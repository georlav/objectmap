package httpclient_test

import (
	"io/ioutil"
	"testing"
	"time"

	"github.com/georlav/objectmap/httpclient"
	"github.com/georlav/objectmap/useragent"
)

func TestHTTPClient_NewRequestFromFile(t *testing.T) {
	testCases := []struct {
		File       string
		Method     string
		BodyLength int
	}{
		{File: "testdata/get.req", Method: "GET", BodyLength: 0},
		{File: "testdata/post.req", Method: "POST", BodyLength: 42},
		{File: "testdata/post-multipart-form.req", Method: "POST", BodyLength: 49},
	}

	hc := httpclient.New(time.Second*10, useragent.Default(), true)

	for _, tc := range testCases {
		rtc := tc
		t.Run("Testing request "+rtc.File, func(t *testing.T) {
			r, err := hc.NewRequestFromFile(rtc.File)
			if err != nil {
				t.Fatalf("Failed to create request from file %s. Reason: %s", rtc.File, err)
			}

			if r.Method != rtc.Method {
				t.Fatalf("Wrong method, expected %s got %s", rtc.Method, r.Method)
			}

			if r.Body != nil {
				b, err := ioutil.ReadAll(r.Body)
				if err != nil {
					t.Fatal(err)
				}

				if len(b) != rtc.BodyLength {
					t.Fatalf("Invalid body length got %d expected %d", len(b), rtc.BodyLength)
				}
			}

		})
	}
}
