package main

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	"github.com/jersonsatoru/lets-go/internal/models"
)

type application struct {
	errorLogger    *log.Logger
	infoLogger     *log.Logger
	snippetModel   *models.SnippetModel
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
}

func NewApplication(
	infoLogger *log.Logger, errorLogger *log.Logger, snippetModel *models.SnippetModel,
	templateCache map[string]*template.Template, formDecoder *form.Decoder,
	sessionManager *scs.SessionManager,
) *application {
	return &application{
		errorLogger:    errorLogger,
		infoLogger:     infoLogger,
		snippetModel:   snippetModel,
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
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

func (app *application) logRequests(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		app.infoLogger.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())
		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func (app *application) panicRecover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "Close")
				app.serverError(w, fmt.Errorf("%s", err))
				return
			}

		}()
		next.ServeHTTP(w, r)
	})
}

func (app *application) decodePostForm(r *http.Request, dst any) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}

	err = app.formDecoder.Decode(&dst, r.PostForm)
	if err != nil {
		var invalidDecoderError *form.InvalidDecoderError
		if errors.As(err, &invalidDecoderError) {
			panic(err)
		}

		return err
	}

	return nil
}
