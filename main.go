package main

import (
	"log"
	"net/http"
)

func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	w.Write([]byte("Hello from quickbits home!"))
}

func showSnippet(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display a specific snippet..."))
}

func createSnippet(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Create a new snippet..."))
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/", showSnippet)
	mux.HandleFunc("/snippet/create", createSnippet)

	// Potential security issue that can be exploited since it uses DefaultServeMux which is accessible by all packages (even 3rd party)
	// Instead use locally scoped servemux
	// http.HandleFunc("/", home)
	// http.HandleFunc("/snippet/", showSnippet)
	// http.HandleFunc("/snippet/create", createSnippet)

	log.Println("Listening on port :4000")
	err := http.ListenAndServe(":4000", nil)
	log.Fatal(err)
}
