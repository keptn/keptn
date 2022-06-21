package go_tests

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/keptn/go-utils/pkg/api/models"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/stretchr/testify/require"
)

const unleashServiceK8sManifest = `---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: unleash-service
  labels:
    keptn.sh/integration-name: unleash-service-integration-name
spec:
  selector:
    matchLabels:
      run: unleash-service
  replicas: 1
  template:
    metadata:
      labels:
        run: unleash-service
        app.kubernetes.io/name: keptn
        app.kubernetes.io/component: unleash-service
        app.kubernetes.io/version: 0.3.2
        keptn.sh/integration-name: unleash-service-integration-name
    spec:
      serviceAccountName: keptn-default
      containers:
        - name: unleash-service
          image: keptncontrib/unleash-service:0.3.2
          ports:
            - containerPort: 8080
          env:
            - name: EVENTBROKER
              value: 'http://localhost:8081/event'
            - name: CONFIGURATION_SERVICE
              value: 'http://configuration-service:8080'
          envFrom:
            - secretRef:
                name: unleash
                optional: true
        - name: distributor
          image: ${distributor-image}
          ports:
            - containerPort: 8080
          resources:
            requests:
              memory: "16Mi"
              cpu: "25m"
            limits:
              memory: "32Mi"
              cpu: "250m"
          env:
            - name: PUBSUB_URL
              value: 'nats://keptn-nats'
            - name: PUBSUB_TOPIC
              value: 'sh.keptn.event.action.triggered'
            - name: PUBSUB_RECIPIENT
              value: '127.0.0.1'
            - name: VERSION
              valueFrom:
                fieldRef:
                  apiVersion: v1
                  fieldPath: 'metadata.labels[''app.kubernetes.io/version'']'
            - name: DISTRIBUTOR_VERSION
              valueFrom:
                fieldRef:
                  fieldPath: metadata.labels['app.kubernetes.io/version']
            - name: INTEGRATION_NAME
              valueFrom:
                fieldRef:
                  apiVersion: v1
                  fieldPath: metadata.labels['keptn.sh/integration-name']
            - name: K8S_POD_NAME
              valueFrom:
                fieldRef:
                  apiVersion: v1
                  fieldPath: metadata.name
            - name: K8S_NAMESPACE
              valueFrom:
                fieldRef:
                  apiVersion: v1
                  fieldPath: metadata.namespace
            - name: K8S_NODE_NAME
              valueFrom:
                fieldRef:
                  apiVersion: v1
                  fieldPath: spec.nodeName
---
apiVersion: v1
kind: Service
metadata:
  name: unleash-service
  labels:
    run: unleash-service
spec:
  ports:
    - port: 8080
      protocol: TCP
  selector:
    run: unleash-service`

const selfHealingShipyard = `apiVersion: "spec.keptn.sh/0.2.2"
kind: "Shipyard"
metadata:
  name: "shipyard-sockshop"
spec:
  stages:
    - name: "production"
      sequences:
        - name: "remediation"
          triggeredOn:
            - event: "production.remediation.finished"
              selector:
                match:
                  evaluation.result: "fail"
          tasks:
            - name: "get-action"
            - name: "action"
            - name: "evaluation"
              triggeredAfter: "1m"`

const remediationFileContent = `apiVersion: spec.keptn.sh/0.1.4
kind: Remediation
metadata:
  name: service-remediation
spec:
  remediations:
    - problemType: Response time degradation
      actionsOnOpen:
      - action: toggle-feature
        name: toggle-feature
        description: Toggle feature flag EnablePromotion to OFF
        value:
          EnablePromotion: "off"`

const defaultUnleashServiceVersion = "master"
const unleashServiceEnvVar = "UNLEASH_SERVICE_VERSION"

type RemediationTriggered struct {
	keptnv2.EventData
	Problem keptnv2.ProblemDetails `json:"problem"`
}

func Test_SelfHealing(t *testing.T) {
	projectName := "self-healing"
	serviceName := "my-service"
	shipyardFilePath, err := CreateTmpShipyardFile(selfHealingShipyard)
	require.Nil(t, err)
	defer os.Remove(shipyardFilePath)

	_, err = ExecuteCommand(fmt.Sprintf("kubectl delete configmap -n %s lighthouse-config-%s", GetKeptnNameSpaceFromEnv(), projectName))
	t.Logf("creating project %s", projectName)
	projectName, err = CreateProject(projectName, shipyardFilePath)
	require.Nil(t, err)

	t.Logf("creating service %s", serviceName)
	_, err = ExecuteCommand(fmt.Sprintf("keptn create service %s --project=%s", serviceName, projectName))
	require.Nil(t, err)

	// trigger a remediation - this should fail because no remediation.yaml is available yet
	t.Log("triggering a remediation with no remediation.yaml")
	remediationFinishedEvent := performRemediation(t, projectName, serviceName)

	require.Equal(t, "shipyard-controller", *remediationFinishedEvent.Source)
	finishedEventData := &keptnv2.EventData{}
	err = keptnv2.Decode(remediationFinishedEvent.Data, finishedEventData)
	require.Nil(t, err)
	require.Equal(t, "shipyard-controller", *remediationFinishedEvent.Source)
	require.Equal(t, keptnv2.StatusErrored, finishedEventData.Status)
	require.Equal(t, keptnv2.ResultFailed, finishedEventData.Result)

	t.Log("adding remediation.yaml file")
	remediationFilePath, err := CreateTmpFile("remediation-*.yaml", remediationFileContent)
	defer os.Remove(remediationFilePath)
	require.Nil(t, err)
	require.NotEmpty(t, remediationFilePath)

	_, err = ExecuteCommand(fmt.Sprintf("keptn add-resource --project=%s --stage=%s --service=%s --resource=%s --resourceUri=remediation.yaml", projectName, "production", serviceName, remediationFilePath))
	require.Nil(t, err)

	t.Log("Installing unleash-service")
	imageName, err := GetImageOfDeploymentContainer("lighthouse-service", "lighthouse-service")
	require.Nil(t, err)
	distributorImage := strings.Replace(imageName, "lighthouse-service", "distributor", 1)
	unleashServiceManifestContent := strings.ReplaceAll(unleashServiceK8sManifest, "${distributor-image}", distributorImage)

	tmpFile, err := CreateTmpFile("unleash-service-*.yaml", unleashServiceManifestContent)
	defer func() {
		if err := os.Remove(tmpFile); err != nil {
			t.Logf("Could not delete file: %v", err)
		}
	}()
	_, err = KubeCtlApplyFromURL(tmpFile)
	require.Nil(t, err)

	err = WaitForPodOfDeployment("unleash-service")
	require.Nil(t, err)

	var uniformServiceIntegration models.Integration
	require.Eventually(t, func() bool {
		uniformServiceIntegration, err = GetIntegrationWithName("unleash-service")
		return err == nil
	}, time.Second*20, time.Second*3)

	require.NotEmpty(t, uniformServiceIntegration.Subscriptions)

	// it seems like the unleash service is not immediately ready after the distributor has registered itself
	// to be safe let's wait a couple of seconds here. This can be removed as soon as we have decoupled our tests from the unleash-service
	<-time.After(15 * time.Second)

	t.Log("remediation.yaml and unleash-service are ready. let's trigger another remediation")
	remediationFinishedEvent = performRemediation(t, projectName, serviceName)

	// inspect the remediation.finished event again
	finishedEventData = &keptnv2.EventData{}
	err = keptnv2.Decode(remediationFinishedEvent.Data, finishedEventData)
	require.Nil(t, err)
	require.Equal(t, "shipyard-controller", *remediationFinishedEvent.Source)
	require.Equal(t, keptnv2.StatusErrored, finishedEventData.Status)
	require.Equal(t, keptnv2.ResultFailed, finishedEventData.Result)

	t.Log("verifying if action.triggered event has been sent")
	var actionTriggered *models.KeptnContextExtendedCE
	require.Eventually(t, func() bool {
		event, err := GetLatestEventOfType(remediationFinishedEvent.Shkeptncontext, projectName, "production", keptnv2.GetTriggeredEventType(keptnv2.ActionTaskName))
		if err != nil || event == nil {
			return false
		}
		actionTriggered = event
		return true
	}, 1*time.Minute, 10*time.Second)

	triggeredEventData := &keptnv2.ActionTriggeredEventData{}
	err = keptnv2.Decode(actionTriggered.Data, triggeredEventData)
	require.Equal(t, keptnv2.ActionInfo{
		Name:        "toggle-feature",
		Action:      "toggle-feature",
		Description: "Toggle feature flag EnablePromotion to OFF",
		Value: map[string]interface{}{
			"EnablePromotion": "off",
		},
	}, triggeredEventData.Action)

	t.Log("verifying if action.finished event has been sent")
	var actionFinished *models.KeptnContextExtendedCE
	require.Eventually(t, func() bool {
		event, err := GetLatestEventOfType(remediationFinishedEvent.Shkeptncontext, projectName, "production", keptnv2.GetFinishedEventType(keptnv2.ActionTaskName))
		if err != nil || event == nil {
			return false
		}
		actionFinished = event
		return true
	}, 1*time.Minute, 10*time.Second)

	finishedEventData = &keptnv2.EventData{}
	err = keptnv2.Decode(actionFinished.Data, finishedEventData)
	require.Equal(t, keptnv2.StatusErrored, finishedEventData.Status)
}

func performRemediation(t *testing.T, projectName string, serviceName string) *models.KeptnContextExtendedCE {
	keptnContext, err := TriggerSequence(projectName, serviceName, "production", "remediation", &RemediationTriggered{
		Problem: keptnv2.ProblemDetails{
			RootCause:    "Response time degradation",
			ProblemTitle: "My Problem",
		},
	})

	require.Nil(t, err)
	require.NotEmpty(t, keptnContext)

	t.Log("waiting for remediation.finished event to be available")
	var remediationFinishedEvent *models.KeptnContextExtendedCE
	require.Eventually(t, func() bool {
		event, err := GetLatestEventOfType(keptnContext, projectName, "production", keptnv2.GetFinishedEventType("production.remediation"))
		if err != nil || event == nil {
			return false
		}
		remediationFinishedEvent = event
		return true
	}, 1*time.Minute, 10*time.Second)
	return remediationFinishedEvent
}
