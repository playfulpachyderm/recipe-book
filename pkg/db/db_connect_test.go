package db_test

import (
	"errors"
	"fmt"
	"io"
	"math/rand"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	. "recipe_book/pkg/db"
)

var seed_sql string

func init() {
	file, err := os.Open("../../sample_data/seed.sql")
	if err != nil {
		panic(err)
	}
	data, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}
	seed_sql = string(data)
}

func get_test_db() DB {
	db_path := "../../sample_data/data/test.db"
	db, err := DBCreate(db_path)
	if errors.Is(err, ErrTargetExists) {
		db, err = DBConnect(db_path)
	} else if err == nil {
		db.DB.MustExec(seed_sql)
	}
	if err != nil {
		panic(err)
	}
	return db
}

func get_food(db DB, id FoodID) Food {
	ret, err := db.GetFoodByID(id)
	if err != nil {
		panic(err)
	}
	return ret
}

func TestCreateAndConnectToDB(t *testing.T) {
	i := rand.Uint32()
	_, err := DBCreate(fmt.Sprintf("../../sample_data/data/random-%d.db", i))
	assert.NoError(t, err)

	_, err = DBConnect(fmt.Sprintf("../../sample_data/data/random-%d.db", i))
	assert.NoError(t, err)
}
