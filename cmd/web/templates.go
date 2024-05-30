package main

import (
	"html/template"
	"path/filepath"
	"time"

	"snippetbox.micypac.io/internal/models"
)


type templateData struct {
	CurrentYear 		int
	Snippet 				*models.Snippet
	Snippets 				[]*models.Snippet
	Form 						any
	Flash 					string
	IsAuthenticated bool
}

func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}


func newTemplateCache() (map[string]*template.Template, error) {
	// Initialize new map to act as the cache
	cache := map[string]*template.Template{}

	/*
		Use the filepath.Glob() func to get a slice of all the filepaths that match the
		pattern "./ui/html/pages/*.tmpl". 
	*/
	pages, err := filepath.Glob("./ui/html/pages/*.tmpl.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		// files := []string{
		// 	"./ui/html/base.tmpl.html",
		// 	"./ui/html/partials/nav.tmpl.html",
		// 	page,
		// }

		// template.Funcs() need a new template set before calling ParseFiles.
		ts:= template.New(name)
		ts = ts.Funcs(functions)

		// Parse the base template file into a template set.
		ts, err := ts.ParseFiles("./ui/html/base.tmpl.html")
		if err != nil {
			return nil, err
		}

		// Call ParseGlob() on this template set to add any partials.
		ts, err = ts.ParseGlob("./ui/html/partials/*.tmpl.html")
		if err != nil {
			return nil, err
		}

		// Call ParseFiles() on this template set to add the page template.
		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}
