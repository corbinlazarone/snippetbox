package main

import (
	"fmt"
	"net/http"
)

// Add secure headers to all incoming requests
func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy", "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")

		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")

		// NOTE: Any code here will execute on the way down the chain.
		// Any early returns will break the chain.

		next.ServeHTTP(w, r)

		// NOTE: Any code here will execute on the way back up the chain.
	})
}

// Will log the IP address, URL, and http method of the user for every incoming request
func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())

		next.ServeHTTP(w, r)
	})
}

// NOTE: Each http request is handler in its on goroutine, so when a pancic occurs in
// a request it will crash that request and not the whole server. -- that is cool

// NOTE: any extra goroutine you add to a request to do some other background processing,
// you will have to do your own panic recover, because as stated above the middleware below
// only handles the request goroutine.

// Show 500 internal server error to user on request panics
func (app *application) recoverFromPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// defered function that will alwasy run in the event of a panic as Go unwinds
		// the stack.
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connetion", "close")

				// show a 500 server error to the user
				app.serverError(w, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (app *application) requireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// if the user is not authenticated redirect them to the login page
		// and return from the middleware chain so that no subsequent middleware
		// is executed.
		if !app.isAuthenticated(r) {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}

		// Otherwise set the "Cache-Control: no-store" header so that pages that require
		// authentication will not be cached by the browser.
		w.Header().Add("Cache-Control", "no-cache")

		// And call the next handler in the chain
		next.ServeHTTP(w, r)
	})
}
