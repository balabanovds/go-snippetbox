package main

import (
	"net/http"
	"testing"
)

func TestPing(t *testing.T) {

	app := newTestApplication()

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	statusCode, _, body := ts.get("/ping")

	if statusCode != http.StatusOK {
		t.Errorf("want %d; got %d", http.StatusOK, statusCode)
	}

	if string(body) != "OK" {
		t.Errorf("want body 'OK', got %s", body)
	}
}
