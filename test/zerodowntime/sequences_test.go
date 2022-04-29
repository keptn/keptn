package zerodowntime

import (
	"fmt"
	"github.com/benbjohnson/clock"
	"github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/go-utils/pkg/lib/keptn"
	"github.com/keptn/go-utils/pkg/lib/v0_2_0"
	testutils "github.com/keptn/keptn/test/go-tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gopkg.in/yaml.v3"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

type TestSuiteSequences struct {
	suite.Suite
	env     *ZeroDowntimeEnv
	project string
}

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

func (suite *TestSuiteSequences) SetupSuite() {

	suite.T().Log("Starting test for sequences")
	// if needed the following line can setup a project eat very clock tick
	suite.createNew()
}

func (suite *TestSuiteSequences) createNew() {
	var err error
	projectName := "zd-sequence" + suite.gedId()
	suite.T().Logf("Creating project with id %s ", projectName)
	suite.project, err = testutils.CreateProject(projectName, suite.env.ShipyardFile)
	suite.Nil(err)
	output, err := testutils.ExecuteCommand(fmt.Sprintf("keptn create service %s --project=%s", "myservice", suite.project))
	suite.Nil(err)
	suite.Contains(output, "created successfully")

	suite.T().Logf("Starting test for project %s ", suite.project)
}

func (suite *TestSuiteSequences) BeforeTest(suiteName, testName string) {
	atomic.AddUint64(&suite.env.FiredSequences, 1)
	suite.T().Log("Running one more test, tot ", suite.env.FiredSequences)
}

//Test_Sequences can be used to test a single run of the sequence test suite
func Test_Sequences(t *testing.T) {
	Env := setEnv(t)

	s := &TestSuiteSequences{
		env: Env,
	}
	suite.Run(t, s)
}

func setEnv(t *testing.T) *ZeroDowntimeEnv {
	Env := SetupZD()
	var err error
	Env.ExistingProject, err = testutils.CreateProject("projectzd", Env.ShipyardFile)
	assert.Nil(t, err)
	_, err = testutils.ExecuteCommand(fmt.Sprintf("keptn create service %s --project=%s", "myservice", Env.ExistingProject))
	assert.Nil(t, err)
	return Env
}

// to perform tests sequentially inside ZD
func Sequences(t *testing.T, env *ZeroDowntimeEnv) {
	var s *TestSuiteSequences
	env.Wg.Add(1)
	wgSequences := &sync.WaitGroup{}
	seqTicker := clock.New().Ticker(sequencesInterval)

Loop:
	for {
		select {
		case <-env.Ctx.Done():
			break Loop
		case <-seqTicker.C:
			s = &TestSuiteSequences{
				env: env,
			}
			wgSequences.Add(1)
			go func() {
				suite.Run(t, s)
				wgSequences.Done()
			}()

		}
	}
	wgSequences.Wait()
	env.Wg.Done()
}

func (suite *TestSuiteSequences) Test_Evaluation() {
	var finished *models.KeptnContextExtendedCE
	defer func() {
		if finished.Data == nil {
			atomic.AddUint64(&suite.env.FailedSequences, 1)
		} else {
			atomic.AddUint64(&suite.env.PassedSequences, 1)
		}
	}()

	suite.T().Log("deleting lighthouse configmap from previous test run")
	_, _ = testutils.ExecuteCommand(fmt.Sprintf("kubectl delete configmap -n %s lighthouse-config-%s", testutils.GetKeptnNameSpaceFromEnv(), suite.project))

	//// now let's add an SLI provider
	suite.T().Log("adding SLI provider")
	_, err := testutils.ExecuteCommand(fmt.Sprintf("kubectl create configmap -n %s lighthouse-config-%s --from-literal=sli-provider=my-sli-provider", testutils.GetKeptnNameSpaceFromEnv(), suite.project))
	suite.Nil(err)
	_, finished = testutils.PerformResourceServiceTest(suite.T(), suite.project, "myservice", true)

}

func (suite *TestSuiteSequences) Test_DeliveryFails() {
	suite.trigger("delivery", nil, false)
}

func (suite *TestSuiteSequences) Test_ExistingEvaluationFails() {
	suite.trigger("evaluation", nil, true)
}

func (suite *TestSuiteSequences) trigger(triggerType string, data keptn.EventProperties, existing bool) {
	project := suite.project
	if existing {
		project = suite.env.ExistingProject
	}

	suite.T().Logf("triggering sequence %s for project %s", triggerType, project)
	// trigger a delivery sequence
	keptnContext, err := testutils.TriggerSequence(project, "myservice", "hardening", triggerType, data)
	suite.Nil(err)
	suite.T().Logf("triggered sequence %s for project %s with context %s", triggerType, project, keptnContext)
	sequence := NewTriggeredSequence(keptnContext, project, triggerType)

	suite.checkSequence(sequence)
}

func (suite *TestSuiteSequences) checkSequence(sequence TriggeredSequence) {

	var sequenceFinishedEvent *models.KeptnContextExtendedCE
	stageSequenceName := fmt.Sprintf("%s.%s", "hardening", sequence.sequenceName)
	var err error

	suite.T().Logf("verifying completion of sequence %s with keptnContext %s in project %s", sequence.sequenceName, sequence.keptnContext, sequence.projectName)
	suite.Eventually(func() bool {
		sequenceFinishedEvent, err = testutils.GetLatestEventOfType(sequence.keptnContext, sequence.projectName, "hardening", v0_2_0.GetFinishedEventType(stageSequenceName))
		if sequenceFinishedEvent == nil || err != nil {
			return false
		}
		atomic.AddUint64(&suite.env.PassedSequences, 1)
		return true
	}, 1*time.Minute, 5*time.Second)

	if sequenceFinishedEvent == nil || err != nil {
		atomic.AddUint64(&suite.env.FailedSequences, 1)
		suite.T().Errorf("sequence %s with keptnContext %s in project %s has NOT been finished", sequence.sequenceName, sequence.keptnContext, sequence.projectName)

	} else {
		suite.T().Logf("sequence %s with keptnContext %s in project %s has been finished", sequence.sequenceName, sequence.keptnContext, sequence.projectName)
	}
}

func GetShipyard() (string, error) {
	shipyard := &v0_2_0.Shipyard{
		ApiVersion: "0.2.3",
		Kind:       "shipyard",
		Metadata:   v0_2_0.Metadata{},
		Spec: v0_2_0.ShipyardSpec{
			Stages: []v0_2_0.Stage{},
		},
	}

	stage := v0_2_0.Stage{
		Name: "hardening",
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

func (suite *TestSuiteSequences) gedId() string {
	atomic.AddUint64(&suite.env.Id, 1)
	return fmt.Sprintf("%d", suite.env.Id)
}

func PrintSequencesResults(t *testing.T, env *ZeroDowntimeEnv) {

	t.Log("-----------------------------------------------")
	t.Log("TOTAL SEQUENCES: ", env.FiredSequences)
	t.Log("TOTAL SUCCESS ", env.PassedSequences)
	t.Log("TOTAL FAILURES ", env.FailedSequences)
	t.Log("-----------------------------------------------")

}
