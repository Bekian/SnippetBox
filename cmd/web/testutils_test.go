package main

import (
	"bytes"
	"html"
	"io"
	"log/slog"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"regexp"
	"testing"
	"time"

	"github.com/Bekian/SnippetBox/internal/models/mocks"

	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
)

// regex expression to extract csrf token
var csrfTokenRX = regexp.MustCompile(`<input type='hidden' name='csrf_token' value='(.+)'>`)

func extractCSRFToken(t *testing.T, body string) string {
	matches := csrfTokenRX.FindStringSubmatch(body)
	if len(matches) < 2 {
		t.Fatal("no csrf token found in body")
	}

	return html.UnescapeString(matches[1])
}

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

// implements post method
// the form parameter is any values to pass to the form body
func (ts *testServer) postForm(t *testing.T, urlPath string, form url.Values) (int, http.Header, string) {
	// call endpoint
	rs, err := ts.Client().PostForm(ts.URL+urlPath, form)
	if err != nil {
		t.Fatal(err)
	}
	// read response
	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	body = bytes.TrimSpace(body)
	// return status, headers, and body
	return rs.StatusCode, rs.Header, string(body)

}
