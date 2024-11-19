package web

import (
	"fmt"
	"net/http"
)

func panic_if(err error) {
	if err != nil {
		panic(err)
	}
}

func (app *Application) error_400(w http.ResponseWriter, r *http.Request, msg string) {
	http.Error(w, fmt.Sprintf("Bad Request\n\n%s", msg), 400)
}

func (app *Application) error_404(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not Found", 404)
}

func (app *Application) error_500(w http.ResponseWriter, r *http.Request, err error) {
	panic("TODO")
}

