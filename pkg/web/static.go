package web

import (
	"embed"
	"net/http"
	"path"
	"runtime"
)

//go:embed "static"
var embedded_files embed.FS

var use_embedded = ""

var this_dir string

func init() {
	_, this_file, _, _ := runtime.Caller(0) // `this_file` is absolute path to this source file
	this_dir = path.Dir(this_file)
}

// Serve static assets, either from the disk (if running in development mode), or from go:embedded files
func (app *Application) ServeStatic(w http.ResponseWriter, r *http.Request) {
	// Static files can be stored in browser cache
	w.Header().Set("Cache-Control", "public, max-age=3600")
	if use_embedded == "true" {
		// Serve directly from the embedded files
		http.FileServer(http.FS(embedded_files)).ServeHTTP(w, r)
	} else {
		// Serve from disk
		http.FileServer(http.Dir(path.Join(this_dir, "static"))).ServeHTTP(w, r)
	}
}
