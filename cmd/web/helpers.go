package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"
)

func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, statusCode int) {
	http.Error(w, http.StatusText(statusCode), statusCode)
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *application) render(w http.ResponseWriter, r *http.Request, name string, data *templateData) {
	ts, ok := app.templateCache[name]
	if !ok {
		app.serverError(w, fmt.Errorf("template %s does not exist", name))
		return
	}

	buf := new(bytes.Buffer)

	app.addDefaultData(data, r)

	err := ts.Execute(buf, data)
	if err != nil {
		app.serverError(w, err)
		return
	}

	buf.WriteTo(w)
}

func (app *application) addDefaultData(td *templateData, r *http.Request) {
	if td == nil {
		td = newTemplateData()
	}

	td.CurrentYear = time.Now().Year()
}

func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}
