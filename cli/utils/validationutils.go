package utils

import (
	"strings"

	"github.com/hashicorp/go-version"
)

// ErrorContains checks if the error message in out contains the text in want.
func ErrorContains(out error, want string) bool {
	if out == nil {
		return want == ""
	}
	if want == "" {
		return false
	}
	return strings.Contains(out.Error(), want)
}

// IsOfficialKeptnVersion checks whether the provided version string follows a
// Keptn version pattern
func IsOfficialKeptnVersion(versionStr string) bool {
	// First check if version ends with time stamp
	_, err := version.NewSemver(versionStr)
	return err == nil
}
