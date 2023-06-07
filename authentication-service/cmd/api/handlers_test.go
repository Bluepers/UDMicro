package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

type RoundTripFunc func(req *http.Request) *http.Response

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

func NewTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: fn,
	}
}
func TestAuthenticate(t *testing.T) {
	jsonToResponse := `
{
	"error": false,
	"message": "A message",
}`

	client := NewTestClient(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewBufferString(jsonToResponse)),
			Header:     make(http.Header),
		}
	})

	testApp.Client = client

	postBody := map[string]interface{}{
		"email":    "email@test.com",
		"password": "verysecret",
	}

	body, _ := json.Marshal(postBody)

	req, _ := http.NewRequest("POST", "/authenticate", bytes.NewReader(body))

	resrec := httptest.NewRecorder()

	handler := http.HandlerFunc(testApp.Authenticate)

	handler.ServeHTTP(resrec, req)

	if resrec.Code != http.StatusAccepted {
		t.Error("Error: expected http.StatusAccepted but got " + strconv.Itoa(resrec.Code))
	}
}
