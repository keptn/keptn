package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/keptn/keptn/cli/pkg/credentialmanager"
	"github.com/keptn/keptn/cli/pkg/logging"
)

func init() {
	logging.InitLoggers(os.Stdout, os.Stdout, os.Stderr)
}

// testShipyard writes a default shipyard file or uses the value from the shipyard parameter.
// It returns a function to delete the shipyard file.
func testShipyard(t *testing.T, shipyardFilePath string, shipyard string) func() {
	if shipyard == "" {
		shipyard = `stages:
  - name: dev
    deployment_strategy: direct
  - name: staging
    deployment_strategy: blue_green_service
  - name: production
    deployment_strategy: blue_green_service`
	}

	ioutil.WriteFile(shipyardFilePath, []byte(shipyard), 0644)

	buf := new(bytes.Buffer)
	rootCmd.SetOutput(buf)

	return func() {
		os.Remove(shipyardFilePath)
	}
}

// TestCreateProjectCmd tests the default use of the create project command
func TestCreateProjectCmd(t *testing.T) {
	credentialmanager.MockAuthCreds = true
	checkEndPointStatusMock = true

	shipyardFilePath := "./shipyard.yaml"
	defer testShipyard(t, shipyardFilePath, "")()

	cmd := fmt.Sprintf("create project sockshop --shipyard=%s --mock", shipyardFilePath)
	_, err := executeActionCommandC(cmd)
	if err != nil {
		t.Errorf(unexpectedErrMsg, err)
	}
}

// TestCreateProjectIncorrectProjectNameCmd tests whether the create project command aborts
// due to a project name with upper case character
func TestCreateProjectIncorrectProjectNameCmd(t *testing.T) {
	credentialmanager.MockAuthCreds = true
	checkEndPointStatusMock = true

	shipyardFilePath := "./shipyard.yaml"
	defer testShipyard(t, shipyardFilePath, "")()

	cmd := fmt.Sprintf("create project Sockshop --shipyard=%s --mock", shipyardFilePath)
	_, err := executeActionCommandC(cmd)

	if !errorContains(err, "contains upper case letter(s) or special character(s)") {
		t.Errorf("missing expected error, but got %v", err)
	}
}

// TestCreateProjectIncorrectProjectNameCmd tests whether the create project command aborts
// due to a stage name, which contains a special character (-)
func TestCreateProjectIncorrectStageNameCmd(t *testing.T) {
	credentialmanager.MockAuthCreds = true
	checkEndPointStatusMock = true

	shipyardFilePath := "./shipyard.yaml"
	shipyardContent := `stages:
- name: dev
  deployment_strategy: direct
- name: staging-projectA
  deployment_strategy: blue_green_service
- name: production
  deployment_strategy: blue_green_service`

	defer testShipyard(t, shipyardFilePath, shipyardContent)()

	cmd := fmt.Sprintf("create project Sockshop --shipyard=%s --mock", shipyardFilePath)
	_, err := executeActionCommandC(cmd)

	if !errorContains(err, "contains upper case letter(s) or special character(s)") {
		t.Errorf("missing expected error, but got %v", err)
	}
}

// TestCreateProjectCmdWithGitMissingParam tests whether the create project command aborts
// due to a missing flag for defining a git upstream
func TestCreateProjectCmdWithGitMissingParam(t *testing.T) {
	credentialmanager.MockAuthCreds = true
	checkEndPointStatusMock = true

	shipyardFilePath := "./shipyard.yaml"
	defer testShipyard(t, shipyardFilePath, "")()

	cmd := fmt.Sprintf("create project sockshop --shipyard=%s --git-user=%s --git-token=%s --mock",
		shipyardFilePath, "user", "token")
	_, err := executeActionCommandC(cmd)

	if !errorContains(err, gitErrMsg) {
		t.Errorf("missing expected error, but got %v", err)
	}
}

// TestCreateProjectCmdWithGitMissingParam tests a successful create project
// command with git upstream parameters
func TestCreateProjectCmdWithGit(t *testing.T) {
	credentialmanager.MockAuthCreds = true
	checkEndPointStatusMock = true

	shipyardFilePath := "./shipyard.yaml"
	defer testShipyard(t, shipyardFilePath, "")()

	cmd := fmt.Sprintf("create project sockshop --shipyard=%s --git-user=%s --git-token=%s --git-remote-url=%s --mock",
		shipyardFilePath, "user", "token", "https://")
	_, err := executeActionCommandC(cmd)

	if err != nil {
		t.Errorf(unexpectedErrMsg, err)
	}
}

func errorContains(out error, want string) bool {
	if out == nil {
		return want == ""
	}
	if want == "" {
		return false
	}
	return strings.Contains(out.Error(), want)
}
