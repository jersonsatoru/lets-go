package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	mux := httprouter.New()
	mux.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))
	mux.HandlerFunc(http.MethodGet, "/snippets/create", app.createSnippetForm)
	mux.HandlerFunc(http.MethodPost, "/snippets/create", app.createSnippet)
	mux.HandlerFunc(http.MethodGet, "/snippets/view/:id", app.viewSnippet)
	mux.HandlerFunc(http.MethodGet, "/", http.HandlerFunc(app.home))
	return app.panicRecover(app.logRequests(securityMiddleware(mux)))
}
