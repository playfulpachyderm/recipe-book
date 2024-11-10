package db_test

import (
	"testing"

	"github.com/go-test/deep"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	. "recipe_book/pkg/db"
)

func TestSaveAndLoadIngredient(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	db := get_test_db()

	// Setup
	recipe := Recipe{
		Name:  "some Recipe",
		Blurb: "Lorem Ispum dolor sit amet consquiter id blah blabh albha blahbla blahblahblh",
		Instructions: RecipeInstructions{
			"instr 1", "isntr 2", "instr3", "ins32gjkifw",
		},
	}
	db.SaveRecipe(&recipe)
	food := Food{
		Name:    "a food",
		Cals:    10,
		Carbs:   1,
		Protein: 2,
	}
	db.SaveFood(&food)

	// Create an ingredient on the recipe
	ingr := Ingredient{
		FoodID:              food.ID,
		QuantityNumerator:   3,
		QuantityDenominator: 2,
		Units:               1, // count
		InRecipeID:          recipe.ID,
		ListOrder:           0,
	}
	assert.Equal(ingr.ID, IngredientID(0))
	db.SaveIngredient(&ingr)
	assert.NotEqual(ingr.ID, IngredientID(0))

	// It should be added to the recipe at position 0
	new_recipe, err := db.GetRecipeByID(recipe.ID)
	assert.NoError(err)
	require.Len(new_recipe.Ingredients, 1)
	new_ingr := new_recipe.Ingredients[0]
	if diff := deep.Equal(ingr, new_ingr); diff != nil {
		t.Error(diff)
	}

	// Modify the ingredient
	ingr.QuantityNumerator = 5
	ingr.QuantityDenominator = 4
	ingr.Units = 2

	// Save it and reload the recipe
	db.SaveIngredient(&ingr)
	new_recipe, err = db.GetRecipeByID(recipe.ID)
	assert.NoError(err)
	require.Len(new_recipe.Ingredients, 1)
	new_ingr = new_recipe.Ingredients[0]
	if diff := deep.Equal(ingr, new_ingr); diff != nil {
		t.Error(diff)
	}

	// Delete and reload-- should be gone
	db.DeleteIngredient(ingr)
	new_recipe, err = db.GetRecipeByID(recipe.ID)
	assert.NoError(err)
	require.Len(new_recipe.Ingredients, 0)
}