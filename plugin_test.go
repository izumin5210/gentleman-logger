package logger

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"gopkg.in/h2non/gentleman.v2"
)

type httpTestContext struct {
	mux    *http.ServeMux
	server *httptest.Server
}

func newHTTPTestContext() *httpTestContext {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	return &httpTestContext{
		mux:    mux,
		server: server,
	}
}

func Test_LoggerRoundTripper(t *testing.T) {
	ctx := newHTTPTestContext()
	defer ctx.server.Close()

	buf := bytes.NewBufferString("")

	ctx.mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		fmt.Fprint(w, "{\"message\": \"pong\"}")
	})

	client := gentleman.New().BaseURL(ctx.server.URL).Use(New(buf))

	_, err := client.Request().Path("/ping").Send()

	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	if got, want := buf.String(), "/ping"; !strings.Contains(got, want) {
		t.Errorf("logged %q, wanna contain path %q", got, want)
	}

	if got, want := buf.String(), "200"; !strings.Contains(got, want) {
		t.Errorf("logged %q, wanna contain status code %q", got, want)
	}

	if got, want := buf.String(), "pong"; !strings.Contains(got, want) {
		t.Errorf("logged %q, wanna contain response body %q", got, want)
	}
}
