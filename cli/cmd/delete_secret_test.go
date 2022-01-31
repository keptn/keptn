package cmd

import (
	"fmt"
	"os"
	"testing"

	"github.com/keptn/keptn/cli/pkg/credentialmanager"

	"github.com/keptn/keptn/cli/pkg/logging"
)

func init() {
	logging.InitLoggers(os.Stdout, os.Stdout, os.Stderr)
}

func TestDeleteSecretCmd(t *testing.T) {
	credentialmanager.MockAuthCreds = true

	cmd := fmt.Sprintf("delete secret mysecret --mock --scope=my-scope")
	_, err := executeActionCommandC(cmd)
	if err != nil {
		t.Errorf(unexpectedErrMsg, err)
	}
}

// TestDeleteSecretUnknownCommand
func TestDeleteSecretUnknownCommand(t *testing.T) {

	cmd := fmt.Sprintf("delete secret mysecret someUnknownCommand")
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

// TestDeleteSecretUnknownParameter
func TestDeleteSecretUnknownParmeter(t *testing.T) {

	cmd := fmt.Sprintf("delete secret mysecret --projectt=mysecret")
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
