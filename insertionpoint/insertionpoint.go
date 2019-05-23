// Package insertionpoint analyzes requests and generates insertion points
package insertionpoint

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"objectmap/httpclient"

	"github.com/pkg/errors"
)

// Point type of insertion
type Point string

// Available points
// nolint:golint
const (
	Cookie    Point = "Cookie"
	Header    Point = "Header"
	Param     Point = "Param"
	BodyParam Point = "Body Param"
)

// InsertionPoint is a point in an http request where a payload can be injected
type InsertionPoint struct {
	Name  string
	Point Point
}

// InsertionPoints slice of InsertionPoint objects
type InsertionPoints []InsertionPoint

// NewInsertionPoint creates a new InsertionPoint
func NewInsertionPoint(
	n string,
	p Point,
) InsertionPoint {
	return InsertionPoint{
		Name:  n,
		Point: p,
	}
}

// NewInsertionPoints create a slice of insertion point objects
func NewInsertionPoints(
	resp *http.Response,
	c httpclient.HTTPClient,
) (iPts InsertionPoints, err error) {
	req := resp.Request
	reqBody, err := c.ReadBody(&req.Body)
	if err != nil {
		return iPts, err
	}

	if err = req.ParseForm(); err != nil {
		return iPts, errors.Wrap(err, "failed to parse request body")
	}
	req.Body = ioutil.NopCloser(bytes.NewReader(reqBody))

	for p := range req.Form {
		iPts = append(iPts, NewInsertionPoint(
			p,
			Param,
		))
	}

	for i := range resp.Request.Header {
		toSkip := map[string]string{"Content-Type": "", "Content-Length": "", "Cookie": "", "Host": ""}
		if _, ok := toSkip[http.CanonicalHeaderKey(i)]; !ok {
			iPts = append(iPts, NewInsertionPoint(
				i,
				Header,
			))
		}
	}

	cookies := append(resp.Request.Cookies(), resp.Cookies()...)
	for _, c := range cookies {
		iPts = append(iPts, NewInsertionPoint(
			c.Name,
			Cookie,
		))
	}

	return iPts, nil
}
