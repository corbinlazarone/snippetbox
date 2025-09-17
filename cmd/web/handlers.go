package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/corbinlazarone/snippetbox/internal/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// only show home on / path
	if r.URL.Path != "/" {
		app.clientError(w, http.StatusNotFound)
		return
	}

	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	// use our templateData holding struct
	data := &templateData{
		Snippets: &snippets,
	}

	// use render helper function to render our template page
	app.render(w, "home.tmpl.html", data, http.StatusOK)
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
	title := "O snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
	expires := 7
	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// redirect the user to the relvant snippet id page
	http.Redirect(w, r, fmt.Sprintf("/snippet/view?id=%d\n", id), http.StatusSeeOther)
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
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
	data := &templateData{
		Snippet: snippet,
	}

	app.render(w, "view.tmpl.html", data, http.StatusOK)
}
