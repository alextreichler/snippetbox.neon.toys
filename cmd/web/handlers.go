package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/alextreichler/snippetbox.neon.toys/internal/models"
	"github.com/alextreichler/snippetbox.neon.toys/internal/validator"
	"github.com/julienschmidt/httprouter"
)

func (a *application) home(w http.ResponseWriter, r *http.Request) {

	snippets, err := a.snippets.Latest()
	if err != nil {
		a.serverError(w, r, err)
	}

	data := a.newTemplateData(r)
	data.Snippets = snippets

	a.render(w, r, http.StatusOK, "home.gotmpl", data)

}

func (a *application) snippetView(w http.ResponseWriter, r *http.Request) {

	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		a.notFound(w)
		return
	}

	snippet, err := a.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			a.notFound(w)
		} else {
			a.serverError(w, r, err)
		}
		return
	}

	data := a.newTemplateData(r)
	data.Snippet = snippet

	a.render(w, r, http.StatusOK, "view.gotmpl", data)
}

func (a *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	data := a.newTemplateData(r)

	data.Form = snippetCreateForm{
		Expires: 365,
	}

	a.render(w, r, http.StatusOK, "create.gotmpl", data)
}

type snippetCreateForm struct {
	Title               string `form:"title"`
	Content             string `form:"content"`
	Expires             int    `form:"expires"`
	validator.Validator `form:"-"`
}

func (a *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	var form snippetCreateForm

	err := a.decodePostForm(r, &form)
	if err != nil {
		a.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Title),
		"title", "This field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 100),
		"title", "This field cannot be more than 100 characters long")
	form.CheckField(validator.NotBlank(form.Content),
		"content", "This field cannot be blank")
	form.CheckField(validator.PermittedValue(form.Expires, 1, 7, 365),
		"expires", "This field must equal 1, 7 or 365")

	if !form.Valid() {
		data := a.newTemplateData(r)
		data.Form = form
		a.render(w, r, http.StatusUnprocessableEntity,
			"create.gotmpl", data)
		return
	}
	id, err := a.snippets.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		a.serverError(w, r, err)
		return
	}

	a.sessionManager.Put(r.Context(), "flash", "Snippet successfully created!")

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id),
		http.StatusSeeOther)
}
