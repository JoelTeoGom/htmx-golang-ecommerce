package handlers

import (
	"html/template"
	"net/http"
)

func HomeHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			tmpl := template.Must(template.ParseFiles("templates/home.html"))
			tmpl.Execute(w, nil)
			return
		}
	}
}
