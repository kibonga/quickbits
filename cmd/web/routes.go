package main

import "net/http"

func (a *app) routes() {
	a.mux.Handle("/handler", &myType{})
	a.mux.HandleFunc("/handler/function", funcHandler)

	a.mux.Handle("/handler/func", http.HandlerFunc(funcHandler))
	a.mux.HandleFunc("/", a.home)
	a.mux.HandleFunc("/snippet", a.showSnippet)
	a.mux.HandleFunc("/snippet/create", a.createSnippet)

	a.mux.Handle("/static/", initFS())
}

func initFS() http.Handler {
	staticFileServer := http.FileServer(http.Dir("../../ui/static"))
	staticFileHandler := http.StripPrefix("/static", staticFileServer)

	return staticFileHandler
}
