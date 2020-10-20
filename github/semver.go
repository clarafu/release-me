package github

import (
	"regexp"
)

var semverRegex = regexp.MustCompile(`(\d+)\.(\d+)\.(\d+)`)

func isPatchRelease(name string) bool {
	segments := semverRegex.FindStringSubmatch(name)
	if len(segments) < 4 {
		return false
	}
	patchNum := segments[3]
	return patchNum != "0"
}

