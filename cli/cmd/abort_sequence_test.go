package cmd

import (
	"testing"
)

// TestAbortSequenceUnknownCommand
func TestAbortSequenceUnknownCommand(t *testing.T) {
	testInvalidInputHelper("abort sequence someUnknownCommand --project=sockshop --keptn-context=djsfjdfdsjjcs", "unknown command \"someUnknownCommand\" for \"keptn abort sequence\"", t)
}

// TestAbortSequenceUnknownParameter
func TestAbortSequenceUnknownParmeter(t *testing.T) {
	testInvalidInputHelper("abort sequence --projectt=sockshop --keptn-context=djsfjdfdsjjcs", "unknown flag: --projectt", t)
}
