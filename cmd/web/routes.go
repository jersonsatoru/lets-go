package main

import "net/http"

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	mux.Handle("/", http.HandlerFunc(app.home))
	mux.HandleFunc("/snippets/create", app.createSnippet)
	mux.HandleFunc("/snippets/view", app.viewSnippet)
	return app.panicRecover(app.logRequests(securityMiddleware(mux)))
}
