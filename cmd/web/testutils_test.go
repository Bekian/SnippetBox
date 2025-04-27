package main

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Bekian/SnippetBox/internal/models/mocks"

	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
)

// helper to create an instance of our application in test
func newTestApplication(t *testing.T) *application {
	// template cache init
	templateCache, err := newTemplateCache()
	if err != nil {
		t.Fatal(err)
	}

	// form decoder
	formDecoder := form.NewDecoder()
	// session manager with in-memory store
	sessionManager := scs.New()
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true

	return &application{
		logger:         slog.New(slog.DiscardHandler),
		snippets:       &mocks.SnippetModel{}, // from mock package
		users:          &mocks.UserModel{},    // from mock package
		templateCache:  templateCache,
		formDecorder:   formDecoder,
		sessionManager: sessionManager,
	}
}

// custom testServer type with a httptest.Server instance
type testServer struct {
	*httptest.Server
}

// helper to create an instance of the above type
func newTestServer(t *testing.T, h http.Handler) *testServer {
	ts := httptest.NewTLSServer(h)

	// cookiejar to hold cookies
	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}
	// response cookies are stored with the client in the jar
	ts.Client().Jar = jar

	// disable redirect following for the client
	// this runs whenever any http status codes in the redirection class 300-399
	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	return &testServer{ts}
}

// implement get method to make GET requests to a given url for testing
// returns status code, headers, and body
func (ts *testServer) get(t *testing.T, urlPath string) (int, http.Header, string) {
	rs, err := ts.Client().Get(ts.URL + urlPath)
	if err != nil {
		t.Fatal(err)
	}

	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	body = bytes.TrimSpace(body)

	return rs.StatusCode, rs.Header, string(body)
}
