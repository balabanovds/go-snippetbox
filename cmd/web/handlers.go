package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/balabanovds/go-snippetbox/pkg/models"
	"github.com/go-chi/chi"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {

	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := &templateData{
		Snippets: snippets,
	}

	app.render(w, r, "home.page.html", data)

}

func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	s, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
			return
		}
		app.serverError(w, err)
		return
	}

	data := &templateData{Snippet: s}

	app.render(w, r, "show.page.html", data)
}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")
	expires := r.PostForm.Get("expires")

	errors := make(map[string]string)

	if strings.TrimSpace(title) == "" {
		errors["title"] = "Field can not be blank"
	} else if utf8.RuneCountInString(title) > 100 {
		errors["title"] = "Field can not be more than 100 characters"
	}

	if strings.TrimSpace(content) == "" {
		errors["content"] = "Field can not be blank"
	}

	if strings.TrimSpace(expires) == "" {
		errors["expires"] = "Field can not be blank"
	} else if expires != "365" && expires != "7" && expires != "1" {
		errors["expires"] = "Field is invalid"
	}

	if len(errors) > 0 {
		fmt.Fprint(w, errors)
		return
	}

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
}

func (app *application) createSnippetForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "create.page.html", nil)
}
