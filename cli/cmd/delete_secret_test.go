package cmd

import (
	"testing"
)

// TestDeleteSecretUnknownCommand
func TestDeleteSecretUnknownCommand(t *testing.T) {
	testInvalidInputHelper("delete secret mysecret someUnknownCommand", "too many arguments set", t)
}

// TestDeleteSecretUnknownParameter
func TestDeleteSecretUnknownParmeter(t *testing.T) {
	testInvalidInputHelper("delete secret mysecret --projectt=mysecret", "unknown flag: --projectt", t)
}
