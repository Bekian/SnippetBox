package main

import (
	"github.com/Bekian/SnippetBox/internal/models"
	"html/template"
	"path/filepath"
)

// dep injection struct for holding template data
type templateData struct {
	CurrentYear int
	Snippet     models.Snippet
	Snippets    []models.Snippet
}

func newTemplateCache() (map[string]*template.Template, error) {
	// map to act as the cache
	cache := map[string]*template.Template{}
	// grab all files that match the following pattern
	pages, err := filepath.Glob("./ui/html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		// get the file's base name
		name := filepath.Base(page)

		// parse base template into a template set
		ts, err := template.ParseFiles("./ui/html/base.tmpl")
		if err != nil {
			return nil, err
		}

		// this loads the partials onto the base template set
		ts, err = ts.ParseGlob("./ui/html/partials/*.tmpl")
		if err != nil {
			return nil, err
		}

		// add the page to the template set
		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		// add template set to the cache map
		// with the base name as the key
		cache[name] = ts
	}

	return cache, nil
}
