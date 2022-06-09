package zero_downtime

import (
	"fmt"
	"github.com/benbjohnson/clock"
	"github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/go-utils/pkg/lib/keptn"
	"github.com/keptn/go-utils/pkg/lib/v0_2_0"
	testutils "github.com/keptn/keptn/test/go-tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
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
	suite.Assert().Contains(suite.project, "zd-sequence")
}

func (suite *TestSuiteSequences) createNew() {
	var err error
	projectName := "zd-sequence" + suite.env.gedId()
	suite.T().Logf("Creating project with id %s ", projectName)
	suite.project, err = testutils.CreateProject(projectName, suite.env.ShipyardFile)
	suite.Nil(err)
	output, err := testutils.ExecuteCommand(fmt.Sprintf("keptn create service %s --project=%s", "myservice", suite.project))
	suite.Require().Nil(err)
	suite.Require().Contains(output, "created successfully")
	suite.T().Logf("Starting test for project %s ", suite.project)
}

func (suite *TestSuiteSequences) BeforeTest(suiteName, testName string) {
	atomic.AddUint64(&suite.env.FiredSequences, 1)
	suite.T().Log("Current time ", time.Now())
	suite.T().Log("Running one more test, total tests: ", suite.env.FiredSequences)
}

//Test_Sequences can be used to test a single run of the sequence test suite
func Test_Sequences(t *testing.T) {
	Env := setSequencesEnv(t)

	s := &TestSuiteSequences{
		env: Env,
	}
	suite.Run(t, s)
}

func setSequencesEnv(t *testing.T) *ZeroDowntimeEnv {
	Env := SetupZD()
	var err error
	Env.ExistingProject, err = testutils.CreateProject("projectzd", Env.ShipyardFile)
	assert.Nil(t, err)
	_, err = testutils.ExecuteCommand(fmt.Sprintf("keptn create service %s --project=%s", "myservice", Env.ExistingProject))
	assert.Nil(t, err)
	return Env
}

// Sequences is used to perform tests sequentially inside the zerodowntime suite
func Sequences(t *testing.T, env *ZeroDowntimeEnv) {
	wgSequences := &sync.WaitGroup{}
	seqTicker := clock.New().Ticker(env.SequencesInterval)
	t.Logf("started Sequence tests")
Loop:
	for {
		select {
		case <-env.quit:
			break Loop
		case <-seqTicker.C:
			wgSequences.Add(1)
			go func() {
				suite.Run(t, &TestSuiteSequences{
					env: env,
				})
				wgSequences.Done()
			}()

		}
	}
	wgSequences.Wait()

}

// Pass : evaluation 100
func (suite *TestSuiteSequences) Test_Evaluation() {
	var finished *models.KeptnContextExtendedCE
	suite.env.failSequence()
	suite.T().Log("deleting lighthouse configmap from previous test run")
	_, _ = testutils.ExecuteCommand(fmt.Sprintf("kubectl delete configmap -n %s lighthouse-config-%s", testutils.GetKeptnNameSpaceFromEnv(), suite.project))

	//// now let's add an SLI provider
	suite.T().Log("adding SLI provider")
	_, err := testutils.ExecuteCommand(fmt.Sprintf("kubectl create configmap -n %s lighthouse-config-%s --from-literal=sli-provider=my-sli-provider", testutils.GetKeptnNameSpaceFromEnv(), suite.project))
	suite.Require().Nil(err)
	_, finished = testutils.PerformResourceServiceTest(suite.T(), suite.project, "myservice", true)
	if finished.Data != nil {
		suite.env.passFailedSequence()
	}
}

// Fails : no task sequence with name delivery found in stage hardening
func (suite *TestSuiteSequences) Test_DeliveryFails() {
	suite.trigger("delivery", nil, false)
}

// Fails : Evaluation 0 event does not contain evaluation timeframe
func (suite *TestSuiteSequences) Test_ExistingEvaluationFails() {
	suite.trigger("evaluation", nil, false)
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
		suite.env.passSequence()
		return true
	}, 1*time.Minute, 5*time.Second)

	if sequenceFinishedEvent == nil || err != nil {
		suite.env.failSequence()
		suite.T().Errorf("sequence %s with keptnContext %s in project %s has NOT been finished", sequence.sequenceName, sequence.keptnContext, sequence.projectName)

	} else {
		suite.T().Logf("sequence %s with keptnContext %s in project %s has been finished", sequence.sequenceName, sequence.keptnContext, sequence.projectName)
	}
}
