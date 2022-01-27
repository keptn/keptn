package cmd

import (
	"fmt"
	"github.com/keptn/keptn/cli/pkg/credentialmanager"
	"os"
	"testing"

	"github.com/keptn/keptn/cli/pkg/logging"
)

func init() {
	logging.InitLoggers(os.Stdout, os.Stdout, os.Stderr)
}

// TestGetStage tests the get stage command
func TestGetStage(t *testing.T) {

	credentialmanager.MockAuthCreds = true

	cmd := fmt.Sprintf("get stage hardening --project=sockshop --mock")
	_, err := executeActionCommandC(cmd)
	if err != nil {
		t.Errorf(unexpectedErrMsg, err)
	}
}
