package main

import (
	"database/sql"
	"flag"
	"kibonga/quickbits/internal/models"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type app struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	bits     *models.BitModel
	flags    *flags
	db       *sql.DB
}

type flags struct {
	dsn      string
	addr     string
	htmlPath string
}

func main() {
	app := initApp()

	app.DB()
	defer app.db.Close()

	srv := app.createServer()

	app.infoLog.Printf("Starting server on %s", srv.Addr)
	err := srv.ListenAndServe()
	app.errorLog.Fatal(err)
}

func initApp() *app {
	return &app{
		infoLog:  log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
		errorLog: log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
		flags:    parseCLIFlags(),
	}
}

func parseCLIFlags() *flags {
	dsn := flag.String("dsn", "web:pass@/quickbits_local?parseTime=true", "Database DSN")
	addr := flag.String("addr", ":4000", "HTTP network address")
	htmlPath := flag.String("htmlPath", "../../", "Path to HTML")

	flag.Parse()

	return &flags{
		addr:     *addr,
		dsn:      *dsn,
		htmlPath: *htmlPath,
	}
}

func (a *app) DB() {
	db, err := connectDb(a.flags.dsn)
	if err != nil {
		a.errorLog.Fatal(err)
	}
	a.bits = &models.BitModel{DB: db}
	a.db = db
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

func (a *app) createServer() *http.Server {
	return &http.Server{
		Addr:     a.flags.addr,
		ErrorLog: a.errorLog,
		Handler:  a.routes(),
	}
}
