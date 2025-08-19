package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {

	port := flag.String("port", ":4000", "HTTP server port number")

	flag.Parse()

	mux := http.NewServeMux()

	// Use the mux.Handle() function to register the file server as the handler for
	// all URL paths that start with "/static/". For matching paths, we strip the
	// "/static" prefix before the request reaches the file server.
	fileServer := http.FileServer(http.Dir("./ui/static"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)

	log.Printf("Server running on %s...\n", *port)
	err := http.ListenAndServe(*port, mux)
	log.Fatal(err)
}
