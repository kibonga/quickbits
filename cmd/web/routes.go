package main

import (
	"net/http"

	"github.com/justinas/alice"
)

func (a *app) routes() http.Handler {

	mux := http.NewServeMux()

	mux.Handle("/handler", &myType{})
	mux.HandleFunc("/handler/function", funcHandler)

	mux.Handle("/handler/func", http.HandlerFunc(funcHandler))
	mux.HandleFunc("/", a.home)
	mux.HandleFunc("/bits/view", a.viewBit)
	mux.HandleFunc("/bits/create", a.createBit)

	// Transaction
	mux.HandleFunc("/bits/update", a.updateBit)

	mux.Handle("/static/", staticFileServer())

	common := alice.New(a.recoverPanic, a.afterMiddleware, a.logRequest, a.beforeMiddleware, a.secureHeaders)

	// headersMiddleware := a.secureHeaders(mux)
	// beforeMiddleware := a.beforeMiddleware(headersMiddleware)
	// loggerMiddleware := a.logRequest(beforeMiddleware)
	// afterMiddleware := a.afterMiddleware(loggerMiddleware)
	// panicMiddleware := a.recoverPanic(afterMiddleware)

	return common.Then(mux)
}

func staticFileServer() http.Handler {
	staticFileServer := http.FileServer(http.Dir("../../ui/static"))
	staticFileHandler := http.StripPrefix("/static", staticFileServer)

	return staticFileHandler
}
