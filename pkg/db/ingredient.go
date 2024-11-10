package db

import (
// "fmt"
)

type IngredientID uint64

type Ingredient struct {
	ID       IngredientID `db:"rowid"`
	FoodID   FoodID       `db:"food_id"`
	RecipeID RecipeID     `db:"recipe_id"`

	QuantityNumerator   int64 `db:"quantity_numerator"`
	QuantityDenominator int64 `db:"quantity_denominator"`
	Units               Units `db:"units"`

	InRecipeID RecipeID `db:"in_recipe_id"`
	ListOrder  int64    `db:"list_order"`
	IsHidden   bool     `db:"is_hidden"`
}

// // Format as string
// func (i Ingredient) String() string {
// 	return fmt.Sprintf("%s(%d)", f.Name, f.ID)
// }
