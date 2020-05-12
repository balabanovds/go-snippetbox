package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newTestApplication() *application {
	return &application{
		errLog:  log.New(ioutil.Discard, "", 0),
		infoLog: log.New(ioutil.Discard, "", 0),
	}
}

type testServer struct {
	*httptest.Server
	t *testing.T
}

func newTestServer(t *testing.T, h http.Handler) *testServer {
	ts := httptest.NewTLSServer(h)
	return &testServer{ts, t}
}

func (ts *testServer) get(path string) (int, http.Header, []byte) {
	rs, err := ts.Client().Get(ts.URL + path)
	ts.handleFatal(err)

	defer rs.Body.Close()
	body, err := ioutil.ReadAll(rs.Body)
	ts.handleFatal(err)

	return rs.StatusCode, rs.Header, body
}

func (ts *testServer) handleFatal(err error) {
	if err != nil {
		ts.t.Fatal(err)
	}
}
