package cmd

import (
	"testing"
)

// TestUpdateSecretUnknownCommand
func TestUpdateSecretUnknownCommand(t *testing.T) {
	testInvalidInputHelper("update secret mysecret someUnknownCommand --from-literal='key2=value2' --scope=my-scope", "too many arguments set", t)
}

// TestUpdateSecretUnknownParameter
func TestUpdateSecretUnknownParmeter(t *testing.T) {
	testInvalidInputHelper("update secret mysecret --froom-literal='key2=value2' --scope=my-scope", "unknown flag: --froom-literal", t)
}
