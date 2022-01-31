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

// TestCreateProjectCmd tests the default use of the create project command
func TestCreateServiceCmd(t *testing.T) {
	credentialmanager.MockAuthCreds = true

	cmd := fmt.Sprintf("create service carts --project=%s --mock", "sockshop")
	_, err := executeActionCommandC(cmd)
	if err != nil {
		t.Errorf(unexpectedErrMsg, err)
	}
}

// TestCreateServiceUnknownCommand
func TestCreateServiceUnknownCommand(t *testing.T) {

	cmd := fmt.Sprintf("create service myservice someUnknownCommand --project=sockshop")
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

// TestCreateServiceUnknownParameter
func TestCreateServiceUnknownParmeter(t *testing.T) {

	cmd := fmt.Sprintf("create service myservice --projectt=sockshop")
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
