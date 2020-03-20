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
func testShipyard(t *testing.T, shipyardFileName string, shipyard string) func() {
	if shipyard == "" {
		shipyard = `stages:
  - name: dev
    deployment_strategy: direct
  - name: staging
    deployment_strategy: blue_green_service
  - name: production
    deployment_strategy: blue_green_service`
	}

	ioutil.WriteFile(shipyardFileName, []byte(shipyard), 0644)

	buf := new(bytes.Buffer)
	rootCmd.SetOutput(buf)

	return func() {
		os.Remove(shipyardFileName)
	}
}

// TestCreateProjectCmd tests the default use of the create project command
func TestCreateProjectCmd(t *testing.T) {
	credentialmanager.MockAuthCreds = true

	shipyardFileName := "shipyard.yaml"
	defer testShipyard(t, shipyardFileName, "")()

	args := []string{
		"create",
		"project",
		"sockshop",
		fmt.Sprintf("--shipyard=%s", shipyardFileName),
		"--mock",
	}
	rootCmd.SetArgs(args)
	err := rootCmd.Execute()

	if err != nil {
		t.Errorf("An error occured: %v", err)
	}
}

// TestCreateProjectIncorrectProjectNameCmd tests whether the create project command aborts
// due to a project name with upper case character
func TestCreateProjectIncorrectProjectNameCmd(t *testing.T) {
	credentialmanager.MockAuthCreds = true

	shipyardFileName := "shipyard.yaml"
	defer testShipyard(t, shipyardFileName, "")()

	args := []string{
		"create",
		"project",
		"Sockshop", // invalid name, only lowercase is allowed
		fmt.Sprintf("--shipyard=%s", shipyardFileName),
	}
	rootCmd.SetArgs(args)
	err := rootCmd.Execute()

	if err != nil {
		if !errorContains(err, "contains upper case letter(s) or special character(s)") {
			t.Errorf("An error occured: %v", err)
		}
	} else {
		t.Fail()
	}
}

// TestCreateProjectIncorrectProjectNameCmd tests whether the create project command aborts
// due to a stage name, which contains a special character (-)
func TestCreateProjectIncorrectStageNameCmd(t *testing.T) {
	credentialmanager.MockAuthCreds = true

	shipyardFileName := "shipyard.yaml"
	shipyardContent := `stages:
- name: dev
  deployment_strategy: direct
- name: staging-projectA
  deployment_strategy: blue_green_service
- name: production
  deployment_strategy: blue_green_service`

	defer testShipyard(t, shipyardFileName, shipyardContent)()

	args := []string{
		"create",
		"project",
		"sockshop",
		fmt.Sprintf("--shipyard=%s", shipyardFileName),
	}
	rootCmd.SetArgs(args)
	err := rootCmd.Execute()

	if err != nil {
		if !errorContains(err, "contains upper case letter(s) or special character(s)") {
			t.Errorf("An error occured: %v", err)
		}
	} else {
		t.Fail()
	}
}

// TestCreateProjectCmdWithGitMissingParam tests whether the create project command aborts
// due to a missing parameters for defining a git upstream
func TestCreateProjectCmdWithGitMissingParam(t *testing.T) {
	credentialmanager.MockAuthCreds = true

	shipyardFileName := "shipyard.yaml"
	defer testShipyard(t, shipyardFileName, "")()

	args := []string{
		"create",
		"project",
		"sockshop",
		fmt.Sprintf("--shipyard=%s", shipyardFileName),
		fmt.Sprintf("--git-user=%s", "user"),
		fmt.Sprintf("--git-token=%s", "token"),
		"--mock",
	}
	rootCmd.SetArgs(args)
	err := rootCmd.Execute()

	if err != nil {
		if !errorContains(err, "For configuring a Git upstream") {
			t.Errorf("An error occured: %v", err)
		}
	} else {
		t.Fail()
	}
}

// TestCreateProjectCmdWithGitMissingParam tests a successful create project
// command with git upstream parameters
func TestCreateProjectCmdWithGit(t *testing.T) {
	credentialmanager.MockAuthCreds = true

	shipyardFileName := "shipyard.yaml"
	defer testShipyard(t, shipyardFileName, "")()

	args := []string{
		"create",
		"project",
		"sockshop",
		fmt.Sprintf("--shipyard=%s", shipyardFileName),
		fmt.Sprintf("--git-user=%s", "user"),
		fmt.Sprintf("--git-token=%s", "token"),
		fmt.Sprintf("--git-remote-url=%s", "https://"),
		"--mock",
	}
	rootCmd.SetArgs(args)
	err := rootCmd.Execute()

	if err != nil {
		t.Fail()
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
