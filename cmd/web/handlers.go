package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"kibonga/quickbits/internal/models"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type myType struct{}

func (t *myType) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is my custom handler"))
}

func funcHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is my custom handler function"))
}

func (a *app) bitsIndex(w http.ResponseWriter, r *http.Request) {
	bits, err := a.bitModel.Latest()
	if err != nil {
		a.serverError(w, err)
		return
	}

	data := a.newTemplateData(r)
	data.Bits = bits

	a.render(w, http.StatusOK, "home.tmpl", data)
}

func (a *app) bitsView(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(params.ByName("id"))

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

	data := a.newTemplateData(r)
	data.Bit = b

	a.render(w, http.StatusOK, "view.tmpl", data)
}

func (a *app) bitsCreate(w http.ResponseWriter, r *http.Request) {
	id, err := a.bitModel.Insert("title", "content", 1)
	if err != nil {
		a.serverError(w, err)
		return
	}

	a.infoLog.Printf("Inserted bit with ID %d\n", id)
	http.Redirect(w, r, fmt.Sprintf("/bit/view?id=%d", id), http.StatusSeeOther)
}

func (a *app) bitsNew(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display bit form for creating new bit"))
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
	http.Redirect(w, r, fmt.Sprintf("/bits/view?id=%d", id), http.StatusSeeOther)
}
