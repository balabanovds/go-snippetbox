package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/balabanovds/go-snippetbox/pkg/models/mysql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golangcollege/sessions"
)

type application struct {
	infoLog       *log.Logger
	errLog        *log.Logger
	session       *sessions.Session
	snippets      *mysql.SnippetModel
	templateCache map[string]*template.Template
}

func main() {
	host := flag.String("host", "localhost", "Host name to start server on")
	port := flag.Int("port", 4000, "Port to start server on")
	dbUser := flag.String("db_user", "", "MySQL user")
	dbPassw := flag.String("db_password", "", "MySQL password")
	dbHost := flag.String("db_host", "localhost", "MySQL host")
	dbPort := flag.Int("db_port", 3306, "MySQL port")
	dbName := flag.String("db_name", "snippetbox", "MySQL app DB name")
	secret := flag.String("secret", "LtcznmGmzys[J,tpmzy", "Secret key")

	flag.Parse()

	addr := *host + ":" + strconv.Itoa(*port)

	// logging
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	if *dbUser == "" || *dbPassw == "" {
		errLog.Print("db user and password should be defined")
		flag.Usage()
		os.Exit(1)
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
		*dbUser, *dbPassw, *dbHost, *dbPort, *dbName)

	db, err := openDB(dsn)
	if err != nil {
		errLog.Fatal(err)
	}
	defer db.Close()

	templateCache, err := newTemplateCache("./ui/templates")
	if err != nil {
		errLog.Fatal(err)
	}

	session := sessions.New([]byte(*secret))
	session.Lifetime = 12 * time.Hour

	app := &application{
		infoLog:       infoLog,
		errLog:        errLog,
		session:       session,
		snippets:      mysql.NewSnippetModel(db),
		templateCache: templateCache,
	}

	tlsConfig := &tls.Config{

		PreferServerCipherSuites: true,
		CurvePreferences:         []tls.CurveID{tls.CurveP256, tls.X25519},
	}

	srv := &http.Server{
		Addr:         addr,
		ErrorLog:     errLog,
		Handler:      app.routes(),
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	infoLog.Printf("serving on https://%s\n", addr)
	errLog.Fatalln(srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem"))
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
