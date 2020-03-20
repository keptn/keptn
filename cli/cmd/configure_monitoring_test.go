package cmd

import (
	"bytes"
	"os"
	"testing"

	"github.com/keptn/keptn/cli/pkg/credentialmanager"
	"github.com/keptn/keptn/cli/pkg/logging"
)

func init() {
	logging.InitLoggers(os.Stdout, os.Stdout, os.Stderr)
}

func TestConfigureMonitoringCmd(t *testing.T) {

	credentialmanager.MockAuthCreds = true
	mocking = true
	buf := new(bytes.Buffer)
	rootCmd.SetOutput(buf)

	args := []string{
		"configure",
		"monitoring",
		"prometheus",
		"--project=sockshop",
		"--service=carts",
	}
	rootCmd.SetArgs(args)
	err := rootCmd.Execute()

	if err != nil {
		t.Errorf("An error occured: %v", err)
	}
}
