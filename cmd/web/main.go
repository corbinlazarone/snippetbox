package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

// Struct to hold our resources that needs to be
// shared across our app.
type application struct {
	errLog  *log.Logger // field to introduce a custom error logger
	infoLog *log.Logger // field to introduce a custom info logger2
}

func main() {
	// port flag
	port := flag.String("port", ":4000", "HTTP server port number")
	flag.Parse()

	// info and error logging
	infoLog := log.New(os.Stdout, "INFO:\t", log.Ldate|log.Ltime|log.Lshortfile)
	errLog := log.New(os.Stderr, "ERROR:\t", log.Ldate|log.Ltime|log.Lshortfile)

	app := &application{
		errLog:  errLog,
		infoLog: infoLog,
	}

	srv := &http.Server{
		Addr:     *port,
		Handler:  app.routes(),
		ErrorLog: errLog,
	}

	infoLog.Printf("Server running on %s...\n", *port)
	err := srv.ListenAndServe()
	errLog.Fatal(err)
}
