// Package httpclient ...
package httpclient

import (
	"io"
	"net/http"
	"net/url"
	"objectmap/useragent"
	"time"

	"github.com/pkg/errors"
)

// HTTPClient keeps client and its configuration
type HTTPClient struct {
	Client         *http.Client
	userAgent      useragent.UserAgent
	followRedirect bool
}

// New creates a new HTTPClient
func New(timeout time.Duration, ua useragent.UserAgent, noFollow bool) HTTPClient {
	hc := HTTPClient{
		userAgent:      ua,
		followRedirect: !noFollow,
		Client: &http.Client{
			Timeout: timeout,
		},
	}

	// nolint:unparam
	hc.Client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		if hc.followRedirect {
			return nil
		}

		return http.ErrUseLastResponse
	}

	return hc
}

// NewRequest creates a new http.Request
func (h *HTTPClient) NewRequest(method string, url *url.URL, body io.Reader) (*http.Request, error) {
	r, err := http.NewRequest(method, url.String(), body)
	if err != nil {
		return r, errors.Errorf("Invalid request, %s", err)
	}
	r.Header.Add("User-Agent", string(h.userAgent))

	return r, nil
}

//// NewRequestFromCLI
//func NewRequestFromCLI(input cli.Input) {
//
//}

// NewRequestFromFile calls RequestFromFile
func (h *HTTPClient) NewRequestFromFile(filepath string) (*http.Request, error) {
	return RequestFromFile(filepath)
}

// ReadBody reads body without draining it
func (h *HTTPClient) ReadBody(body *io.ReadCloser) ([]byte, error) {
	return ReadBody(body)
}
