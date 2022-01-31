package cmd

import (
	"fmt"
	"testing"
)

// TestDeleteServiceUnknownCommand
func TestDeleteServiceUnknownCommand(t *testing.T) {

	cmd := fmt.Sprintf("delete secret myservice someUnknownCommand")
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

// TestDeleteServiceUnknownParameter
func TestDeleteServiceUnknownParmeter(t *testing.T) {

	cmd := fmt.Sprintf("delete secret myservice --projectt=myservice")
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
