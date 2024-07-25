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
	Form          any
	Flash         string
}

func createTemplateCache(htmlPath string) (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob(fmt.Sprintf("%s%s", htmlPath, "ui/html/pages/*.tmpl"))
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		tmpl := template.New(name).Funcs(functions)

		tmpl, err := tmpl.ParseFiles(fmt.Sprintf("%s%s", htmlPath, "ui/html/base.tmpl"))
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
		Flash:         a.sessionManager.PopString(r.Context(), "flash"),
	}
}

func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}
