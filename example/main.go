package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/fschnko/httpcli"
)

var ca = flag.String("ca", "", "custom certificate authority")

func main() {

	cli, err := httpcli.New(httpcli.TLSTransport(*ca))
	if err != nil {
		log.Fatalf("new client: %s", err)
	}

	state := new(State)

	req, err := httpcli.NewRequest(
		httpcli.RequestMethod(http.MethodGet),
		httpcli.RequestURL(httpcli.SchemeHTTP, "127.0.0.1:8080", "/state"),
		httpcli.RequestUserAgent("awesome agent"),
	)
	if err != nil {
		log.Fatalf("new request: %s", err)
	}

	err = httpcli.Execute(cli, req,
		HandleFailJSONResponse(),
		httpcli.HandleSuccessJSONResponse(state),
	)
	if err != nil {
		log.Fatalf("execute API call: %s", err)
	}

	fmt.Printf("%+v\n", state)
}

// HandleFailJSONResponse handles JSON fail response.
func HandleFailJSONResponse() httpcli.ResponseHandlerFunc {
	return func(resp *http.Response) error {
		if resp.StatusCode >= 400 && resp.StatusCode <= 599 {
			buf, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return fmt.Errorf("read response body: %s", err)
			}
			apierr := new(Error)
			err = json.Unmarshal(buf, apierr)
			if err != nil {
				return fmt.Errorf("unmarshal json: %s", err)
			}
			return apierr
		}
		return nil
	}
}

// State represents value parameter in json.
type State struct {
	Value string `json:"value"`
}

// Error represents error parameter in json.
type Error struct {
	Text string `json:"error"`
}

// Error returns text of the error.
func (e *Error) Error() string {
	return e.Text
}
