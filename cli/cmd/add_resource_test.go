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

func setup(t *testing.T) {
	t.Setenv("MOCK_SERVER", "http://some-valid-url.com")
	credentialmanager.MockAuthCreds = true

	*addResourceCmdParams.AllStages = false
	*addResourceCmdParams.Stage = ""
	*addResourceCmdParams.Service = ""
	*addResourceCmdParams.Project = ""
	*addResourceCmdParams.Resource = ""
	*addResourceCmdParams.ResourceURI = ""
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

	setup(t)

	resourceFileName := "testResource.txt"
	defer testResource(t, resourceFileName, "")()

	cmd := fmt.Sprintf("add-resource --project=%s --stage=%s --service=%s --resource=%s "+
		"--resourceUri=%s --mock", "sockshop", "dev", "carts", resourceFileName, "resource/"+resourceFileName)
	_, err := executeActionCommandC(cmd)
	if err != nil {
		t.Errorf(unexpectedErrMsg, err)
	}
}

// TestAddResourceToProjectStageAndAllStages tests that using --stage and --all-stages together doesn't work
func TestAddResourceToProjectStageAndAllStages(t *testing.T) {

	setup(t)

	resourceFileName := "testResource.txt"
	defer testResource(t, resourceFileName, "")()

	cmd := fmt.Sprintf("add-resource --project=%s --stage=%s --all-stages --service=%s --resource=%s "+
		"--resourceUri=%s --mock", "sockshop", "dev", "carts", resourceFileName, "resource/"+resourceFileName)
	_, err := executeActionCommandC(cmd)
	if err.Error() != "Cannot use --stage and --all-stages at the same time" {
		t.Errorf(unexpectedErrMsg, err)
	}
}

// TestAddResourceToProjectAndStage tests that using --project and --stage (without --service) works
func TestAddResourceToProjectAndStage(t *testing.T) {

	setup(t)

	resourceFileName := "testResource.txt"
	defer testResource(t, resourceFileName, "")()

	cmd := fmt.Sprintf("add-resource --project=%s --stage=%s --resource=%s "+
		"--resourceUri=%s --mock", "sockshop", "dev", resourceFileName, "resource/"+resourceFileName)
	_, err := executeActionCommandC(cmd)
	if err != nil {
		t.Errorf(unexpectedErrMsg, err)
	}
}

// TestAddResourceAllStages tests that using --all-stages works
func TestAddResourceAllStages(t *testing.T) {

	setup(t)

	resourceFileName := "testResource.txt"
	defer testResource(t, resourceFileName, "")()

	cmd := fmt.Sprintf("add-resource --all-stages --project=%s --resource=%s "+
		"--resourceUri=%s --mock", "sockshop", resourceFileName, "resource/"+resourceFileName)
	_, err := executeActionCommandC(cmd)
	if err != nil {
		t.Errorf(unexpectedErrMsg, err)
	}
}

// TestAddResourceToProjectServiceAllStages tests that using --project, --service and --all-stages works
func TestAddResourceToProjectServiceAllStages(t *testing.T) {

	setup(t)

	resourceFileName := "testResource.txt"
	defer testResource(t, resourceFileName, "")()

	cmd := fmt.Sprintf("add-resource --project=%s --service=%s --all-stages --resource=%s "+
		"--resourceUri=%s --mock", "sockshop", "carts", resourceFileName, "resource/"+resourceFileName)
	_, err := executeActionCommandC(cmd)
	if err != nil {
		t.Errorf(unexpectedErrMsg, err)
	}
}

// TestAddResourceToProjectServiceStage tests that using --project, --service and --stage works
func TestAddResourceToProjectServiceStage(t *testing.T) {

	setup(t)

	resourceFileName := "testResource.txt"
	defer testResource(t, resourceFileName, "")()

	cmd := fmt.Sprintf("add-resource --project=%s --service=%s --stage=%s --resource=%s "+
		"--resourceUri=%s --mock", "sockshop", "carts", "dev", resourceFileName, "resource/"+resourceFileName)
	_, err := executeActionCommandC(cmd)
	if err != nil {
		t.Errorf(unexpectedErrMsg, err)
	}
}

// TestAddResourceToProjectStage
func TestAddResourceToProjectStage(t *testing.T) {

	setup(t)

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

	setup(t)

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

	setup(t)

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
	testInvalidInputHelper("add-resource --project=sockshop --stage=dev --service=carts --resource=testResource.txt "+
		"-- resourceUri=resource/testResource.txt --mock", "accepts 0 arg(s), received 2", t)
}

// TestAddResourceUnknownCommand
func TestAddResourceUnknownCommand(t *testing.T) {
	testInvalidInputHelper("add-resource someUnknownCommand --project=sockshop --stage=dev --service=carts --resource=testResource.txt "+
		"--resourceUri=resource/testResource.txt --mock", "accepts 0 arg(s), received 1", t)
}

// TestAddResourceUnknownParameter
func TestAddResourceUnknownParmeter(t *testing.T) {
	testInvalidInputHelper("add-resource --projectt=sockshop --stage=dev --service=carts --resource=testResource.txt "+
		"--resourceUri=resource/testResource.txt --mock", "unknown flag: --projectt", t)
}
