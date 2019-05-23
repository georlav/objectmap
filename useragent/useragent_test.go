package useragent_test

import (
	"objectmap/useragent"
	"strings"
	"testing"
)

func TestDefault(t *testing.T) {
	ua := useragent.Default()

	if !strings.Contains(string(ua), "ObjectMap") {
		t.Fatal("Invalid Agent")
	}
}

func TestRandom(t *testing.T) {
	ua := useragent.Random()

	found := false
	for _, a := range useragent.All() {
		if a == ua {
			found = true
		}
	}

	if !found {
		t.Fatal("Unknown Agent")
	}

}
