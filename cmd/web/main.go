package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"

	ms_sql "pl1x/pkg/models/mssql"

	_ "net/url"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/microsoft/go-mssqldb"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	//snippets *ms_sql.SnippetModel
	snippets *ms_sql.SnippetModel
}

func main() {
	addr := flag.String("addr", ":4000", "Net-addr HTTP")
	dsn := flag.String("dsn", "server=localhost\\SQLExpress;user id=Web;password=hackOFF1;database=master;app name=MyAppName;", "connection string")

	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openMSSQL(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
		//snippets: &ms_sql.SnippetModel{DB: db},
		snippets: &ms_sql.SnippetModel{DB: db},
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

func openMSSQL(dsn string) (*sql.DB, error) {
	db, err := sql.Open("sqlserver", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func openMYSQL(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
