package main

import (
	"errors"
	"fmt"
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

	bits, err := a.bits.Latest()
	if err != nil {
		a.serverError(w, err)
		return
	}

	for _, b := range bits {
		fmt.Fprintf(w, "%+v\n", b)
	}

	// files := []string{
	// 	fmt.Sprintf("%s%s", a.flags.htmlPath, "ui/html/home.page.tmpl"),
	// 	fmt.Sprintf("%s%s", a.flags.htmlPath, "ui/html/footer.partial.tmpl"),
	// 	fmt.Sprintf("%s%s", a.flags.htmlPath, "ui/html/base.layout.tmpl"),
	// }

	// ts, err := template.ParseFiles(files...)
	// if err != nil {
	// 	a.serverError(w, err)
	// 	return
	// }

	// err = ts.Execute(w, nil)
	// if err != nil {
	// 	a.serverError(w, err)
	// 	return
	// }
}

func (a *app) showBit(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))

	if err != nil || id < 1 {
		a.notFound(w)
		return
	}

	b, err := a.bits.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			a.notFound(w)
			return
		}
		a.serverError(w, err)
	}

	fmt.Fprintf(w, "%+v", b)
}

func (a *app) createBit(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.Header().Set("Allow", http.MethodPost)
		a.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	id, err := a.bits.Insert("title", "content", 1)
	if err != nil {
		a.serverError(w, err)
		return
	}

	a.infoLog.Printf("Inserted bit with ID %d\n", id)
	http.Redirect(w, r, fmt.Sprintf("/bit/view?id=%d", id), http.StatusSeeOther)
}
