package main

import "net/http"

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

	headersMiddleware := a.secureHeaders(mux)
	beforeMiddleware := a.beforeMiddleware(headersMiddleware)
	loggerMiddleware := a.logRequest(beforeMiddleware)
	afterMiddleware := a.afterMiddleware(loggerMiddleware)
	panicMiddleware := a.recoverPanic(afterMiddleware)

	return panicMiddleware
}

func staticFileServer() http.Handler {
	staticFileServer := http.FileServer(http.Dir("../../ui/static"))
	staticFileHandler := http.StripPrefix("/static", staticFileServer)

	return staticFileHandler
}
