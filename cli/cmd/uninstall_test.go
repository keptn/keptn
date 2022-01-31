package cmd

import (
	"fmt"
	"testing"
)

// TestUninstallUnknownCommand
func TestUninstallUnknownCommand(t *testing.T) {

	cmd := fmt.Sprintf("uninstall someUnknownCommand")
	_, err := executeActionCommandC(cmd)
	if err == nil {
		t.Errorf("Expected an error")
	}

	got := err.Error()
	expected := "unknown command \"someUnknownCommand\" for \"keptn uninstall\""
	if got != expected {
		t.Errorf("Expected %q, got %q", expected, got)
	}
}

// TestUninstallUnknownParameter
func TestUninstallUnknownParmeter(t *testing.T) {

	cmd := fmt.Sprintf("uninstall --project=sockshop")
	_, err := executeActionCommandC(cmd)
	if err == nil {
		t.Errorf("Expected an error")
	}

	got := err.Error()
	expected := "unknown flag: --project"
	if got != expected {
		t.Errorf("Expected %q, got %q", expected, got)
	}
}
