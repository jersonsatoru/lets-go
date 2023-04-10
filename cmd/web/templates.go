package main

import (
	"html/template"
	"path/filepath"
	"time"

	"github.com/jersonsatoru/lets-go/internal/models"
)

type TemplateData struct {
	Snippet  *models.Snippet
	Snippets []*models.Snippet
	Year     int
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

func humanDate(date time.Time) string {
	return date.Format("02 Jan 2006 at 15:04")
}

func newTemplateCache() (map[string]*template.Template, error) {
	files, err := filepath.Glob("./ui/html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}

	templates := make(map[string]*template.Template)

	for _, file := range files {
		base := filepath.Base(file)

		ts, err := template.New(base).Funcs(functions).ParseFiles("./ui/html/base.tmpl")
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob("./ui/html/partials/*.tmpl")
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseFiles(file)
		if err != nil {
			return nil, err
		}

		templates[base] = ts
	}

	return templates, nil
}
