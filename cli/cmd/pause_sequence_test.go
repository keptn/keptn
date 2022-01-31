package cmd

import (
	"fmt"
	"testing"
)

// TestPauseSequenceUnknownCommand
func TestPauseSequenceUnknownCommand(t *testing.T) {

	cmd := fmt.Sprintf("pause sequence someUnknownCommand --project=sockshop --keptn-context=djsfjdfdsjjcs")
	_, err := executeActionCommandC(cmd)
	if err == nil {
		t.Errorf("Expected an error")
	}

	got := err.Error()
	expected := "unknown command \"someUnknownCommand\" for \"keptn pause sequence\""
	if got != expected {
		t.Errorf("Expected %q, got %q", expected, got)
	}
}

// TestPauseSequenceUnknownParameter
func TestPauseSequenceUnknownParmeter(t *testing.T) {

	cmd := fmt.Sprintf("pause sequence --projectt=sockshop --keptn-context=djsfjdfdsjjcs")
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
