package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"savemanager/config"
	"savemanager/router"
	"strings"

	"github.com/go-chi/chi/v5"
)

func main() {
	println("Starting server...")
	config.GetConfig()
	mux, err := router.InitRouter()
	if err != nil {
		log.Fatal(err.Error())
	}
	workDir, _ := os.Getwd()
	// mux.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	root := http.Dir(filepath.Join(workDir, "static"))
	FileServer(mux, "/static/", root)
	http.ListenAndServe(":5263", mux)
}

func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
	println("Save Manager online !")
}
