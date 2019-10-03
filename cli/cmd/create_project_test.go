package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/keptn/keptn/cli/pkg/logging"
	"github.com/keptn/keptn/cli/utils"
	"github.com/keptn/keptn/cli/utils/credentialmanager"
)

func init() {
	logging.InitLoggers(os.Stdout, os.Stdout, os.Stderr)
}

func TestCreateProjectCmd(t *testing.T) {

	credentialmanager.MockAuthCreds = true

	// Write temporary shipyardTest.yml file
	const tmpShipyardFileName = "shipyardTest.yml"
	shipYardContent := `stages: 
- name: dev
  deployment_strategy: direct
- name: staging
  deployment_strategy: blue_green_service
- name: production
  deployment_strategy: blue_green_service`

	ioutil.WriteFile(tmpShipyardFileName, []byte(shipYardContent), 0644)

	buf := new(bytes.Buffer)
	rootCmd.SetOutput(buf)

	args := []string{
		"create",
		"project",
		"sockshop",
		fmt.Sprintf("--shipyard=%s", tmpShipyardFileName),
		"--mock",
	}
	rootCmd.SetArgs(args)
	err := rootCmd.Execute()

	// Delete temporary shipyard.yml file
	os.Remove(tmpShipyardFileName)

	if err != nil {
		t.Errorf("An error occured %v", err)
	}
}

func TestCreateProjectIncorrectProjectNameCmd(t *testing.T) {

	credentialmanager.MockAuthCreds = true

	// Write temporary shipyardTest.yml file
	const tmpShipyardFileName = "shipyardTest.yml"
	shipYardContent := `stages:
- name: dev
  deployment_strategy: direct
- name: staging
  deployment_strategy: blue_green_service
- name: production
  deployment_strategy: blue_green_service`

	ioutil.WriteFile(tmpShipyardFileName, []byte(shipYardContent), 0644)

	buf := new(bytes.Buffer)
	rootCmd.SetOutput(buf)

	args := []string{
		"create",
		"project",
		"Sockshop", // invalid name, only lowercase is allowed
		fmt.Sprintf("--shipyard=%s", tmpShipyardFileName),
	}
	rootCmd.SetArgs(args)
	err := rootCmd.Execute()

	// Delete temporary shipyard.yml file
	os.Remove(tmpShipyardFileName)

	if err != nil {
		if !utils.ErrorContains(err, "Project name contains invalid characters or is not well-formed.") {
			t.Errorf("An error occured: %v", err)
		}
	} else {
		t.Fail()
	}
}

func TestCreateProjectIncorrectStageNameCmd(t *testing.T) {

	credentialmanager.MockAuthCreds = true

	// Write temporary shipyardTest.yml file
	const tmpShipyardFileName = "shipyardTest.yml"
	shipYardContent := `stages:
- name: dev
  deployment_strategy: direct
- name: staging-team1
  deployment_strategy: blue_green_service
- name: production
  deployment_strategy: blue_green_service`

	ioutil.WriteFile(tmpShipyardFileName, []byte(shipYardContent), 0644)

	buf := new(bytes.Buffer)
	rootCmd.SetOutput(buf)

	args := []string{
		"create",
		"project",
		"sockshop",
		fmt.Sprintf("--shipyard=%s", tmpShipyardFileName),
	}
	rootCmd.SetArgs(args)
	err := rootCmd.Execute()

	// Delete temporary shipyard.yml file
	os.Remove(tmpShipyardFileName)

	if err != nil {
		if !utils.ErrorContains(err, "Stage staging-team1 contains invalid characters or is not well-formed.") {
			t.Errorf("An error occured: %v", err)
		}
	} else {
		t.Fail()
	}
}

func TestCreateProjectCmdWithGitMissingParam(t *testing.T) {

	credentialmanager.MockAuthCreds = true

	// Write temporary shipyardTest.yml file
	const tmpShipyardFileName = "shipyardTest.yml"
	shipYardContent := `stages: 
- name: dev
  deployment_strategy: direct
- name: staging
  deployment_strategy: blue_green_service
- name: production
  deployment_strategy: blue_green_service`

	ioutil.WriteFile(tmpShipyardFileName, []byte(shipYardContent), 0644)

	buf := new(bytes.Buffer)
	rootCmd.SetOutput(buf)

	args := []string{
		"create",
		"project",
		"sockshop",
		fmt.Sprintf("--shipyard=%s", tmpShipyardFileName),
		fmt.Sprintf("--git-user=%s", "user"),
		fmt.Sprintf("--git-token=%s", "token"),
		"--mock",
	}
	rootCmd.SetArgs(args)
	err := rootCmd.Execute()

	// Delete temporary shipyard.yml file
	os.Remove(tmpShipyardFileName)

	if err != nil {
		if !utils.ErrorContains(err, "For configuring a Git upstream") {
			t.Errorf("An error occured: %v", err)
		}
	} else {
		t.Fail()
	}
}

func TestCreateProjectCmdWithGit(t *testing.T) {

	credentialmanager.MockAuthCreds = true

	// Write temporary shipyardTest.yml file
	const tmpShipyardFileName = "shipyardTest.yml"
	shipYardContent := `stages: 
- name: dev
  deployment_strategy: direct
- name: staging
  deployment_strategy: blue_green_service
- name: production
  deployment_strategy: blue_green_service`

	ioutil.WriteFile(tmpShipyardFileName, []byte(shipYardContent), 0644)

	buf := new(bytes.Buffer)
	rootCmd.SetOutput(buf)

	args := []string{
		"create",
		"project",
		"sockshop",
		fmt.Sprintf("--shipyard=%s", tmpShipyardFileName),
		fmt.Sprintf("--git-user=%s", "user"),
		fmt.Sprintf("--git-token=%s", "token"),
		fmt.Sprintf("--git-remote-url=%s", "https://"),
		"--mock",
	}
	rootCmd.SetArgs(args)
	err := rootCmd.Execute()

	// Delete temporary shipyard.yml file
	os.Remove(tmpShipyardFileName)

	if err != nil {
		t.Fail()
	}
}
