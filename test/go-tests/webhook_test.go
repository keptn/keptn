package go_tests

import (
	"fmt"
	"github.com/keptn/go-utils/pkg/api/models"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/stretchr/testify/require"
	"os"
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
        - name: "mysequence"
          tasks:
            - name: "mytask"`

const webhookYaml = `apiVersion: webhookconfig.keptn.sh/v1alpha1
kind: WebhookConfig
metadata:
  name: webhook-configuration
spec:
  webhooks:
    - type: "sh.keptn.event.mytask.triggered"
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

	// now, let's add an webhook.yaml file to our service
	webhookFilePath, err := CreateTmpFile("webhook.yaml", webhookYaml)
	require.Nil(t, err)
	defer os.Remove(webhookFilePath)

	t.Log("Adding webhook.yaml to our service")
	_, err = ExecuteCommand(fmt.Sprintf("keptn add-resource --project=%s --service=%s --resource=%s --resourceUri=webhook.yaml --all-stages", projectName, serviceName, webhookFilePath))

	require.Nil(t, err)

	t.Logf("triggering sequence %s in stage %s", sequencename, stageName)
	keptnContextID, _ := TriggerSequence(projectName, serviceName, stageName, sequencename, nil)

	var taskFinishedEvent *models.KeptnContextExtendedCE
	require.Eventually(t, func() bool {
		taskFinishedEvent, err = GetLatestEventOfType(keptnContextID, projectName, stageName, keptnv2.GetFinishedEventType("mytask"))
		if err != nil || taskFinishedEvent == nil {
			return false
		}
		return true
	}, 30*time.Second, 3*time.Second)

	require.NotNil(t, taskFinishedEvent)

	decodedEvent := map[string]interface{}{}

	err = keptnv2.EventDataAs(*taskFinishedEvent, &decodedEvent)

	require.Nil(t, err)

	// check if the requests have been executed and yielded some results
	require.NotNil(t, decodedEvent["sh.keptn.event.mytask.triggered"])

}
