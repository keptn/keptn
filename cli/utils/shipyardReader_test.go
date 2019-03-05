package utils

import (
	"os"
	"testing"
)

func init() {
	Init(os.Stdout, os.Stdout, os.Stderr)
}

func TestShipyardReader(t *testing.T) {

	data := `
stages: 
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

	stages := UnmarshalStages([]byte(data))

	if len(stages) != 3 {
		t.Fatal("invalid number of stages")
	}
	if len(stages[1]) != 8 {
		t.Fatal("invalid number of parameters for first stage")
	}

}
