package cmd

import (
	"bytes"
	"os"
	"testing"

	"github.com/keptn/keptn/cli/utils"
	"github.com/keptn/keptn/cli/utils/credentialmanager"
)

func init() {
	utils.Init(os.Stdout, os.Stdout, os.Stderr)
}

func TestCreateProjectCmd(t *testing.T) {

	credentialmanager.MockCreds = true

	buf := new(bytes.Buffer)
	rootCmd.SetOutput(buf)

	args := []string{
		"create",
		"project",
		"sockshop",
		"shipyard.yml",
	}
	rootCmd.SetArgs(args)
	err := rootCmd.Execute()

	if err != nil {
		t.Errorf("An error occured %v", err)
	}
}
