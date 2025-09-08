package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/corbinlazarone/snippetbox/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Struct to hold our resources that needs to be
// shared across our app.
type application struct {
	errLog   *log.Logger          // field to introduce a custom error logger
	infoLog  *log.Logger          // field to introduce a custom info logger
	snippets *models.SnippetModel // used so our handlers.go can see out model
}

func main() {
	// we can easily change the port at runtime with the -port flag
	port := flag.String("port", ":4000", "HTTP server port number")

	// we can easily change databases at runtime with the -db flag
	datasource := flag.String("db", "YOUR_DB_URL", "my postgres db url")

	flag.Parse()

	// info and error logging
	infoLog := log.New(os.Stdout, "INFO:\t", log.Ldate|log.Ltime|log.Lshortfile)
	errLog := log.New(os.Stderr, "ERROR:\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*datasource)
	if err != nil {
		errLog.Fatal(err)
	}

	defer db.Close()

	app := &application{
		errLog:  errLog,
		infoLog: infoLog,
		snippets: &models.SnippetModel{
			DB: db,
		},
	}

	srv := &http.Server{
		Addr:     *port,
		Handler:  app.routes(),
		ErrorLog: errLog,
	}

	infoLog.Printf("Server running on %s...\n", *port)
	err = srv.ListenAndServe()
	errLog.Fatal(err)
}

func openDB(dataSource string) (*pgxpool.Pool, error) {
	ctx := context.Background()
	conn, err := pgxpool.New(ctx, dataSource)

	if err != nil {
		return nil, err
	}

	// pinging the connection to see if it was successful
	if err = conn.Ping(ctx); err != nil {
		return nil, err
	}

	return conn, nil
}
