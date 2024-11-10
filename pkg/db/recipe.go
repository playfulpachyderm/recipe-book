package db

import (
	"database/sql/driver"
	"fmt"
	"strings"
)

type RecipeID uint64

type RecipeInstructions []string

// Join the instructions with 0x1F, the "Unit Separator" ASCII character
func (ri RecipeInstructions) Value() (driver.Value, error) {
	return strings.Join(ri, "\x1F"), nil
}

// Split the stored string by "Unit Separator" characters
func (ri *RecipeInstructions) Scan(src interface{}) error {
	val, is_ok := src.(string)
	if !is_ok {
		return fmt.Errorf("incompatible type for RecipeInstructions list: %#v", src)
	}
	*ri = RecipeInstructions(strings.Split(val, "\x1F"))
	return nil
}

type Recipe struct {
	ID           RecipeID           `db:"rowid"`
	Name         string             `db:"name"`
	Blurb        string             `db:"blurb"`
	Instructions RecipeInstructions `db:"instructions"`

	Ingredients []Ingredient

	ComputedFoodID FoodID `db:"computed_food_id"`
}

func (db *DB) SaveRecipe(r *Recipe) {
	if r.ID == RecipeID(0) {
		// Do create
		result, err := db.DB.NamedExec(`
			insert into recipes (name, blurb, instructions)
	                     values (:name, :blurb, :instructions)
	                on conflict do update
	                        set name=:name,
	                            blurb=:blurb,
	                            instructions=:instructions
		`, r)
		if err != nil {
			panic(err)
		}

		// Update the ID
		id, err := result.LastInsertId()
		if err != nil {
			panic(err)
		}
		r.ID = RecipeID(id)
	} else {
		// Do update
		result, err := db.DB.NamedExec(`
			update recipes set name=:name, blurb=:blurb, instructions=:instructions where rowid = :rowid
		`, r)
		if err != nil {
			panic(err)
		}
		count, err := result.RowsAffected()
		if err != nil {
			panic(err)
		}
		if count != 1 {
			panic(fmt.Errorf("Got recipe with ID (%d), so attempted update, but it doesn't exist", r.ID))
		}
	}
	// TODO: recompute the computed_food
}

func (db *DB) GetRecipeByID(id RecipeID) (ret Recipe, err error) {
	err = db.DB.Get(&ret, `
		select rowid, name, blurb, instructions
	      from recipes
	     where rowid = ?
	`, id)
	if err != nil {
		return Recipe{}, fmt.Errorf("fetching recipe with ID %d: %w", id, err)
	}

	// Load the ingredients
	err = db.DB.Select(&ret.Ingredients, `
		select rowid, ifnull(food_id, 0) food_id, ifnull(recipe_id, 0) recipe_id, quantity_numerator, quantity_denominator, units,
		       in_recipe_id, list_order, is_hidden
		  from ingredients
	     where in_recipe_id = ?
	     order by list_order asc
    `, id)
	return
}