package main

import (
	"net/http"

	"github.com/Bekian/SnippetBox/ui"

	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	// router
	mux := http.NewServeMux()

	// serve static files
	mux.Handle("GET /static/", http.FileServerFS(ui.Files))

	mux.HandleFunc("GET /ping", ping)

	// middleware chain for session manager logic
	dynamic := alice.New(app.sessionManager.LoadAndSave, noSurf, app.authenticate)

	/// server app routes with middleware chain
	// routes that do not require auth
	mux.Handle("GET /{$}", dynamic.ThenFunc(app.home))
	mux.Handle("GET /snippet/view/{id}", dynamic.ThenFunc(app.snippetView))
	// user signup and login routes
	mux.Handle("GET /user/signup", dynamic.ThenFunc(app.userSignup))
	mux.Handle("POST /user/signup", dynamic.ThenFunc(app.userSignupPost))
	mux.Handle("GET /user/login", dynamic.ThenFunc(app.userLogin))
	mux.Handle("POST /user/login", dynamic.ThenFunc(app.userLoginPost))

	// auth protected routes
	protected := dynamic.Append(app.requireAuthentication)

	mux.Handle("GET /snippet/create", protected.ThenFunc(app.snippetCreate))
	mux.Handle("POST /snippet/create", protected.ThenFunc(app.snippetCreatePost))
	mux.Handle("POST /user/logout", protected.ThenFunc(app.userLogoutPost))

	// middleware chain obj using alice
	standard := alice.New(app.recoverPanic, app.logRequest, commonHeaders)

	return standard.Then(mux)
}
