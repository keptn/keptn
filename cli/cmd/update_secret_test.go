package cmd

import (
	"fmt"
	"testing"
)

// TestUpdateSecretUnknownCommand
func TestUpdateSecretUnknownCommand(t *testing.T) {

	cmd := fmt.Sprintf("update secret mysecret someUnknownCommand --from-literal='key2=value2' --scope=my-scope")
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

// TestUpdateSecretUnknownParameter
func TestUpdateSecretUnknownParmeter(t *testing.T) {

	cmd := fmt.Sprintf("update secret mysecret --froom-literal='key2=value2' --scope=my-scope")
	_, err := executeActionCommandC(cmd)
	if err == nil {
		t.Errorf("Expected an error")
	}

	got := err.Error()
	expected := "unknown flag: --froom-literal"
	if got != expected {
		t.Errorf("Expected %q, got %q", expected, got)
	}
}
