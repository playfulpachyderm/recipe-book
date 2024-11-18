package db

type Units uint64

const (
	COUNT = Units(iota + 1) // Start at 1 to match SQLite ID column
	GRAMS
	LBS
	OZ
	ML
	CUPS
	TSP
	TBSP
	FLOZ
)

var names = []string{"", "count", "grams", "pounds", "ounces", "milliliters", "cups", "teaspoons", "tablespoons", "fluid ounces"}
var abbreviations = []string{"", "", "g", "lbs", "oz", "mL", "cups", "tsp", "tbsp", "fl-oz"}

func (u Units) Name() string {
	return names[u]
}
func (u Units) Abbreviation() string {
	return abbreviations[u]
}

func (u Units) Of(f Food, n float32) Ingredient {
	switch u {
	case COUNT:
		return Ingredient{FoodID: f.ID, Quantity: n, Units: u, Food: f}
	case GRAMS:
		return Ingredient{FoodID: f.ID, Quantity: n / f.Mass, Units: u, Food: f}
	case LBS:
		return Ingredient{FoodID: f.ID, Quantity: n * 454 / f.Mass, Units: u, Food: f}
	case OZ:
		return Ingredient{FoodID: f.ID, Quantity: n * 28 / f.Mass, Units: u, Food: f}
	case ML:
		return Ingredient{FoodID: f.ID, Quantity: n * f.Density / f.Mass, Units: u, Food: f}
	case CUPS:
		return Ingredient{FoodID: f.ID, Quantity: n * f.Density * 250 / f.Mass, Units: u, Food: f}
	case TSP:
		return Ingredient{FoodID: f.ID, Quantity: n * f.Density * 5 / f.Mass, Units: u, Food: f}
	case TBSP:
		return Ingredient{FoodID: f.ID, Quantity: n * f.Density * 15 / f.Mass, Units: u, Food: f}
	case FLOZ:
		return Ingredient{FoodID: f.ID, Quantity: n * f.Density * 30 / f.Mass, Units: u, Food: f}
	default:
		panic(u)
	}
}

func (u Units) Portion(r Recipe, n float32) Ingredient {
	f := r.ComputeFood()
	switch u {
	case COUNT:
		return Ingredient{RecipeID: r.ID, Quantity: n, Units: u, Food: f}
	case GRAMS:
		return Ingredient{RecipeID: r.ID, Quantity: n / f.Mass, Units: u, Food: f}
	case LBS:
		return Ingredient{RecipeID: r.ID, Quantity: n * 454 / f.Mass, Units: u, Food: f}
	case OZ:
		return Ingredient{RecipeID: r.ID, Quantity: n * 28 / f.Mass, Units: u, Food: f}
	case ML:
		return Ingredient{RecipeID: r.ID, Quantity: n * f.Density / f.Mass, Units: u, Food: f}
	case CUPS:
		return Ingredient{RecipeID: r.ID, Quantity: n * f.Density * 250 / f.Mass, Units: u, Food: f}
	case TSP:
		return Ingredient{RecipeID: r.ID, Quantity: n * f.Density * 5 / f.Mass, Units: u, Food: f}
	case TBSP:
		return Ingredient{RecipeID: r.ID, Quantity: n * f.Density * 15 / f.Mass, Units: u, Food: f}
	case FLOZ:
		return Ingredient{RecipeID: r.ID, Quantity: n * f.Density * 30 / f.Mass, Units: u, Food: f}
	default:
		panic(u)
	}
}
