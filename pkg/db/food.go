package db

import (
	"fmt"
)

type FoodID uint64

type Food struct {
	ID   FoodID `db:"rowid"`
	Name string `db:"name"`

	Cals    float32 `db:"cals"`
	Carbs   float32 `db:"carbs"`
	Protein float32 `db:"protein"`
	Fat     float32 `db:"fat"`
	Sugar   float32 `db:"sugar"`
	Alcohol float32 `db:"alcohol"`

	Water float32 `db:"water"`

	Potassium  float32 `db:"potassium"`
	Calcium    float32 `db:"calcium"`
	Sodium     float32 `db:"sodium"`
	Magnesium  float32 `db:"magnesium"`
	Phosphorus float32 `db:"phosphorus"`
	Iron       float32 `db:"iron"`
	Zinc       float32 `db:"zinc"`

	Mass      float32 `db:"mass"` // In grams
	Price     float32 `db:"price"`
	Density   float32 `db:"density"`
	CookRatio float32 `db:"cook_ratio"`
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
