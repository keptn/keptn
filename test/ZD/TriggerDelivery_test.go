package ZD

import (
	"fmt"
	"github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/go-utils/pkg/lib/keptn"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"

	//"github.com/anandvarma/namegen"
	"github.com/keptn/go-utils/pkg/lib/v0_2_0"
	testutils "github.com/keptn/keptn/test/go-tests"
	"gopkg.in/yaml.v3"
	"sync/atomic"
)

var shipyardFile, _ = getShipyard()

type TestEvaluation struct {
	suite.Suite
	project string
}

//type TriggeredSequences struct {
//	sequences []TriggeredSequence
//	mutex     *sync.Mutex
//}

//func NewTriggeredSequences() *TriggeredSequences {
//	return &TriggeredSequences{
//		sequences: []TriggeredSequence{},
//		mutex:     new(sync.Mutex),
//	}
//}
//
//func (ts *TriggeredSequences) Add(s TriggeredSequence) {
//	ts.mutex.Lock()
//	defer ts.mutex.Unlock()
//	ts.sequences = append(ts.sequences, s)
//}

type TriggeredSequence struct {
	keptnContext string
	projectName  string
	sequenceName string
}

func NewTriggeredSequence(keptnContext string, projectName string, seqName string) TriggeredSequence {
	return TriggeredSequence{
		keptnContext: keptnContext,
		projectName:  projectName,
		sequenceName: seqName,
	}
}

func (suite *TestEvaluation) SetupSuite() {
	var err error
	projectName := "zd" + gedId()
	suite.project, err = testutils.CreateProject(projectName, shipyardFile)
	suite.Nil(err)

	output, err := testutils.ExecuteCommand(fmt.Sprintf("keptn create service %s --project=%s", "myservice", suite.project))
	suite.Nil(err)
	suite.Contains(output, "created successfully")

	suite.T().Logf("Starting test for project %s ", suite.project)
}

func Test_Run(t *testing.T) {
	wg.Add(2)
	atomic.AddUint64(&FiredSequences, 2)
	suite.Run(t, new(TestEvaluation))
}

func (suite *TestEvaluation) Test_EvaluationFails() {
	suite.trigger("evaluation", nil)
}

func (suite *TestEvaluation) Test_DeliveryFails() {
	suite.trigger("delivery", nil)
}

func (suite *TestEvaluation) trigger(triggerType string, data keptn.EventProperties) {

	// trigger a delivery sequence
	keptnContext, err := testutils.TriggerSequence(suite.project, "myservice", "dev", triggerType, data)
	suite.Nil(err)
	//
	sequence := NewTriggeredSequence(keptnContext, suite.project, triggerType)
	//sequences.Add(sequence)

	go suite.checkSequence(keptnContext, sequence)
}

func (suite *TestEvaluation) checkSequence(keptnContext string, sequence TriggeredSequence) {

	defer wg.Done()
	var sequenceFinishedEvent *models.KeptnContextExtendedCE
	stageSequenceName := fmt.Sprintf("%s.%s", "dev", sequence.sequenceName)
	var err error
	suite.Eventually(func() bool {
		sequenceFinishedEvent, err = testutils.GetLatestEventOfType(keptnContext, sequence.projectName, "dev", v0_2_0.GetFinishedEventType(stageSequenceName))
		if sequenceFinishedEvent == nil || err != nil {
			return false
		}
		atomic.AddUint64(&PassedSequences, 1)
		return true
	}, 30*time.Second, 5*time.Second)
	if sequenceFinishedEvent == nil || err != nil {
		suite.T().Logf("Sequence %s in stage %s has not been finished", keptnContext, sequence.projectName)
		atomic.AddUint64(&FailedSequences, 1)
	}

	suite.T().Logf("Finished %s in %s ", sequence.sequenceName, sequence.projectName)
}

//func CheckSequence(sequence TriggeredSequence, wg sync.WaitGroup, t *testing.T) {
//	defer wg.Done()
//	t.Logf("Checking Sequence %s in project %s", sequence.keptnContext, sequence.projectName)
//	var sequenceFinishedEvent *models.KeptnContextExtendedCE
//	stageSequenceName := fmt.Sprintf("%s.%s", "dev", sequence.sequenceName)
//	var err error
//	assert.Eventually(t, func() bool {
//		sequenceFinishedEvent, err = testutils.GetLatestEventOfType(sequence.keptnContext, sequence.projectName, "dev", v0_2_0.GetFinishedEventType(stageSequenceName))
//		if sequenceFinishedEvent == nil || err != nil {
//			return false
//		}
//		atomic.AddUint64(&PassedSequences, 1)
//
//		return true
//	}, 10*time.Second, 5*time.Second)
//	if sequenceFinishedEvent == nil || err != nil {
//
//		t.Logf("Sequence %s in stage %s has not been finished", sequence.keptnContext, "dev")
//		atomic.AddUint64(&FailedSequences, 1)
//
//	}
//}

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
