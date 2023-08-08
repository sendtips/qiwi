package qiwi

import "fmt"

// AppVersion for version string of main application version
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
