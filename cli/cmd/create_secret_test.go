package cmd

import (
	"testing"
)

// TestCreateSecretUnknownCommand
func TestCreateSecretUnknownCommand(t *testing.T) {
	testInvalidInputHelper("create secret mysecret someUnknownCommand --from-literal='key1=value1' --scope=myscope", "too many arguments set", t)
}

// TestCreateSecretUnknownParameter
func TestCreateSecretUnknownParmeter(t *testing.T) {
	testInvalidInputHelper("create secret mysecret--from-literal='key1=value1' --projectt=myscope", "unknown flag: --projectt", t)
}
