package cmd

import (
	"os"
	"testing"

	"github.com/keptn/keptn/cli/pkg/logging"
	"github.com/keptn/keptn/cli/utils/credentialmanager"
)

func init() {
	logging.InitLoggers(os.Stdout, os.Stdout, os.Stderr)
}

func TestDeleteProjectCmd(t *testing.T) {

	credentialmanager.MockAuthCreds = true

	args := []string{
		"delete",
		"project",
		"sockshop",
		"--mock",
	}
	rootCmd.SetArgs(args)
	err := rootCmd.Execute()

	if err != nil {
		t.Errorf("An error occured %v", err)
	}
}
