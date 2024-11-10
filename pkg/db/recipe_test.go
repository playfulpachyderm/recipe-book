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
