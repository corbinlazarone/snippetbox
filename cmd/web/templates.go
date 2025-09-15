package main

import "github.com/corbinlazarone/snippetbox/internal/models"

// templateData will act as a holding structure for
// any dynamic data we want to pass to our html templates.
type templateData struct {
	Snippet *models.Snippet
}
