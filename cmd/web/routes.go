package main

import (
	"net/http"

	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	// router
	mux := http.NewServeMux()
	// serve static files
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	// middleware chain for session manager logic
	dynamic := alice.New(app.sessionManager.LoadAndSave)
	// server app routes with middleware chain
	mux.Handle("GET /{$}", dynamic.ThenFunc(app.home))
	mux.Handle("GET /snippet/view/{id}", dynamic.ThenFunc(app.snippetView))
	mux.Handle("GET /snippet/create", dynamic.ThenFunc(app.snippetCreate))
	mux.Handle("POST /snippet/create", dynamic.ThenFunc(app.snippetCreatePost))

	// middleware chain obj using alice
	standard := alice.New(app.recoverPanic, app.logRequest, commonHeaders)

	return standard.Then(mux)
}
