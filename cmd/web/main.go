package main

import (
	"context"
	"crypto/tls"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/pgxstore"
	"github.com/alexedwards/scs/v2"
	"github.com/corbinlazarone/snippetbox/internal/models"
	"github.com/go-playground/form/v4"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Struct to hold our resources that needs to be
// shared across our app.
type application struct {
	errLog         *log.Logger                   // field to introduce a custom error logger
	infoLog        *log.Logger                   // field to introduce a custom info logger
	snippets       *models.SnippetModel          // used so our handlers.go can see our model
	templateCache  map[string]*template.Template // holds our cached templates
	formDecoder    *form.Decoder                 // used so our handerls.go can auto parse forms
	sessionManager *scs.SessionManager
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

	// initialize a new tempalte cache
	templateCache, err := newTemplateCache()
	if err != nil {
		errLog.Fatal(err) // log the error with customer logger and kill the program
	}

	db, err := openDB(*datasource)
	if err != nil {
		errLog.Fatal(err)
	}

	defer db.Close() // make sure the db connection is closed right after opening it.

	// Initialze a new decoder instance.
	formDecoder := form.NewDecoder()

	// Initialize a new session manager with scs.New(), then configured it to use our Postgresql database
	// as the session store, and set a lifetime of 12 hours.
	sessionManager := scs.New()
	sessionManager.Store = pgxstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour

	// Cookie will only be sent by a user's web browser when a
	// HTTPS connection is being used. (Will not be send over
	// and HTTP connection).
	sessionManager.Cookie.Secure = true

	app := &application{
		errLog:  errLog,
		infoLog: infoLog,
		snippets: &models.SnippetModel{
			DB: db,
		},
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
	}

	// Initialize a tls.Config struct to hold the non-default TLS settings we
	// want the server to use. In this case the only thing that we're changing
	// is the curve preferences value, so that only elliptic curves with
	// assembly implementations are used.
	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	srv := &http.Server{
		Addr:         *port,
		Handler:      app.routes(),
		ErrorLog:     errLog,
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	infoLog.Printf("Server running on %s...\n", *port)
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
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
