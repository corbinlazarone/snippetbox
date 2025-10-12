package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/corbinlazarone/snippetbox/internal/models"
	"github.com/corbinlazarone/snippetbox/internal/validator"
	"github.com/julienschmidt/httprouter"
)

// NOTE: Struct field must be public in order to be read by
// the html/template package when rendering the template.

// NOTE: the struct tags tells the decoder how to map HTML form values into
// different struct fields. For example, here we are telling the decoder to
// store the value from the HTML form input with the name "title" in the Title
// field.
type snippetCreateForm struct {
	Title               string     `form:"title"`
	Content             string     `form:"content"`
	Expires             int        `form:"expires"`
	validator.Validator `form:"-"` // tells the from decoder to ignore this field
}

type userSignUpForm struct {
	Name                string     `form:"name"`
	Email               string     `form:"email"`
	Password            string     `form:"password"`
	validator.Validator `form:"-"` // tells the from decoder to ignore this field
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	// use our templateData holding struct
	data := app.newTemplateData(r)
	data.Snippets = &snippets

	// use render helper function to render our template page
	app.render(w, "home.tmpl.html", data, http.StatusOK)
}

// renders the html for our snippet create form
func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

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
	var form snippetCreateForm

	// we pass a pointer to our form to the Decoder and the request and
	// it will fill out our struct that holds the form values with the values
	// from the HTML form.
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Check that the title field is not blank.
	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")

	// Check the the title field is not more than 100 characters long.
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "This field cannot be more than 100 characters long")

	// Check that the content value isn't blank
	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")

	// Check that the expires value matches one of our permitted values (1, 7 or
	// 365).
	form.CheckField(validator.PermittedInt(form.Expires, 1, 7, 365), "expires", "This field must equal 1, 7 or 365")

	// If any errors in our map than re render the create.tmpl.html page
	// with a 422 status code error.
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, "create.tmpl.html", data, http.StatusUnprocessableEntity)
		return
	}

	id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Use the Put() method to add a string value and the corresponding key to
	// session data.
	app.sessionManager.Put(r.Context(), "flash", "Snippet created successfully created!")

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
	data := app.newTemplateData(r)
	data.Snippet = snippet

	app.render(w, "view.tmpl.html", data, http.StatusOK)
}

// userSignUp() renders the user sign up html form.
func (app *application) userSignUp(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userSignUpForm{}

	app.render(w, "signup.tmpl.html", data, http.StatusOK)
}

// userSignUpPost() handles the user sign up post request and validates the form data and
// creates a new user in the database.
func (app *application) userSignUpPost(w http.ResponseWriter, r *http.Request) {
	var form userSignUpForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Check if name field is not blank
	form.CheckField(validator.NotBlank(form.Name), "name", "This field can not be blank")

	// Check if email field is not blank
	form.CheckField(validator.NotBlank(form.Email), "email", "This field can not be blank")

	// Check if email field is correct email format
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")

	// Check if password is not blank
	form.CheckField(validator.NotBlank(form.Password), "password", "This field can not be blank")

	// Check if password is at least 8 characters long
	form.CheckField(validator.MinChars(form.Password, 8), "password", "This field must be at least 8 characters long")

	// If an errors in our map than re render the signup.tmpl.html page
	// with a 422 status code error.
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, "signup.tmpl.html", data, http.StatusUnprocessableEntity)
		return
	}

	err = app.users.Insert(form.Name, form.Email, form.Password)
	if err != nil {
		// check for duplicate email error
		if errors.Is(err, models.ErrDuplicateEmail) {
			// add form error
			form.AddFieldError("email", "Email address is already in use")

			// re render the signup page with the new field error
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, "signup.tmpl.html", data, http.StatusUnprocessableEntity)
		}
	}

	app.sessionManager.Put(r.Context(), "flash", "Your signup was successful. Please log in.")

	// redirect the user to the login page
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)

}

func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Display a html from for loggin in a user...")
}

func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Authenticate and login the user...")
}

func (app *application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Logout the user...")
}
