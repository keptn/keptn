package cmd

import (
	"testing"
)

// TestStatusUnknownCommand
func TestStatusUnknownCommand(t *testing.T) {
	testInvalidInputHelper("status someUnknownCommand", "unknown command \"someUnknownCommand\" for \"keptn status\"", t)
}

// TestStatusUnknownParameter
func TestStatusUnknownParmeter(t *testing.T) {
	testInvalidInputHelper("status --projectt=sockshop", "unknown flag: --projectt", t)
}
