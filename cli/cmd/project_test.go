package cmd

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	"github.com/keptn/keptn/cli/utils"
	"github.com/keptn/keptn/cli/utils/credentialmanager"
)

func init() {
	utils.Init(os.Stdout, os.Stdout, os.Stderr)
}

func TestCreateProjectCmd(t *testing.T) {

	credentialmanager.MockCreds = true

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
		"sockshop20",
		tmpShipyardFileName,
	}
	rootCmd.SetArgs(args)
	err := rootCmd.Execute()

	// Delete temporary shipyard.yml file
	os.Remove(tmpShipyardFileName)

	if err != nil {
		t.Errorf("An error occured %v", err)
	}
}
