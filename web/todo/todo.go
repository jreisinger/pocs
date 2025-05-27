package todo

type Todo struct {
	ID   int
	Text string
	Done bool
}

var Todos = []Todo{
	{ID: 1, Text: "Learn templ", Done: false},
	{ID: 2, Text: "Learn htmx", Done: false},
	{ID: 3, Text: "Learn CSS", Done: false},
}
