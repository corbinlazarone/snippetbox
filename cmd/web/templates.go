package main

import (
	"html/template"
	"path/filepath"

	"github.com/corbinlazarone/snippetbox/internal/models"
)

// templateData will act as a holding structure for
// any dynamic data we want to pass to our html templates.
type templateData struct {
	Snippet  *models.Snippet
	Snippets *[]models.Snippet
}

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

		// parse our base template file
		tmplSet, err := template.ParseFiles("./ui/html/base.tmpl.html")
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
