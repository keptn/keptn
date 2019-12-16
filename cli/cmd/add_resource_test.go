package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/keptn/keptn/cli/utils/credentialmanager"

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

	resourceFileName := "testResource.txt"
	defer testResource(t, resourceFileName, "")()

	buf := new(bytes.Buffer)
	rootCmd.SetOutput(buf)

	args := []string{
		"add-resource",

		fmt.Sprintf("--project=%s", "sockshop"),
		fmt.Sprintf("--stage=%s", "dev"),
		fmt.Sprintf("--service=%s", "carts"),
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

// TestAddResourceToProjectStage
func TestAddResourceToProjectStage(t *testing.T) {

	credentialmanager.MockAuthCreds = true

	resourceFileName := "testResource.txt"
	defer testResource(t, resourceFileName, "")()

	buf := new(bytes.Buffer)
	rootCmd.SetOutput(buf)

	args := []string{
		"add-resource",

		fmt.Sprintf("--project=%s", "sockshop"),
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

// TestAddResourceToProject
func TestAddResourceToProject(t *testing.T) {

	credentialmanager.MockAuthCreds = true

	resourceFileName := "testResource.txt"
	defer testResource(t, resourceFileName, "")()

	buf := new(bytes.Buffer)
	rootCmd.SetOutput(buf)

	args := []string{
		"add-resource",

		fmt.Sprintf("--project=%s", "sockshop"),
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

// TestAddResourceToProjectService
func TestAddResourceToProjectService(t *testing.T) {

	credentialmanager.MockAuthCreds = true

	resourceFileName := "testResource.txt"
	defer testResource(t, resourceFileName, "")()

	buf := new(bytes.Buffer)
	rootCmd.SetOutput(buf)

	args := []string{
		"add-resource",

		fmt.Sprintf("--project=%s", "sockshop"),
		fmt.Sprintf("--service=%s", "carts"),
		fmt.Sprintf("--resource=%s", resourceFileName),
		fmt.Sprintf("--resourceUri=%s", "resource/"+resourceFileName),
		"--mock",
	}
	rootCmd.SetArgs(args)
	err := rootCmd.Execute()

	if err == nil {
		t.Errorf("No error occured: %v", err)
	}
}
