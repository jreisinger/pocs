package main

import (
	"bytes"
	"html/template"
	"net/http"
)

// convertLinks converts [PageName] to <a href="/view/PageName>PageName</a>"
func convertLinks(body []byte) template.HTML {
	var escapedBody = new(bytes.Buffer)
	template.HTMLEscape(escapedBody, body)
	converted := linkRegex.ReplaceAllFunc(escapedBody.Bytes(), func(match []byte) []byte {
		pageName := match[1 : len(match)-1] // remove the square brackets
		return []byte(`<a href="/view/` + string(pageName) + `">` + string(pageName) + `</a>`)
	})
	return template.HTML(converted)
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	p.BodyWithLinks = convertLinks(p.Body)
	renderTemplate(w, "view", *p)
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", *p)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := r.PostFormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+p.Title, http.StatusFound)
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2]) // The title is the second subexpression.
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/view/FrontPage", http.StatusFound)
}
