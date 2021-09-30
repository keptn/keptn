package go_tests

import (
	"fmt"
	"github.com/keptn/go-utils/pkg/api/models"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/shipyard-controller/common"
	"github.com/keptn/keptn/shipyard-controller/operations"
	"github.com/stretchr/testify/require"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"
)

const webhookShipyard = `--- 
apiVersion: "spec.keptn.sh/0.2.3"
kind: Shipyard
metadata:
  name: "shipyard-echo-service"
spec:
  stages:
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
        - name: "mysequence"
          tasks:
            - name: "mytask"`

const webhookYaml = `apiVersion: webhookconfig.keptn.sh/v1alpha1
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
        - "curl http://shipyard-controller:8080/v1/project{{.unknownKey}}"
    - type: "sh.keptn.event.unallowedtask.triggered"
      subscriptionID: ${unallowedtask-sub-id}
      sendFinished: true
      envFrom: 
        - name: "secretKey"
          secretRef:
            name: "my-webhook-k8s-secret"
            key: "my-key"
      requests:
        - "curl http://kubernetes.default.svc.cluster.local:443/v1"
    - type: "sh.keptn.event.mytask.triggered"
      subscriptionID: ${mytask-sub-id}
      sendFinished: true
      envFrom: 
        - name: "secretKey"
          secretRef:
            name: "my-webhook-k8s-secret"
            key: "my-key"
      requests:
        - "curl --header 'x-token: {{.env.secretKey}}' http://shipyard-controller:8080/v1/project/{{.data.project}}"
        - "curl http://shipyard-controller:8080/v1/project/{{.data.project}}/stage/{{.data.stage}}"`

const webhookWithDisabledFinishedEventsYaml = `apiVersion: webhookconfig.keptn.sh/v1alpha1
kind: WebhookConfig
metadata:
  name: webhook-configuration
spec:
  webhooks:
    - type: "sh.keptn.event.othertask.triggered"
      subscriptionID: ${othertask-sub-id}
      sendFinished: false
      envFrom: 
        - name: "secretKey"
          secretRef:
            name: "my-webhook-k8s-secret"
            key: "my-key"
      requests:
        - "curl http://shipyard-controller:8080/v1/project{{.unknownKey}}"
        - "curl http://shipyard-controller:8080/v1/project{{.unknownKey}}"
    - type: "sh.keptn.event.unallowedtask.triggered"
      subscriptionID: ${unallowedtask-sub-id}
      sendFinished: false
      envFrom: 
        - name: "secretKey"
          secretRef:
            name: "my-webhook-k8s-secret"
            key: "my-key"
      requests:
        - "curl http://shipyard-controller:8080/v1/project/{{.data.project}}/stage/{{.data.stage}}"
        - "curl http://kubernetes.default.svc.cluster.local:443/v1"
    - type: "sh.keptn.event.mytask.triggered"
      subscriptionID: ${mytask-sub-id}
      sendFinished: false
      envFrom: 
        - name: "secretKey"
          secretRef:
            name: "my-webhook-k8s-secret"
            key: "my-key"
      requests:
        - "curl --header 'x-token: {{.env.secretKey}}' http://shipyard-controller:8080/v1/project/{{.data.project}}"
        - "curl http://shipyard-controller:8080/v1/project/{{.data.project}}/stage/{{.data.stage}}"`

func Test_Webhook(t *testing.T) {
	projectName := "webhooks"
	serviceName := "myservice"
	stageName := "dev"
	sequencename := "mysequence"

	shipyardFilePath, err := CreateTmpShipyardFile(webhookShipyard)
	require.Nil(t, err)
	defer os.Remove(shipyardFilePath)

	t.Logf("creating project %s", projectName)
	err = CreateProject(projectName, shipyardFilePath, true)
	require.Nil(t, err)

	t.Logf("creating service %s", serviceName)
	output, err := ExecuteCommand(fmt.Sprintf("keptn create service %s --project=%s", serviceName, projectName))

	require.Nil(t, err)
	require.Contains(t, output, "created successfully")

	// create a secret that should be referenced in the webhook
	_, err = ApiPOSTRequest("/secrets/v1/secret", map[string]interface{}{
		"name": "my-webhook-k8s-secret",
		//		"scope": "keptn-default-scope",
		"data": map[string]string{
			"my-key": "my-value",
		},
	})
	require.Nil(t, err)

	// create subscriptions for the webhook-service
	taskTypes := []string{"mytask", "othertask", "unallowedtask", "unknowntask"}

	webhookYamlWithSubscriptionIDs := webhookYaml
	for _, taskType := range taskTypes {
		subscriptionID, err := CreateSubscription(t, "webhook-service", models.EventSubscription{
			Event: keptnv2.GetTriggeredEventType(taskType),
			Filter: models.EventSubscriptionFilter{
				Projects: []string{projectName},
			},
		})
		require.Nil(t, err)

		subscriptionPlaceholder := fmt.Sprintf("${%s-sub-id}", taskType)
		webhookYamlWithSubscriptionIDs = strings.Replace(webhookYamlWithSubscriptionIDs, subscriptionPlaceholder, subscriptionID, -1)
	}

	require.Nil(t, err)

	// wait some time to make sure the webhook service has pulled the updated subscription
	<-time.After(20 * time.Second) // sorry :(

	// now, let's add an webhook.yaml file to our service
	webhookFilePath, err := CreateTmpFile("webhook.yaml", webhookYamlWithSubscriptionIDs)
	require.Nil(t, err)
	defer os.Remove(webhookFilePath)

	t.Log("Adding webhook.yaml to our service")
	_, err = ExecuteCommand(fmt.Sprintf("keptn add-resource --project=%s --service=%s --resource=%s --resourceUri=webhook/webhook.yaml --all-stages", projectName, serviceName, webhookFilePath))

	require.Nil(t, err)

	triggerSequenceAndVerifyTaskFinishedEvent := func(sequencename, taskFinishedType string, verify func(t *testing.T, decodedEvent map[string]interface{})) {
		t.Logf("triggering sequence %s in stage %s", sequencename, stageName)
		keptnContextID, _ := TriggerSequence(projectName, serviceName, stageName, sequencename, nil)

		var taskFinishedEvent *models.KeptnContextExtendedCE
		require.Eventually(t, func() bool {
			taskFinishedEvent, err = GetLatestEventOfType(keptnContextID, projectName, stageName, keptnv2.GetFinishedEventType(taskFinishedType))
			if err != nil || taskFinishedEvent == nil {
				return false
			}
			return true
		}, 30*time.Second, 3*time.Second)

		require.NotNil(t, taskFinishedEvent)

		decodedEvent := map[string]interface{}{}

		err = keptnv2.EventDataAs(*taskFinishedEvent, &decodedEvent)
		require.Nil(t, err)

		verify(t, decodedEvent)
	}

	triggerSequenceAndVerifyTaskFinishedEvent(sequencename, "mytask", func(t *testing.T, decodedEvent map[string]interface{}) {
		require.NotNil(t, decodedEvent["mytask"])
	})

	// Now, trigger another sequence that tries to execute a webhook with a reference to an unknown variable - this should fail
	sequencename = "othersequence"

	triggerSequenceAndVerifyTaskFinishedEvent(sequencename, "othertask", func(t *testing.T, decodedEvent map[string]interface{}) {
		// check the result - this time it should be set to fail because an unknown Key was referenced in the webhook
		require.Equal(t, string(keptnv2.ResultFailed), decodedEvent["result"])
		require.Nil(t, decodedEvent["othertask"])
	})

	// Now, trigger another sequence that tries to execute a webhook with a call to the kubernetes API - this one should fail as well
	sequencename = "unallowedsequence"

	triggerSequenceAndVerifyTaskFinishedEvent(sequencename, "unallowedtask", func(t *testing.T, decodedEvent map[string]interface{}) {
		// check the result - this time it should be set to fail because an unknown Key was referenced in the webhook
		require.Equal(t, string(keptnv2.ResultFailed), decodedEvent["result"])
		require.Nil(t, decodedEvent["unallowedtask"])
	})

	// Now, trigger another sequence that contains a task for which we don't have a webhook configured - this one should fail as well
	sequencename = "sequencewithunknowntask"

	triggerSequenceAndVerifyTaskFinishedEvent(sequencename, "unknowntask", func(t *testing.T, decodedEvent map[string]interface{}) {
		// check the result - this time it should be set to fail because an unknown Key was referenced in the webhook
		require.Equal(t, string(keptnv2.ResultFailed), decodedEvent["result"])
		require.Nil(t, decodedEvent["unknowntask"])
	})
}

func Test_WebhookWithDisabledFinishedEvents(t *testing.T) {
	projectName := "webhooks-no-finish"
	serviceName := "myservice"
	stageName := "dev"
	sequencename := "mysequence"

	shipyardFilePath, err := CreateTmpShipyardFile(webhookShipyard)
	require.Nil(t, err)
	defer os.Remove(shipyardFilePath)

	t.Logf("creating project %s", projectName)
	err = CreateProject(projectName, shipyardFilePath, true)
	require.Nil(t, err)

	t.Logf("creating service %s", serviceName)
	output, err := ExecuteCommand(fmt.Sprintf("keptn create service %s --project=%s", serviceName, projectName))

	require.Nil(t, err)
	require.Contains(t, output, "created successfully")

	// create a secret that should be referenced in the webhook
	_, err = ApiPOSTRequest("/secrets/v1/secret", map[string]interface{}{
		"name": "my-webhook-k8s-secret",
		//		"scope": "keptn-default-scope",
		"data": map[string]string{
			"my-key": "my-value",
		},
	})
	require.Nil(t, err)

	// create subscriptions for the webhook-service
	taskTypes := []string{"mytask", "othertask", "unallowedtask", "unknowntask"}

	webhookYamlWithSubscriptionIDs := webhookWithDisabledFinishedEventsYaml
	for _, taskType := range taskTypes {
		subscriptionID, err := CreateSubscription(t, "webhook-service", models.EventSubscription{
			Event: keptnv2.GetTriggeredEventType(taskType),
			Filter: models.EventSubscriptionFilter{
				Projects: []string{projectName},
			},
		})
		require.Nil(t, err)

		subscriptionPlaceholder := fmt.Sprintf("${%s-sub-id}", taskType)
		webhookYamlWithSubscriptionIDs = strings.Replace(webhookYamlWithSubscriptionIDs, subscriptionPlaceholder, subscriptionID, -1)
	}

	require.Nil(t, err)

	// wait some time to make sure the webhook service has pulled the updated subscription
	<-time.After(20 * time.Second) // sorry :(

	// now, let's add an webhook.yaml file to our service
	webhookFilePath, err := CreateTmpFile("webhook.yaml", webhookYamlWithSubscriptionIDs)
	require.Nil(t, err)
	defer os.Remove(webhookFilePath)

	t.Log("Adding webhook.yaml to our service")
	_, err = ExecuteCommand(fmt.Sprintf("keptn add-resource --project=%s --service=%s --resource=%s --resourceUri=webhook/webhook.yaml --all-stages", projectName, serviceName, webhookFilePath))

	require.Nil(t, err)

	triggerSequenceAndVerifyStartedEvents := func(sequencename, taskName string, nrExpected int) string {
		t.Logf("triggering sequence %s in stage %s", sequencename, stageName)
		keptnContextID, _ := TriggerSequence(projectName, serviceName, stageName, sequencename, nil)

		var taskStartedEvents []*models.KeptnContextExtendedCE
		require.Eventually(t, func() bool {
			taskStartedEvents, err = GetEventsOfType(keptnContextID, projectName, stageName, keptnv2.GetStartedEventType(taskName))
			if err != nil || taskStartedEvents == nil || len(taskStartedEvents) != nrExpected {
				return false
			}
			return true
		}, 30*time.Second, 3*time.Second)
		return keptnContextID
	}

	// verify that two .started events have been sent
	keptnContextID := triggerSequenceAndVerifyStartedEvents(sequencename, "mytask", 2)

	<-time.After(5 * time.Second)

	// verify that no .finished event has been sent for 'mytask'
	taskFinishedEvent, err := GetLatestEventOfType(keptnContextID, projectName, stageName, keptnv2.GetFinishedEventType("mytask"))

	require.Nil(t, taskFinishedEvent)

	t.Log("verified desired state, aborting sequence")
	resp, err := ApiPOSTRequest(fmt.Sprintf("/controlPlane/v1/sequence/%s/%s/control", projectName, keptnContextID), operations.SequenceControlCommand{
		State: common.AbortSequence,
		Stage: "",
	})
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.Response().StatusCode)

	// Now, trigger another sequence that tries to execute a webhook with a reference to an unknown variable - this should fail
	sequencename = "othersequence"

	keptnContextID = triggerSequenceAndVerifyStartedEvents(sequencename, "othertask", 2)

	// verify that we have received two .finished events with the status set to fail
	var taskFinishedEvents []*models.KeptnContextExtendedCE
	require.Eventually(t, func() bool {
		taskFinishedEvents, err = GetEventsOfType(keptnContextID, projectName, stageName, keptnv2.GetFinishedEventType("othertask"))
		if err != nil {
			t.Logf("got error: %s. will try again in a few seconds", err.Error())
			return false
		} else if taskFinishedEvents == nil {
			t.Log("did not receive any .finished events")
			return false
		} else if len(taskFinishedEvents) != 2 {
			t.Logf("received %d .finished events, but expected 2", len(taskFinishedEvents))
			return false
		}
		return true
	}, 30*time.Second, 2*time.Second)

	for _, event := range taskFinishedEvents {
		decodedEvent := map[string]interface{}{}

		err = keptnv2.EventDataAs(*event, &decodedEvent)
		require.Nil(t, err)
		require.Equal(t, string(keptnv2.ResultFailed), decodedEvent["result"])
		require.NotEmpty(t, string(keptnv2.ResultFailed), decodedEvent["message"])
	}

	// Now, trigger another sequence that tries to execute a webhook with a call to the kubernetes API - this one should fail as well
	sequencename = "unallowedsequence"

	keptnContextID = triggerSequenceAndVerifyStartedEvents(sequencename, "unallowedtask", 2)

	<-time.After(5 * time.Second)
	// verify that we have received one .finished events with the status set to fail
	require.Eventually(t, func() bool {
		taskFinishedEvents, err = GetEventsOfType(keptnContextID, projectName, stageName, keptnv2.GetFinishedEventType("unallowedtask"))
		if err != nil {
			t.Logf("got error: %s. will try again in a few seconds", err.Error())
			return false
		} else if taskFinishedEvents == nil {
			t.Log("did not receive any .finished events")
			return false
		} else if len(taskFinishedEvents) != 1 {
			t.Logf("received %d .finished events, but expected 1", len(taskFinishedEvents))
			return false
		}
		return true
	}, 10*time.Second, 2*time.Second)

	for _, event := range taskFinishedEvents {
		decodedEvent := map[string]interface{}{}

		err = keptnv2.EventDataAs(*event, &decodedEvent)
		require.Nil(t, err)
		require.Equal(t, string(keptnv2.ResultFailed), decodedEvent["result"])
		require.NotEmpty(t, string(keptnv2.ResultFailed), decodedEvent["message"])
	}
}
