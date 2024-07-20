package main

import "net/http"

func (a *app) routes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("/handler", &myType{})
	mux.HandleFunc("/handler/function", funcHandler)

	mux.Handle("/handler/func", http.HandlerFunc(funcHandler))
	mux.HandleFunc("/", a.home)
	mux.HandleFunc("/snippet", a.showSnippet)
	mux.HandleFunc("/snippet/create", a.createSnippet)

	mux.Handle("/static/", initFS())

	return mux
}

func initFS() http.Handler {
	staticFileServer := http.FileServer(http.Dir("../../ui/static"))
	staticFileHandler := http.StripPrefix("/static", staticFileServer)

	return staticFileHandler
}
