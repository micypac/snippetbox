package main

import (
	"net/http"
	"testing"

	"snippetbox.micypac.io/internal/assert"
)

// func TestPing(t *testing.T) {
// 	// Create a new instance of our app struct. Include mock loggers(which discard
// 	// anything written to them). These are needed by the logRequest and recoverPanic middlewares.
// 	// Running without these 2 dependencies will result in a panic.
// 	app := &application{
// 		errorLog: log.New(io.Discard, "", 0),
// 		infoLog: log.New(io.Discard, "", 0),
// 	}

// 	// Use httptest.NewTLSServer() func to create a new test server, passing in the value returned by
// 	// our app.routes() method as the handler for the server. This starts a HTTPS server w/c listens
// 	// on a random port of your local machine for the duration of the test.
// 	// Defer ts.Close() so the server is shutdown when it finishes.
// 	ts := httptest.NewTLSServer(app.routes())
// 	defer ts.Close()

// 	// The network address the test sever is listening on is contained in ts.URL field.
// 	// Use this along ts.Client().Get() method to make a GET /ping request against the test server.
// 	resp, err := ts.Client().Get(ts.URL + "/ping")
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	// Check response status code and body.
// 	assert.Equal(t, resp.StatusCode, http.StatusOK)

// 	defer resp.Body.Close()

// 	body, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	bytes.TrimSpace(body)

// 	assert.Equal(t, string(body), "OK")
// }

func TestPing(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	code, _, body := ts.get(t, "/ping")

	assert.Equal(t, code, http.StatusOK)
	assert.Equal(t, body, "OK")
}
