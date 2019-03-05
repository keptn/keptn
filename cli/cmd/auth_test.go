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
	utils.Init(os.Stdout, os.Stdout, os.Stderr)
}

// TestAuthCmd tests the auth command. Therefore, this test assumes a file named "endPoint.txt" containing
// the endpoint and api-token in the executing directory.
func TestAuthCmd(t *testing.T) {

	credentialmanager.MockCreds = true

	endPoint, apiToken, err := credentialmanager.GetCreds()
	if err != nil {
		t.Error(err)
		return
	}

	buf := new(bytes.Buffer)
	rootCmd.SetOutput(buf)

	args := []string{
		"auth",
		fmt.Sprintf("--endpoint=%s", endPoint),
		fmt.Sprintf("--api-token=%s", apiToken),
	}
	rootCmd.SetArgs(args)
	err = rootCmd.Execute()

	if err != nil {
		t.Errorf("An error occured %v", err)
	}
}
