// Package cli takes user input form cli
package cli

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/georlav/objectmap/httpclient"

	"github.com/pkg/errors"

	"github.com/urfave/cli"
)

const (
	// AppName the application name
	AppName = "ObjectMap"
)

// Input keeps console input variables and the cli app instance
type Input struct {
	App                *cli.App
	URL                *url.URL
	URLScheme          string
	Method             string
	Body               io.ReadCloser
	RequestFile        string
	Request            *http.Request
	RequestRetries     int
	ConcurrentRequests int
	FollowRedirect     bool
	Timeout            int64
	UserAgent          string
	RandomUserAgent    bool
	Packages           bool
	PackagesRop        bool
	Banner             bool
	Verbosity          int
}

// NewInput creates a new ConsoleIn object
func NewInput() (ci Input, err error) {
	ci.App = newCliApp()
	ci.App.Flags = newCliAppFlags(&ci)
	ci.App.HideHelp = true
	ci.App.HideVersion = true

	ci.App.Before = func(c *cli.Context) error {
		// if no flags, show help
		if c.NumFlags() == 0 {
			cli.ShowAppHelpAndExit(c, 0)
		}

		if c.Bool("help") {
			cli.ShowAppHelpAndExit(c, 0)
		}

		// convert method name to uppercase
		if c.String("method") != "" {
			if err = c.Set("method", strings.ToUpper(c.String("method"))); err != nil {
				return err
			}

		}

		return nil
	}

	ci.App.Action = func(c *cli.Context) error {
		// if flags url or request should be present
		if c.String("url") == "" && c.String("request") == "" {
			return cli.NewExitError("You need to specify a target (use --url or --request or --help)", 1)
		}

		// set the request body
		if c.String("body") != "" {
			ci.Body = ioutil.NopCloser(
				strings.NewReader(
					c.String("body"),
				),
			)
		}

		// set target via --url param
		if u := c.String("url"); u != "" {
			// If url has no scheme add http
			if !strings.HasPrefix(u, "http") {
				u = "http://" + u
			}

			ci.URL, err = url.Parse(u)
			if err != nil {
				return errors.Wrapf(
					err,
					"Invalid value '%s' for flag --url",
					u,
				)
			}
		}

		// Set target via --request param
		if ci.RequestFile != "" {
			ci.Request, err = httpclient.RequestFromFile(ci.RequestFile)
			if err != nil {
				return errors.Wrap(err, "Unable to load request from file")
			}

			ci.Request.URL.Scheme = ci.URLScheme

			// nolint:govet
			b, err := httpclient.ReadBody(&ci.Request.Body)
			if err != nil {
				return err
			}

			ci.Body = ioutil.NopCloser(bytes.NewReader(b))
		}

		return nil
	}

	err = ci.App.Run(os.Args)
	if err != nil {
		return ci, err
	}

	err = validateParsedFlags(ci, Validator{})
	if err != nil {
		return ci, err
	}

	return ci, nil
}

func newCliApp() *cli.App {
	app := cli.NewApp()
	app.Name = AppName
	app.Usage = ""
	app.Author = "georlav"
	app.Description = "Object Injection Vulnerability Scanner"
	app.UsageText = AppName + " --url https://example.com [options]"
	app.Version = "0.0.1"
	app.Commands = []cli.Command{}

	return app
}

func newCliAppFlags(ci *Input) []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "url, u",
			Usage: "Target url",
		},
		cli.StringFlag{
			Name:        "url-scheme, us",
			Usage:       "Set the URL scheme [http, https]",
			Destination: &ci.URLScheme,
			Value:       "http",
		},
		cli.StringFlag{
			Name:  "method, m",
			Value: http.MethodGet,
			Usage: fmt.Sprintf(
				"Set the HTTP request method, supported methods are %v",
				[]string{"GET", "POST", "PUT", "PATCH", "DELETE"},
			),
			Destination: &ci.Method,
		},
		cli.StringFlag{
			Name:  "body",
			Usage: "Set the request body",
		},
		cli.StringFlag{
			Name:        "request, r",
			Usage:       "Load target from request file",
			Destination: &ci.RequestFile,
		},
		cli.IntFlag{
			Name:        "request-concurrency, rc",
			Usage:       "Set the number of concurrent requests",
			Destination: &ci.ConcurrentRequests,
			Value:       1,
		},
		cli.IntFlag{
			Name:        "request-retries, rr",
			Usage:       "Set number of retries on request failure",
			Destination: &ci.RequestRetries,
			Value:       2,
		},
		cli.BoolFlag{
			Name:        "no-follow, nf",
			Usage:       "Do not follow http redirects (default: follows)",
			Destination: &ci.FollowRedirect,
		},
		cli.Int64Flag{
			Name:        "timeout, t",
			Usage:       "Set the max timeout limit in seconds for http requests",
			Value:       10,
			Destination: &ci.Timeout,
		},
		cli.StringFlag{
			Name:        "user-agent",
			Value:       "ObjectMap/1.0",
			Usage:       "Set client user agent",
			Destination: &ci.UserAgent,
		},
		cli.BoolFlag{
			Name:        "random-agent",
			Usage:       "Set client to use a random user agent",
			Destination: &ci.RandomUserAgent,
		},
		// todo:impl
		//cli.BoolFlag{
		//	Name:        "packages, p",
		//	Usage:       "If target is vulnerable try to enumerate application packages",
		//	Destination: &ci.Packages,
		//},
		//cli.BoolFlag{
		//	Name:        "packages-rop, vp",
		//	Usage:       "If target is vulnerable try to enumerate packages that can be used to start a POP chain",
		//	Destination: &ci.PackagesRop,
		//},
		cli.BoolTFlag{
			Name:        "banner, b",
			Usage:       "Retrieve server banner",
			Destination: &ci.Banner,
		},
		cli.IntFlag{
			Name:        "verbose, v",
			Usage:       "Set the verbosity level [1-5]",
			Value:       4,
			Destination: &ci.Verbosity,
		},
		cli.BoolFlag{
			Name:  "help, h",
			Usage: "Show help",
		},
	}
}

func validateParsedFlags(ci Input, v Validator) (err error) {
	// Validate http method
	if err = v.HTTPMethod(ci.Method); err != nil {
		return err
	}

	// validate verbosity
	if err = v.VerboseLevel(ci.Verbosity); err != nil {
		return err
	}

	return nil
}
