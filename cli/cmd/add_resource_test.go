package cmd

import (
	"bytes"
	"fmt"
	"github.com/keptn/keptn/cli/utils/credentialmanager"
	"io/ioutil"
	"os"
	"testing"

	"github.com/keptn/keptn/cli/pkg/logging"
)

func init() {
	logging.InitLoggers(os.Stdout, os.Stdout, os.Stderr)
}

// testResource writes a default file
func testResource(t *testing.T, fileName string, fileContent string) func() {
	if fileContent == "" {
		fileContent = `This is a test file`
	}

	ioutil.WriteFile(fileName, []byte(fileContent), 0644)

	buf := new(bytes.Buffer)
	rootCmd.SetOutput(buf)

	return func() {
		os.Remove(fileName)
	}
}

// TestAuthCmd tests the auth command. Therefore, this test assumes a file "~/keptn/.keptnmock" containing
// the endpoint and api-token.
func TestAddResource(t *testing.T) {

	credentialmanager.MockAuthCreds = true

	resourceFileName := "testResource.txt"
	defer testResource(t, resourceFileName, "")()

	buf := new(bytes.Buffer)
	rootCmd.SetOutput(buf)

	args := []string{
		"add-resource",

		fmt.Sprintf("--project=%s", "sockshop"),
		fmt.Sprintf("--service=%s", "carts"),
		fmt.Sprintf("--stage=%s", "dev"),
		fmt.Sprintf("--resource=%s", resourceFileName),
		fmt.Sprintf("--resourceUri=%s", "resource/"+resourceFileName),
		"--mock",
	}
	rootCmd.SetArgs(args)
	err := rootCmd.Execute()

	if err != nil {
		t.Errorf("An error occured: %v", err)
	}
}
