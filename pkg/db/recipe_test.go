package db_test

import (
	"testing"

	"github.com/go-test/deep"
	"github.com/stretchr/testify/assert"

	. "recipe_book/pkg/db"
)

func TestRecipeSaveAndLoad(t *testing.T) {
	assert := assert.New(t)
	db := get_test_db()

	recipe := Recipe{
		Name:  "some Recipe",
		Blurb: "Lorem Ispum dolor sit amet consquiter id blah blabh albha blahbla blahblahblh",
		Instructions: RecipeInstructions{
			"instr 1", "isntr 2", "instr3", "ins32gjkifw",
		},
	}
	assert.Equal(recipe.ID, RecipeID(0))
	db.SaveRecipe(&recipe)
	assert.NotEqual(recipe.ID, RecipeID(0))
	new_recipe, err := db.GetRecipeByID(recipe.ID)
	assert.NoError(err)
	if diff := deep.Equal(recipe, new_recipe); diff != nil {
		t.Error(diff)
	}

	// Modify it
	recipe.Name = "some recipe 2"
	recipe.Blurb = "another blurb"
	recipe.Instructions = RecipeInstructions{"i1", "i2", "i3"}

	// Save it and reload
	db.SaveRecipe(&recipe)
	new_recipe, err = db.GetRecipeByID(recipe.ID)
	assert.NoError(err)
	if diff := deep.Equal(recipe, new_recipe); diff != nil {
		t.Error(diff)
	}
}

func TestRecipeComputeFood(t *testing.T) {
	assert := assert.New(t)
	f1 := Food{0, "", 1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0, 11.0, 12.0, 13.0, 14.0, 15.0, 16.0, 0, 0}
	f2 := Food{0, "", 16.5, 15.5, 14.5, 13.5, 12.5, 11.5, 10.5, 9.5, 8.5, 7.5, 6.5, 5.5, 4.5, 3.5, 2.5, 1.5, 0, 0}

	recipe := Recipe{Ingredients: []Ingredient{
		{Quantity: 1, Food: &f1},
		{Quantity: 1, Food: &f2},
	}}
	computed_food := recipe.ComputeFood()
	assert.Equal(computed_food.Cals, float32(17.5))
	assert.Equal(computed_food.Carbs, float32(17.5))
	assert.Equal(computed_food.Protein, float32(17.5))
	assert.Equal(computed_food.Fat, float32(17.5))
	assert.Equal(computed_food.Sugar, float32(17.5))
	assert.Equal(computed_food.Alcohol, float32(17.5))
	assert.Equal(computed_food.Water, float32(17.5))
	assert.Equal(computed_food.Potassium, float32(17.5))
	assert.Equal(computed_food.Calcium, float32(17.5))
	assert.Equal(computed_food.Sodium, float32(17.5))
	assert.Equal(computed_food.Magnesium, float32(17.5))
	assert.Equal(computed_food.Phosphorus, float32(17.5))
	assert.Equal(computed_food.Iron, float32(17.5))
	assert.Equal(computed_food.Zinc, float32(17.5))
	assert.Equal(computed_food.Mass, float32(17.5))
	assert.Equal(computed_food.Price, float32(17.5))

	recipe2 := Recipe{Ingredients: []Ingredient{
		{Quantity: 1.5, Food: &f1},
		{Quantity: 0.5, Food: &f2},
	}}
	computed_food2 := recipe2.ComputeFood()
	assert.Equal(computed_food2.Cals, float32(9.75))
	assert.Equal(computed_food2.Carbs, float32(10.75))
	assert.Equal(computed_food2.Protein, float32(11.75))
	assert.Equal(computed_food2.Fat, float32(12.75))
	assert.Equal(computed_food2.Sugar, float32(13.75))
	assert.Equal(computed_food2.Alcohol, float32(14.75))
	assert.Equal(computed_food2.Water, float32(15.75))
	assert.Equal(computed_food2.Potassium, float32(16.75))
	assert.Equal(computed_food2.Calcium, float32(17.75))
	assert.Equal(computed_food2.Sodium, float32(18.75))
	assert.Equal(computed_food2.Magnesium, float32(19.75))
	assert.Equal(computed_food2.Phosphorus, float32(20.75))
	assert.Equal(computed_food2.Iron, float32(21.75))
	assert.Equal(computed_food2.Zinc, float32(22.75))
	assert.Equal(computed_food2.Mass, float32(23.75))
	assert.Equal(computed_food2.Price, float32(24.75))
}
