package cli

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

// Validator handle
type Validator struct{}

// URL validator
func (v Validator) URL(u string) error {
	pu, err := url.ParseRequestURI(u)
	if err != nil {
		return errors.Wrapf(
			err,
			"Invalid value '%s' for flag --url",
			u,
		)
	}

	// Empty Host
	if pu.Host == "" {
		return errors.Errorf("Invalid value '%s' for flag --url, no host", u)
	}

	// Invalid Scheme
	if !strings.HasPrefix(pu.Scheme, "http") && !strings.HasPrefix(pu.Scheme, "https") {
		return errors.Errorf("Invalid value '%s' for flag --url, url should start with http/https not %s", pu.Host, pu.Scheme)
	}

	return nil
}

// HTTPMethod validator
func (v Validator) HTTPMethod(method string) error {
	httpMethods := map[string]struct{}{
		http.MethodGet:    {},
		http.MethodPost:   {},
		http.MethodPut:    {},
		http.MethodPatch:  {},
		http.MethodDelete: {},
	}

	_, ok := httpMethods[method]
	if !ok {
		return errors.Errorf(
			"Invalid value '%s' for flag -method, valid values are %s",
			method,
			mapKeys(httpMethods),
		)
	}

	return nil
}

// VerboseLevel validator
func (v Validator) VerboseLevel(level int) error {
	verboseLevels := map[string]struct{}{
		"1": {},
		"2": {},
		"3": {},
		"4": {},
		"5": {},
	}

	_, ok := verboseLevels[strconv.Itoa(level)]
	if !ok {
		return errors.Errorf(
			"Invalid value '%d' for flag -verbose, valid values are %s",
			level,
			mapKeys(verboseLevels),
		)
	}

	return nil
}

func mapKeys(m map[string]struct{}) []string {
	// nolint:prealloc
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}

	return keys
}
