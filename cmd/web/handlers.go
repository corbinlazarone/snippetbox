package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/corbinlazarone/snippetbox/internal/models"
	"github.com/julienschmidt/httprouter"
)

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

	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")

	expires, err := strconv.Atoi(r.PostForm.Get("expires"))
	if err != nil {
		app.serverError(w, err)
		return
	}

	id, err := app.snippets.Insert(title, content, expires)
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
