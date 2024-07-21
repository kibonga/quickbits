package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
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
