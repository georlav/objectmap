package httpclient_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/georlav/objectmap/httpclient"
)

func TestReadBody(t *testing.T) {
	// nolint[:errcheck, gosec]
	req, _ := http.NewRequest(http.MethodGet, "", strings.NewReader("param1=hello"))
	// nolint[:errcheck, gosec]
	b1, _ := httpclient.ReadBody(&req.Body)
	// nolint[:errcheck, gosec]
	b2, _ := ioutil.ReadAll(req.Body)

	// nolint:gosimple
	if bytes.Compare(b1, b2) != 0 {
		t.Fatal("Invalid body, expected the request bodies to be equal")
	}
}
