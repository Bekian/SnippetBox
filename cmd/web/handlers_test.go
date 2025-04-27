package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Bekian/SnippetBox/internal/assert"
)

func TestPing(t *testing.T) {
	// recorder to record http responses
	rr := httptest.NewRecorder()

	// send dummy req
	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// call ping
	ping(rr, r)
	// get the result from the above request
	rs := rr.Result()
	// ensure status code is what we expect
	assert.Equal(t, rs.StatusCode, http.StatusOK)

	// check that the response body is also what we expect
	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	body = bytes.TrimSpace(body)
	assert.Equal(t, string(body), "OK")

}
