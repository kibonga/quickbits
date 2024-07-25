package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (a *app) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		a.notFound(w)
	})

	// router.HandlerFunc(http.MethodGet, "/static/", http.StripPrefix("/static", fileServer).(http.HandlerFunc))
	fileServer := http.FileServer(http.Dir("../../ui/static"))
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))

	dynamic := alice.New(a.sessionManager.LoadAndSave)

	// Exercise purpose
	// router.Handler(http.MethodGet, "/", dynamic.Then(http.HandlerFunc(a.bitsIndex)))
	// router.HandlerFunc(http.MethodGet, "/bits/view/:id", http.HandlerFunc(dynamic.ThenFunc(a.bitView).ServeHTTP))
	// router.HandlerFunc(http.MethodGet, "/bits/create", dynamic.Then(http.HandlerFunc(a.bitCreateForm)).ServeHTTP)
	// router.Handler(http.MethodPost, "/bits/create", dynamic.ThenFunc(a.bitCreate))

	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(a.bitsIndex))
	router.Handler(http.MethodGet, "/bits/view/:id", dynamic.ThenFunc(a.bitsView))
	router.Handler(http.MethodGet, "/bits/create", dynamic.ThenFunc(a.bitsCreateForm))
	router.Handler(http.MethodPost, "/bits/create", dynamic.ThenFunc(a.bitsCreate))

	// Common middlewares
	common := alice.New(a.recoverPanic, a.afterMiddleware, a.logRequest, a.beforeMiddleware, a.secureHeaders)
	return common.Then(router)
}

func (a *app) routesMux() http.Handler {

	mux := http.NewServeMux()

	mux.Handle("/handler", &myType{})
	mux.HandleFunc("/handler/function", funcHandler)

	mux.Handle("/handler/func", http.HandlerFunc(funcHandler))
	mux.HandleFunc("/", a.bitsIndex)
	mux.HandleFunc("/bits/view", a.bitsView)
	mux.HandleFunc("/bits/create", a.bitsCreate)
	mux.HandleFunc("/bits/create/form", a.bitsCreateForm)

	// Transaction
	mux.HandleFunc("/bits/update", a.bitsUpdate)

	mux.Handle("/static/", staticFileHandler())

	common := alice.New(a.recoverPanic, a.afterMiddleware, a.logRequest, a.beforeMiddleware, a.secureHeaders)

	// headersMiddleware := a.secureHeaders(mux)
	// beforeMiddleware := a.beforeMiddleware(headersMiddleware)
	// loggerMiddleware := a.logRequest(beforeMiddleware)
	// afterMiddleware := a.afterMiddleware(loggerMiddleware)
	// panicMiddleware := a.recoverPanic(afterMiddleware)

	return common.Then(mux)
}

func staticFileHandler() http.Handler {
	staticFileServer := http.FileServer(http.Dir("../../ui/static"))
	staticFileHandler := http.StripPrefix("/static", staticFileServer)

	return staticFileHandler
}
