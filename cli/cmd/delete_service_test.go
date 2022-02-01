package cmd

import (
	"testing"
)

// TestDeleteServiceUnknownCommand
func TestDeleteServiceUnknownCommand(t *testing.T) {
	testInvalidInputHelper("delete secret myservice someUnknownCommand", "too many arguments set", t)
}

// TestDeleteServiceUnknownParameter
func TestDeleteServiceUnknownParmeter(t *testing.T) {
	testInvalidInputHelper("delete secret myservice --projectt=myservice", "unknown flag: --projectt", t)
}
