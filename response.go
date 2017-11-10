package httpcli

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

// The ResponseHandlerFunc type is an adapter
// that allows to use functions to handle HTTP response.
// Any package can define its own response handle functions.
type ResponseHandlerFunc func(resp *http.Response) error

// Execute executes the http request
// and calls all the handlers for the response until an error occurs.
func Execute(client *http.Client, req *http.Request, handlers ...ResponseHandlerFunc) error {

	resp, err := client.Do(req)
	if resp != nil {
		defer func(c io.Closer) {
			ignoreErr(c.Close())
		}(resp.Body)
	}
	if err != nil {
		return fmt.Errorf("send request: %s", err)
	}
	for _, handle := range handlers {
		err = handle(resp)
		if err != nil {
			return fmt.Errorf("handle: %s", err)
		}
	}
	return nil
}

// HandleSuccessJSONResponse returns the ResponseHandlerFunc
// to read the JSON response data, if the StatusCode is in the range 200-299.
func HandleSuccessJSONResponse(v interface{}) ResponseHandlerFunc {
	return func(resp *http.Response) error {
		if resp.StatusCode >= 200 && resp.StatusCode <= 299 && v != nil {
			buf, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return fmt.Errorf("read response body: %s", err)
			}

			err = json.Unmarshal(buf, v)
			if err != nil {
				return fmt.Errorf("unmarshal json: %s", err)
			}
		}
		return nil
	}
}

// HandleFailJSONResponse returns the ResponseHandlerFunc
// to read the JSON response data, if the StatusCode is in the range 400-599.
func HandleFailJSONResponse(v interface{}) ResponseHandlerFunc {
	return func(resp *http.Response) error {
		if resp.StatusCode >= 400 && resp.StatusCode <= 599 && v != nil {
			buf, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return fmt.Errorf("read response body: %s", err)
			}

			err = json.Unmarshal(buf, v)
			if err != nil {
				return fmt.Errorf("unmarshal json: %s", err)
			}
		}
		return nil
	}
}

// ignore helps ignore errors
func ignoreErr(err error) {
}
