package web

import (
	"context"
	"errors"
	"fmt"
	"github.com/a-h/templ"
	"net/http"
	"strconv"
	"strings"

	. "recipe_book/pkg/db"

	"recipe_book/pkg/web/tpl/pages"
)

// Router: `/ingredients`
func (app *Application) Recipes(w http.ResponseWriter, r *http.Request) {
	app.traceLog.Printf("Recipe: %s", r.URL.Path)
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")

	if parts[0] == "" {
		// Index page
		switch r.Method {
		case "GET":
			app.RecipesIndex(w, r)
		}
		return
	}

	// Details page
	recipe_id, err := strconv.Atoi(parts[0])
	if err != nil {
		app.error_400(w, r, fmt.Sprintf("invalid ID: %s", parts[0]))
		return
	}
	recipe, err := app.DB.GetRecipeByID(RecipeID(recipe_id))
	if errors.Is(err, ErrNotInDB) {
		app.error_404(w, r)
		return
	} else if err != nil {
		panic(err)
	}
	app.RecipeDetail(recipe, w, r)
}

// Handler: `GET /recipes`
func (app *Application) RecipesIndex(w http.ResponseWriter, r *http.Request) {
	foods := app.DB.GetAllRecipes()
	err := pages.Base("Ingredients").Render(
		templ.WithChildren(
			context.Background(),
			pages.RecipesIndex(foods),
		),
		w)
	panic_if(err)
}

// // Handler: `POST /ingredients`
// func (app *Application) IngredientCreate(w http.ResponseWriter, r *http.Request) {
// 	var food Food
// 	data, err := io.ReadAll(r.Body)
// 	if err != nil {
// 		panic(err)
// 	}
// 	err = json.Unmarshal(data, &food)
// 	if err != nil {
// 		app.ErrorLog.Print(err)
// 		panic(err)
// 	}
// 	app.DB.SaveFood(&food)

// 	http.Redirect(w, r, fmt.Sprintf("/ingredients/%d", food.ID), 303)
// }

// Handler: `GET /ingredients/:id`
func (app *Application) RecipeDetail(recipe Recipe, w http.ResponseWriter, r *http.Request) {
	err := pages.Base(recipe.Name).Render(
		templ.WithChildren(
			context.Background(),
			templ.Join(pages.RecipeDetail(recipe)),
		),
		w)
	panic_if(err)
}
