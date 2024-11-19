package web

import (
	"fmt"
	"net/http"
	"time"
)

// Middlewares are wrappers around `http.Handler`s
type Middleware func(http.Handler) http.Handler

func (app *Application) WithMiddlewares() http.Handler {
	var ret http.Handler = app
	for i := range app.Middlewares {
		ret = app.Middlewares[i](ret)
	}
	return ret
}

func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Referrer-Policy", "same-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")

		next.ServeHTTP(w, r)
	})
}

func (app *Application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Since(t)

		app.accessLog.Printf("%s - %s %s %s\t%s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI(), duration)
	})
}

func (app *Application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				app.error_500(w, r, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}
