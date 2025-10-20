package main

import (
	"html/template"
	"net/http"
	"path/filepath"
	"time"

	"github.com/corbinlazarone/snippetbox/internal/models"
)

// templateData will act as a holding structure for
// any dynamic data we want to pass to our html templates.
type templateData struct {
	CurrentYear     int
	Snippet         *models.Snippet
	Snippets        *[]models.Snippet
	Form            any
	Flash           string
	IsAuthenticated bool
}

// initialize the templateData struct with a current year
func (app *application) newTemplateData(r *http.Request) *templateData {
	return &templateData{
		CurrentYear: time.Now().Year(),

		// add flash message to template data if it exists
		Flash: app.sessionManager.PopString(r.Context(), "flash"),

		// add the authentication status to the template data.
		IsAuthenticated: app.isAuthenticated(r),
	}
}

// NOTE: Custome template function like the one below can only
// return one value, unless of course you return a value and an
// error.

// returns a formatted date
func humanReadableDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humanReadableDate,
}

// newTemplateCache() parses all our html pages when the app starts
// and stores them in map (cache) so we can refrence them in our app, instead
// of parsing and reading the template html page everytime the user goes to a new page.
func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	// Grab a slice for all our page tempaltes
	matches, err := filepath.Glob("./ui/html/pages/*.tmpl.html")
	if err != nil {
		return nil, err
	}

	for _, page := range matches {
		// grab the name of the page for the key in our map.
		// name will be the last element of the path ex: home.tmpl.html
		name := filepath.Base(page)

		// NOTE: Customer template functions must be registerd before the template is parsed.
		// This means we have to create a empty template set using template.New() to register the
		// template.FuncMap and then parse the file.

		// parse our base template file
		tmplSet, err := template.New(name).Funcs(functions).ParseFiles("./ui/html/base.tmpl.html")
		if err != nil {
			return nil, err
		}

		// parse all the partials template files b/c we will add more later
		tmplSet, err = tmplSet.ParseGlob("./ui/html/partials/*.tmpl.html")
		if err != nil {
			return nil, err
		}

		// parse our page template
		tmplSet, err = tmplSet.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		cache[name] = tmplSet
	}

	return cache, nil
}
