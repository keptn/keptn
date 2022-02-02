package cmd

import (
	"testing"
)

// TestGenerateCESpecUnknownCommand
func TestGenerateCESpecUnknownCommand(t *testing.T) {
	testInvalidInputHelper("generate cloud-events-spec someUnknownCommand", "unknown command \"someUnknownCommand\" for \"keptn generate cloud-events-spec\"", t)
}

// TestGenerateCESpecUnknownParameter
func TestGenerateCESpecUnknownParmeter(t *testing.T) {
	testInvalidInputHelper("generate cloud-events-spec --project=sockshop", "unknown flag: --project", t)
}
