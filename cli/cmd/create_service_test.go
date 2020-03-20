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

	project := "sockshop"

	args := []string{
		"create",
		"service",
		"carts",
		fmt.Sprintf("--project=%s", project),
		"--mock",
	}
	rootCmd.SetArgs(args)
	err := rootCmd.Execute()

	if err != nil {
		t.Errorf("An error occured: %v", err)
	}
}
