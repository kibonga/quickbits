package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.Handle("/handler", &myType{})
	mux.HandleFunc("/handler/function", funcHandler)
	mux.Handle("/handler/func", http.HandlerFunc(funcHandler))
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet", showSnippet)
	mux.HandleFunc("/snippet/create", createSnippet)

	staticFileServer := http.FileServer(http.Dir("../../ui/static"))
	staticFileHandler := http.StripPrefix("/static", staticFileServer)

	mux.Handle("/static/", staticFileHandler)

	log.Println("Starting server on :4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
