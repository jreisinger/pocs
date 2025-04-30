package handlers

import (
	"html/template"
	"net/http"
)

var tmpl = template.Must(template.ParseFiles("internal/templates/base.html"))

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	err := tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func LissajousHandler(w http.ResponseWriter, r *http.Request) {
	lissajous(w)
}

func DynamicContentHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`<img src="http://localhost:8080/lissajous" alt="lissajous">`))
}
