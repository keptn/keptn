package cmd

import (
	"fmt"
	"testing"
)

// TestGenerateCESpecUnknownCommand
func TestGenerateCESpecUnknownCommand(t *testing.T) {

	cmd := fmt.Sprintf("generate cloud-events-spec someUnknownCommand")
	_, err := executeActionCommandC(cmd)
	if err == nil {
		t.Errorf("Expected an error")
	}

	got := err.Error()
	expected := "unknown command \"someUnknownCommand\" for \"keptn generate cloud-events-spec\""
	if got != expected {
		t.Errorf("Expected %q, got %q", expected, got)
	}
}

// TestGenerateCESpecUnknownParameter
func TestGenerateCESpecUnknownParmeter(t *testing.T) {

	cmd := fmt.Sprintf("generate cloud-events-spec --project=sockshop")
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
