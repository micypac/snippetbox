package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"snippetbox.micypac.io/internal/assert"
)

func TestSecureHeaders(t *testing.T) {
	
	// Initialize a new ResponseRecorder and dummy http.Request.
	rr := httptest.NewRecorder()
	
	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a mock HTTP handler that we can pass to our secureHeaders middleware,
	// which writes a 200 status code and an "OK" response body.
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	// Pass the mock HTTP handler to our secureHeaders mid. Our middleware *returns* a http.Handler,
	// we can call its ServeHTTP(), passing in the ResponseRecorder and dummy request to execute it.
	secureHeaders(next).ServeHTTP(rr, r)

	// Call the Result() method on the ResponseRecorder.
	rs := rr.Result()

	// Check the middleware has correctly sets the response headers.
	want := "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com"
	got := rs.Header.Get("Content-Security-Policy")
	assert.Equal(t, got, want)

	want = "origin-when-cross-origin"
	got = rs.Header.Get("Referrer-Policy")
	assert.Equal(t, got, want)

	want = "nosniff"
	got = rs.Header.Get("X-Content-Type-Options")
	assert.Equal(t, got, want)

	want = "deny"
	got = rs.Header.Get("X-Frame-Options")
	assert.Equal(t, got, want)

	want = "0"
	got = rs.Header.Get("X-XSS-Protection")
	assert.Equal(t, got, want)

	// Check the middleware has correctly called the next handler in line and the resp status code
	// and body are as expected.
	assert.Equal(t, rs.StatusCode, http.StatusOK)

	defer rs.Body.Close()

	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}

	bytes.TrimSpace(body)

	assert.Equal(t, string(body), "OK")

}
