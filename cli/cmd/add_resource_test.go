package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/keptn/keptn/cli/pkg/credentialmanager"

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

// TestAddResourceToProjectStageService
func TestAddResourceToProjectStageService(t *testing.T) {

	credentialmanager.MockAuthCreds = true
	checkEndPointStatusMock = true

	resourceFileName := "testResource.txt"
	defer testResource(t, resourceFileName, "")()

	cmd := fmt.Sprintf("add-resource --project=%s --stage=%s --service=%s --resource=%s "+
		"--resourceUri=%s --mock", "sockshop", "dev", "carts", resourceFileName, "resource/"+resourceFileName)
	_, err := executeActionCommandC(cmd)
	if err != nil {
		t.Errorf(unexpectedErrMsg, err)
	}
}

// TestAddResourceToProjectStage
func TestAddResourceToProjectStage(t *testing.T) {

	credentialmanager.MockAuthCreds = true
	checkEndPointStatusMock = true

	resourceFileName := "testResource.txt"
	defer testResource(t, resourceFileName, "")()

	cmd := fmt.Sprintf("add-resource --project=%s --stage=%s --resource=%s "+
		"--resourceUri=%s --mock", "sockshop", "dev", resourceFileName, "resource/"+resourceFileName)
	_, err := executeActionCommandC(cmd)
	if err != nil {
		t.Errorf(unexpectedErrMsg, err)
	}
}

// TestAddResourceToProject
func TestAddResourceToProject(t *testing.T) {

	credentialmanager.MockAuthCreds = true
	checkEndPointStatusMock = true

	resourceFileName := "testResource.txt"
	defer testResource(t, resourceFileName, "")()

	cmd := fmt.Sprintf("add-resource --project=%s --resource=%s "+
		"--resourceUri=%s --mock", "sockshop", resourceFileName, "resource/"+resourceFileName)
	_, err := executeActionCommandC(cmd)
	if err != nil {
		t.Errorf(unexpectedErrMsg, err)
	}
}

// TestAddResourceToProjectService
func TestAddResourceToProjectService(t *testing.T) {

	credentialmanager.MockAuthCreds = true
	checkEndPointStatusMock = true

	*addResourceCmdParams.Stage = ""

	resourceFileName := "testResource.txt"
	defer testResource(t, resourceFileName, "")()

	cmd := fmt.Sprintf("add-resource --project=%s --service=%s --resource=%s "+
		"--resourceUri=%s --mock", "sockshop", "carts", resourceFileName, "resource/"+resourceFileName)
	_, err := executeActionCommandC(cmd)
	if err.Error() != "Flag 'stage' is missing" {
		t.Errorf(unexpectedErrMsg, err)
	}
}

// TestAddResourceWhenArgsArePresent
func TestAddResourceWhenArgsArePresent(t *testing.T) {

	credentialmanager.MockAuthCreds = true
	checkEndPointStatusMock = true

	resourceFileName := "testResource.txt"
	defer testResource(t, resourceFileName, "")()

	cmd := fmt.Sprintf("add-resource --project=%s --stage=%s --service=%s --resource=%s "+
		"-- resourceUri=%s --mock", "sockshop", "dev", "carts", resourceFileName, "resource/"+resourceFileName)
	_, err := executeActionCommandC(cmd)
	if err == nil {
		t.Errorf("Expected an error")
	}
	got := err.Error()
	expected := "accepts 0 arg(s), received 2"
	if got != expected {
		t.Errorf("Expected %q, got %q", expected, got)
	}
}
