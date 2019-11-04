package cmd

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/keptn/keptn/cli/pkg/logging"
	"github.com/keptn/keptn/cli/utils/credentialmanager"
)

func init() {
	logging.InitLoggers(os.Stdout, os.Stdout, os.Stderr)
}

// TestEvaluationStart tests the evaluation.start command.
func TestEvaluationStart(t *testing.T) {

	credentialmanager.MockAuthCreds = true
	buf := new(bytes.Buffer)
	rootCmd.SetOutput(buf)

	args := []string{
		"send",
		"event",
		"evaluation.start",
		fmt.Sprintf("--project=%s", "sockshop"),
		fmt.Sprintf("--stage=%s", "hardening"),
		fmt.Sprintf("--service=%s", "carts"),
		fmt.Sprintf("--timeframe=%s", "5m"),
		"--mock",
	}
	rootCmd.SetArgs(args)
	err := rootCmd.Execute()

	if err != nil {
		t.Errorf("An error occured: %v", err)
	}
}

func TestEvaluationStartWrongFormat(t *testing.T) {

	credentialmanager.MockAuthCreds = true
	buf := new(bytes.Buffer)
	rootCmd.SetOutput(buf)

	args := []string{
		"send",
		"event",
		"evaluation.start",
		fmt.Sprintf("--project=%s", "sockshop"),
		fmt.Sprintf("--stage=%s", "hardening"),
		fmt.Sprintf("--service=%s", "carts"),
		fmt.Sprintf("--timeframe=%s", "5h"),
		"--mock",
	}
	rootCmd.SetArgs(args)
	err := rootCmd.Execute()

	if err == nil {
		t.Error("An error occured: expect an error due to wrong time frame format")
	}
}
