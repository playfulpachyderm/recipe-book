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

// Save the recipe.  New recipes will have their ID back-filled from the DB.
//
// Automatically updates the computed food, creating one if it's a new recipe and back-filling the
// `computed_food_id` foreign key to it.
func (db *DB) SaveRecipe(r *Recipe) {
	if r.ID == RecipeID(0) {
		// Do create

		// Create the computed food
		computed_food := Food{Name: r.Name}
		db.SaveFood(&computed_food)
		r.ComputedFoodID = computed_food.ID

		// Create the recipe
		result, err := db.DB.NamedExec(`
			insert into recipes (name, blurb, instructions, computed_food_id)
	                     values (:name, :blurb, :instructions, :computed_food_id)
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
	for i := range r.Ingredients {
		r.Ingredients[i].InRecipeID = r.ID
		db.SaveIngredient(&r.Ingredients[i])
	}
	// Update the computed food
	computed_food := r.ComputeFood()
	db.SaveFood(&computed_food)
}

func (db *DB) GetRecipeByID(id RecipeID) (ret Recipe, err error) {
	err = db.DB.Get(&ret, `
		select rowid, name, blurb, instructions, computed_food_id
	      from recipes
	     where rowid = ?
	`, id)
	if err != nil {
		return Recipe{}, fmt.Errorf("fetching recipe with ID %d: %w", id, err)
	}

	// Load the ingredients
	err = db.DB.Select(&ret.Ingredients, `
		select rowid, ifnull(food_id, 0) food_id, ifnull(recipe_id, 0) recipe_id, quantity, units,
		       in_recipe_id, list_order, is_hidden
		  from ingredients
	     where in_recipe_id = ?
	     order by list_order asc
    `, id)
	if err != nil {
		panic(err)
	}
	for i := range ret.Ingredients {
		if ret.Ingredients[i].FoodID != FoodID(0) {
			// ingredient is a food
			ret.Ingredients[i].Food, err = db.GetFoodByID(ret.Ingredients[i].FoodID)
		} else {
			// ingredient is a food; i.Food is the ComputedFood of the Recipe
			var computed_food_id FoodID
			err = db.DB.Get(&computed_food_id, `select computed_food_id from recipes where rowid = ?`, ret.Ingredients[i].RecipeID)
			if err != nil {
				panic(err)
			}
			ret.Ingredients[i].Food, err = db.GetFoodByID(computed_food_id)
		}
		if err != nil {
			panic(err)
		}
	}
	return
}

func (r Recipe) ComputeFood() Food {
	// If r.ComputedFoodID is 0, so should be the ID of returned Food
	ret := Food{ID: r.ComputedFoodID, Name: r.Name}
	for _, ingr := range r.Ingredients {
		ret.Cals += ingr.Quantity * ingr.Food.Cals
		ret.Carbs += ingr.Quantity * ingr.Food.Carbs
		ret.Protein += ingr.Quantity * ingr.Food.Protein
		ret.Fat += ingr.Quantity * ingr.Food.Fat
		ret.Sugar += ingr.Quantity * ingr.Food.Sugar
		ret.Alcohol += ingr.Quantity * ingr.Food.Alcohol
		ret.Water += ingr.Quantity * ingr.Food.Water
		ret.Potassium += ingr.Quantity * ingr.Food.Potassium
		ret.Calcium += ingr.Quantity * ingr.Food.Calcium
		ret.Sodium += ingr.Quantity * ingr.Food.Sodium
		ret.Magnesium += ingr.Quantity * ingr.Food.Magnesium
		ret.Phosphorus += ingr.Quantity * ingr.Food.Phosphorus
		ret.Iron += ingr.Quantity * ingr.Food.Iron
		ret.Zinc += ingr.Quantity * ingr.Food.Zinc
		ret.Mass += ingr.Quantity * ingr.Food.Mass
		ret.Price += ingr.Quantity * ingr.Food.Price
	}
	return ret
}
