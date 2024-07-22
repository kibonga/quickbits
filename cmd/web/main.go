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
	bitModel *models.BitModel
	cliFlags *cliFlags
	db       *sql.DB
}

type cliFlags struct {
	dsn      string
	addr     string
	htmlPath string
}

func main() {
	app := initApp()

	app.db = app.connectDb()
	defer app.db.Close()

	app.addBitModel()

	srv := app.createServer()

	app.infoLog.Printf("Starting server on %s", srv.Addr)
	err := srv.ListenAndServe()
	app.errorLog.Fatal(err)
}

func initApp() *app {
	return &app{
		infoLog:  log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
		errorLog: log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
		cliFlags: parseCLIFlags(),
		db:       nil,
		bitModel: nil,
	}
}

func parseCLIFlags() *cliFlags {
	dsn := flag.String("dsn", "web:pass@/quickbits_local?parseTime=true", "Database DSN")
	addr := flag.String("addr", ":4000", "HTTP network address")
	htmlPath := flag.String("htmlPath", "../../", "Path to HTML")

	flag.Parse()

	return &cliFlags{
		addr:     *addr,
		dsn:      *dsn,
		htmlPath: *htmlPath,
	}
}

func (a *app) connectDb() *sql.DB {
	db, err := openMysqlConn(a.cliFlags.dsn)
	if err != nil {
		a.errorLog.Fatal(err)
	}

	return db
}

func openMysqlConn(dsn string) (*sql.DB, error) {
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
		Addr:     a.cliFlags.addr,
		ErrorLog: a.errorLog,
		Handler:  a.routes(),
	}
}

func (a *app) addBitModel() {
	bitModel, err := models.CreateBitModel(a.db)
	if err != nil {
		a.errorLog.Fatal(err)
	}

	a.bitModel = bitModel
}
