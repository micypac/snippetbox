package main

import (
	"net/http"
	"net/url"
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


func TestSnippetView(t *testing.T) {
	// Create a new instance of application struct that uses mocked dependencies.
	app := newTestApplication(t)

	// Establish a new test server for running e2e tests.
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	tests := []struct{
		name string
		urlPath string
		wantCode int
		wantBody string
	}{
		{
			name: "Valid ID",
			urlPath: "/snippet/view/1",
			wantCode: http.StatusOK,
			wantBody: "An old silent pond...",
		},
		{
			name: "Non-Existent ID",
			urlPath: "/snippet/view/2",
			wantCode: http.StatusNotFound,
		},
		{
			name: "Negative ID",
			urlPath: "/snippet/view/-1",
			wantCode: http.StatusNotFound,
		},
		{
			name: "Decimal ID",
			urlPath: "/snippet/view/1.23",
			wantCode: http.StatusNotFound,
		},
		{
			name: "String ID",
			urlPath: "/snippet/view/foo",
			wantCode: http.StatusNotFound,
		},
		{
			name: "Empty ID",
			urlPath: "/snippet/view/",
			wantCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, _, body := ts.get(t, tt.urlPath)

			assert.Equal(t, code, tt.wantCode)

			if tt.wantBody != "" {
				assert.StringContains(t, body, tt.wantBody)
			}
		})
	}
}


func TestSnippetCreate(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())

	defer ts.Close()

	t.Run("Unathenticated", func(t *testing.T) {
		code, headers, _ := ts.get(t, "/snippet/create")

		assert.Equal(t, code, http.StatusSeeOther)
		assert.Equal(t, headers.Get("Location"), "/user/login")
	})

	t.Run("Authenticated", func(t *testing.T) {
		// Make a GET /user/login and extract the CSRF token.
		_, _, body := ts.get(t, "/user/login")
		csrfToken := extractCSRFToken(t, body)

		// Make a POST /user/login using the extracted csrf token
		// and mock user model credentials.
		form := url.Values{}
		form.Add("email", "alice@example.com")
		form.Add("password", "pa$$word")
		form.Add("csrf_token", csrfToken)
		ts.postForm(t, "/user/login", form)

		// Check that the authenticated user is shown the create snippet page
		code, _, body := ts.get(t, "/snippet/create")

		assert.Equal(t, code, http.StatusOK)
		assert.StringContains(t, body, `<form action="/snippet/create" method="POST">`)
	})
}

func TestUserSignup(t *testing.T) {
	// Create the app struct containing the mocked dependencies and 
	// set up the test server for running e2e test.
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	// Make a GET /user/signup request and extract the CSRF token from the resp body.
	_, _, body := ts.get(t, "/user/signup")
	validCSRFToken := extractCSRFToken(t, body)

	// t.Logf("CSRF token is: %q", validCSRFToken)

	const (
		validName = "Bob"
		validPassword = "validPa$$word"
		validEmail = "bob@example.com"
		formTag = `<form action="/user/signup" method="POST" novalidate>`
	)

	tests := []struct{
		name 					string
		userName 			string
		userEmail 		string
		userPassword 	string
		csrfToken 		string
		wantCode 			int
		wantFormTag 	string
	}{
		{
			name:					"Valid Submission",
			userName: 		validName,
			userEmail:		validEmail,
			userPassword: validPassword,
			csrfToken: 		validCSRFToken,
			wantCode: 		http.StatusSeeOther,
		},
		{
			name:					"Invalid CSRF Token",
			userName: 		validName,
			userEmail:		validEmail,
			userPassword: validPassword,
			csrfToken: 		"wrongToken",
			wantCode: 		http.StatusBadRequest,
		},
		{
			name:					"Empty name",
			userName: 		"",
			userEmail:		validEmail,
			userPassword: validPassword,
			csrfToken: 		validCSRFToken,
			wantCode: 		http.StatusUnprocessableEntity,
			wantFormTag:  formTag,
		},
		{
			name:					"Empty email",
			userName: 		validName,
			userEmail:		"",
			userPassword: validPassword,
			csrfToken: 		validCSRFToken,
			wantCode: 		http.StatusUnprocessableEntity,
			wantFormTag:  formTag,
		},
		{
			name:					"Empty password",
			userName: 		validName,
			userEmail:		validEmail,
			userPassword: "",
			csrfToken: 		validCSRFToken,
			wantCode: 		http.StatusUnprocessableEntity,
			wantFormTag:  formTag,
		},
		{
			name:					"Invalid email",
			userName: 		validName,
			userEmail:		"bob@example.",
			userPassword: validPassword,
			csrfToken: 		validCSRFToken,
			wantCode: 		http.StatusUnprocessableEntity,
			wantFormTag:  formTag,
		},
		{
			name:					"Short password",
			userName: 		validName,
			userEmail:		validEmail,
			userPassword: "pa$$",
			csrfToken: 		validCSRFToken,
			wantCode: 		http.StatusUnprocessableEntity,
			wantFormTag:  formTag,
		},
		{
			name:					"Duplicate email",
			userName: 		validName,
			userEmail:		"dupe@example.com",
			userPassword: validPassword,
			csrfToken: 		validCSRFToken,
			wantCode: 		http.StatusUnprocessableEntity,
			wantFormTag:  formTag,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			form := url.Values{}
			form.Add("name", tt.userName)
			form.Add("email", tt.userEmail)
			form.Add("password", tt.userPassword)
			form.Add("csrf_token", tt.csrfToken)

			code, _, body := ts.postForm(t, "/user/signup", form)

			assert.Equal(t, code, tt.wantCode)

			if tt.wantFormTag != "" {
				assert.StringContains(t, body, tt.wantFormTag)
			}
		})
	}
}
