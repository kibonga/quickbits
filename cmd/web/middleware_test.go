package main

import (
	"bytes"
	"io"
	"kibonga/quickbits/internal/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSecureHeaders(t *testing.T) {
	rr := httptest.NewRecorder()

	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	a := &app{}
	a.secureHeaders(next).ServeHTTP(rr, r)

	rs := rr.Result()

	expectedVal := "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com"
	assert.Equal(t, expectedVal, rs.Header.Get("Content-Security-Policy"))

	expectedVal = "origin-when-cross-origin"
	assert.Equal(t, expectedVal, rs.Header.Get("Referrer-Policy"))

	expectedVal = "nosniff"
	assert.Equal(t, expectedVal, rs.Header.Get("X-Content-Type-Options"))

	expectedVal = "deny"
	assert.Equal(t, expectedVal, rs.Header.Get("X-Frame-Options"))

	expectedVal = "0"
	assert.Equal(t, expectedVal, rs.Header.Get("X-XSS-Protection"))

	assert.Equal(t, rs.StatusCode, http.StatusOK)

	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}

	bytes.TrimSpace(body)
	assert.Equal(t, string(body), "OK")
}
