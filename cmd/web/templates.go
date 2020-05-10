package main

import (
	"github.com/balabanovds/go-snippetbox/pkg/forms"
	"html/template"
	"path/filepath"

	"github.com/balabanovds/go-snippetbox/pkg/models"
)

type templateData struct {
	CSRFToken     string
	CurrentYear   int
	Authenticated bool
	Name          string
	Flash         string
	Form          *forms.Form
	Snippet       *models.Snippet
	Snippets      []*models.Snippet
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateData() *templateData {
	return &templateData{}
}

func newTemplateCache(dir string) (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob(filepath.Join(dir, "*.page.tmpl"))
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob(filepath.Join(dir, "*.layout.tmpl"))
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob(filepath.Join(dir, "*.partial.tmpl"))
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}
