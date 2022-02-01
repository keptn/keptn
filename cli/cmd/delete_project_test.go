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

func TestDeleteProjectCmd(t *testing.T) {
	credentialmanager.MockAuthCreds = true

	cmd := fmt.Sprintf("delete project %s --mock", "sockshop")
	_, err := executeActionCommandC(cmd)
	if err != nil {
		t.Errorf(unexpectedErrMsg, err)
	}
}

// TestDeleteProjectUnknownCommand
func TestDeleteProjectUnknownCommand(t *testing.T) {
	testInvalidInputHelper("delete project sockshop someUnknownCommand", "too many arguments set", t)
}

// TestDeleteProjectUnknownParameter
func TestDeleteProjectUnknownParmeter(t *testing.T) {
	testInvalidInputHelper("delete project sockshop --projectt=sockshop", "unknown flag: --projectt", t)
}
