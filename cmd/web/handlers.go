package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"kibonga/quickbits/internal/models"
	"net/http"
	"strconv"
)

type myType struct{}

func (t *myType) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is my custom handler"))
}

func funcHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is my custom handler function"))
}

func (a *app) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	bits, err := a.bitModel.Latest()
	if err != nil {
		a.serverError(w, err)
		return
	}

	files := []string{
		fmt.Sprintf("%s%s", a.cliFlags.htmlPath, "ui/html/base.layout.tmpl"),
		fmt.Sprintf("%s%s", a.cliFlags.htmlPath, "ui/html/home.page.tmpl"),
		fmt.Sprintf("%s%s", a.cliFlags.htmlPath, "ui/html/footer.partial.tmpl"),
	}

	tmpl, err := template.ParseFiles(files...)
	if err != nil {
		a.serverError(w, err)
		return
	}

	data := &templateData{
		Bits: bits,
	}

	if err := tmpl.ExecuteTemplate(w, "base", data); err != nil {
		a.serverError(w, err)
	}
}

func (a *app) viewBit(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))

	if err != nil || id < 1 {
		a.notFound(w)
		return
	}

	b, err := a.bitModel.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			a.notFound(w)
			return
		}
		a.serverError(w, err)
		return
	}

	files := []string{
		fmt.Sprintf("%s%s", a.cliFlags.htmlPath, "ui/html/base.layout.tmpl"),
		fmt.Sprintf("%s%s", a.cliFlags.htmlPath, "ui/html/footer.partial.tmpl"),
		fmt.Sprintf("%s%s", a.cliFlags.htmlPath, "ui/html/view.page.tmpl"),
	}

	tmpl, err := template.ParseFiles(files...)
	if err != nil {
		a.serverError(w, err)
		return
	}

	data := &templateData{
		Bit: b,
	}

	if err := tmpl.ExecuteTemplate(w, "base", data); err != nil {
		a.serverError(w, err)
		return
	}
}

func (a *app) createBit(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.Header().Set("Allow", http.MethodPost)
		a.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	id, err := a.bitModel.Insert("title", "content", 1)
	if err != nil {
		a.serverError(w, err)
		return
	}

	a.infoLog.Printf("Inserted bit with ID %d\n", id)
	http.Redirect(w, r, fmt.Sprintf("/bit/view?id=%d", id), http.StatusSeeOther)
}

func (a *app) updateBit(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" && r.Method != "PATCH" {
		w.Header().Set("Allow", http.MethodPut+http.MethodPatch)
		a.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		a.notFound(w)
		return
	}

	bit := &models.UpdateBit{}
	if json.NewDecoder(r.Body).Decode(&bit) != nil {
		a.clientError(w, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	err = a.bitModel.Update(id, bit)

	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			a.notFound(w)
			return
		}
		a.serverError(w, err)
		return
	}

	a.infoLog.Printf("updated bit with ID %d\n", id)
	http.Redirect(w, r, fmt.Sprintf("/bit/view?id=%d", id), http.StatusSeeOther)
}
