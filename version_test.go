package qiwi

import (
	"fmt"
	"testing"
)

func TestVersion(t *testing.T) {
	s := fmt.Sprintf("%s/%s", userAgent, Version)

	if s != version() {
		t.Error("Version string is invalid")
	}

	AppVersion = "somever/1.0"
	s = fmt.Sprintf("%s (%s/%s)", AppVersion, userAgent, Version)

	if s != version() {
		t.Error("Version string is invalid")
	}
}
