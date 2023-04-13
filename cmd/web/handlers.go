package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"runtime/debug"
	"strconv"

	"github.com/jersonsatoru/lets-go/internal/models"
	"github.com/jersonsatoru/lets-go/internal/validator"
	"github.com/julienschmidt/httprouter"
)

type CreateSnippetForm struct {
	Title               string `form:"title"`
	Content             string `form:"content"`
	Expires             int    `form:"expires"`
	validator.Validator `form:"-"`
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	snippets, err := app.snippetModel.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	fmt.Println(r.Header.Get("x-correlation-id"))
	templateData := app.newTemplateData(r)
	templateData.Snippets = snippets
	app.render(w, http.StatusOK, "home.tmpl", templateData)
}

func (app *application) createSnippetForm(w http.ResponseWriter, r *http.Request) {
	templateData := app.newTemplateData(r)
	templateData.Form = CreateSnippetForm{
		Expires: 365,
	}

	app.render(w, http.StatusOK, "create.tmpl", templateData)
}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	var form CreateSnippetForm
	err := app.decodePostForm(r, form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank 11")
	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Content, 100), "content", "This field cannot be greater than 100")
	form.CheckField(validator.PermittedInt(form.Expires, 1, 7, 365), "expires", "Should be either 1, 7, 365")

	if len(form.FieldErrors) > 0 {
		templateData := app.newTemplateData(r)
		templateData.Form = form

		app.render(w, http.StatusOK, "create.tmpl", templateData)
		return
	}

	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	id, err := app.snippetModel.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		app.errorLogger.Output(2, string(debug.Stack()))
		app.clientError(w, http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippets/view/%d", id), http.StatusSeeOther)
}

func (app *application) viewSnippet(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	snippet, err := app.snippetModel.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	templateData := app.newTemplateData(r)
	templateData.Snippet = snippet
	app.render(w, http.StatusOK, "view.tmpl", templateData)
}
