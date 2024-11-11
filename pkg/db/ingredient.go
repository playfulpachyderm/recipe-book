package db

import (
	"fmt"
)

type IngredientID uint64

type Ingredient struct {
	ID       IngredientID `db:"rowid"`
	FoodID   FoodID       `db:"food_id"`
	RecipeID RecipeID     `db:"recipe_id"`

	Quantity float32
	Units    UnitsID `db:"units"`

	InRecipeID RecipeID `db:"in_recipe_id"`
	ListOrder  int64    `db:"list_order"`
	IsHidden   bool     `db:"is_hidden"`

	Food *Food
}

// // Format as string
// func (i Ingredient) String() string {
// 	return fmt.Sprintf("%s(%d)", f.Name, f.ID)
// }

func (db *DB) SaveIngredient(i *Ingredient) {
	if i.ID == IngredientID(0) {
		// Do create
		result, err := db.DB.NamedExec(`
			insert into ingredients
			            (food_id, recipe_id, quantity, units, in_recipe_id, list_order, is_hidden)
			     values (nullif(:food_id, 0), nullif(:recipe_id, 0), :quantity, :units, :in_recipe_id,
			                :list_order, :is_hidden)
		`, i)
		if err != nil {
			panic(err)
		}

		// Update the ID
		id, err := result.LastInsertId()
		if err != nil {
			panic(err)
		}
		i.ID = IngredientID(id)
	} else {
		// Do update
		result, err := db.DB.NamedExec(`
			update ingredients
			   set food_id=nullif(:food_id, 0),
			       recipe_id=nullif(:recipe_id, 0),
			       quantity=:quantity,
			       units=:units,
			       list_order=:list_order,
			       is_hidden=:is_hidden
			 where rowid = :rowid
		`, i)
		if err != nil {
			panic(err)
		}
		count, err := result.RowsAffected()
		if err != nil {
			panic(err)
		}
		if count != 1 {
			panic(fmt.Errorf("Got ingredient with ID (%d), so attempted update, but it doesn't exist", i.ID))
		}
	}
}

func (db *DB) DeleteIngredient(i Ingredient) {
	result, err := db.DB.Exec(`delete from ingredients where rowid = ?`, i.ID)
	if err != nil {
		panic(err)
	}
	count, err := result.RowsAffected()
	if err != nil {
		panic(err)
	}
	if count != 1 {
		panic(fmt.Errorf("tried to delete ingredient with ID (%d) but it doesn't exist", i.ID))
	}
}

