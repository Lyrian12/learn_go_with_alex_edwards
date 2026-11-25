package main

import (
	"errors"
	"fmt"
	"gojourney/internal/models"
	"gojourney/internal/validator"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
)

type snippetCreateForm struct {
	Title   string
	Content string
	Expires string
	validator.Validator
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	app.renderCreateForm(w, &snippetCreateForm{Expires: "365"}, http.StatusOK)
}

func (app *application) renderCreateForm(w http.ResponseWriter, form *snippetCreateForm, status int) {
	files := []string{
		"./ui/html/base.tmpl",
		"./ui/html/partials/nav.tmpl",
		"./ui/html/pages/create.tmpl",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(status)

	err = ts.ExecuteTemplate(w, "base", &templateData{Form: form})
	if err != nil {
		app.errorLog.Println(err.Error())
		app.serverError(w, err)
	}
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {

	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	files := []string{
		"./ui/html/base.tmpl",
		"./ui/html/partials/nav.tmpl",
		"./ui/html/pages/home.tmpl",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := &templateData{
		Snippets: snippets,
	}

	w.Header().Set("Content-Type", "text/html")

	err = ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		app.errorLog.Println(err.Error())
		app.serverError(w, err)
	}
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}
	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}
	files := []string{
		"./ui/html/base.tmpl",
		"./ui/html/partials/nav.tmpl",
		"./ui/html/pages/view.tmpl",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := &templateData{
		Snippet: snippet,
		Flash:   app.sessionManager.PopString(r.Context(), "flash"),
	}

	w.Header().Set("Content-Type", "text/html")

	err = ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		app.errorLog.Println(err.Error())
		app.serverError(w, err)
	}
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := &snippetCreateForm{
		Title:   r.PostForm.Get("title"),
		Content: r.PostForm.Get("content"),
		Expires: r.PostForm.Get("expires"),
	}

	expires, err := strconv.Atoi(form.Expires)
	if err != nil {
		form.Check(false, "expires", "Choose a valid expiration period.")
	} else {
		form.Check(validator.PermittedInt(expires, 1, 7, 365), "expires", "Choose a valid expiration period.")
	}

	form.Check(validator.NotBlank(form.Title), "title", "This field is required.")
	form.Check(validator.MaxChars(strings.TrimSpace(form.Title), 50), "title", "Title must be 50 characters or fewer.")
	form.Check(validator.NotBlank(form.Content), "content", "This field is required.")

	if !form.Valid() {
		app.renderCreateForm(w, form, http.StatusUnprocessableEntity)
		return
	}

	title := strings.TrimSpace(form.Title)
	id, err := app.snippets.Insert(title, form.Content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.sessionManager.Put(r.Context(), "flash", "Your snippet was saved successfully.")

	http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
}
