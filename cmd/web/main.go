package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/balabanovds/go-snippetbox/pkg/models/mysql"
	"log"
	"net/http"
	"os"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

type application struct {
	infoLog  *log.Logger
	errLog   *log.Logger
	snippets *mysql.SnippetModel
}

func main() {
	host := flag.String("host", "localhost", "Host name to start server on")
	port := flag.Int("port", 4000, "Port to start server on")
	dbUser := flag.String("db_user", "", "MySQL user")
	dbPassw := flag.String("db_password", "", "MySQL password")
	dbHost := flag.String("db_host", "localhost", "MySQL host")
	dbPort := flag.Int("db_port", 3306, "MySQL port")
	dbName := flag.String("db_name", "snippetbox", "MySQL app DB name")

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

	infoLog.Println(dsn)
	db, err := openDB(dsn)
	if err != nil {
		errLog.Fatal(err)
	}
	defer db.Close()

	app := &application{
		infoLog:  infoLog,
		errLog:   errLog,
		snippets: mysql.NewSnippetModel(db),
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
