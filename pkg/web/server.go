package web

import (
	"crypto/tls"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	pkg_db "recipe_book/pkg/db"
)

type Application struct {
	accessLog *log.Logger
	traceLog  *log.Logger
	InfoLog   *log.Logger
	ErrorLog  *log.Logger

	Middlewares []Middleware

	DB pkg_db.DB
}

func NewApp(db pkg_db.DB) Application {
	ret := Application{
		accessLog: log.New(os.Stdout, "ACCESS\t", log.Ldate|log.Ltime),
		traceLog:  log.New(os.Stdout, "TRACE\t", log.Ldate|log.Ltime),
		InfoLog:   log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
		ErrorLog:  log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),

		DB: db,
	}
	ret.Middlewares = []Middleware{
		secureHeaders,
		ret.logRequest,
		ret.recoverPanic,
	}
	return ret
}

// Manual router implementation.
// I don't like the weird matching behavior of http.ServeMux, and it's not hard to write by hand.
func (app *Application) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	app.traceLog.Printf("base handler: %s", r.URL.Path)
	parts := strings.Split(r.URL.Path, "/")[1:]
	switch parts[0] {
	case "static":
		http.StripPrefix("/static", http.HandlerFunc(app.ServeStatic)).ServeHTTP(w, r)
	case "ingredients":
		http.StripPrefix("/ingredients", http.HandlerFunc(app.Ingredients)).ServeHTTP(w, r)
	case "recipes":
		http.StripPrefix("/recipes", http.HandlerFunc(app.Recipes)).ServeHTTP(w, r)
	default:
		app.error_404(w, r)
		return
	}
}

func (app *Application) Run(address string, should_auto_open bool) {
	srv := &http.Server{
		Addr:     address,
		ErrorLog: app.ErrorLog,
		Handler:  app.WithMiddlewares(),
		TLSConfig: &tls.Config{
			CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
		},
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	app.InfoLog.Printf("Starting server on %s", address)

	if should_auto_open {
		go func(url string) {
			var cmd *exec.Cmd
			switch runtime.GOOS {
			case "darwin": // macOS
				cmd = exec.Command("open", url)
			case "windows":
				cmd = exec.Command("cmd", "/c", "start", url)
			default: // Linux and others
				cmd = exec.Command("xdg-open", url)
			}
			if err := cmd.Run(); err != nil {
				log.Printf("Failed to open homepage: %s", err.Error())
			}
		}("http://" + address)
	}
	err := srv.ListenAndServe()
	app.ErrorLog.Fatal(err)
}
