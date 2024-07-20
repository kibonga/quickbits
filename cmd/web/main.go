package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

type app struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	mux      *http.ServeMux
}

func initApp() *app {
	return &app{
		infoLog:  log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
		errorLog: log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime),
		mux:      http.NewServeMux(),
	}
}

func (a *app) createServer() *http.Server {
	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()

	return &http.Server{
		Addr:     *addr,
		ErrorLog: a.errorLog,
		Handler:  a.mux,
	}
}

func (a *app) initHandlers() {
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

func main() {
	app := initApp()
	app.initHandlers()

	srv := app.createServer()

	app.infoLog.Printf("Starting server on %s", srv.Addr)
	err := srv.ListenAndServe()
	app.errorLog.Fatal(err)
}
