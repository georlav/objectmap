package cli_test

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	st := m.Run()
	os.Exit(st)
}

func TestNewCInput(t *testing.T) {
}
