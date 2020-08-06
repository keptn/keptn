package cmd

import (
	"fmt"
	"github.com/keptn/keptn/cli/pkg/credentialmanager"
	"github.com/keptn/keptn/cli/pkg/logging"
	"os"
	"testing"
)

func init() {
	logging.InitLoggers(os.Stdout, os.Stdout, os.Stderr)
}

func TestGetEventCmdEptyInput(t *testing.T) {
	credentialmanager.MockAuthCreds = true

	cmd := fmt.Sprintf("get event --project=%s --mock",
		"sockshop")
	_, err := executeActionCommandC(cmd)

	if !errorContains(err, "please provide an event type as an argument") {
		t.Errorf("missing expected error, but got %v", err)
	}
}

func TestGetEventNoProject(t *testing.T) {
	credentialmanager.MockAuthCreds = true

	cmd := fmt.Sprintf("get event sh.keptn.events.problem --mock")
	_, err := executeActionCommandC(cmd)
	if !errorContains(err, "required flag(s) \"project\" not set") {
		t.Errorf("missing expected error, but got %v", err)
	}
}

// TestGetEvent tests the get event command
func TestGetEvent(t *testing.T) {
	credentialmanager.MockAuthCreds = true

	cmd := fmt.Sprintf("get event sh.keptn.events.problem --project=%s --mock", "sockshop")
	_, err := executeActionCommandC(cmd)
	if err != nil {
		t.Errorf(unexpectedErrMsg, err)
	}
}

