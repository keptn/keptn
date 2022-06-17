package cmd

import (
	"testing"
)

// TestUninstallUnknownCommand
func TestUninstallUnknownCommand(t *testing.T) {
	testInvalidInputHelper("uninstall someUnknownCommand", "unknown command \"someUnknownCommand\" for \"keptn uninstall\"", t)
}

// TestUninstallUnknownParameter
func TestUninstallUnknownParmeter(t *testing.T) {
	testInvalidInputHelper("uninstall --project=sockshop", "unknown flag: --project", t)
}

// TestUninstallDeprecated
func TestUninstallDeprecated(t *testing.T) {
	Version = "0.16.0"
	testInvalidInputHelper("uninstall", "this command is deprecated", t)
}
