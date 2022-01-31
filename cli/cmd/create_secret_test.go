package cmd

import (
	"fmt"
	"testing"
)

// TestCreateSecretUnknownCommand
func TestCreateSecretUnknownCommand(t *testing.T) {

	cmd := fmt.Sprintf("create secret mysecret someUnknownCommand --from-literal='key1=value1' --scope=myscope")
	_, err := executeActionCommandC(cmd)
	if err == nil {
		t.Errorf("Expected an error")
	}

	got := err.Error()
	expected := "too many arguments set"
	if got != expected {
		t.Errorf("Expected %q, got %q", expected, got)
	}
}

// TestCreateSecretUnknownParameter
func TestCreateSecretUnknownParmeter(t *testing.T) {

	cmd := fmt.Sprintf("create secret mysecret--from-literal='key1=value1' --projectt=myscope")
	_, err := executeActionCommandC(cmd)
	if err == nil {
		t.Errorf("Expected an error")
	}

	got := err.Error()
	expected := "unknown flag: --projectt"
	if got != expected {
		t.Errorf("Expected %q, got %q", expected, got)
	}
}
