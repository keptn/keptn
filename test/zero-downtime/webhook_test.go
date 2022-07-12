package zero_downtime

import (
	"fmt"
	"github.com/benbjohnson/clock"
	"github.com/keptn/go-utils/pkg/api/models"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	testutils "github.com/keptn/keptn/test/go-tests"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
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

const webhookConfig = `apiVersion: webhookconfig.keptn.sh/v1beta1
kind: WebhookConfig
metadata:
  name: webhook-configuration
spec:
  webhooks:
    - type: "sh.keptn.event.othertask.triggered"
      subscriptionID: ${othertask-sub-id}
      sendFinished: true
      envFrom: 
        - name: "secretKey"
          secretRef:
            name: "my-webhook-k8s-secret"
            key: "my-key"
      requests:
        - url: http://shipyard-controller:8080/v1/project{{.unknownKey}}
          method: GET
    - type: "sh.keptn.event.failedtask.triggered"
      subscriptionID: ${failedtask-sub-id}
      sendFinished: true
      requests:
        - url: http://shipyard-controller:8080/v1/some-unknown-api
          method: GET
    - type: "sh.keptn.event.unallowedtask.triggered"
      subscriptionID: ${unallowedtask-sub-id}
      sendFinished: true
      envFrom: 
        - name: "secretKey"
          secretRef:
            name: "my-webhook-k8s-secret"
            key: "my-key"
      requests:
        - url: http://kubernetes.default.svc.cluster.local:443/v1
          method: GET
    - type: "sh.keptn.event.loopback.triggered"
      subscriptionID: ${loopback-sub-id}
      sendFinished: true
      requests:
        - url: http://localhost:8080
          method: GET
    - type: "sh.keptn.event.loopback2.triggered"
      subscriptionID: ${loopback2-sub-id}
      sendFinished: true
      requests:
        - url: http://127.0.0.1:8080
          method: GET
    - type: "sh.keptn.event.loopback3.triggered"
      subscriptionID: ${loopback3-sub-id}
      sendFinished: true
      requests:
        - url: http://[::1]:8080
          method: GET
    - type: "sh.keptn.event.mytask.finished"
      subscriptionID: ${mytask-finished-sub-id}
      sendFinished: true
      envFrom: 
        - name: "secretKey"
          secretRef:
            name: "my-webhook-k8s-secret"
            key: "my-key"
      requests:
        - url: http://shipyard-controller:8080/v1/some-unknown-api
          method: GET
    - type: "sh.keptn.event.mytask.triggered"
      subscriptionID: ${mytask-sub-id}
      sendFinished: true
      envFrom: 
        - name: "secretKey"
          secretRef:
            name: "my-webhook-k8s-secret"
            key: "my-key"
      requests:
        - url: http://keptn.sh
          method: GET
          headers:
            - key: x-token
              value: "{{.env.secretKey}}"
        - url: http://keptn.sh
          method: GET`
const webhookConfigMap = "keptn-webhook-config"

// 1 Pass 7 Fails
func (suite *TestSuiteWebhook) Test_Webhook() {
	suite.env.failSequence()
	t := suite.T()
	oldConfig, err := testutils.GetFromConfigMap(testutils.GetKeptnNameSpaceFromEnv(), webhookConfigMap, func(data map[string]string) string {
		return data["denyList"]
	})
	require.Nil(t, err)
	testutils.PutConfigMapDataVal(testutils.GetKeptnNameSpaceFromEnv(), webhookConfigMap, "denyList", "kubernetes")
	defer testutils.PutConfigMapDataVal(testutils.GetKeptnNameSpaceFromEnv(), webhookConfigMap, "denyList", oldConfig)

	projectName := "webhooks" + suite.env.gedId()
	serviceName := "myservice"
	projectName, shipyardFilePath := testutils.CreateWebhookProject(t, projectName, serviceName)
	defer testutils.DeleteFile(t, shipyardFilePath)
	stageName := "dev"
	sequencename := "mysequence"
	taskName := "mytask"

	// create subscriptions for the webhook-service
	webhookYamlWithSubscriptionIDs := webhookConfig
	subscriptionID, err := testutils.CreateSubscription(t, "webhook-service", models.EventSubscription{
		Event: keptnv2.GetTriggeredEventType(taskName),
		Filter: models.EventSubscriptionFilter{
			Projects: []string{projectName},
			Stages:   []string{stageName},
		},
	})
	require.Nil(t, err)

	webhookYamlWithSubscriptionIDs = strings.Replace(webhookYamlWithSubscriptionIDs, "${mytask-sub-id}", subscriptionID, -1)

	// create a second subscription that overlaps with the previously created one
	subscriptionID2, err := testutils.CreateSubscription(t, "webhook-service", models.EventSubscription{
		Event: keptnv2.GetTriggeredEventType(taskName),
		Filter: models.EventSubscriptionFilter{
			Projects: []string{projectName},
			Stages:   []string{stageName, "otherstage"},
		},
	})
	require.Nil(t, err)

	webhookYamlWithSubscriptionIDs = strings.Replace(webhookYamlWithSubscriptionIDs, "${mytask-sub-2-id}", subscriptionID2, -1)

	// wait some time to make sure the webhook service has pulled the updated subscription
	<-time.After(20 * time.Second) // sorry :(

	// now, let's add a webhook.yaml file to our service
	webhookFilePath, err := testutils.CreateTmpFile("webhook.yaml", webhookYamlWithSubscriptionIDs)
	require.Nil(t, err)
	defer testutils.DeleteFile(t, webhookFilePath)

	t.Log("Adding webhook.yaml to our service")
	_, err = testutils.ExecuteCommand(fmt.Sprintf("keptn add-resource --project=%s --service=%s --resource=%s --resourceUri=webhook/webhook.yaml --all-stages", projectName, serviceName, webhookFilePath))

	require.Nil(t, err)

	t.Logf("triggering sequence %s in stage %s", sequencename, stageName)
	keptnContextID, _ := testutils.TriggerSequence(projectName, serviceName, stageName, sequencename, nil)

	var taskFinishedEvent []*models.KeptnContextExtendedCE
	require.Eventually(t, func() bool {
		taskFinishedEvent, err = testutils.GetEventsOfType(keptnContextID, projectName, stageName, keptnv2.GetFinishedEventType(taskName))
		if err != nil || taskFinishedEvent == nil || len(taskFinishedEvent) != 2 {
			return false
		}
		return true
	}, 30*time.Second, 3*time.Second)
	//if test returns then it's passed
	suite.env.passFailedSequence()
}
