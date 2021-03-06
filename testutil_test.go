package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"maunium.net/go/mautrix/id"
)

// newResponse creates a new HTTP response with the given data.
func newResponse(statusCode int, body string) *http.Response {
	return &http.Response{
		StatusCode: statusCode,
		Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
	}
}

// matrixTripper mocks out RoundTrip and calls a registered handler instead.
type matrixTripper struct {
	handlers map[string]func(req *http.Request) (*http.Response, error)
}

func newMatrixTripper() *matrixTripper {
	return &matrixTripper{
		handlers: make(map[string]func(req *http.Request) (*http.Response, error)),
	}
}

func (rt *matrixTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	key := req.Method + " " + req.URL.Path
	if handler, ok := rt.handlers[key]; ok {
		return handler(req)
	}
	for strMatch, handler := range rt.handlers {
		// try to match key with wildcard handlers
		if strMatch[len(strMatch)-1] == '*' && strings.HasPrefix(key, strMatch[:len(strMatch)-1]) {
			return handler(req)
		}
	}
	panic(fmt.Sprintf(
		"RoundTrip: Unhandled request: %s\nHandlers: %d",
		key, len(rt.handlers),
	))
}

func (rt *matrixTripper) Handle(method, path string, handler func(req *http.Request) (*http.Response, error)) {
	key := method + " " + path
	if _, exists := rt.handlers[key]; exists {
		panic(fmt.Sprintf("Test handler with key %s already exists", key))
	}
	rt.handlers[key] = handler
}

func (rt *matrixTripper) HandlePOSTFilter(userID id.UserID) {
	rt.Handle("POST", "/_matrix/client/r0/user/"+userID.String()+"/filter",
		func(req *http.Request) (*http.Response, error) {
			return newResponse(200, `{
				"filter_id":"abcdef"
			}`), nil
		},
	)
}

func (rt *matrixTripper) ClearHandlers() {
	for k := range rt.handlers {
		delete(rt.handlers, k)
	}
}
