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

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(app.recoverPanic)
	//r.Use(app.logRequest)
	r.Use(secureHeaders)
	r.Use(app.session.Enable)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/", app.home)
	r.Group(func(r chi.Router) {
		r.Route("/snippet", func(r chi.Router) {
			r.Get("/{id}", app.showSnippet)
			r.Get("/create", app.createSnippetForm)
			r.Post("/create", app.createSnippet)
		})
	})

	filesDir := http.Dir("./ui/static")
	fileServer(r, "/static", filesDir)

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
