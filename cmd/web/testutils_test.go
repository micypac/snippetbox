package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"testing"
	"time"

	"snippetbox.micypac.io/internal/models/mocks"

	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
)

// Create a newTestAppication helper that return an instance of app struct containing
// mocked dependencies.
func newTestApplication(t *testing.T) *application {
	// Create instance of template cache.
	templateCache, err := newTemplateCache()
	if err != nil {
		t.Fatal(err)
	}

	// Create instance for form decoder.
	formDecoder := form.NewDecoder()

	// Create session manager instance.
	sessionManager := scs.New()
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true

	return &application{
		errorLog: log.New(io.Discard, "", 0),
		infoLog: log.New(io.Discard, "", 0),
		snippets: &mocks.SnippetModel{},
		users: &mocks.UserModel{},
		templateCache: templateCache,
		formDecoder: formDecoder,
		sessionManager: sessionManager,
	}
}

// Define a custom testServer type which embeds a httptest.Server instance.
type testServer struct {
	*httptest.Server
}


// Create a newTestServer he;per which initializes and returns a new instance of 
// custom testServer type.
func newTestServer(t *testing.T, h http.Handler) *testServer {
	// Initialize a new test server.
	ts := httptest.NewTLSServer(h)

	// Initialize a new cookie jar.
	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}

	// Add the cookie jar to the test server client. Any response cookies will now be stored
	// and sent with subsequent requests when using the client.
	ts.Client().Jar = jar

	// Disable redirect-following for the test server client by setting a custom CheckRedirect func.
	// This func will be called whenever a 3xx response is received by the client, and by always
	// returning a http.ErrUseLastResponse error it forces the client to immediately return a received
	// response.
	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	return &testServer{ts}
}


// Implement a get() method on custom testServer type. This makes GET request to a 
// given url path using the test server client, and returns the response status code, headers, & body.
func (ts *testServer) get(t *testing.T, urlPath string) (int, http.Header, string) {
	res, err := ts.Client().Get(ts.URL + urlPath)
	if err != nil {
		t.Fatal(err)
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	bytes.TrimSpace(body)

	return res.StatusCode, res.Header, string(body)
}
