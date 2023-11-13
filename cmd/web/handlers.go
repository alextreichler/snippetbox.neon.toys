package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)


func (a *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		a.notFound(w)
		http.NotFound(w,r)
		return
	}

	files := []string{
		"./ui/html/pages/home.gohtml",
		"./ui/html/base.gohtml",
		"./ui/html/partials/nav.gohtml",
	}
	
	ts, err := template.ParseFiles(files...)
	if err != nil {
		a.logger.Error(err.Error(),
			"method", r.Method,
			"uri",r.URL.RequestURI())
		a.serverError(w,r,err)
		return
	}
	err = ts.ExecuteTemplate(w,"base", nil)
	if err != nil {
		a.logger.Error(err.Error(),
			"method", r.Method,
			"uri", r.URL.RequestURI())
		a.serverError(w,r,err)
	}

}





func (a *application) snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		a.notFound(w)
		return
	}
	
	fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
}

func (a *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow",http.MethodPost)
		a.clientError(w,http.StatusNotFound)
		return
	}

	title := "O snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
	expires := 7

	id, err := a.snippets.Insert(title, content, expires)
	if err != nil {
		a.serverError(w,r,err)
		return
	}
	
	http.Redirect(w, r, fmt.Sprintf("/snippet/view?id=%d", id), http.StatusSeeOther)
}
