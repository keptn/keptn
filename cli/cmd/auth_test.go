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

// This test assumes that the endpoint and api token are already saved as credentials
func TestValidAuthCmd(t *testing.T) {

	buf := new(bytes.Buffer)
	rootCmd.SetOutput(buf)

	endPoint, apiToken, err := credentialmanager.GetCreds()
	if err != nil {
		t.Errorf("An error occured %v", err)
	}

	args := []string{
		"auth",
		fmt.Sprintf("--endpoint=%s", endPoint),
		fmt.Sprintf("--api-token=%s", apiToken),
	}
	rootCmd.SetArgs(args)
	err = authCmd.Execute()

	if err != nil {
		t.Errorf("An error occured %v", err)
	}
}
