package db_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	. "recipe_book/pkg/db"
)

func TestUnitsNamesAndAbbrevs(t *testing.T) {
	assert := assert.New(t)
	assert.Equal("count", COUNT.Name())
	assert.Equal("cups", CUPS.Name())
	assert.Equal("pounds", LBS.Name())

	assert.Equal("tsp", TSP.Abbreviation())
	assert.Equal("fl-oz", FLOZ.Abbreviation())
}

func TestUnitsOfFoods(t *testing.T) {
	assert := assert.New(t)
	db := get_test_db()

	avocado := get_food(db, 20)
	onion := get_food(db, 28)
	eggs := get_food(db, 8)
	oats := get_food(db, 2)

	ingr := COUNT.Of(avocado, 3)
	assert.Equal(ingr.FoodID, FoodID(20))
	assert.Equal(ingr.Quantity, float32(3.0))
	assert.Equal(ingr.Units, COUNT)

	ingr = GRAMS.Of(avocado, 50)
	assert.Equal(ingr.Quantity, float32(0.5))
	assert.Equal(ingr.Units, GRAMS)

	assert.Equal(GRAMS.Of(eggs, 106).Quantity, float32(2))
	assert.InDelta(LBS.Of(onion, 2).Quantity, float32(4.127), 0.001)
	assert.InDelta(OZ.Of(onion, 2).Quantity, float32(0.254), 0.001)
	assert.InDelta(ML.Of(onion, 200).Quantity, float32(0.909), 0.001)
	assert.InDelta(ML.Of(oats, 200).Quantity, float32(0.88), 0.001)
	assert.InDelta(CUPS.Of(oats, 2).Quantity, float32(2.2), 0.001)
	assert.InDelta(TSP.Of(avocado, 1.5).Quantity, float32(0.075), 0.001)
	assert.InDelta(TBSP.Of(avocado, 1.5).Quantity, float32(0.225), 0.001)
}
