package utils

import (
	"log"
	"regexp"
	"strings"
)

// ValidateK8sName valides if a given name starts with lowercase letter and only contains lowercase letters and -
func ValidateK8sName(svcName string) bool {
	reg, err := regexp.Compile("[a-z][a-zA-Z0-9/-]+")
	if err != nil {
		log.Fatal(err)
	}
	processedString := reg.FindString(svcName)
	return len(processedString) == len(svcName)
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
