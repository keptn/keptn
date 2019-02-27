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
	const tmpShipyardFileName = "shipyardTest.tyml"
	shipYardContent := `stages: 
- name: dev
  deployment_strategy: direct
  deployment_operator: jenkins-operator, slack
  test_strategy: functional
  test_operator: neotys_operator
  validation_operator: keptn.monspec-evaluator
  remediation_handler: // TBD    
  next: staging
- name: staging
  deployment_strategy: service_blue/green
  deployment_operator: jenkins-operator, slack
  test_strategy: continous_performance
  test_operator: neotys_operator
  validation_operator: keptn.monspec-evaluator
  remediation_handler: rollback
  next: production
- name: production
  deployment_strategy: application blue/green
  deployment_operator: jenkins-operator, slack
  test_strategy: production
  test_operator: neotys_operator
  validation_strategy: production
  remediation_handler: rollback`

	ioutil.WriteFile(tmpShipyardFileName, []byte(shipYardContent), 0644)

	buf := new(bytes.Buffer)
	rootCmd.SetOutput(buf)

	args := []string{
		"create",
		"project",
		"sockshop",
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
