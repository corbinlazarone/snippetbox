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
	infoLog *log.Logger // field to introduce a custom info logger
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

	mux := http.NewServeMux()
	// Use the mux.Handle() function to register the file server as the handler for
	// all URL paths that start with "/static/". For matching paths, we strip the
	// "/static" prefix before the request reaches the file server.
	fileServer := http.FileServer(http.Dir("./ui/static"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet/view", app.snippetView)
	mux.HandleFunc("/snippet/create", app.snippetCreate)

	srv := &http.Server{
		Addr:     *port,
		Handler:  mux,
		ErrorLog: errLog,
	}

	infoLog.Printf("Server running on %s...\n", *port)
	err := srv.ListenAndServe()
	errLog.Fatal(err)
}
