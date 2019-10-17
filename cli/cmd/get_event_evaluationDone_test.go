package cmd

import (
	"os"
	"testing"

	"github.com/keptn/keptn/cli/pkg/logging"
)

func init() {
	logging.InitLoggers(os.Stdout, os.Stdout, os.Stderr)
}

// TestEvaluationDone tests the evaluation-done command.
func TestEvaluationDone(t *testing.T) {
	/*
		credentialmanager.MockAuthCreds = true
		buf := new(bytes.Buffer)
		rootCmd.SetOutput(buf)

		args := []string{
			"get",
			"event",
			"evaluation-done",
			fmt.Sprintf("--keptn-context=%s", "??"),
			"--mock",
		}
		rootCmd.SetArgs(args)
		err := rootCmd.Execute()

		if err != nil {
			t.Errorf("An error occured %v", err)
		}
	*/
}
