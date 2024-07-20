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
		errorLog: log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
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

func main() {
	app := initApp()
	app.routes()

	srv := app.createServer()

	app.infoLog.Printf("Starting server on %s", srv.Addr)
	err := srv.ListenAndServe()
	app.errorLog.Fatal(err)
}
