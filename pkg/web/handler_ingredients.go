package web

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/a-h/templ"
	"io"
	"net/http"
	"strconv"
	"strings"

	. "recipe_book/pkg/db"

	"recipe_book/pkg/web/tpl/pages"
)

// Router: `/ingredients`
func (app *Application) Ingredients(w http.ResponseWriter, r *http.Request) {
	app.traceLog.Printf("Ingredient: %s", r.URL.Path)
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")

	if parts[0] == "" {
		// Index page
		switch r.Method {
		case "GET":
			app.IngredientsIndex(w, r)
		case "POST":
			app.IngredientCreate(w, r)
		}
		return
	}

	// Details page
	food_id, err := strconv.Atoi(parts[0])
	if err != nil {
		app.error_400(w, r, fmt.Sprintf("invalid ID: %s", parts[0]))
		return
	}
	food, err := app.DB.GetFoodByID(FoodID(food_id))
	if errors.Is(err, ErrNotInDB) {
		app.error_404(w, r)
		return
	} else if err != nil {
		panic(err)
	}
	app.IngredientDetail(food, w, r)
}

// Handler: `GET /ingredients`
func (app *Application) IngredientsIndex(w http.ResponseWriter, r *http.Request) {
	foods := app.DB.GetAllBaseFoods()
	err := pages.Base("Ingredients").Render(
		templ.WithChildren(
			context.Background(),
			pages.IngredientsIndex(foods),
		),
		w)
	panic_if(err)
}

// Handler: `POST /ingredients`
func (app *Application) IngredientCreate(w http.ResponseWriter, r *http.Request) {
	var food Food
	data, err := io.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, &food)
	if err != nil {
		app.ErrorLog.Print(err)
		panic(err)
	}
	app.DB.SaveFood(&food)

	http.Redirect(w, r, fmt.Sprintf("/ingredients/%d", food.ID), 303)
}

// Handler: `GET /ingredients/:id`
func (app *Application) IngredientDetail(food Food, w http.ResponseWriter, r *http.Request) {
	// If it's a POST request, update the food
	if r.Method == "POST" {
		data, err := io.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}
		err = json.Unmarshal(data, &food)
		if err != nil {
			app.ErrorLog.Print(err)
			panic(err)
		}
		// Save the updated food
		app.DB.SaveFood(&food)
		app.traceLog.Printf("POST Ingredient Detail: %#v", food)
	}

	err := pages.Base(fmt.Sprintf("Ingredient: %s", food.Name)).Render(
		templ.WithChildren(
			context.Background(),
			pages.IngredientDetail(food),
		),
		w)
	panic_if(err)
}
