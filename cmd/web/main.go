package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"

	_ "net/url"

	_ "github.com/denisenkom/go-mssqldb"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

func main() {
	addr := flag.String("addr", ":4000", "Net-addr HTTP")
	dsn := flag.String("dsn", "server=localhost\\SQLExpress;user id=Administrator;database=master;app name=MyAppName", "connection string")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := sql.Open("sqlserver", *dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Start server at %s", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}
