package qiwi

import "fmt"

// AppVersion sets version for main application.
var AppVersion string

const (
	// Version string.
	Version   string = "1.0.0"
	userAgent string = "Sendtips-QIWI-Go" // UserAgent name
)

func version() (s string) {
	s = fmt.Sprintf("%s/%s", userAgent, Version)

	if AppVersion != "" {
		s = fmt.Sprintf("%s (%s)", AppVersion, s)
	}

	return s
}
