package main

import (
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func (app *application) routes() http.Handler {

	r := chi.NewRouter()

	// r.Use(app.recoverPanic, app.logRequest, secureHeaders)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)
	r.Use(secureHeaders)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/", app.home)
	r.Route("/snippet", func(r chi.Router) {
		r.Get("/{id}", app.showSnippet)
		r.Get("/create", app.createSnippet)
		r.Post("/create", app.createSnippetForm)
	})

	filesDir := http.Dir("./ui/static")
	fileServer(r, "/static", filesDir)
	// fileServer := http.FileServer(http.Dir("./ui/static"))
	// r.Get("/static/", func(w http.ResponseWriter, r *http.Request) {
	// 	fs := http.StripPrefix("/static", fileServer)
	// 	fs.ServeHTTP(w, r)
	// })

	return r
}

func fileServer(r chi.Router, path string, root http.FileSystem) {
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
}
