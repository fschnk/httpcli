// Package httpcli it's helper for creating net/http Request, Client
// and executing a request with response processing by user handlers.
package httpcli

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

// The OptionFunc type is an adapter
// that allows to use functions to configure the HTTP client.
// Any package can define its own client configuration functions.
type OptionFunc func(*http.Client) error

// New creates a new http client and configures it using the OptionFunc functions.
func New(options ...OptionFunc) (*http.Client, error) {

	cli := &http.Client{
		Timeout: 20 * time.Second, // 20 sec
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   15 * time.Second,  // 15 sec
				KeepAlive: 300 * time.Second, //  6 min
				DualStack: true,
			}).DialContext,
			MaxIdleConns:          1,
			MaxIdleConnsPerHost:   1,
			IdleConnTimeout:       20 * time.Second, // 20 sec
			TLSHandshakeTimeout:   10 * time.Second, // 10 sec
			ExpectContinueTimeout: 1 * time.Second,  // 1 sec
			DisableCompression:    false,
			DisableKeepAlives:     false,
		},
	}

	for _, opt := range options {
		err := opt(cli)
		if err != nil {
			return nil, fmt.Errorf("apply option: %s", err)
		}
	}

	return cli, nil
}

// TLSTransport returns the OptionFunc function
// to install the client TLS transport with a custom CA.
func TLSTransport(CA string) OptionFunc {
	return func(cli *http.Client) error {
		rootCertPEM, err := ioutil.ReadFile(CA)
		if err != nil {
			return fmt.Errorf("read certificate file: %s", err)
		}

		certPool := x509.NewCertPool()
		ok := certPool.AppendCertsFromPEM(rootCertPEM)
		if !ok {
			return errors.New("can't append root CA")
		}

		t, ok := cli.Transport.(*http.Transport)
		if !ok {
			cli.Transport = &http.Transport{
				TLSClientConfig: &tls.Config{RootCAs: certPool},
			}
			return nil
		}

		t.TLSClientConfig = &tls.Config{RootCAs: certPool}

		cli.Transport = t
		return nil
	}
}
