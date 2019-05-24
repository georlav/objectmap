package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/georlav/objectmap/cli"
	"github.com/georlav/objectmap/httpclient"
	"github.com/georlav/objectmap/insertionpoint"
	"github.com/georlav/objectmap/payload"
	"github.com/georlav/objectmap/pentest"
	"github.com/georlav/objectmap/useragent"

	"github.com/gookit/color"

	"github.com/olekukonko/tablewriter"

	"github.com/sirupsen/logrus"
)

const (
	// AppName the application name
	AppName = "objectMap"
)

func main() {
	cout := cli.NewOutput()

	ci, err := cli.NewInput()
	if err != nil {
		cout.Fatal(err)
	}

	// Configure logger
	cout.SetLevel(logrus.Level(ci.Verbosity))
	cout.Debugf("Starting %s", AppName)

	// Setup HTTP client
	c := httpclient.New(time.Second*time.Duration(ci.Timeout), func() useragent.UserAgent {
		ua := useragent.Default()
		if ci.RandomUserAgent {
			ua = useragent.Random()
		}

		return ua
	}(), ci.FollowRedirect)

	// If request is from file use it, else use target url to create one
	if ci.Request == nil {
		ci.Request, err = c.NewRequest(ci.Method, ci.URL, ci.Body)
		if err != nil {
			cout.Fatal(err)
		}
	}

	if ci.Request.Header.Get("Content-Type") == "" && ci.Body != nil {
		cout.Warning("Request has body but Content-Type is empty, this might lead to fewer tests")
		if strings.Contains(http.MethodPost+http.MethodPut+http.MethodPatch, ci.Request.Method) {
			cout.Info("Setting Content-Type to application/x-www-form-urlencoded")
			ci.Request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		}
	}

	// Create a copy of request body
	reqBody, err := httpclient.ReadBody(&ci.Request.Body)
	if err != nil {
		cout.Fatal(err)
	}

	// Initial Request
	resp, err := c.Client.Do(ci.Request)
	if err != nil {
		cout.Fatal(err)
	}
	defer func() {
		if err = resp.Body.Close(); err != nil {
			cout.Fatal(err)
		}
	}()

	// set request body
	resp.Request.Body = ioutil.NopCloser(bytes.NewReader(reqBody))

	// On request error exit, a valid response is required to continue
	if resp.StatusCode >= http.StatusBadRequest || resp.StatusCode < http.StatusOK {
		if resp, rErr := httputil.DumpResponse(resp, true); err == rErr {
			cout.Debug(string(resp))
		}
		cout.Fatalf("Please provide a valid http resource, got %s from %s", resp.Status, resp.Request.URL.String())
	}

	// Calculate insertion points
	cout.Info("Calculating insertion points")
	insertions, err := insertionpoint.NewInsertionPoints(resp, c)
	if err != nil {
		cout.Fatal(err)
	}
	cout.Infof("Found %d insertion points", len(insertions))
	cout.Debug(insertions)

	// load required payloads
	payloads := payload.Payloads{
		payload.NewObjectInjection(),
		payload.NewJavaDeserialization(),
	}

	// Generate tests
	ptCH, _ := pentest.NewPenTests(resp, insertions, payloads, c, ci.ConcurrentRequests)
	if err != nil {
		cout.Fatal(err)
	}

	// Create a wait group of ConcurrentRequests size
	var wg sync.WaitGroup
	wg.Add(ci.ConcurrentRequests)

	report := tablewriter.NewWriter(os.Stdout)
	report.SetHeader([]string{"Insertion Point", "Vulnerability", "Status"})
	results := make(map[string][]string)
	reqCounter := 0

	for i := 1; i <= ci.ConcurrentRequests; i++ {
		go func() {
			for {
				pt, more := <-ptCH
				if !more {
					wg.Done()
					break
				}
			Retry:
				reqCounter++

				dreq, err := httputil.DumpRequestOut(pt.Request, true)
				if err != nil {
					cout.Error(err)
					continue
				}
				cout.Debug(pt.Info)
				cout.Debug(string(dreq))

				resp, err := c.Client.Do(pt.Request)
				if err != nil {
					cout.Errorf("Failed to do request, %s (Attempts: %d)", err, pt.Retries)
					if pt.Retries <= ci.RequestRetries+1 {
						pt.Retries++
						goto Retry
					}
					continue
				}
				// nolint[:golint, errcheck]
				defer resp.Body.Close()

				dresp, err := httputil.DumpResponse(resp, true)
				if err != nil {
					cout.Error(err)
					continue
				}
				cout.Debug(string(dresp))

				status := color.Green.Sprint(`Clean`)
				if bytes.Contains(dresp, pt.ExpectInResp) {
					status = color.Error.Sprint(`Vulnerable`)
				}

				key := base64.StdEncoding.EncodeToString([]byte(
					fmt.Sprint(pt.Info.InsertionPoint.Point, pt.Info.InsertionPoint.Name, pt.Info.VulnName),
				))
				if _, ok := results[key]; !ok || (ok && strings.Contains(status, "Vulnerable")) {
					results[key] = []string{
						fmt.Sprintf(`%s[%s]`, pt.Info.InsertionPoint.Point, pt.Info.InsertionPoint.Name),
						pt.Info.VulnName,
						status,
					}
				}

			}
		}()
	}

	wg.Wait()
	for i := range results {
		report.Append(results[i])
	}
	report.SetFooter([]string{"", "Total Requests", fmt.Sprint(reqCounter)})
	report.Render()
}
