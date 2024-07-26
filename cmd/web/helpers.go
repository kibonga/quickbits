package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/go-playground/form/v4"
)

func (a *app) httpError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (a *app) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())

	a.errorLog.Output(2, trace)

	a.httpError(w, http.StatusInternalServerError)
}

func (a *app) clientError(w http.ResponseWriter, status int) {
	a.httpError(w, status)
}

func (a *app) notFound(w http.ResponseWriter) {
	a.clientError(w, http.StatusNotFound)
}

func (a *app) render(w http.ResponseWriter, status int, page string, data *templateData) {
	tmpl, ok := a.templateCache[page]

	buff := &bytes.Buffer{}

	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		a.serverError(w, err)
		return
	}

	if err := tmpl.ExecuteTemplate(buff, "base", data); err != nil {
		a.serverError(w, err)
	}

	w.WriteHeader(status)

	buff.WriteTo(w)
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

func (a *app) isUserAuthenticated(ctx context.Context) bool {
	return a.sessionManager.Exists(ctx, "authUserId")
}
