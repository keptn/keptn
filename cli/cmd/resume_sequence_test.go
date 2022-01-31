package cmd

import (
	"fmt"
	"testing"
)

// TestResumeSequenceUnknownCommand
func TestResumeSequenceUnknownCommand(t *testing.T) {

	cmd := fmt.Sprintf("resume sequence someUnknownCommand --project=sockshop --keptn-context=djsfjdfdsjjcs")
	_, err := executeActionCommandC(cmd)
	if err == nil {
		t.Errorf("Expected an error")
	}

	got := err.Error()
	expected := "unknown command \"someUnknownCommand\" for \"keptn resume sequence\""
	if got != expected {
		t.Errorf("Expected %q, got %q", expected, got)
	}
}

// TestResumeSequenceUnknownParameter
func TestResumeSequenceUnknownParmeter(t *testing.T) {

	cmd := fmt.Sprintf("resume sequence --projectt=sockshop --keptn-context=djsfjdfdsjjcs")
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
