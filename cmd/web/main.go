package main

import (
	"database/sql"
	"flag"
	"html/template"
	"kibonga/quickbits/internal/models"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	_ "github.com/go-sql-driver/mysql"
)

type app struct {
	errorLog       *log.Logger
	infoLog        *log.Logger
	bitModel       *models.BitModel
	cliFlags       *cliFlags
	db             *sql.DB
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
}

type cliFlags struct {
	dsn      string
	addr     string
	htmlPath string
}

func main() {
	cliFlags := parseCLIFlags()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openMysqlConn(cliFlags.dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	bitModel, err := models.CreateBitModel(db)
	if err != nil {
		errorLog.Fatal(err)
	}

	tmplCache, err := createTemplateCache(cliFlags.htmlPath)
	if err != nil {
		errorLog.Fatal(err)
	}

	sessionManager := scs.New()
	sessionManager.Store = &mysqlstore.MySQLStore{}
	sessionManager.Lifetime = 12 * time.Hour

	app := &app{
		errorLog:       errorLog,
		infoLog:        infoLog,
		bitModel:       bitModel,
		cliFlags:       cliFlags,
		db:             db,
		templateCache:  tmplCache,
		formDecoder:    form.NewDecoder(),
		sessionManager: sessionManager,
	}

	srv := &http.Server{
		Addr:     cliFlags.addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Starting server on %s", srv.Addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
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
