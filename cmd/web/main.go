package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type app struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

func main() {
	app := initApp()

	srv := app.createServer()
	app.infoLog.Printf("address: %s", srv.Addr)

	dsn := "web2:pass@tcp(quickbits_mysql:3306)/quickbits?parseTime=true"
	db, err := connectDb(dsn)
	if err != nil {
		app.errorLog.Fatal(err)
	}

	defer db.Close()

	app.infoLog.Printf("Starting server on %s", srv.Addr)
	err = srv.ListenAndServe()
	app.errorLog.Fatal(err)
}

func initApp() *app {
	return &app{
		infoLog:  log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
		errorLog: log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

func (a *app) createServer() *http.Server {
	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()

	return &http.Server{
		Addr:     *addr,
		ErrorLog: a.errorLog,
		Handler:  a.routes(),
	}
}

func connectDb(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
