package zero_downtime

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/keptn/go-utils/pkg/api/models"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	testutils "github.com/keptn/keptn/test/go-tests"
	keptnkubeutils "github.com/keptn/kubernetes-utils/pkg"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
	"strings"
	"sync"
	"testing"
	"time"
)

const zdShipyard = `apiVersion: "spec.keptn.sh/0.2.0"
kind: "Shipyard"
metadata:
  name: "shipyard-sockshop"
spec:
  stages:
    - name: "hardening"
      sequences:
        - name: "remediation"
          tasks:
            - name: "action"
            - name: "approval"
              properties:
                pass: "automatic"
                warning: "automatic"
            - name: "evaluation"
              properties:
                timeframe: "5m"
        - name: "evaluation"
          tasks:
            - name: "evaluation"
            - name: "approval"
              properties:
                pass: "automatic"
                warning: "automatic"`

const webhookYaml = `apiVersion: webhookconfig.keptn.sh/v1alpha1
kind: WebhookConfig
metadata:
  name: webhook-configuration
spec:
  webhooks:
    - type: sh.keptn.event.action.triggered
      requests:
        - >-
          curl --request GET http://shipyard-controller:8080/v1/project
      subscriptionID: ${action-sub-id}
      sendFinished: true
      sendStarted: true`

const remediationYaml = `apiVersion: spec.keptn.sh/0.1.4
kind: Remediation
metadata:
  name: remediation-configuration
spec:
  remediations: 
  - problemType: "default"
    actionsOnOpen:
    - name: Execute webhook
      action: webhook
      description: Execute a nice webhook`

const sloYaml = `---
spec_version: '0.1.0'
comparison:
  compare_with: "single_result"
  include_result_with_score: "pass"
  aggregate_function: avg
objectives:
  - sli: test-metric
    pass:
      - criteria:
          - "<=4"
    warning:
      - criteria:
          - "<=5"
total_score:
  pass: "51"
  warning: "20"`

func TestEvaluationsWithApproval(t *testing.T) {
	images := []string{"0.15.1-dev.202205240824", "0.15.1-dev.202205240902"}
	services := []string{"api-service", "shipyard-controller", "resource-service", "lighthouse-service", "approval-service", "webhook-service", "remediation-service", "mongodb-datastore"}
	//services := []string{"mongodb-datastore"}

	project := "a-zd-test"
	stage := "hardening"
	service := "myservice"

	shipyardFile, err := testutils.CreateTmpShipyardFile(zdShipyard)
	require.Nil(t, err)
	defer os.Remove(shipyardFile)

	t.Logf("Creating project %s", project)
	project, err = testutils.CreateProject(project, shipyardFile)
	require.Nil(t, err)

	t.Logf("creating service %s", service)
	_, err = testutils.ExecuteCommand(fmt.Sprintf("keptn create service %s --project=%s", service, project))

	taskTypes := []string{"action"}

	t.Logf("Setting up subscription")
	webhookYamlWithSubscriptionIDs := webhookYaml
	webhookYamlWithSubscriptionIDs = getWebhookYamlWithSubscriptionIDs(t, taskTypes, project, webhookYamlWithSubscriptionIDs)

	// wait some time to make sure the webhook service has pulled the updated subscription
	<-time.After(20 * time.Second) // sorry :(

	t.Logf("Adding webhook")
	// now, let's add an webhook.yaml file to our service
	webhookFilePath, err := testutils.CreateTmpFile("webhook.yaml", webhookYamlWithSubscriptionIDs)
	require.Nil(t, err)
	defer func() {
		err := os.Remove(webhookFilePath)
		if err != nil {
			t.Logf("Could not delete tmp file: %s", err.Error())
		}

	}()

	t.Log("Adding webhook.yaml to our service")
	_, err = testutils.ExecuteCommand(fmt.Sprintf("keptn add-resource --project=%s --service=%s --resource=%s --resourceUri=webhook/webhook.yaml --all-stages", project, service, webhookFilePath))

	remediationFilePath, err := testutils.CreateTmpFile("remediation.yaml", remediationYaml)
	t.Log("Adding remediation.yaml to our service")
	_, err = testutils.ExecuteCommand(fmt.Sprintf("keptn add-resource --project=%s --service=%s --resource=%s --resourceUri=remediation.yaml --all-stages", project, service, remediationFilePath))

	t.Log("deleting lighthouse configmap from previous test run")
	testutils.ExecuteCommandf("kubectl delete configmap -n %s lighthouse-config-%s", testutils.GetKeptnNameSpaceFromEnv(), project)

	t.Log("adding SLI provider")
	_, err = testutils.ExecuteCommand(fmt.Sprintf("kubectl create configmap -n %s lighthouse-config-%s --from-literal=sli-provider=my-sli-provider", testutils.GetKeptnNameSpaceFromEnv(), project))
	require.Nil(t, err)

	sloFilePath, err := testutils.CreateTmpFile("slo.yaml", sloYaml)
	t.Log("Adding slo.yaml to our service")
	_, err = testutils.ExecuteCommand(fmt.Sprintf("keptn add-resource --project=%s --service=%s --resource=%s --resourceUri=slo.yaml --all-stages", project, service, sloFilePath))

	ctx, cancel := context.WithCancel(context.Background())

	for _, svc := range services {
		go func(service string) {
			err := updateImageOfService(ctx, t, service, images)
			if err != nil {
				t.Logf("%v", err)
			}
		}(svc)
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go startSLIRetrieval(t, project, stage, service)
	doEvaluations(project, stage, service)

	cancel()

	require.Eventually(t, func() bool {
		states := &models.SequenceStates{}
		t.Log("Checking if all sequences are completed")
		resp, err := testutils.ApiGETRequest("/controlPlane/v1/sequence/"+project+"?state=started", 3)
		if err != nil {
			return false
		}
		err = resp.ToJSON(states)
		if err != nil {
			return false
		}

		if states.TotalCount != 0 {
			t.Logf("Currently there are still %d open sequences", states.TotalCount)
			return false
		}
		t.Logf("All sequences completed!")
		return true
	}, 10*time.Minute, 10*time.Second)
}

func doEvaluations(project, stage, service string) {
	for i := 0; i < 3; i++ {
		nrEvaluations := 0
		go func() {
			//_, err := triggerEvaluation("podtatohead", "hardening", "helloservice")
			_, err := triggerRemediation(project, stage, service)
			if err != nil {
				nrEvaluations++
			}
		}()

		<-time.After(5 * time.Second)
	}
}

func triggerEvaluation(projectName, stageName, serviceName string) (string, error) {
	cliResp, err := testutils.ExecuteCommand(fmt.Sprintf("keptn trigger evaluation --project=%s --stage=%s --service=%s --timeframe=5m", projectName, stageName, serviceName))

	if err != nil {
		return "", err
	}
	var keptnContext string
	split := strings.Split(cliResp, "\n")
	for _, line := range split {
		if strings.Contains(line, "ID of") {
			splitLine := strings.Split(line, ":")
			if len(splitLine) == 2 {
				keptnContext = strings.TrimSpace(splitLine[1])
			}
		}
	}
	return keptnContext, err
}

func triggerRemediation(projectName, stageName, serviceName string) (string, error) {
	source := "golang-test"
	eventData := keptnv2.EventData{}
	eventType := keptnv2.GetTriggeredEventType(stageName + ".remediation")
	eventData.SetProject(projectName)
	eventData.SetService(serviceName)
	eventData.SetStage(stageName)

	resp, err := testutils.ApiPOSTRequest("/v1/event", models.KeptnContextExtendedCE{
		Contenttype:        "application/json",
		Data:               eventData,
		ID:                 uuid.NewString(),
		Shkeptnspecversion: "0.2.0",
		Source:             &source,
		Specversion:        "1.0",
		Type:               &eventType,
	}, 0)

	eventContext := &models.EventContext{}
	err = resp.ToJSON(eventContext)
	if err != nil {
		return "", err
	}
	return *eventContext.KeptnContext, nil
}

func updateImageOfService(ctx context.Context, t *testing.T, service string, images []string) error {
	clientset, err := keptnkubeutils.GetClientset(false)

	if err != nil {
		return err
	}

	i := 0
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			nextImage := images[i%len(images)]
			get, err := clientset.AppsV1().Deployments(testutils.GetKeptnNameSpaceFromEnv()).Get(context.TODO(), service, v1.GetOptions{})
			if err != nil {
				break
			}

			imageWithTag := get.Spec.Template.Spec.Containers[0].Image
			split := strings.Split(imageWithTag, ":")
			updatedImage := fmt.Sprintf("%s:%s", split[0], nextImage)

			get.Spec.Template.Spec.Containers[0].Image = updatedImage

			t.Logf("upgrading %s to %s", service, updatedImage)
			_, err = clientset.AppsV1().Deployments(testutils.GetKeptnNameSpaceFromEnv()).Update(context.TODO(), get, v1.UpdateOptions{})
			if err != nil {
				break
			}

			require.Eventually(t, func() bool {
				pods, err := clientset.CoreV1().Pods(testutils.GetKeptnNameSpaceFromEnv()).List(context.TODO(), v1.ListOptions{LabelSelector: "app.kubernetes.io/name=" + service})
				if err != nil {
					return false
				}

				if int32(len(pods.Items)) != 1 {
					// make sure only one pod is running
					return false
				}

				for _, item := range pods.Items {
					if len(item.Spec.Containers) == 0 {
						continue
					}
					if item.Spec.Containers[0].Image == updatedImage {
						return true
					}
				}
				return false
			}, 3*time.Minute, 10*time.Second)
			<-time.After(5 * time.Second)
			i++
		}
	}
}

func startSLIRetrieval(t *testing.T, project, stage, service string) {
	for {
		<-time.After(3 * time.Second)
		if err := reportSLIValues(project, stage, service); err != nil {
			t.Logf("Error while SLI retrieval: %v", err)
		}
	}
}

func reportSLIValues(project string, stage string, service string) error {
	resp, err := testutils.ApiGETRequest(fmt.Sprintf("/mongodb-datastore/event?project=%s&stage=%s&service=%s&type=sh.keptn.event.get-sli.triggered", project, stage, service), 3)
	if err != nil {
		return err
	}
	events := &models.Events{}
	if err := resp.ToJSON(events); err != nil {
		return err
	}

	if len(events.Events) == 0 {
		return nil
	}

	sliFinishedEventType := keptnv2.GetFinishedEventType(keptnv2.GetSLITaskName)
	source := "golang-test"
	for _, sliTriggeredEvent := range events.Events {
		_, err := testutils.ApiPOSTRequest("/v1/event", models.KeptnContextExtendedCE{
			Contenttype: "application/json",
			Data: keptnv2.GetSLIFinishedEventData{
				EventData: keptnv2.EventData{
					Project: project,
					Stage:   stage,
					Service: service,
					Status:  keptnv2.StatusSucceeded,
					Result:  keptnv2.ResultPass,
				},
				GetSLI: keptnv2.GetSLIFinished{
					IndicatorValues: []*keptnv2.SLIResult{
						{
							Metric:  "test-metric",
							Value:   1,
							Success: true,
						},
					},
				},
			},
			ID:                 uuid.NewString(),
			Shkeptnspecversion: "0.2.0",
			Source:             &source,
			Specversion:        "1.0",
			Shkeptncontext:     sliTriggeredEvent.Shkeptncontext,
			Triggeredid:        sliTriggeredEvent.ID,
			Type:               &sliFinishedEventType,
		}, 0)

		if err != nil {
			continue
		}
	}
	return nil
}

func getWebhookYamlWithSubscriptionIDs(t *testing.T, taskTypes []string, projectName string, webhookYamlWithSubscriptionIDs string) string {
	for _, taskType := range taskTypes {
		eventType := keptnv2.GetTriggeredEventType(taskType)
		if strings.HasSuffix(taskType, "-finished") {
			eventType = keptnv2.GetFinishedEventType(strings.TrimSuffix(taskType, "-finished"))
		}
		subscriptionID, err := testutils.CreateSubscription(t, "webhook-service", models.EventSubscription{
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
