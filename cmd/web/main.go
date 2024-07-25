package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"fmt"
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
	rootPath string
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

	tmplCache, err := createTemplateCache(cliFlags.rootPath)
	if err != nil {
		errorLog.Fatal(err)
	}

	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
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

	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	srv := &http.Server{
		Addr:         cliFlags.addr,
		ErrorLog:     errorLog,
		Handler:      app.routes(),
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	infoLog.Printf("Starting server on %s", srv.Addr)
	certPath := fmt.Sprintf("%s%s", app.cliFlags.rootPath, "tls/cert.pem")
	pkPath := fmt.Sprintf("%s%s", app.cliFlags.rootPath, "tls/key.pem")
	err = srv.ListenAndServeTLS(certPath, pkPath)
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
		rootPath: *htmlPath,
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
