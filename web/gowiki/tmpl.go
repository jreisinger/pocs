package main

import (
	"html/template"
	"net/http"
	"regexp"
)

var (
	templates = make(map[string]*template.Template)
	validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")
	linkRegex = regexp.MustCompile(`\[(\w+)\]`) // Regular expression to match [PageName]
)

func init() {
	templates["view"] = template.Must(template.ParseFiles("tmpl/view.html", "tmpl/base.html"))
	templates["edit"] = template.Must(template.ParseFiles("tmpl/edit.html", "tmpl/base.html"))
}

func renderTemplate(w http.ResponseWriter, name string, p Page) {
	err := templates[name].ExecuteTemplate(w, "base", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
