package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"kibonga/quickbits/internal/models"
	validator "kibonga/quickbits/internal/validators"
	"net/http"
	"strconv"

	"github.com/go-playground/form/v4"
	"github.com/julienschmidt/httprouter"
)

type myType struct{}

type bitCreateForm struct {
	Title               string `form:"title"`
	Content             string `form:"content"`
	ExpiresAt           int    `form:"expires"`
	validator.Validator `form:"-"`
}

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

func (a *app) bitView(w http.ResponseWriter, r *http.Request) {
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

func (a *app) bitCreate(w http.ResponseWriter, r *http.Request) {
	// bit := &models.InsertBit{}
	// if err := json.NewDecoder(r.Body).Decode(&bit); err != nil {
	// 	fmt.Println(err)
	// 	a.clientError(w, http.StatusBadRequest)
	// 	return
	// }
	// defer r.Body.Close()

	form := &bitCreateForm{}
	if err := a.decodePostForm(r, form); err != nil {
		a.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Title), "title", "Title is required")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "Title must be less than 100 characters")
	form.CheckField(validator.NotBlank(form.Content), "content", "Content is required")
	form.CheckField(validator.MaxChars(form.Content, 1000), "content", "Content must be less than 1000 characters")
	form.CheckField(validator.PermittedInt(form.ExpiresAt, 1, 7, 365), "expires", "Expires must be either 1, 7 or 365 days")

	if !form.Valid() {
		data := a.newTemplateData(r)
		data.Form = form
		a.render(w, http.StatusUnprocessableEntity, "create.tmpl", data)
		return
	}

	bit := &models.InsertBit{
		Title:     form.Title,
		Content:   form.Content,
		DaysValid: form.ExpiresAt,
	}
	id, err := a.bitModel.Insert(bit)

	if err != nil {
		a.serverError(w, err)
		return
	}

	a.infoLog.Printf("Inserted bit with ID %d\n", id)
	http.Redirect(w, r, fmt.Sprintf("/bits/view/%d", id), http.StatusSeeOther)
}

func (a *app) bitCreateForm(w http.ResponseWriter, r *http.Request) {
	// w.Write([]byte("Display bit form for creating new bit"))
	data := a.newTemplateData(r)
	data.Form = &bitCreateForm{
		ExpiresAt: 365,
	}
	a.render(w, http.StatusOK, "create.tmpl", data)
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
	http.Redirect(w, r, fmt.Sprintf("/bits/view/%d", id), http.StatusSeeOther)
}

func (a *app) decodePostForm(r *http.Request, dst any) error {
	if err := r.ParseForm(); err != nil {
		return err
	}

	if err := a.formDecoder.Decode(dst, r.PostForm); err != nil {
		var invalidDecoderError *form.InvalidDecoderError

		if errors.As(err, &invalidDecoderError) {
			panic(err)
		}

		return err
	}

	return nil
}
