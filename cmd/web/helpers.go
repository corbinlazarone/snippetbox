package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/go-playground/form/v4"
)

// The second parameter here, destination, is the target destination that we want
// to decode the form data into.
func (app *application) decodePostForm(r *http.Request, destination any) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}

	err = app.formDecoder.Decode(destination, r.PostForm)

	if err != nil {
		// NOTE: If we try to use an invalid target destination, the Decode() method
		// will return an error with the type *form.InvalidDecoderError.We use
		// errors.As() to check for this and raise a panic rather than returning
		// the error.
		var invalidDecodeError *form.InvalidDecoderError

		if errors.As(err, &invalidDecodeError) {
			panic(err)
		}

		// For all other errors, return as normal
		return err
	}
	return nil
}

func (app *application) render(w http.ResponseWriter, pageName string, data *templateData, statusCode int) {
	tmplSet, ok := app.templateCache[pageName]

	if !ok {
		// the page tempalte is not in the map
		err := fmt.Errorf("%s not found in template cache\n", pageName)
		app.serverError(w, err)
		return
	}

	// Understanding the problem:
	// --------------------------
	// With our templates now handling dynamic data rendering,
	// we need to handle runtime errors where the data may be null,
	// and the tempalte doesn't render correctly.
	// To do this we will first write the template to a temp buffer
	// and if that is successful then we know that html template doesn't have
	// any errors in it, therefor we can render the page to the user.
	// if the write to the buffer fails we know their is an error, so
	// we will show a server error to the user.

	buf := new(bytes.Buffer)

	// write the tempalte to the buffer instead of straight to
	// the http.ResponseWriter
	err := tmplSet.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// write out the provided statusCode
	w.WriteHeader(statusCode)

	// write the contents of the buffer out
	buf.WriteTo(w)
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
