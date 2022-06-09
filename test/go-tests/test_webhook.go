package go_tests

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/keptn/go-utils/pkg/api/models"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/stretchr/testify/require"
)

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

const webhookYamlAlpha = `apiVersion: webhookconfig.keptn.sh/v1alpha1
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
    - type: "sh.keptn.event.failedtask.triggered"
      subscriptionID: ${failedtask-sub-id}
      sendFinished: true
      requests:
        - "curl http://shipyard-controller:8080/v1/some-unknown-api"
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
    - type: "sh.keptn.event.loopback.triggered"
      subscriptionID: ${loopback-sub-id}
      sendFinished: true
      requests:
        - "curl http://localhost:8080"
    - type: "sh.keptn.event.loopback2.triggered"
      subscriptionID: ${loopback2-sub-id}
      sendFinished: true
      requests:
        - "curl http://127.0.0.1:8080"
    - type: "sh.keptn.event.loopback3.triggered"
      subscriptionID: ${loopback3-sub-id}
      sendFinished: true
      requests:
        - "curl http://[::1]:8080"
    - type: "sh.keptn.event.mytask.finished"
      subscriptionID: ${mytask-finished-sub-id}
      sendFinished: true
      envFrom: 
        - name: "secretKey"
          secretRef:
            name: "my-webhook-k8s-secret"
            key: "my-key"
      requests:
        - "curl http://shipyard-controller:8080/v1/some-unknown-api"
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

const WebhookYamlBeta = `apiVersion: webhookconfig.keptn.sh/v1beta1
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

const webhookWithDisabledFinishedEventsYamlAlpha = `apiVersion: webhookconfig.keptn.sh/v1alpha1
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

const webhookWithDisabledFinishedEventsYamlBeta = `apiVersion: webhookconfig.keptn.sh/v1beta1
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
        - url: http://keptn.sh{{.unknownKey}}
          method: GET
        - url: http://keptn.sh{{.unknownKey}}
          method: GET
    - type: "sh.keptn.event.unallowedtask.triggered"
      subscriptionID: ${unallowedtask-sub-id}
      sendFinished: false
      envFrom: 
        - name: "secretKey"
          secretRef:
            name: "my-webhook-k8s-secret"
            key: "my-key"
      requests:
        - url: http://keptn.sh
          method: GET
        - url: http://kubernetes.default.svc.cluster.local:443/v1
          method: GET
    - type: "sh.keptn.event.mytask.triggered"
      subscriptionID: ${mytask-sub-id}
      sendFinished: false
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

const webhookWithDisabledStartedEventsYamlAlpha = `apiVersion: webhookconfig.keptn.sh/v1alpha1
kind: WebhookConfig
metadata:
  name: webhook-configuration
spec:
  webhooks:
    - type: "sh.keptn.event.mytask.triggered"
      subscriptionID: ${mytask-sub-id}
      sendStarted: false
      sendFinished: false
      envFrom: 
        - name: "secretKey"
          secretRef:
            name: "my-webhook-k8s-secret"
            key: "my-key"
      requests:
        - "curl --header 'x-token: {{.env.secretKey}}' http://shipyard-controller:8080/v1/project/{{.data.project}}"`

const webhookWithDisabledStartedEventsYamlBeta = `apiVersion: webhookconfig.keptn.sh/v1beta1
kind: WebhookConfig
metadata:
  name: webhook-configuration
spec:
  webhooks:
    - type: "sh.keptn.event.mytask.triggered"
      subscriptionID: ${mytask-sub-id}
      sendStarted: false
      sendFinished: false
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
            value: "{{.env.secretKey}}"`

const webhookWithOverlappingSubscriptionsYamlAlpha = `apiVersion: webhookconfig.keptn.sh/v1alpha1
kind: WebhookConfig
metadata:
  name: webhook-configuration
spec:
  webhooks:
    - type: "sh.keptn.event.mytask.triggered"
      subscriptionID: ${mytask-sub-id}
      sendFinished: true
      requests:
        - "curl http://shipyard-controller:8080/v1/project/{{.data.project}}/stage/{{.data.stage}}"
    - type: "sh.keptn.event.mytask.triggered"
      subscriptionID: ${mytask-sub-2-id}
      sendFinished: true
      envFrom: 
        - name: "secretKey"
          secretRef:
            name: "my-webhook-k8s-secret"
            key: "my-key"
      requests:
        - "curl --header 'x-token: {{.env.secretKey}}' http://shipyard-controller:8080/v1/project/{{.data.project}}"
        - "curl http://shipyard-controller:8080/v1/project/{{.data.project}}/stage/{{.data.stage}}"`

const webhookWithOverlappingSubscriptionsYamlBeta = `apiVersion: webhookconfig.keptn.sh/v1beta1
kind: WebhookConfig
metadata:
  name: webhook-configuration
spec:
  webhooks:
    - type: "sh.keptn.event.mytask.triggered"
      subscriptionID: ${mytask-sub-id}
      sendFinished: true
      requests:
        - url: http://shipyard-controller:8080/v1/project/{{.data.project}}/stage/{{.data.stage}}
          method: GET
    - type: "sh.keptn.event.mytask.triggered"
      subscriptionID: ${mytask-sub-2-id}
      sendFinished: true
      envFrom: 
        - name: "secretKey"
          secretRef:
            name: "my-webhook-k8s-secret"
            key: "my-key"
      requests:
        - url: http://shipyard-controller:8080/v1/project/{{.data.project}}
          method: GET
          headers:
            - key: x-token
              value: "{{.env.secretKey}}"
        - url: http://shipyard-controller:8080/v1/project/{{.data.project}}/stage/{{.data.stage}}
          method: GET`

const failwebhookSimpleYamlBeta = `apiVersion: webhookconfig.keptn.sh/v1beta1
kind: WebhookConfig
metadata:
  name: webhook-configuration
spec:
  webhooks:
    - type: "sh.keptn.event.mytask.triggered"
      subscriptionID: ${mytask-sub-id}
      sendStarted: true
      sendFinished: true
      requests:
        - url: http://shipyard-controller:8080/v1/project/{{.data.project}}/stage/{{.data.stage}}
          method: GET`

const webhookSimpleYamlAlpha = `apiVersion: webhookconfig.keptn.sh/v1alpha1
kind: WebhookConfig
metadata:
  name: webhook-configuration
spec:
  webhooks:
    - type: "sh.keptn.event.mytask.triggered"
      subscriptionID: ${mytask-sub-id}
      sendStarted: true
      sendFinished: true
      requests:
        - "curl http://keptn.sh"`

const webhookSimpleYamlBeta = `apiVersion: webhookconfig.keptn.sh/v1beta1
kind: WebhookConfig
metadata:
  name: webhook-configuration
spec:
  webhooks:
    - type: "sh.keptn.event.mytask.triggered"
      subscriptionID: ${mytask-sub-id}
      sendStarted: true
      sendFinished: true
      requests:
        - url: http://keptn.sh
          method: GET`

const webhookSimpleYamlBetaAPI = `apiVersion: webhookconfig.keptn.sh/v1beta1
kind: WebhookConfig
metadata:
  name: webhook-configuration
spec:
  webhooks:
    - type: "sh.keptn.event.mytask.triggered"
      subscriptionID: ${mytask-sub-id}
      sendStarted: true
      sendFinished: true
      envFrom: 
        - name: "tokensecretKey"
          secretRef:
            name: "my-webhook-k8s-secret-token"
            key: "x-token"
      requests:
        - url: http://shipyard-controller:8080/v1/project
          method: GET
          headers:
            - key: x-token
              value: "{{.env.tokensecretKey}}"`

const WebhookConfigMap = "keptn-webhook-config"

func CreateWebhookProject(t *testing.T, projectName, serviceName string) (string, string) {

	shipyardFilePath, err := CreateTmpShipyardFile(webhookShipyard)
	require.Nil(t, err)

	t.Logf("creating project %s", projectName)
	projectName, err = CreateProject(projectName, shipyardFilePath)
	require.Nil(t, err)

	t.Logf("creating service %s", serviceName)
	output, err := ExecuteCommand(fmt.Sprintf("keptn create service %s --project=%s", serviceName, projectName))

	require.Nil(t, err)
	require.Contains(t, output, "created successfully")

	// create a secret that should be referenced in the webhook yaml
	_, err = ApiPOSTRequest("/secrets/v1/secret", map[string]interface{}{
		"name":  "my-webhook-k8s-secret",
		"scope": "keptn-webhook-service",
		"data": map[string]string{
			"my-key": "my-value",
		},
	}, 3)
	require.Nil(t, err)
	return projectName, shipyardFilePath
}

func Test_Webhook_Beta(t *testing.T) {
	projectName := "webhooks-b"
	serviceName := "myservice"
	projectName, shipyardFilePath := CreateWebhookProject(t, projectName, serviceName)
	defer deleteFile(t, shipyardFilePath)
	Test_Webhook(t, WebhookYamlBeta, projectName, serviceName)
}

func Test_Webhook_Alpha(t *testing.T) {
	projectName := "webhooks-a"
	serviceName := "myservice"
	projectName, shipyardFilePath := CreateWebhookProject(t, projectName, serviceName)
	defer deleteFile(t, shipyardFilePath)
	Test_Webhook(t, webhookYamlAlpha, projectName, serviceName)
}

func Test_Webhook_Beta_API(t *testing.T) {
	projectName := "webhooks-b-api"
	oldConfig, err := GetFromConfigMap(GetKeptnNameSpaceFromEnv(), WebhookConfigMap, func(data map[string]string) string {
		return data["denyList"]
	})
	require.Nil(t, err)

	//temporary enabling all communication
	PutConfigMapDataVal(GetKeptnNameSpaceFromEnv(), WebhookConfigMap, "denyList", "kubernetes")
	defer PutConfigMapDataVal(GetKeptnNameSpaceFromEnv(), WebhookConfigMap, "denyList", oldConfig)

	api, err := NewAPICaller()
	require.Nil(t, err)

	// create a secret that should be referenced in the webhook yaml
	_, err = ApiPOSTRequest("/secrets/v1/secret", map[string]interface{}{
		"name":  "my-webhook-k8s-secret-token",
		"scope": "keptn-webhook-service",
		"data": map[string]string{
			"x-token": api.token,
		},
	}, 3)
	require.Nil(t, err)

	Test_WebhookConfigAtStageLevel(t, webhookSimpleYamlBetaAPI, projectName)
}

func Test_Webhook(t *testing.T, webhookYaml string, projectName, serviceName string) {

	stageName := "dev"
	sequencename := "mysequence"

	// create subscriptions for the webhook-service
	taskTypes := []string{"mytask", "mytask-finished", "othertask", "unallowedtask", "unknowntask", "failedtask", "loopback", "loopback2", "loopback3"}

	webhookYamlWithSubscriptionIDs := webhookYaml
	webhookYamlWithSubscriptionIDs = getWebhookYamlWithSubscriptionIDs(t, taskTypes, projectName, webhookYamlWithSubscriptionIDs)

	// wait some time to make sure the webhook service has pulled the updated subscription
	<-time.After(20 * time.Second) // sorry :(

	// now, let's add an webhook.yaml file to our service
	webhookFilePath, err := CreateTmpFile("webhook.yaml", webhookYamlWithSubscriptionIDs)
	require.Nil(t, err)
	defer func() {
		err := os.Remove(webhookFilePath)
		if err != nil {
			t.Logf("Could not delete tmp file: %s", err.Error())
		}

	}()

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
		}, 60*time.Second, 5*time.Second)

		require.NotNil(t, taskFinishedEvent)

		decodedEvent := map[string]interface{}{}

		err = keptnv2.EventDataAs(*taskFinishedEvent, &decodedEvent)
		require.Nil(t, err)

		verify(t, decodedEvent)

		// verify that no <task>.finished.finished event is sent
		finishedFinishedEvent, err := GetLatestEventOfType(keptnContextID, projectName, stageName, keptnv2.GetFinishedEventType("mytask.finished"))
		require.Nil(t, err)
		require.Nil(t, finishedFinishedEvent)
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

	// Now, trigger another sequence that tries to execute a webhook with a call to the localhost - this one should fail as well
	sequencename = "loopbacksequence"

	triggerSequenceAndVerifyTaskFinishedEvent(sequencename, "loopback", func(t *testing.T, decodedEvent map[string]interface{}) {
		// check the result - this time it should be set to fail because an unknown Key was referenced in the webhook
		require.Equal(t, string(keptnv2.ResultFailed), decodedEvent["result"])
		require.Nil(t, decodedEvent["loopback"])
	})

	// Now, trigger another sequence that tries to execute a webhook with a call to the 127.0.0.1 - this one should fail as well
	sequencename = "loopbacksequence2"

	triggerSequenceAndVerifyTaskFinishedEvent(sequencename, "loopback2", func(t *testing.T, decodedEvent map[string]interface{}) {
		// check the result - this time it should be set to fail because an unknown Key was referenced in the webhook
		require.Equal(t, string(keptnv2.ResultFailed), decodedEvent["result"])
		require.Nil(t, decodedEvent["loopback"])
	})

	// Now, trigger another sequence that tries to execute a webhook with a call to the 127.0.0.1 - this one should fail as well
	sequencename = "loopbacksequence3"

	triggerSequenceAndVerifyTaskFinishedEvent(sequencename, "loopback3", func(t *testing.T, decodedEvent map[string]interface{}) {
		// check the result - this time it should be set to fail because an unknown Key was referenced in the webhook
		require.Equal(t, string(keptnv2.ResultFailed), decodedEvent["result"])
		require.Nil(t, decodedEvent["loopback"])
	})

	// Now, trigger another sequence that contains a task for which we don't have a webhook configured - this one should fail as well
	sequencename = "sequencewithunknowntask"

	triggerSequenceAndVerifyTaskFinishedEvent(sequencename, "unknowntask", func(t *testing.T, decodedEvent map[string]interface{}) {
		// check the result - this time it should be set to fail because an unknown Key was referenced in the webhook
		require.Equal(t, string(keptnv2.ResultFailed), decodedEvent["result"])
		require.Nil(t, decodedEvent["unknowntask"])
	})

	// Now, trigger another sequence that contains a task which results in a HTTP error status
	sequencename = "failedsequence"

	triggerSequenceAndVerifyTaskFinishedEvent(sequencename, "failedtask", func(t *testing.T, decodedEvent map[string]interface{}) {
		// check the result - this time it should be set to fail because an unknown Key was referenced in the webhook
		require.Equal(t, string(keptnv2.ResultFailed), decodedEvent["result"])
		require.Nil(t, decodedEvent["failedtask"])
	})
}

func Test_Webhook_OverlappingSubscriptions_Beta(t *testing.T) {
	projectName := "webhooks-subscription-overlap"
	serviceName := "myservice"
	projectName, shipyardFilePath := CreateWebhookProject(t, projectName, serviceName)
	defer deleteFile(t, shipyardFilePath)
	Test_Webhook_OverlappingSubscriptions(t, webhookWithOverlappingSubscriptionsYamlBeta, projectName, serviceName)
}

func Test_Webhook_OverlappingSubscriptions_Alpha(t *testing.T) {
	projectName := "webhooks-subscription-overlap-a"
	serviceName := "myservice"
	projectName, shipyardFilePath := CreateWebhookProject(t, projectName, serviceName)
	defer deleteFile(t, shipyardFilePath)
	Test_Webhook_OverlappingSubscriptions(t, webhookWithOverlappingSubscriptionsYamlAlpha, projectName, serviceName)
}

func Test_Webhook_OverlappingSubscriptions(t *testing.T, webhookWithOverlappingSubscriptionsYaml string, projectName, serviceName string) {

	stageName := "dev"
	sequencename := "mysequence"
	taskName := "mytask"

	// create subscriptions for the webhook-service
	webhookYamlWithSubscriptionIDs := webhookWithOverlappingSubscriptionsYaml
	subscriptionID, err := CreateSubscription(t, "webhook-service", models.EventSubscription{
		Event: keptnv2.GetTriggeredEventType(taskName),
		Filter: models.EventSubscriptionFilter{
			Projects: []string{projectName},
			Stages:   []string{stageName},
		},
	})
	require.Nil(t, err)

	webhookYamlWithSubscriptionIDs = strings.Replace(webhookYamlWithSubscriptionIDs, "${mytask-sub-id}", subscriptionID, -1)

	// create a second subscription that overlaps with the previously created one
	subscriptionID2, err := CreateSubscription(t, "webhook-service", models.EventSubscription{
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
	webhookFilePath, err := CreateTmpFile("webhook.yaml", webhookYamlWithSubscriptionIDs)
	require.Nil(t, err)
	defer deleteFile(t, webhookFilePath)

	t.Log("Adding webhook.yaml to our service")
	_, err = ExecuteCommand(fmt.Sprintf("keptn add-resource --project=%s --service=%s --resource=%s --resourceUri=webhook/webhook.yaml --all-stages", projectName, serviceName, webhookFilePath))

	require.Nil(t, err)

	t.Logf("triggering sequence %s in stage %s", sequencename, stageName)
	keptnContextID, _ := TriggerSequence(projectName, serviceName, stageName, sequencename, nil)

	var taskFinishedEvent []*models.KeptnContextExtendedCE
	require.Eventually(t, func() bool {
		taskFinishedEvent, err = GetEventsOfType(keptnContextID, projectName, stageName, keptnv2.GetFinishedEventType(taskName))
		if err != nil || taskFinishedEvent == nil || len(taskFinishedEvent) != 2 {
			return false
		}
		return true
	}, 30*time.Second, 3*time.Second)

}

func Test_WebhookFailInternalAddress_Beta(t *testing.T) {
	projectName := "webhooks-fail-internal-host-b"
	serviceName := "myservice"
	Test_WebhookConfigAtProjectLevel(t, failwebhookSimpleYamlBeta, projectName, serviceName, false)
}

func Test_WebhookConfigAtProjectLevel_Alpha(t *testing.T) {
	projectName := "webhooks-config-project-a"
	serviceName := "myservice"
	Test_WebhookConfigAtProjectLevel(t, webhookSimpleYamlAlpha, projectName, serviceName, true)
}

func Test_WebhookConfigAtProjectLevel_Beta(t *testing.T) {
	projectName := "webhooks-config-project-b"
	serviceName := "myservice"
	Test_WebhookConfigAtProjectLevel(t, webhookSimpleYamlBeta, projectName, serviceName, true)
}

func Test_WebhookConfigAtProjectLevel(t *testing.T, webhookSimpleYaml string, projectName, serviceName string, pass bool) {
	stageName := "dev"
	sequencename := "mysequence"
	projectName, shipyardFile := CreateWebhookProject(t, projectName, serviceName)
	defer deleteFile(t, shipyardFile)
	simpleWebhookTest(t, pass, stageName, projectName, serviceName, sequencename, "mytask", func(t *testing.T, projectName, webhookFilePath string) {
		t.Log("Adding webhook.yaml to our service")
		_, err := ExecuteCommand(fmt.Sprintf("keptn add-resource --project=%s --resource=%s --resourceUri=webhook/webhook.yaml", projectName, webhookFilePath))

		require.Nil(t, err)
	}, webhookSimpleYaml)
}

func Test_WebhookConfigAtStageLevel_Alpha(t *testing.T) {
	projectName := "webhooks-config-stage-a"
	Test_WebhookConfigAtStageLevel(t, webhookSimpleYamlAlpha, projectName)
}

func Test_WebhookConfigAtStageLevel_Beta(t *testing.T) {
	projectName := "webhooks-config-stage-b"
	Test_WebhookConfigAtStageLevel(t, webhookSimpleYamlBeta, projectName)
}

func Test_WebhookConfigAtStageLevel(t *testing.T, webhookSimpleYaml string, projectName string) {
	serviceName := "myservice"
	stageName := "dev"
	sequencename := "mysequence"
	projectName, shipyardFile := CreateWebhookProject(t, projectName, serviceName)
	defer deleteFile(t, shipyardFile)
	simpleWebhookTest(t, true, stageName, projectName, serviceName, sequencename, "mytask", func(t *testing.T, projectName, webhookFilePath string) {
		t.Log("Adding webhook.yaml to our service")
		_, err := ExecuteCommand(fmt.Sprintf("keptn add-resource --project=%s --stage=%s --resource=%s --resourceUri=webhook/webhook.yaml", projectName, stageName, webhookFilePath))

		require.Nil(t, err)
	}, webhookSimpleYaml)
}

func Test_WebhookConfigAtServiceLevel_Alpha(t *testing.T) {
	Test_WebhookConfigAtServiceLevel(t, webhookSimpleYamlAlpha)
}

func Test_WebhookConfigAtServiceLevel_Beta(t *testing.T) {
	Test_WebhookConfigAtServiceLevel(t, webhookSimpleYamlBeta)
}

func Test_WebhookConfigAtServiceLevel(t *testing.T, webhookSimpleYaml string) {
	projectName := "webhooks-config-service"
	serviceName := "myservice"
	stageName := "dev"
	sequencename := "mysequence"
	projectName, shipyardFile := CreateWebhookProject(t, projectName, serviceName)
	defer deleteFile(t, shipyardFile)
	simpleWebhookTest(t, true, stageName, projectName, serviceName, sequencename, "mytask", func(t *testing.T, projectName, webhookFilePath string) {
		t.Log("Adding webhook.yaml to our service")
		_, err := ExecuteCommand(fmt.Sprintf("keptn add-resource --project=%s --service=%s --resource=%s --resourceUri=webhook/webhook.yaml --all-stages", projectName, serviceName, webhookFilePath))

		require.Nil(t, err)
	}, webhookSimpleYaml)
}

// simpleWebhookTest triggers a sequence and checks whether a started and finished event is sent for the given task
func simpleWebhookTest(t *testing.T, pass bool, stageName, projectName, serviceName, sequencename, taskname string, addConfigFunc func(t *testing.T, projectName, webhookFilePath string), webhookSimpleYaml string) {

	// create subscriptions for the webhook-service
	taskTypes := []string{"mytask"}

	webhookYamlWithSubscriptionIDs := webhookSimpleYaml
	webhookYamlWithSubscriptionIDs = getWebhookYamlWithSubscriptionIDs(t, taskTypes, projectName, webhookYamlWithSubscriptionIDs)

	// wait some time to make sure the webhook service has pulled the updated subscription
	<-time.After(20 * time.Second) // sorry :(

	// now, let's add a webhook.yaml file to our service
	webhookFilePath, err := CreateTmpFile("webhook.yaml", webhookYamlWithSubscriptionIDs)
	require.Nil(t, err)
	defer deleteFile(t, webhookFilePath)

	addConfigFunc(t, projectName, webhookFilePath)

	t.Logf("triggering sequence %s in stage %s", sequencename, stageName)
	keptnContextID, _ := TriggerSequence(projectName, serviceName, stageName, sequencename, nil)

	var taskFinishedEvents []*models.KeptnContextExtendedCE
	t.Logf("Checking for started event of task %s", taskname)
	require.Eventually(t, func() bool {
		taskFinishedEvents, err = GetEventsOfType(keptnContextID, projectName, stageName, keptnv2.GetStartedEventType(taskname))
		if err != nil {
			t.Logf("got error: %s. will try again in a few seconds", err.Error())
			return false
		} else if taskFinishedEvents == nil {
			t.Log("did not receive any .started events")
			return false
		} else if len(taskFinishedEvents) != 1 {
			t.Logf("received %d .started events, but expected %d", len(taskFinishedEvents), 1)
			return false
		}
		return true
	}, 30*time.Second, 3*time.Second)

	t.Logf("Checking for finished event of task %s", taskname)
	require.Eventually(t, func() bool {
		taskFinishedEvents, err = GetEventsOfType(keptnContextID, projectName, stageName, keptnv2.GetFinishedEventType(taskname))
		if err != nil {
			t.Logf("got error: %s. will try again in a few seconds", err.Error())
			return false
		} else if taskFinishedEvents == nil {
			t.Log("did not receive any .started events")
			return false
		} else if len(taskFinishedEvents) != 1 {
			t.Logf("received %d .started events, but expected %d", len(taskFinishedEvents), 1)
			return false
		}
		return true
	}, 30*time.Second, 3*time.Second)

	decodedEvent := map[string]interface{}{}

	err = keptnv2.EventDataAs(*taskFinishedEvents[0], &decodedEvent)
	require.Nil(t, err)
	if pass {
		require.Equal(t, string(keptnv2.ResultPass), decodedEvent["result"])
	} else {
		require.Equal(t, string(keptnv2.ResultFailed), decodedEvent["result"])
	}
}

func Test_WebhookWithDisabledFinishedEvents_Beta(t *testing.T) {
	serviceName := "myservice"
	projectName, shipyardFilePath := CreateWebhookProject(t, "webhooks-no-finish-b", serviceName)
	defer deleteFile(t, shipyardFilePath)
	Test_WebhookWithDisabledFinishedEvents(t, webhookWithDisabledFinishedEventsYamlBeta, projectName, serviceName)
}
func Test_WebhookWithDisabledFinishedEvents_Alpha(t *testing.T) {
	serviceName := "myservice"
	projectName, shipyardFilePath := CreateWebhookProject(t, "webhooks-no-finish-a", serviceName)
	defer deleteFile(t, shipyardFilePath)
	Test_WebhookWithDisabledFinishedEvents(t, webhookWithDisabledFinishedEventsYamlAlpha, projectName, serviceName)
}

func Test_WebhookWithDisabledFinishedEvents(t *testing.T, webhookWithDisabledFinishedEventsYaml string, projectName, serviceName string) {

	stageName := "dev"
	sequencename := "mysequence"

	// create subscriptions for the webhook-service
	taskTypes := []string{"mytask", "othertask", "unallowedtask", "unknowntask"}

	webhookYamlWithSubscriptionIDs := webhookWithDisabledFinishedEventsYaml
	webhookYamlWithSubscriptionIDs = getWebhookYamlWithSubscriptionIDs(t, taskTypes, projectName, webhookYamlWithSubscriptionIDs)

	// wait some time to make sure the webhook service has pulled the updated subscription
	<-time.After(20 * time.Second) // sorry :(

	// now, let's add an webhook.yaml file to our service
	webhookFilePath, err := CreateTmpFile("webhook.yaml", webhookYamlWithSubscriptionIDs)
	require.Nil(t, err)
	defer deleteFile(t, webhookFilePath)

	t.Log("Adding webhook.yaml to our service")
	_, err = ExecuteCommand(fmt.Sprintf("keptn add-resource --project=%s --service=%s --resource=%s --resourceUri=webhook/webhook.yaml --all-stages", projectName, serviceName, webhookFilePath))

	require.Nil(t, err)

	triggerSequenceAndVerifyStartedEvents := func(sequencename, taskName string, nrExpected int) string {
		t.Logf("triggering sequence %s in stage %s", sequencename, stageName)
		keptnContextID, _ := TriggerSequence(projectName, serviceName, stageName, sequencename, nil)

		var taskStartedEvents []*models.KeptnContextExtendedCE
		require.Eventually(t, func() bool {
			taskStartedEvents, err = GetEventsOfType(keptnContextID, projectName, stageName, keptnv2.GetStartedEventType(taskName))
			if err != nil {
				t.Logf("got error: %s. will try again in a few seconds", err.Error())
				return false
			} else if taskStartedEvents == nil {
				t.Log("did not receive any .started events")
				return false
			} else if len(taskStartedEvents) != nrExpected {
				t.Logf("received %d .started events, but expected %d", len(taskStartedEvents), nrExpected)
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
	resp, err := ApiPOSTRequest(fmt.Sprintf("/controlPlane/v1/sequence/%s/%s/control", projectName, keptnContextID), models.SequenceControlCommand{
		State: models.AbortSequence,
		Stage: "",
	}, 3)
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

func Test_WebhookWithDisabledStartedEvents_Beta(t *testing.T) {
	projectName := "webhooks-no-started"
	serviceName := "myservice"
	projectName, shipyardFilePath := CreateWebhookProject(t, projectName, serviceName)
	defer deleteFile(t, shipyardFilePath)
	Test_WebhookWithDisabledStartedEvents(t, webhookWithDisabledStartedEventsYamlBeta, projectName, serviceName)
}

func Test_WebhookWithDisabledStartedEvents_Alpha(t *testing.T) {
	projectName := "webhooks-no-started-a"
	serviceName := "myservice"
	projectName, shipyardFilePath := CreateWebhookProject(t, projectName, serviceName)
	defer deleteFile(t, shipyardFilePath)
	Test_WebhookWithDisabledStartedEvents(t, webhookWithDisabledStartedEventsYamlAlpha, projectName, serviceName)
}

func Test_WebhookWithDisabledStartedEvents(t *testing.T, webhookWithDisabledStartedEventsYaml string, projectName, serviceName string) {

	stageName := "dev"
	sequencename := "mysequence"

	// create subscriptions for the webhook-service
	taskTypes := []string{"mytask"}

	webhookYamlWithSubscriptionIDs := webhookWithDisabledStartedEventsYaml
	webhookYamlWithSubscriptionIDs = getWebhookYamlWithSubscriptionIDs(t, taskTypes, projectName, webhookYamlWithSubscriptionIDs)

	// wait some time to make sure the webhook service has pulled the updated subscription
	<-time.After(20 * time.Second) // sorry :(

	// now, let's add an webhook.yaml file to our service
	webhookFilePath, err := CreateTmpFile("webhook.yaml", webhookYamlWithSubscriptionIDs)
	require.Nil(t, err)
	defer func() {
		err := os.Remove(webhookFilePath)
		if err != nil {
			t.Logf("Could not delete tmp file: %s", err.Error())
		}
	}()

	t.Log("Adding webhook.yaml to our service")
	_, err = ExecuteCommand(fmt.Sprintf("keptn add-resource --project=%s --service=%s --resource=%s --resourceUri=webhook/webhook.yaml --all-stages", projectName, serviceName, webhookFilePath))

	require.Nil(t, err)

	keptnContextID, _ := TriggerSequence(projectName, serviceName, stageName, sequencename, nil)

	<-time.After(5 * time.Second)
	// verify that no .started events have been sent for 'mytask'
	taskStartedEvent, err := GetLatestEventOfType(keptnContextID, projectName, stageName, keptnv2.GetStartedEventType("mytask"))

	require.Nil(t, taskStartedEvent)

	// verify that no .finished event has been sent for 'mytask'
	taskFinishedEvent, err := GetLatestEventOfType(keptnContextID, projectName, stageName, keptnv2.GetFinishedEventType("mytask"))

	require.Nil(t, taskFinishedEvent)

	t.Log("verified desired state, aborting sequence")
	resp, err := ApiPOSTRequest(fmt.Sprintf("/controlPlane/v1/sequence/%s/%s/control", projectName, keptnContextID), models.SequenceControlCommand{
		State: models.AbortSequence,
		Stage: "",
	}, 3)
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.Response().StatusCode)
}

func getWebhookYamlWithSubscriptionIDs(t *testing.T, taskTypes []string, projectName string, webhookYamlWithSubscriptionIDs string) string {
	for _, taskType := range taskTypes {
		eventType := keptnv2.GetTriggeredEventType(taskType)
		if strings.HasSuffix(taskType, "-finished") {
			eventType = keptnv2.GetFinishedEventType(strings.TrimSuffix(taskType, "-finished"))
		}
		subscriptionID, err := CreateSubscription(t, "webhook-service", models.EventSubscription{
			Event: eventType,
			Filter: models.EventSubscriptionFilter{
				Projects: []string{projectName},
			},
		})
		require.Nil(t, err)

		subscriptionPlaceholder := fmt.Sprintf("${%s-sub-id}", taskType)
		webhookYamlWithSubscriptionIDs = strings.Replace(webhookYamlWithSubscriptionIDs, subscriptionPlaceholder, subscriptionID, -1)
	}
	return webhookYamlWithSubscriptionIDs
}

func deleteFile(t *testing.T, shipyardFilePath string) {
	func() {
		err := os.Remove(shipyardFilePath)
		if err != nil {
			t.Logf("Could not delete tmp file: %s", err.Error())
		}
	}()
}
