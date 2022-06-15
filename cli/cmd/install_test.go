package cmd

import (
	"testing"
)

// TestInstallUnknownCommand
func TestInstallUnknownCommand(t *testing.T) {
	testInvalidInputHelper("install someUnknownCommand", "unknown command \"someUnknownCommand\" for \"keptn install\"", t)
}

// TestInstallUnknownParameter
func TestInstallUnknownParmeter(t *testing.T) {
	testInvalidInputHelper("install --project=sockshop", "unknown flag: --project", t)
}

// TestInstallUnknownParameter
func TestInstallDeprecated(t *testing.T) {
	testInvalidInputHelper("install --hide-sensitive-data ", MsgDeprecatedUseHelm, t)
}
