package db_test

import (
	"errors"
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"

	. "recipe_book/pkg/db"
)

func get_test_db() DB {
	db_path := "../../sample_data/data/test.db"
	db, err := DBCreate(db_path)
	if errors.Is(err, ErrTargetExists) {
		db, err = DBConnect(db_path)
	}
	if err != nil {
		panic(err)
	}
	return db
}

func TestCreateAndConnectToDB(t *testing.T) {
	i := rand.Uint32()
	_, err := DBCreate(fmt.Sprintf("../../sample_data/data/random-%d.db", i))
	assert.NoError(t, err)

	_, err = DBConnect(fmt.Sprintf("../../sample_data/data/random-%d.db", i))
	assert.NoError(t, err)
}
