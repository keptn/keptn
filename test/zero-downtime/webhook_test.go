package zero_downtime

import (
	"fmt"
	"github.com/benbjohnson/clock"
	testutils "github.com/keptn/keptn/test/go-tests"
	"github.com/stretchr/testify/require"
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
	t.Logf("started Webhook tests")
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

const webhookShipyard = `--- 
apiVersion: "spec.keptn.sh/0.2.3"
kind: Shipyard
metadata:
  name: "shipyard-echo-service"
spec:
  stages:
    - name: "otherstage"
    - name: "dev"
      sequences:
        - name: "othersequence"
          tasks:
            - name: "othertask"
        - name: "sequencewithunknowntask"
          tasks:
            - name: "unknowntask"
        - name: "unallowedsequence"
          tasks:
            - name: "unallowedtask"
        - name: "failedsequence"
          tasks:
            - name: "failedtask"
        - name: "loopbacksequence"
          tasks:
            - name: "loopback"
        - name: "loopbacksequence2"
          tasks:
            - name: "loopback2"
        - name: "loopbacksequence3"
          tasks:
            - name: "loopback3"
        - name: "mysequence"
          tasks:
            - name: "mytask"`

// 1 Pass 7 Fails
func (suite *TestSuiteWebhook) Test_Webhook() {
	projectName := "webhooks" + suite.env.gedId()
	serviceName := "myservice"

	//test considered failed by default so that we can use require
	suite.env.failSequence()

	shipyardFilePath, err := testutils.CreateTmpShipyardFile(webhookShipyard)
	require.Nil(suite.T(), err)

	suite.T().Logf("creating project %s", projectName)
	projectName, err = testutils.CreateProject(projectName, shipyardFilePath)
	require.Nil(suite.T(), err)

	suite.T().Logf("creating service %s", serviceName)
	output, err := testutils.ExecuteCommand(fmt.Sprintf("keptn create service %s --project=%s", serviceName, projectName))

	require.Nil(suite.T(), err)
	require.Contains(suite.T(), output, "created successfully")

	// create a secret that should be referenced in the webhook yaml
	_, _ = testutils.ApiPOSTRequest("/secrets/v1/secret", map[string]interface{}{
		"name":  "my-webhook-k8s-secret",
		"scope": "keptn-webhook-service",
		"data": map[string]string{
			"my-key": "my-value",
		},
	}, 3)

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
