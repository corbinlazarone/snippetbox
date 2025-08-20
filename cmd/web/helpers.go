package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

func (app *application) serverError(w http.ResponseWriter, err error) {
	// debug.Stack() returns the stack trace of the error.
	msg := fmt.Sprintf("%s\n%s\n", err.Error(), debug.Stack())

	// log the error using our customer logger.
	// Output(2, msg) says go down 2 on the stack trace
	// to display where the serverError helper was called
	// not where it was logged i.e right here.
	app.errLog.Output(2, msg)

	// send a 500 http error back to the user
	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
}

// Will make this return custom error messages later.
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}
