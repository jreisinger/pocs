package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"webstats"
)

var urls []string

func main() {
	if len(os.Args[1:]) == 0 {
		log.Fatal("supply at least one URL")
	}
	urls = os.Args[1:]

	http.HandleFunc("/", statsPageHandler)

	addr := "localhost:8080"
	log.Printf("starting server at http://%s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

var tmpl = template.Must(template.New("stats").Funcs(template.FuncMap{
	"humanReadableSize": webstats.HumanReadableSize,
}).Parse(`
<html>
<body>
<h1>Website Fetch Stats</h1>
<table border='1'>
<tr>
<th>URL</th>
<th>Time</th>
<th>Size</th>
<th>Error</th>
</tr>
{{range .}}
<tr>
<td>{{.URL}}</td>
<td>{{printf "%.3fs" .FetchTime.Seconds}}</td>
<td align="right">{{humanReadableSize .FetchSize}}</td>
<td>{{.Err}}</td>
</tr>
{{end}}
</table>
</body>
</html>
`))

func statsPageHandler(w http.ResponseWriter, r *http.Request) {
	stats := webstats.Get(urls)
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	if err := tmpl.Execute(w, stats); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
