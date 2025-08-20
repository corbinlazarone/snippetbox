package main

import (
	"fmt"
	"net/http"
	"strconv"
	"text/template"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// only show home on / path
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	files := []string{
		"./ui/html/base.tmpl.html", // first element must be our base template
		"./ui/html/partials/nav.tmpl.html",
		"./ui/html/pages/home.tmpl.html",
	}

	// Read the template file into a template set
	tmplSet, err := template.ParseFiles(files...) // "files..." unpacks the slice so each element is passed individually
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		app.errLog.Println(err)
		return
	}

	// Read the template set into a response body
	err = tmplSet.ExecuteTemplate(w, "base", nil)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		app.errLog.Println(err)
		return
	}
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Write([]byte("Snippet Create\n"))
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 { // convert id to int and makes sure its greater than 1
		http.NotFound(w, r)
		app.errLog.Println(err)
		return
	}

	fmt.Fprintf(w, "Found snipped with id: %d\n", id)
}
