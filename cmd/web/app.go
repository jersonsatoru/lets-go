package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/jersonsatoru/lets-go/internal/models"
)

type application struct {
	errorLogger   *log.Logger
	infoLogger    *log.Logger
	snippetModel  *models.SnippetModel
	templateCache map[string]*template.Template
}

func NewApplication(infoLogger *log.Logger, errorLogger *log.Logger, snippetModel *models.SnippetModel, templateCache map[string]*template.Template) *application {
	return &application{
		errorLogger:   errorLogger,
		infoLogger:    infoLogger,
		snippetModel:  snippetModel,
		templateCache: templateCache,
	}
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLogger.Output(2, trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) render(w http.ResponseWriter, status int, templateName string, templateParameters *TemplateData) {
	buf := new(bytes.Buffer)

	err := app.templateCache[templateName].ExecuteTemplate(buf, "base", templateParameters)

	if err != nil {
		app.serverError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	buf.WriteTo(w)
}

func (app *application) newTemplateData(r *http.Request) *TemplateData {
	return &TemplateData{
		Year: time.Now().Year(),
	}
}
