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

// TestEvaluationFinishedGetEvent tests the get event evaluation.finished command
func TestEvaluationFinishedGetEvent(t *testing.T) {

	credentialmanager.MockAuthCreds = true

	cmd := fmt.Sprintf("get event evaluation.finished --keptn-context=%s --mock", "8929e5e5-3826-488f-9257-708bfa974909")
	_, err := executeActionCommandC(cmd)
	if err != nil {
		t.Errorf(unexpectedErrMsg, err)
	}
}

// TestGetEventEvaluationFinishedUnknownCommand
func TestGetEventEvaluationFinishedUnknownCommand(t *testing.T) {
	testInvalidInputHelper("get event evaluation.finished someUnknownCommand", "unknown command \"someUnknownCommand\" for \"keptn get event evaluation.finished\"", t)
}

// TestGetEventEvaluationFinishedUnknownParameter
func TestGetEventEvaluationFinishedUnknownParmeter(t *testing.T) {
	testInvalidInputHelper("get event evaluation.finished --projectt=sockshop", "unknown flag: --projectt", t)
}
