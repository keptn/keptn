package zerodowntime

import (
	"fmt"
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

//Test_Sequences can be used to test a single run of the sequence test suite
func Test_Webhook(t *testing.T) {

	s := &TestSuiteWebhook{
		env: SetupZD(),
	}
	suite.Run(t, s)
}

//Zero Downtime Test
func TestWebhookZD(t *testing.T) {

	env := SetupZD()
	t.Run("Rolling Upgrade", func(t2 *testing.T) {
		t2.Parallel()
		RollingUpgrade(t2, env)
	})

	t.Run("API", func(t2 *testing.T) {
		t2.Parallel()
		APIs(t2, env)
	})

	t.Run("Webhook", func(t2 *testing.T) {
		t2.Parallel()
		Webhook(t2, env)
	})

}

//This performs tests sequentially inside ZD
func Webhook(t *testing.T, env *ZeroDowntimeEnv) {
	var s *TestSuiteWebhook
	env.Wg.Add(1)
	wgSequences := &sync.WaitGroup{}
	seqTicker := clock.New().Ticker(sequencesInterval)

Loop:
	for {
		select {
		case <-env.Ctx.Done():
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
	env.Wg.Done()

}

func (suite *TestSuiteWebhook) Test_Webhook() {
	projectName := "webhooks" + suite.gedId()
	serviceName := "myservice"
	projectName, shipyardFilePath := testutils.CreateWebhookProject(suite.T(), projectName, serviceName)
	defer func() {
		err := os.Remove(shipyardFilePath)
		if err != nil {
			suite.T().Logf("Could not delete tmp file: %s", err.Error())
		}
	}()
	testutils.Test_Webhook(suite.T(), testutils.WebhookYamlBeta, projectName, serviceName)
}

func (suite *TestSuiteWebhook) gedId() string {
	atomic.AddUint64(&suite.env.Id, 1)
	return fmt.Sprintf("%d", suite.env.Id)
}
