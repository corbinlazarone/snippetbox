package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/corbinlazarone/snippetbox/internal/models"
	"github.com/julienschmidt/httprouter"
)

// NOTE: Struct field must be public in order to be read by
// the html/template package when rendering the template.
type snippetCreateForm struct {
	Title       string
	Content     string
	Expires     int
	FieldErrors map[string]string
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	// use our templateData holding struct
	data := app.newTemplateData()
	data.Snippets = &snippets

	// use render helper function to render our template page
	app.render(w, "home.tmpl.html", data, http.StatusOK)
}

// renders the html for our snippet create form
func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData()

	// Initialize a new createSnippetForm instance and pass it to the template.
	// Notice how this is also a great opportunity to set any default or
	// 'initial' values for the form --- here we set the initial value for the
	// snippet expiry to 365 days.
	data.Form = snippetCreateForm{
		Expires: 365,
	}

	app.render(w, "create.tmpl.html", data, http.StatusOK)
}

// creates the submitted snippet to the database
func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {

	// Chapter 8.2 - parsing the form

	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	expires, err := strconv.Atoi(r.PostForm.Get("expires"))
	if err != nil {
		app.serverError(w, err)
		return
	}

	form := snippetCreateForm{
		Title:       r.PostForm.Get("title"),
		Content:     r.PostForm.Get("content"),
		Expires:     expires,
		FieldErrors: map[string]string{}, // A map to hold any validation errors for the form fields

	}

	// Check that the title field is not blank and is not more than 100
	// characters long.
	if strings.TrimSpace(form.Title) == "" {
		form.FieldErrors["title"] = "This field cannot be blank"
		// we use the utf8.RuneCountInString() b/c we want to count the number of characters
	} else if utf8.RuneCountInString(form.Title) > 100 {
		form.FieldErrors["title"] = "This field cannot be more than 100 characters long"
	}

	// Check that the content value isn't blank
	if strings.TrimSpace(form.Content) == "" {
		form.FieldErrors["content"] = "This field cannot be blank"
	}

	// Check that the expires value matches one of our permitted values (1, 7 or
	// 365).
	if form.Expires != 1 && form.Expires != 7 && form.Expires != 365 {
		form.FieldErrors["expires"] = "This field must equal 1, 7 or 365"
	}

	// If any errors in our map than re render the create.tmpl.html page
	// with a 422 status code error.
	if len(form.FieldErrors) > 0 {
		data := app.newTemplateData()
		data.Form = form
		app.render(w, "create.tmpl.html", data, http.StatusUnprocessableEntity)
		return
	}

	id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// redirect the user to the relvant snippet id page
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	// when httprouter is parsing a request, the values of any named parameters
	// will be stored in the request context. We can use the ParamsFromContext()
	// function to retrive the slice containing these parameter names and values.
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 { // convert id to int and makes sure its greater than 1
		app.clientError(w, http.StatusNotFound)
		return
	}
	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.clientError(w, http.StatusNotFound)
		} else {
			app.serverError(w, err)
		}
		return
	}

	// fix new lines
	snippet.Content = strings.ReplaceAll(snippet.Content, "\\n", "\n")

	// use our templateData holding struct
	data := app.newTemplateData()
	data.Snippet = snippet

	app.render(w, "view.tmpl.html", data, http.StatusOK)
}
