package cmd

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/keptn/keptn/cli/utils"
	"github.com/keptn/keptn/cli/utils/credentialmanager"
)

func init() {
	utils.InitLoggers(os.Stdout, os.Stdout, os.Stderr)
}

// TestNewArtifact tests the new-artifact command.
func TestNewArtifact(t *testing.T) {

	credentialmanager.MockCreds = true

	buf := new(bytes.Buffer)
	rootCmd.SetOutput(buf)

	args := []string{
		"new-artifact",
		fmt.Sprintf("--service=%s", "carts"),
		fmt.Sprintf("--stage=%s", "dev"),
		fmt.Sprintf("--docker-image=%s", "keptn/carts:latest"),
	}
	rootCmd.SetArgs(args)
	err := rootCmd.Execute()

	if err != nil {
		t.Errorf("An error occured %v", err)
	}
}
