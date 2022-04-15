package ZD

import (
	"fmt"
	"github.com/stretchr/testify/suite"
	"testing"

	//"github.com/anandvarma/namegen"
	"github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/go-utils/pkg/lib/v0_2_0"
	testutils "github.com/keptn/keptn/test/go-tests"
	"gopkg.in/yaml.v3"
	"sync/atomic"
	"time"
)

var shipyardFile, _ = getShipyard()

type TestEvaluation struct {
	id          string
	projectName string
	serviceName string
	stageName   string
	suite.Suite
}

func (suite *TestEvaluation) SetupSuite() {
	var err error
	suite.id = gedId()
	suite.serviceName = "myservice"
	suite.stageName = "dev"

	projectName := "zd" + suite.id
	suite.projectName, err = testutils.CreateProject(projectName, shipyardFile)
	suite.Nil(err)

	output, err := testutils.ExecuteCommand(fmt.Sprintf("keptn create service %s --project=%s", suite.serviceName, suite.projectName))
	suite.Nil(err)
	suite.Contains(output, "created successfully")
}

func Test_Evaluation(t *testing.T) {
	suite.Run(t, new(TestEvaluation))
}

func (suite *TestEvaluation) Test_EvaluationFails() {
	wg.Add(1)

	suite.T().Logf("Starting test evaluation for project %s ", suite.projectName)

	var keptnContext string
	// trigger a delivery sequence
	keptnContext, err := testutils.TriggerSequence(suite.projectName, suite.serviceName, suite.stageName, "evaluation", nil)
	suite.Nil(err)
	go suite.checkSequence(keptnContext, "evaluation")

}

func (suite *TestEvaluation) Test_DeliveryFails() {

	wg.Add(1)

	suite.T().Logf("Starting test Delivery for project %s ", suite.projectName)

	var keptnContext string
	// trigger a delivery sequence
	keptnContext, err := testutils.TriggerSequence(suite.projectName, suite.serviceName, suite.stageName, "delivery", nil)
	suite.Nil(err)
	go suite.checkSequence(keptnContext, "delivery")

}

func (suite *TestEvaluation) checkSequence(keptnContext string, sequenceType string) {

	defer wg.Done()
	var sequenceFinishedEvent *models.KeptnContextExtendedCE
	stageSequenceName := fmt.Sprintf("%s.%s", suite.stageName, sequenceType)
	var err error
	suite.Eventually(func() bool {
		sequenceFinishedEvent, err = testutils.GetLatestEventOfType(keptnContext, suite.projectName, suite.stageName, v0_2_0.GetFinishedEventType(stageSequenceName))
		if sequenceFinishedEvent == nil || err != nil {
			return false
		}
		atomic.AddUint64(&PassedSequences, 1)
		return true
	}, 30*time.Second, 5*time.Second)
	if sequenceFinishedEvent == nil || err != nil {
		suite.T().Logf("Sequence %s in stage %s has not been finished", keptnContext, suite.stageName)
		atomic.AddUint64(&FailedSequences, 1)
	}

}

func getShipyard() (string, error) {
	shipyard := &v0_2_0.Shipyard{
		ApiVersion: "0.2.3",
		Kind:       "shipyard",
		Metadata:   v0_2_0.Metadata{},
		Spec: v0_2_0.ShipyardSpec{
			Stages: []v0_2_0.Stage{},
		},
	}

	stage := v0_2_0.Stage{
		Name: "dev",
		Sequences: []v0_2_0.Sequence{
			{
				Name: "hooks",
				Tasks: []v0_2_0.Task{
					{
						Name: "mytask",
					},
				},
			},
		},
	}

	shipyard.Spec.Stages = append(shipyard.Spec.Stages, stage)

	shipyardFileContent, _ := yaml.Marshal(shipyard)

	return testutils.CreateTmpShipyardFile(string(shipyardFileContent))
}

func gedId() string {
	atomic.AddUint64(&Id, 1)
	return fmt.Sprintf("%d", Id)
}
