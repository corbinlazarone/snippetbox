package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

func (app *application) render(w http.ResponseWriter, pageName string, data *templateData, statusCode int) {
	tmplSet, ok := app.templateCache[pageName]

	if !ok {
		// the page tempalte is not in the map
		err := fmt.Errorf("%s not found in template cache\n", pageName)
		app.serverError(w, err)
		return
	}

	// write out the provided statusCode
	w.WriteHeader(statusCode)

	err := tmplSet.ExecuteTemplate(w, "base", data)
	if err != nil {
		app.serverError(w, err)
		return
	}
}

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
