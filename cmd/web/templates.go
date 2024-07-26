package main

import (
	"html/template"
	"io/fs"
	"kibonga/quickbits/internal/models"
	"kibonga/quickbits/ui"
	"net/http"
	"path/filepath"
	"time"

	"github.com/justinas/nosurf"
)

type templateData struct {
	Bit                 *models.Bit
	Bits                []*models.Bit
	CopyrightYear       int
	Form                any
	Flash               string
	IsUserAuthenticated bool
	CSRFToken           string
}

func createTemplateCache(htmlPath string) (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	// pages, err := filepath.Glob(fmt.Sprintf("%s%s", htmlPath, "ui/html/pages/*.tmpl"))
	pages, err := fs.Glob(ui.Files, "html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {

		name := filepath.Base(page)

		patterns := []string{
			"html/base.tmpl",
			"html/partials/*.tmpl",
			page,
		}

		tmpl, err := template.New(name).Funcs(functions).ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}

		cache[name] = tmpl
	}

	return cache, nil
}

func (a *app) newTemplateData(r *http.Request) *templateData {
	return &templateData{
		CopyrightYear:       time.Now().Year(),
		Flash:               a.sessionManager.PopString(r.Context(), "flash"),
		IsUserAuthenticated: a.isUserAuthenticated(r),
		CSRFToken:           nosurf.Token(r),
	}
}

func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}
