package utils

import (
	"log"
	"regexp"
	"strings"
)

// ValidateK8sName valides if a given name starts with lowercase letter and only contains lowercase letters and numbers.
func ValidateK8sName(name string) bool {
	reg, err := regexp.Compile("[a-z][a-z0-9]+")
	if err != nil {
		log.Fatal(err)
	}
	processedString := reg.FindString(name)
	return len(processedString) == len(name)
}

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
func IsOfficialKeptnVersion(version string) bool {
	// First check if version ends with time stamp
	reg, _ := regexp.Compile(`(\-[0-9]{8}\.[0-9]{4})$`)
	if reg.MatchString(version) {
		return false
	}

	reg, _ = regexp.Compile(`([0-9]+\.){2}(([0-9]+\.[[:alpha:]])|([0-9]+))`)
	return reg.MatchString(version)
}
