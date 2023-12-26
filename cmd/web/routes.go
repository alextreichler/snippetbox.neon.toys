package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (a *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		a.notFound(w)
	})
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))

	dynamic := alice.New(a.sessionManager.LoadAndSave)

	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(a.home))
	router.Handler(http.MethodGet, "/snippet/view/:id", dynamic.ThenFunc(a.snippetView))
	router.Handler(http.MethodGet, "/snippet/create", dynamic.ThenFunc(a.snippetCreate))
	router.Handler(http.MethodPost, "/snippet/create", dynamic.ThenFunc(a.snippetCreatePost))

	router.Handler(http.MethodGet, "/user/signup", dynamic.ThenFunc(a.userSignup))
	router.Handler(http.MethodPost, "/user/signup", dynamic.ThenFunc(a.userSignupPost))
	router.Handler(http.MethodGet, "/user/login", dynamic.ThenFunc(a.userLogin))
	router.Handler(http.MethodPost, "/user/login", dynamic.ThenFunc(a.userLoginPost))
	router.Handler(http.MethodPost, "/user/logout", dynamic.ThenFunc(a.userLogoutPost))

	standard := alice.New(a.recoverPanic, a.logRequest, secureHeaders)

	return standard.Then(router)
}
