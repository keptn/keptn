package cmd

import (
	"fmt"
	"testing"
)

// TestDeleteSecretUnknownCommand
func TestDeleteSecretUnknownCommand(t *testing.T) {

	cmd := fmt.Sprintf("delete secret mysecret someUnknownCommand")
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

// TestDeleteSecretUnknownParameter
func TestDeleteSecretUnknownParmeter(t *testing.T) {

	cmd := fmt.Sprintf("delete secret mysecret --projectt=mysecret")
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
