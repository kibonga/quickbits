package main

import "net/http"

func (a *app) routes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("/handler", &myType{})
	mux.HandleFunc("/handler/function", funcHandler)

	mux.Handle("/handler/func", http.HandlerFunc(funcHandler))
	mux.HandleFunc("/", a.home)
	mux.HandleFunc("/bits/view", a.viewBit)
	mux.HandleFunc("/bits/create", a.createBit)

	// Transaction
	mux.HandleFunc("/bits/update", a.updateBit)

	mux.Handle("/static/", initFS())

	return mux
}

func initFS() http.Handler {
	staticFileServer := http.FileServer(http.Dir("../../ui/static"))
	staticFileHandler := http.StripPrefix("/static", staticFileServer)

	return staticFileHandler
}
