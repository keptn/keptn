package cmd

import (
	"bytes"
	"os"
	"testing"

	"github.com/keptn/keptn/cli/pkg/logging"
	"github.com/keptn/keptn/cli/utils/credentialmanager"
)

func init() {
	logging.InitLoggers(os.Stdout, os.Stdout, os.Stderr)
}

// TestOnboardService tests the onboard service command.
func TestOnboardServiceWrongPath(t *testing.T) {

	credentialmanager.MockAuthCreds = true
	buf := new(bytes.Buffer)
	rootCmd.SetOutput(buf)

	args := []string{
		"onboard",
		"service",
		"carts",
		"--project=sockshop",
		"--chart=cartsX",
	}
	rootCmd.SetArgs(args)
	err := rootCmd.Execute()

	if err == nil {
		t.Errorf("Expected error event, but no one received.")
	}

	expected := "Provided Helm chart does not exist"
	if err.Error() != expected {
		t.Errorf("Error actual = %v, and Expected = %v.", err, expected)
	}
}
