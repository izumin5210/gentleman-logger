package logger

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"gopkg.in/h2non/gentleman.v2/plugin"

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
	cases := []struct {
		createPlugin func(out io.Writer) plugin.Plugin
	}{
		{
			createPlugin: func(out io.Writer) plugin.Plugin {
				return New(out)
			},
		},
		{
			createPlugin: func(out io.Writer) plugin.Plugin {
				return FromLogger(log.New(out, "[gentleman] ", log.LstdFlags))
			},
		},
	}

	var (
		path       = "/ping"
		queryKey   = "foo"
		queryValue = "bar"
		reqBody    = "{\"baz\": \"qux\"}"
		respBody   = "{\"message\": \"pong\"}"
		status     = 201
	)

	ctx := newHTTPTestContext()
	defer ctx.server.Close()

	ctx.mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		fmt.Fprint(w, respBody)
	})

	for _, c := range cases {
		buf := bytes.NewBufferString("")

		client := gentleman.New().BaseURL(ctx.server.URL).Use(c.createPlugin(buf))
		_, err := client.Post().
			Path("/ping").
			SetQuery(queryKey, queryValue).
			BodyString(reqBody).
			Send()

		if err != nil {
			t.Fatalf("Unexpected error %v", err)
		}

		if got, want := buf.String(), path; !strings.Contains(got, want) {
			t.Errorf("logged %q, wanna contain path %q", got, want)
		}

		if got, want := buf.String(), fmt.Sprintf("%s=%s", queryKey, queryValue); !strings.Contains(got, want) {
			t.Errorf("logged %q, wanna contain query %q", got, want)
		}

		if got, want := buf.String(), reqBody; !strings.Contains(got, want) {
			t.Errorf("logged %q, wanna contain request body %q", got, want)
		}

		if got, want := buf.String(), fmt.Sprint(status); !strings.Contains(got, want) {
			t.Errorf("logged %q, wanna contain status code %q", got, want)
		}

		if got, want := buf.String(), respBody; !strings.Contains(got, want) {
			t.Errorf("logged %q, wanna contain response body %q", got, want)
		}
	}
}
