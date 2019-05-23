package cli

import (
	"fmt"
	"testing"
)

func TestValidator_Url(t *testing.T) {
	testCases := []struct {
		input    string
		hasError bool
	}{
		{"target.com", true},
		{"www.target.com", true},
		{"target", true},
		{"https://", true},
		{"http://www.target.com", false},
		{"http://target.com", false},
		{"htt://target.com", true},
	}

	vld := Validator{}

	for _, testCase := range testCases {
		tc := testCase

		t.Run("Testing URL: "+tc.input, func(t *testing.T) {
			hasErr := false
			output := vld.URL(tc.input)
			if output != nil {
				hasErr = true
			}

			if hasErr != tc.hasError {
				t.Fatalf("Unexpected validation result, expected %t got %t", tc.hasError, hasErr)
			}
		})
	}
}

func TestValidator_VerboseLevel(t *testing.T) {
	testCases := []struct {
		input    int
		hasError bool
	}{
		{10, true},
		{1, false},
		{-1, true},
	}

	vld := Validator{}

	for _, testCase := range testCases {
		tc := testCase

		t.Run(fmt.Sprintf("Verbose level %d", tc.input), func(t *testing.T) {
			result := false

			output := vld.VerboseLevel(tc.input)
			if output != nil {
				result = true
			}

			if result != tc.hasError {
				t.Fatalf("Unexpected validation result, expected %t got %t", tc.hasError, result)
			}
		})
	}
}
