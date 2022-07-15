package go_tests

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"testing"

	"github.com/keptn/go-utils/pkg/api/models"
	"github.com/stretchr/testify/require"
)

const webhookConfigMigrationAlpha = `apiVersion: webhookconfig.keptn.sh/v1alpha1
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
        - "curl --data '{\"email\":\"test@example.com\", \"name\": [\"Boolean\", \"World\"]}' -H \"Accept-Charset: utf-8\" -H 'Content-Type: application/json' https://httpbin.org/post --some-random-options -YYY"
    - type: "sh.keptn.event.failedtask.triggered"
      subscriptionID: ${failedtask-sub-id}
      sendFinished: true
      requests:
        - curl http://local:8080 {{.data.project}} {{.env.mysecret}}
    - type: "sh.keptn.event.unallowedtask.triggered"
      subscriptionID: ${unallowedtask-sub-id}
      sendFinished: true
      envFrom: 
        - name: "secretKey"
          secretRef:
            name: "my-webhook-k8s-secret"
            key: "my-key"
      requests:
        - curl http://local:8080 {{.data.project}} {{.env.mysecret}}`

const webhookConfigMigrationBeta = `apiVersion: webhookconfig.keptn.sh/v1beta1
kind: WebhookConfig
metadata:
    name: webhook-configuration
spec:
    webhooks:
        - type: sh.keptn.event.othertask.triggered
          subscriptionID: ${othertask-sub-id}
          sendFinished: true
          envFrom:
            - secretRef:
                key: my-key
                name: my-webhook-k8s-secret
              name: secretKey
          requests:
            - url: https://httpbin.org/post
              method: GET
              headers:
                - key: Accept-Charset
                  value: utf-8
                - key: Content-Type
                  value: application/json
              payload: '{"email":"test@example.com", "name": ["Boolean", "World"]}'
              options: --some-random-options -YYY
        - type: sh.keptn.event.failedtask.triggered
          subscriptionID: ${failedtask-sub-id}
          sendFinished: true
          envFrom: []
          requests:
            - url: http://local:8080
              method: GET
              options: '{{.data.project}} {{.env.mysecret}}'
        - type: sh.keptn.event.unallowedtask.triggered
          subscriptionID: ${unallowedtask-sub-id}
          sendFinished: true
          envFrom:
            - secretRef:
                key: my-key
                name: my-webhook-k8s-secret
              name: secretKey
          requests:
            - url: http://local:8080
              method: GET
              options: '{{.data.project}} {{.env.mysecret}}'
`
const webhookURI = "/%252Fwebhook%252Fwebhook.yaml"

func Test_Webhook_Migrator(t *testing.T) {
	projectName := "webhook-migration"
	serviceName := "myservice"
	shipyardFilePath, err := CreateTmpShipyardFile(webhookShipyard)
	require.Nil(t, err)
	defer func() {
		err := os.Remove(shipyardFilePath)
		if err != nil {
			t.Logf("Could not delete tmp file: %s", err.Error())
		}

	}()

	webhookFilePath, err := CreateTmpFile("webhook.yaml", webhookConfigMigrationAlpha)
	require.Nil(t, err)
	defer func() {
		err := os.Remove(webhookFilePath)
		if err != nil {
			t.Logf("Could not delete tmp file: %s", err.Error())
		}

	}()

	ctx, closeInternalKeptnAPI := context.WithCancel(context.Background())
	defer closeInternalKeptnAPI()
	internalKeptnAPI, err := GetInternalKeptnAPI(ctx, "service/resource-service", "8888", "8080")
	require.Nil(t, err)

	for i := 1; i <= 3; i++ {
		fillProjectWithWebhooks(t, shipyardFilePath, webhookFilePath, projectName+fmt.Sprint(i), serviceName+fmt.Sprint(i))
		checkWebhooksHaveCorrectWersion(t, internalKeptnAPI, "keptn-"+projectName+fmt.Sprint(i), serviceName+fmt.Sprint(i), webhookConfigMigrationAlpha)
	}

	t.Logf("Executing dry-run migration for project keptn-%s1", projectName)
	output, err := ExecuteCommandf("keptn migrate-webhooks --dry-run --project=keptn-%s1 -y", projectName)
	require.Nil(t, err)
	require.Contains(t, output, webhookConfigMigrationBeta)

	t.Logf("Checking if all projects still contain Alpha version")
	for i := 1; i <= 3; i++ {
		checkWebhooksHaveCorrectWersion(t, internalKeptnAPI, "keptn-"+projectName+fmt.Sprint(i), serviceName+fmt.Sprint(i), webhookConfigMigrationAlpha)
	}

	t.Logf("Executing migration for project keptn-%s1", projectName)
	_, err = ExecuteCommandf("keptn migrate-webhooks --project=keptn-%s1 -y", projectName)
	require.Nil(t, err)

	t.Logf("Checking if all webhooks contain the right version")
	checkWebhooksHaveCorrectWersion(t, internalKeptnAPI, "keptn-"+projectName+"1", serviceName+"1", webhookConfigMigrationBeta)
	checkWebhooksHaveCorrectWersion(t, internalKeptnAPI, "keptn-"+projectName+"2", serviceName+"2", webhookConfigMigrationAlpha)
	checkWebhooksHaveCorrectWersion(t, internalKeptnAPI, "keptn-"+projectName+"3", serviceName+"3", webhookConfigMigrationAlpha)

	t.Logf("Executing migration for all projects")
	_, err = ExecuteCommandf("keptn migrate-webhooks -y")
	require.Nil(t, err)

	t.Logf("Checking if all webhooks in all projects were migrated")
	for i := 1; i <= 3; i++ {
		checkWebhooksHaveCorrectWersion(t, internalKeptnAPI, "keptn-"+projectName+fmt.Sprint(i), serviceName+fmt.Sprint(i), webhookConfigMigrationBeta)
	}
}

func fillProjectWithWebhooks(t *testing.T, shipyardFilePath string, webhookFilePath string, projectName string, serviceName string) {
	t.Logf("creating project %s", projectName)
	projectName, err := CreateProject(projectName, shipyardFilePath)
	require.Nil(t, err)

	t.Logf("creating service %s", serviceName)
	output, err := ExecuteCommandf("keptn create service %s --project=%s", serviceName, projectName)
	require.Nil(t, err)
	require.Contains(t, output, "created successfully")

	t.Logf("adding webhook resource for project %s", projectName)
	output, err = ExecuteCommandf("keptn add-resource --project=%s --resource=%s --resourceUri=webhook/webhook.yaml", projectName, webhookFilePath)
	require.Nil(t, err)

	t.Logf("adding webhook resource for stage dev for project %s", projectName)
	output, err = ExecuteCommandf("keptn add-resource --project=%s --resource=%s --resourceUri=webhook/webhook.yaml --stage=dev", projectName, webhookFilePath)
	require.Nil(t, err)

	t.Logf("adding webhook resource for stage otherstage for project %s", projectName)
	output, err = ExecuteCommandf("keptn add-resource --project=%s --resource=%s --resourceUri=webhook/webhook.yaml --stage=otherstage", projectName, webhookFilePath)
	require.Nil(t, err)

	t.Logf("adding webhook resource for all services in all stages for project %s", projectName)
	output, err = ExecuteCommandf("keptn add-resource --project=%s --service=%s --resource=%s --resourceUri=webhook/webhook.yaml --all-stages", projectName, serviceName, webhookFilePath)
	require.Nil(t, err)
}

func checkWebhooksHaveCorrectWersion(t *testing.T, internalKeptnAPI *APICaller, projectName string, serviceName string, webhookConfig string) {
	t.Logf("Checking if webhook config has required format for project")
	resp, err := internalKeptnAPI.Get(basePath+"/"+projectName+"/resource"+webhookURI, 3)
	require.Nil(t, err)
	require.Equal(t, 200, resp.Response().StatusCode)

	t.Logf("Checking body of the received response")
	resource := models.Resource{}
	err = resp.ToJSON(&resource)
	require.Nil(t, err)
	rawDecoded, err := base64.StdEncoding.DecodeString(resource.ResourceContent)
	require.Nil(t, err)
	require.Equal(t, webhookConfig, string(rawDecoded))

	t.Logf("Checking if webhook config has required format for stage dev")
	resp, err = internalKeptnAPI.Get(basePath+"/"+projectName+"/stage/dev/resource"+webhookURI, 3)
	require.Nil(t, err)
	require.Equal(t, 200, resp.Response().StatusCode)

	t.Logf("Checking body of the received response")
	resource = models.Resource{}
	err = resp.ToJSON(&resource)
	require.Nil(t, err)
	rawDecoded, err = base64.StdEncoding.DecodeString(resource.ResourceContent)
	require.Nil(t, err)
	require.Equal(t, webhookConfig, string(rawDecoded))

	t.Logf("Checking if webhook config has required format for stage otherstage")
	resp, err = internalKeptnAPI.Get(basePath+"/"+projectName+"/stage/otherstage/resource"+webhookURI, 3)
	require.Nil(t, err)
	require.Equal(t, 200, resp.Response().StatusCode)

	t.Logf("Checking body of the received response")
	resource = models.Resource{}
	err = resp.ToJSON(&resource)
	require.Nil(t, err)
	rawDecoded, err = base64.StdEncoding.DecodeString(resource.ResourceContent)
	require.Nil(t, err)
	require.Equal(t, webhookConfig, string(rawDecoded))

	t.Logf("Checking if webhook config has required format for stage dev for service")
	resp, err = internalKeptnAPI.Get(basePath+"/"+projectName+"/stage/dev/service/"+serviceName+"/resource"+webhookURI, 3)
	require.Nil(t, err)
	require.Equal(t, 200, resp.Response().StatusCode)

	t.Logf("Checking body of the received response")
	resource = models.Resource{}
	err = resp.ToJSON(&resource)
	require.Nil(t, err)
	rawDecoded, err = base64.StdEncoding.DecodeString(resource.ResourceContent)
	require.Nil(t, err)
	require.Equal(t, webhookConfig, string(rawDecoded))

	t.Logf("Checking if webhook config has required format for stage otherstage for service")
	resp, err = internalKeptnAPI.Get(basePath+"/"+projectName+"/stage/otherstage/service/"+serviceName+"/resource"+webhookURI, 3)
	require.Nil(t, err)
	require.Equal(t, 200, resp.Response().StatusCode)

	t.Logf("Checking body of the received response")
	resource = models.Resource{}
	err = resp.ToJSON(&resource)
	require.Nil(t, err)
	rawDecoded, err = base64.StdEncoding.DecodeString(resource.ResourceContent)
	require.Nil(t, err)
	require.Equal(t, webhookConfig, string(rawDecoded))
}
