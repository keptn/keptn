package cmd

import (
	"bytes"
	"log"
	"os"
	"testing"

	"github.com/keptn/keptn/cli/pkg/logging"
	"github.com/keptn/keptn/cli/utils/credentialmanager"
)

func init() {
	logging.InitLoggers(os.Stdout, os.Stdout, os.Stderr)
}

// TestOnboardService tests the onboard service command.
func TestOnboardService(t *testing.T) {

	credentialmanager.MockAuthCreds = true
	buf := new(bytes.Buffer)
	rootCmd.SetOutput(buf)

	args := []string{
		"onboard",
		"service",
		"carts",
		"--project=sockshop",
		"--chart=carts",
	}
	rootCmd.SetArgs(args)
	err := rootCmd.Execute()

	if err != nil {
		log.Fatalf("An error occured %v", err)
	}
}
