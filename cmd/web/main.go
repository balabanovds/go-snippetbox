package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

type application struct {
	infoLog *log.Logger
	errLog  *log.Logger
}

func main() {
	host := flag.String("host", "localhost", "Host name to start server on")
	port := flag.Int("port", 4000, "Port to start server on")
	dsn := flag.String("dsn", "snippet:feropl09@tcp(balabanov.sknt.ru:43306)/snippetbox?parseTime=true",
		"MySQL data source name")
	flag.Parse()

	addr := *host + ":" + strconv.Itoa(*port)

	// logging
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dsn)
	if err != nil {
		errLog.Fatal(err)
	}

	defer db.Close()

	app := &application{
		infoLog: infoLog,
		errLog:  errLog,
	}

	srv := &http.Server{
		Addr:     addr,
		ErrorLog: errLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("serving on http://%s\n", addr)
	errLog.Fatalln(srv.ListenAndServe())
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
