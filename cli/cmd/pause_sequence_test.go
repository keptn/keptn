package cmd

import (
	"testing"
)

// TestPauseSequenceUnknownCommand
func TestPauseSequenceUnknownCommand(t *testing.T) {
	testInvalidInputHelper("pause sequence someUnknownCommand --project=sockshop --keptn-context=djsfjdfdsjjcs", "unknown command \"someUnknownCommand\" for \"keptn pause sequence\"", t)
}

// TestPauseSequenceUnknownParameter
func TestPauseSequenceUnknownParmeter(t *testing.T) {
	testInvalidInputHelper("pause sequence --projectt=sockshop --keptn-context=djsfjdfdsjjcs", "unknown flag: --projectt", t)
}
