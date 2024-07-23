package main

import (
	"fmt"
	"html/template"
	"kibonga/quickbits/internal/models"
	"net/http"
	"path/filepath"
	"time"
)

type templateData struct {
	Bit           *models.Bit
	Bits          []*models.Bit
	CopyrightYear int
}

func createTemplateCache(htmlPath string) (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob(fmt.Sprintf("%s%s", htmlPath, "ui/html/pages/*.tmpl"))
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		tmpl, err := template.ParseFiles(fmt.Sprintf("%s%s", htmlPath, "ui/html/base.tmpl"))
		if err != nil {
			return nil, err
		}

		tmpl, err = tmpl.ParseGlob(fmt.Sprintf("%s%s", htmlPath, "ui/html/partials/*.tmpl"))
		if err != nil {
			return nil, err
		}

		tmpl, err = tmpl.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		cache[name] = tmpl
	}

	return cache, nil
}

func (a *app) newTemplateData(r *http.Request) *templateData {
	return &templateData{
		CopyrightYear: time.Now().Year(),
	}
}
