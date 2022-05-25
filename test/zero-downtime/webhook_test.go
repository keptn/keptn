package zero_downtime

import (
	"github.com/benbjohnson/clock"
	testutils "github.com/keptn/keptn/test/go-tests"
	"github.com/stretchr/testify/suite"
	"os"
	"sync"
	"sync/atomic"
	"testing"
)

type TestSuiteWebhook struct {
	suite.Suite
	env *ZeroDowntimeEnv
}

func (suite *TestSuiteWebhook) SetupSuite() {
	suite.T().Log("Starting test for webhook")

}

func (suite *TestSuiteWebhook) BeforeTest(suiteName, testName string) {
	atomic.AddUint64(&suite.env.FiredSequences, 1)
	suite.T().Log("Running one more test, tot ", suite.env.FiredSequences)
}

// Test_Webhook can be used to test a single run of the test suite
func Test_Webhook(t *testing.T) {

	s := &TestSuiteWebhook{
		env: SetupZD(),
	}
	suite.Run(t, s)
}

// Webhook performs tests sequentially inside the zerodowntime suite
func Webhook(t *testing.T, env *ZeroDowntimeEnv) {
	var s *TestSuiteWebhook
	wgSequences := &sync.WaitGroup{}
	seqTicker := clock.New().Ticker(env.SequencesInterval)
Loop:
	for {
		select {
		case <-env.quit:
			break Loop
		case <-seqTicker.C:
			s = &TestSuiteWebhook{
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

}

// 1 Pass 7 Fails
func (suite *TestSuiteWebhook) Test_Webhook() {
	projectName := "webhooks" + suite.env.gedId()
	serviceName := "myservice"

	//test considered failed by default so that we can use require
	suite.env.failSequence()

	projectName, shipyardFilePath := testutils.CreateWebhookProject(suite.T(), projectName, serviceName)
	defer func() {
		err := os.Remove(shipyardFilePath)
		if err != nil {
			suite.T().Logf("Could not delete tmp file: %s", err.Error())
		}
	}()
	testutils.Test_Webhook(suite.T(), testutils.WebhookYamlBeta, projectName, serviceName)

	//if test returns then it's passed
	suite.env.passFailedSequence()
}
