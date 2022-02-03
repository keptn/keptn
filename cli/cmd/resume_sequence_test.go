package cmd

import (
	"testing"
)

// TestResumeSequenceUnknownCommand
func TestResumeSequenceUnknownCommand(t *testing.T) {
	testInvalidInputHelper("resume sequence someUnknownCommand --project=sockshop --keptn-context=djsfjdfdsjjcs", "unknown command \"someUnknownCommand\" for \"keptn resume sequence\"", t)
}

// TestResumeSequenceUnknownParameter
func TestResumeSequenceUnknownParmeter(t *testing.T) {
	testInvalidInputHelper("resume sequence --projectt=sockshop --keptn-context=djsfjdfdsjjcs", "unknown flag: --projectt", t)
}
