package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"kibonga/quickbits/internal/models"
	validator "kibonga/quickbits/internal/validators"
	"net/http"
	"strconv"

	formDecoder "github.com/go-playground/form/v4"
	"github.com/julienschmidt/httprouter"
)

type myType struct{}

type bitCreateForm struct {
	Title               string `form:"title"`
	Content             string `form:"content"`
	ExpiresAt           int    `form:"expires"`
	validator.Validator `form:"-"`
}

type userSignupForm struct {
	Name                string `form:"name"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

type userLoginForm struct {
	Email               string `form:"email"`
	Password            string `form:"password"`
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

func (a *app) bitsCreatePost(w http.ResponseWriter, r *http.Request) {
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

	a.sessionManager.Put(r.Context(), "flash", "Bit created successfully!")

	a.infoLog.Printf("Inserted bit with ID %d\n", id)
	http.Redirect(w, r, fmt.Sprintf("/bits/view/%d", id), http.StatusSeeOther)
}

func (a *app) bitsCreate(w http.ResponseWriter, r *http.Request) {
	// w.Write([]byte("Display bit form for creating new bit"))
	data := a.newTemplateData(r)
	data.Form = &bitCreateForm{
		ExpiresAt: 365,
	}
	a.render(w, http.StatusOK, "create.tmpl", data)
}

func (a *app) bitsUpdate(w http.ResponseWriter, r *http.Request) {
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

func (a *app) userSignup(w http.ResponseWriter, r *http.Request) {
	data := a.newTemplateData(r)
	data.Form = &userSignupForm{}
	a.render(w, http.StatusOK, "signup.tmpl", data)
}

func (a *app) userSignupPost(w http.ResponseWriter, r *http.Request) {
	form := &userSignupForm{}
	a.decodePostForm(r, form)

	form.CheckField(validator.NotBlank(form.Name), "name", "Name is required")
	form.CheckField(validator.MaxChars(form.Name, 255), "name", "Name must be less than 255 characters long")
	form.CheckField(validator.NotBlank(form.Email), "email", "Email is required")
	form.CheckField(validator.ValidEmail(form.Email), "email", "Invalid email format")
	form.CheckField(validator.NotBlank(form.Password), "password", "Password is required")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "Password must contain at least 8 characters")
	form.CheckField(validator.MaxChars(form.Password, 30), "password", "Password cannot have more than 30 characters")

	data := a.newTemplateData(r)
	data.Form = form

	if !form.Valid() {
		a.render(w, http.StatusUnprocessableEntity, "signup.tmpl", data)
		return
	}

	u := &models.UserSignupModel{
		Name:     form.Name,
		Email:    form.Email,
		Password: form.Password,
	}

	id, err := a.userModel.Insert(u)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.Validator.AddFieldError("email", "Email already in use")
			a.render(w, http.StatusUnprocessableEntity, "signup.tmpl", data)
			return
		}
		a.serverError(w, err)
		return
	}

	a.sessionManager.Put(r.Context(), "flash", "Registered successfully")
	a.infoLog.Printf("Inserted user with id=%d", id)

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (a *app) userLogin(w http.ResponseWriter, r *http.Request) {
	data := a.newTemplateData(r)
	data.Form = &userLoginForm{}
	a.render(w, http.StatusOK, "login.tmpl", data)
}

func (a *app) userLoginPost(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Println("failed to parse form")
	}

	form := &userLoginForm{}
	if err := formDecoder.NewDecoder().Decode(form, r.PostForm); err != nil {
		var invalidDecoderError *formDecoder.InvalidDecoderError
		if errors.As(err, &invalidDecoderError) {
			panic(err)
		}
		a.serverError(w, err)
		return
	}

	// Validate payload
	form.CheckField(validator.NotBlank(form.Email), "email", "Email is required")
	form.CheckField(validator.ValidEmail(form.Email), "email", "Invalid email format")
	form.CheckField(validator.NotBlank(form.Password), "password", "Password is required")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "Password must contain at least 8 characters")
	form.CheckField(validator.MaxChars(form.Password, 30), "password", "Password cannot have more than 30 characters")

	// Init template data
	data := a.newTemplateData(r)
	data.Form = form

	// Handle validation errors
	if !form.Valid() {
		a.render(w, http.StatusUnprocessableEntity, "login.tmpl", data)
		return
	}

	// Create user model data
	u := &models.UserLoginModel{
		Email:   form.Email,
		Pasword: form.Password,
	}

	// Perform login (check if user exists and if password is valid)
	id, err := a.userModel.Auth(u)
	// Handle invalid creds
	if err != nil {
		if errors.Is(err, models.ErrInvalidCreds) {
			// a.sessionManager.Put(r.Context(), "flash", "Email or password is invalid")
			a.infoLog.Printf("user with email=%s failed to login", form.Email)
			form.GeneralErrors = append(form.GeneralErrors, "Password or email are invalid")
			a.render(w, http.StatusUnprocessableEntity, "login.tmpl", data)
			return
		}
		a.serverError(w, err)
		return
	}

	a.infoLog.Printf("user with id=%d logged in", id)

	if err := a.sessionManager.RenewToken(r.Context()); err != nil {
		a.serverError(w, err)
		return
	}

	// Setup session (auth success)
	a.sessionManager.Put(r.Context(), "authUserId", id)
	a.sessionManager.Put(r.Context(), "flash", "Welcome back")

	// Redirect
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (a *app) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	// a.sessionManager.Remove(r.Context(), "")
	if err := a.sessionManager.RenewToken(r.Context()); err != nil {
		a.serverError(w, err)
		return
	}
	a.sessionManager.Remove(r.Context(), "authUserId")
	a.sessionManager.Put(r.Context(), "flash", "Successfully logged out")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
