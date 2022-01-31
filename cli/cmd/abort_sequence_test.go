package cmd

import (
	"fmt"
	"testing"
)

// TestAbortSequenceUnknownCommand
func TestAbortSequenceUnknownCommand(t *testing.T) {

	cmd := fmt.Sprintf("abort sequence someUnknownCommand --project=sockshop --keptn-context=djsfjdfdsjjcs")
	_, err := executeActionCommandC(cmd)
	if err == nil {
		t.Errorf("Expected an error")
	}

	got := err.Error()
	expected := "unknown command \"someUnknownCommand\" for \"keptn abort sequence\""
	if got != expected {
		t.Errorf("Expected %q, got %q", expected, got)
	}
}

// TestAbortSequenceUnknownParameter
func TestAbortSequenceUnknownParmeter(t *testing.T) {

	cmd := fmt.Sprintf("abort sequence --projectt=sockshop --keptn-context=djsfjdfdsjjcs")
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
