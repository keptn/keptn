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

// TestAuthCmd tests the auth command. Therefore, this test assumes a file "~/keptn/.keptnmock" containing
// the endpoint and api-token.
func TestAuthCmd(t *testing.T) {

	credentialmanager.MockAuthCreds = true

	endPoint, apiToken, err := credentialmanager.GetCreds()
	if err != nil {
		t.Error(err)
		return
	}

	buf := new(bytes.Buffer)
	rootCmd.SetOutput(buf)

	args := []string{
		"auth",

		fmt.Sprintf("--endpoint=%s", endPoint.String()),
		fmt.Sprintf("--api-token=%s", apiToken),
		"--mock",
	}
	rootCmd.SetArgs(args)
	err = rootCmd.Execute()

	if err != nil {
		t.Errorf("An error occured: %v", err)
	}
}
