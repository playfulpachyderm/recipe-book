package db

type UnitsID uint64

type Units struct {
	ID           UnitsID `db:"rowid"`
	Name         string  `db:"name"`
	Abbreviation string  `db:"abbreviation"`
}
