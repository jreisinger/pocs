package main

import (
	"html/template"
	"os"
	"path/filepath"
)

const dataDir = "data"

type Page struct {
	Title         string
	Body          []byte
	BodyWithLinks template.HTML
}

func (p *Page) save() error {
	path := filepath.Join(dataDir, p.Title+".txt")
	return os.WriteFile(path, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	path := filepath.Join(dataDir, title+".txt")
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: b}, nil
}
