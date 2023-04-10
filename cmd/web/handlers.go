package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"runtime/debug"
	"strconv"

	"github.com/jersonsatoru/lets-go/internal/models"
)

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

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	title := "The Snail"
	content := "Impressive work by Mark Jr"
	expires := 7

	id, err := app.snippetModel.Insert(title, content, expires)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		app.errorLogger.Output(2, string(debug.Stack()))
		app.clientError(w, http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippets/view?id=%d", id), http.StatusSeeOther)
}

func (app *application) viewSnippet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
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
