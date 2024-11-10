package db_test

import (
	"testing"

	"github.com/go-test/deep"
	"github.com/stretchr/testify/assert"

	. "recipe_book/pkg/db"
)

func TestFoodSaveAndLoad(t *testing.T) {
	assert := assert.New(t)
	db := get_test_db()
	food := Food{
		Name:       "some food",
		Cals:       1.0,
		Carbs:      2.0,
		Protein:    3.0,
		Fat:        4.0,
		Sugar:      5.0,
		Alcohol:    6.0,
		Water:      7.0,
		Potassium:  8.0,
		Calcium:    9.0,
		Sodium:     10.0,
		Magnesium:  11.0,
		Phosphorus: 12.0,
		Iron:       13.0,
		Zinc:       14.0,
		Mass:       15.0,
		Price:      16.0,
		Density:    17.0,
		CookRatio:  18.0,
	}
	assert.Equal(food.ID, FoodID(0))
	db.SaveFood(&food)
	assert.NotEqual(food.ID, FoodID(0))
	new_food, err := db.GetFoodByID(food.ID)
	assert.NoError(err)
	if diff := deep.Equal(food, new_food); diff != nil {
		t.Error(diff)
	}

	// Modify it
	food.Name = "another food"
	food.Cals = food.Cals + 9.2
	food.Carbs = food.Carbs + 9.2
	food.Protein = food.Protein + 9.2
	food.Fat = food.Fat + 9.2
	food.Sugar = food.Sugar + 9.2
	food.Alcohol = food.Alcohol + 9.2
	food.Water = food.Water + 9.2
	food.Potassium = food.Potassium + 9.2
	food.Calcium = food.Calcium + 9.2
	food.Sodium = food.Sodium + 9.2
	food.Phosphorus = food.Phosphorus + 9.2
	food.Iron = food.Iron + 9.2
	food.Zinc = food.Zinc + 9.2
	food.Mass = food.Mass + 9.2
	food.Price = food.Price + 9.2
	food.Density = food.Density + 9.2
	food.CookRatio = food.CookRatio + 9.2

	// Save it and reload it
	db.SaveFood(&food)
	new_food, err = db.GetFoodByID(food.ID)
	assert.NoError(err)
	if diff := deep.Equal(food, new_food); diff != nil {
		t.Error(diff)
	}
}
