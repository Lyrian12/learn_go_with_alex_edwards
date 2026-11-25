package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	fileserver := http.FileServer(http.Dir("./ui/static"))
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileserver))
	router.HandlerFunc(http.MethodGet, "/", app.home)
	router.HandlerFunc(http.MethodGet, "/snippet/:id", app.snippetView)

	router.HandlerFunc(http.MethodGet, "/create", app.snippetCreate)
	router.HandlerFunc(http.MethodPost, "/create", app.snippetCreatePost)
	return app.sessionManager.LoadAndSave(app.recoverPanic(app.logRequest(secureHeaders(router))))

}
