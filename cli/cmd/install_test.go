package cmd

import (
	"testing"
)

// TestInstallUnknownCommand
func TestInstallUnknownCommand(t *testing.T) {
	testInvalidInputHelper("install someUnknownCommand", "unknown command \"someUnknownCommand\" for \"keptn install\"", t)
}

// TestInstallUnknownParameter
func TestInstallUnknownParameter(t *testing.T) {
	testInvalidInputHelper("install --project=sockshop", "unknown flag: --project", t)
}

// TestInstallDeprecated
func TestInstallDeprecated(t *testing.T) {
	Version = "0.16.0"
	testInvalidInputHelper("install --hide-sensitive-data ", "this command is deprecated", t)
}
