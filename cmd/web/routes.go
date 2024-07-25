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

	router.HandlerFunc(http.MethodGet, "/", a.bitsIndex)
	router.HandlerFunc(http.MethodGet, "/bits/view/:id", a.bitView)
	router.HandlerFunc(http.MethodGet, "/bits/create", a.bitCreateForm)
	router.HandlerFunc(http.MethodPost, "/bits/create", a.bitCreate)

	middlewares := alice.New(a.recoverPanic, a.afterMiddleware, a.logRequest, a.beforeMiddleware, a.secureHeaders)
	return middlewares.Then(router)
}

func (a *app) routesMux() http.Handler {

	mux := http.NewServeMux()

	mux.Handle("/handler", &myType{})
	mux.HandleFunc("/handler/function", funcHandler)

	mux.Handle("/handler/func", http.HandlerFunc(funcHandler))
	mux.HandleFunc("/", a.bitsIndex)
	mux.HandleFunc("/bits/view", a.bitView)
	mux.HandleFunc("/bits/create", a.bitCreate)
	mux.HandleFunc("/bits/create/form", a.bitCreateForm)

	// Transaction
	mux.HandleFunc("/bits/update", a.updateBit)

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
