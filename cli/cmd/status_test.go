package cmd

import (
	"fmt"
	"testing"
)

// TestStatusUnknownCommand
func TestStatusUnknownCommand(t *testing.T) {

	cmd := fmt.Sprintf("status someUnknownCommand")
	_, err := executeActionCommandC(cmd)
	if err == nil {
		t.Errorf("Expected an error")
	}

	got := err.Error()
	expected := "unknown command \"someUnknownCommand\" for \"keptn status\""
	if got != expected {
		t.Errorf("Expected %q, got %q", expected, got)
	}
}

// TestStatusUnknownParameter
func TestStatusUnknownParmeter(t *testing.T) {

	cmd := fmt.Sprintf("status --projectt=sockshop")
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
