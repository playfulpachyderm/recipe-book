package db

import (
	"database/sql"
	"errors"
	"fmt"
)

type FoodID uint64

type Food struct {
	ID   FoodID `db:"rowid" json:"id"`
	Name string `db:"name" json:"name"`

	Cals    float32 `db:"cals" json:"cals,string"`
	Carbs   float32 `db:"carbs" json:"carbs,string"`
	Protein float32 `db:"protein" json:"protein,string"`
	Fat     float32 `db:"fat" json:"fat,string"`
	Sugar   float32 `db:"sugar" json:"sugar,string"`
	Alcohol float32 `db:"alcohol" json:"alcohol,string"`

	Water float32 `db:"water" json:"water,string"`

	Potassium  float32 `db:"potassium" json:"potassium,string"`
	Calcium    float32 `db:"calcium" json:"calcium,string"`
	Sodium     float32 `db:"sodium" json:"sodium,string"`
	Magnesium  float32 `db:"magnesium" json:"magnesium,string"`
	Phosphorus float32 `db:"phosphorus" json:"phosphorus,string"`
	Iron       float32 `db:"iron" json:"iron,string"`
	Zinc       float32 `db:"zinc" json:"zinc,string"`

	Mass      float32 `db:"mass" json:"mass,string"` // In grams
	Price     float32 `db:"price" json:"price,string"`
	Density   float32 `db:"density" json:"density,string"`
	CookRatio float32 `db:"cook_ratio" json:"cook_ratio,string"`
}

// Format as string
func (f Food) String() string {
	return fmt.Sprintf("%s(%d)", f.Name, f.ID)
}

func (db *DB) SaveFood(f *Food) {
	if f.ID == FoodID(0) {
		// Do create
		result, err := db.DB.NamedExec(`
			insert into foods (name, cals, carbs, protein, fat, sugar, alcohol, water, potassium, calcium, sodium,
			                   magnesium, phosphorus, iron, zinc, mass, price, density, cook_ratio)
			           values (:name, :cals, :carbs, :protein, :fat, :sugar, :alcohol, :water, :potassium, :calcium,
			                   :sodium, :magnesium, :phosphorus, :iron, :zinc, :mass, :price, :density, :cook_ratio)
		`, f)
		if err != nil {
			panic(err)
		}
		// Update the ID if necessary
		id, err := result.LastInsertId()
		if err != nil {
			panic(err)
		}
		f.ID = FoodID(id)
	} else {
		// Do update
		result, err := db.DB.NamedExec(`
			update foods
	           set name=:name,
	               cals=:cals,
	               carbs=:carbs,
	               protein=:protein,
	               fat=:fat,
	               sugar=:sugar,
	               alcohol=:alcohol,
	               water=:water,
	               potassium=:potassium,
	               calcium=:calcium,
	               sodium=:sodium,
	               magnesium=:magnesium,
	               phosphorus=:phosphorus,
	               iron=:iron,
	               zinc=:zinc,
	               mass=:mass,
	               price=:price,
	               density=:density,
	               cook_ratio=:cook_ratio
	         where rowid = :rowid
		`, f)
		if err != nil {
			panic(err)
		}
		count, err := result.RowsAffected()
		if err != nil {
			panic(err)
		}
		if count != 1 {
			panic(fmt.Errorf("Got food with ID (%d), so attempted update, but it doesn't exist", f.ID))
		}
	}
}

func (db *DB) GetFoodByID(id FoodID) (ret Food, err error) {
	err = db.DB.Get(&ret, `
		select rowid, name, cals, carbs, protein, fat, sugar, alcohol, water, potassium, calcium, sodium,
			   magnesium, phosphorus, iron, zinc, mass, price, density, cook_ratio
		  from foods
		 where rowid = ?
	`, id)
	if errors.Is(err, sql.ErrNoRows) {
		return Food{}, ErrNotInDB
	}
	return
}

func (db *DB) GetAllBaseFoods() []Food {
	var ret []Food
	err := db.DB.Select(&ret, `
		select rowid, name, cals, carbs, protein, fat, sugar, alcohol, water, potassium, calcium, sodium,
			   magnesium, phosphorus, iron, zinc, mass, price, density, cook_ratio
		  from foods
		 where rowid not in (select computed_food_id from recipes)
	`)
	if err != nil {
		panic(err)
	}
	return ret
}
