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

// TestEvaluationDoneGetEvent tests the evaluation-done command
func TestEvaluationDoneGetEvent(t *testing.T) {

	credentialmanager.MockAuthCreds = true

	cmd := fmt.Sprintf("get event evaluation-done --keptn-context=%s --mock", "8929e5e5-3826-488f-9257-708bfa974909")
	_, err := executeActionCommandC(cmd)
	if err != nil {
		t.Errorf("unexpected error, got '%v'", err)
	}
}
