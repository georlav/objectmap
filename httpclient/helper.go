package httpclient

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
)

// RequestFromFile loads request data from file and creates a new http.Request
// request file should be in HTTP/1.x wire representation format
func RequestFromFile(filepath string) (*http.Request, error) {
	// nolint:gosec
	file, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	file = append(file, []byte("\r\n\r\n")...)

	rr, err := http.ReadRequest(bufio.NewReader(bytes.NewReader(file)))
	if err != nil {
		return nil, err
	}

	// Create outgoing request
	req, err := http.NewRequest(rr.Method, "http://"+rr.Host+rr.RequestURI, rr.Body)
	if err != nil {
		return nil, err
	}

	// Copy Headers
	for k := range rr.Header {
		req.Header.Set(k, rr.Header.Get(k))
	}

	return req, err
}

// ReadBody read body without draining it
func ReadBody(body *io.ReadCloser) ([]byte, error) {
	if *body == nil {
		return nil, nil
	}

	rBody, err := ioutil.ReadAll(*body)
	if err != nil {
		return nil, nil
	}

	*body = ioutil.NopCloser(bytes.NewReader(rBody))

	return rBody, nil
}
