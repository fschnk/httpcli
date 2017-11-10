package httpcli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

// Request protcol schemes.
const (
	SchemeHTTP  = "http"
	SchemeHTTPS = "https"
)

// The RequestOptionFunc type is an adapter
// that allows to use functions to configure the HTTP request.
// Any package can define its own request configuration functions.
type RequestOptionFunc func(*http.Request) error

// NewRequest creates new http request
// and configures it using the RequestOptionFunc functions.
func NewRequest(options ...RequestOptionFunc) (*http.Request, error) {
	req := &http.Request{
		Header: http.Header{},
	}

	for _, opt := range options {
		err := opt(req)
		if err != nil {
			return nil, fmt.Errorf("apply option: %s", err)
		}
	}
	return req, nil
}

// RequestMethod returns the RequestOptionFunc function
// to set the request method.
func RequestMethod(method string) RequestOptionFunc {
	return func(req *http.Request) error {
		req.Method = method
		return nil
	}
}

// RequestURL returns the RequestOptionFunc function
// to set the request URL.
func RequestURL(scheme, host, path string) RequestOptionFunc {
	return func(req *http.Request) error {
		req.URL = &url.URL{
			Scheme: scheme,
			Host:   host,
			Path:   path,
		}
		req.Host = host
		return nil
	}
}

// RequestJSONBody returns the RequestOptionFunc function
// to place v as JSON data and set "Content-Type" to "application/json".
func RequestJSONBody(v interface{}) RequestOptionFunc {
	return func(req *http.Request) error {
		buf, err := json.Marshal(v)
		if err != nil {
			return fmt.Errorf("unmarshal json data: %s", err)
		}

		req.Body = ioutil.NopCloser(bytes.NewBuffer(buf))

		req.Header.Set("Content-Type", "application/json")

		return nil
	}
}

// RequestUserAgent returns the RequestOptionFunc function
// to set request "User-Agent".
func RequestUserAgent(agent string) RequestOptionFunc {
	return func(req *http.Request) error {
		req.Header.Set("User-Agent", agent)
		return nil
	}
}
