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

//const encodedAlpha = "YXBpVmVyc2lvbjogd2ViaG9va2NvbmZpZy5rZXB0bi5zaC92MWFscGhhMQpraW5kOiBXZWJob29rQ29uZmlnCm1ldGFkYXRhOgogIG5hbWU6IHdlYmhvb2stY29uZmlndXJhdGlvbgpzcGVjOgogIHdlYmhvb2tzOgogICAgLSB0eXBlOiAic2gua2VwdG4uZXZlbnQub3RoZXJ0YXNrLnRyaWdnZXJlZCIKICAgICAgc3Vic2NyaXB0aW9uSUQ6ICR7b3RoZXJ0YXNrLXN1Yi1pZH0KICAgICAgc2VuZEZpbmlzaGVkOiB0cnVlCiAgICAgIGVudkZyb206IAogICAgICAgIC0gbmFtZTogInNlY3JldEtleSIKICAgICAgICAgIHNlY3JldFJlZjoKICAgICAgICAgICAgbmFtZTogIm15LXdlYmhvb2stazhzLXNlY3JldCIKICAgICAgICAgICAga2V5OiAibXkta2V5IgogICAgICByZXF1ZXN0czoKICAgICAgICAtICJjdXJsIC0tZGF0YSAne1wiZW1haWxcIjpcInRlc3RAZXhhbXBsZS5jb21cIiwgXCJuYW1lXCI6IFtcIkJvb2xlYW5cIiwgXCJXb3JsZFwiXX0nIC1IIFwiQWNjZXB0LUNoYXJzZXQ6IHV0Zi04XCIgLUggJ0NvbnRlbnQtVHlwZTogYXBwbGljYXRpb24vanNvbicgaHR0cHM6Ly9odHRwYmluLm9yZy9wb3N0IC0tc29tZS1yYW5kb20tb3B0aW9ucyAtWVlZIgogICAgLSB0eXBlOiAic2gua2VwdG4uZXZlbnQuZmFpbGVkdGFzay50cmlnZ2VyZWQiCiAgICAgIHN1YnNjcmlwdGlvbklEOiAke2ZhaWxlZHRhc2stc3ViLWlkfQogICAgICBzZW5kRmluaXNoZWQ6IHRydWUKICAgICAgcmVxdWVzdHM6CiAgICAgICAgLSBjdXJsIGh0dHA6Ly9sb2NhbDo4MDgwIHt7LmRhdGEucHJvamVjdH19IHt7LmVudi5teXNlY3JldH19CiAgICAtIHR5cGU6ICJzaC5rZXB0bi5ldmVudC51bmFsbG93ZWR0YXNrLnRyaWdnZXJlZCIKICAgICAgc3Vic2NyaXB0aW9uSUQ6ICR7dW5hbGxvd2VkdGFzay1zdWItaWR9CiAgICAgIHNlbmRGaW5pc2hlZDogdHJ1ZQogICAgICBlbnZGcm9tOiAKICAgICAgICAtIG5hbWU6ICJzZWNyZXRLZXkiCiAgICAgICAgICBzZWNyZXRSZWY6CiAgICAgICAgICAgIG5hbWU6ICJteS13ZWJob29rLWs4cy1zZWNyZXQiCiAgICAgICAgICAgIGtleTogIm15LWtleSIKICAgICAgcmVxdWVzdHM6CiAgICAgICAgLSBjdXJsIGh0dHA6Ly9sb2NhbDo4MDgwIHt7LmRhdGEucHJvamVjdH19IHt7LmVudi5teXNlY3JldH19"
//const encodedBeta = "YXBpVmVyc2lvbjogd2ViaG9va2NvbmZpZy5rZXB0bi5zaC92MWJldGExCmtpbmQ6IFdlYmhvb2tDb25maWcKbWV0YWRhdGE6CiAgbmFtZTogd2ViaG9vay1jb25maWd1cmF0aW9uCnNwZWM6CiAgd2ViaG9va3M6CiAgICAtIHR5cGU6ICJzaC5rZXB0bi5ldmVudC5vdGhlcnRhc2sudHJpZ2dlcmVkIgogICAgICBzdWJzY3JpcHRpb25JRDogJHtvdGhlcnRhc2stc3ViLWlkfQogICAgICBzZW5kRmluaXNoZWQ6IHRydWUKICAgICAgZW52RnJvbTogCiAgICAgICAgLSBuYW1lOiAic2VjcmV0S2V5IgogICAgICAgICAgc2VjcmV0UmVmOgogICAgICAgICAgICBuYW1lOiAibXktd2ViaG9vay1rOHMtc2VjcmV0IgogICAgICAgICAgICBrZXk6ICJteS1rZXkiCiAgICAgIHJlcXVlc3RzOgogICAgICAgIC0gdXJsOiBodHRwczovL2h0dHBiaW4ub3JnL3Bvc3QKICAgICAgICAgIG1ldGhvZDogR0VUCiAgICAgICAgICBwYXlsb2FkOiAie1wiZW1haWxcIjpcInRlc3RAZXhhbXBsZS5jb21cIiwgXCJuYW1lXCI6IFtcIkJvb2xlYW5cIiwgXCJXb3JsZFwiXX0iCiAgICAgICAgICBvcHRpb25zOiAiLS1zb21lLXJhbmRvbS1vcHRpb25zIC1ZWVkiCiAgICAgICAgICBoZWFkZXJzOgogICAgICAgICAgICAtIGtleTogIkNvbnRlbnQtVHlwZSIKICAgICAgICAgICAgICB2YWx1ZTogImFwcGxpY2F0aW9uL2pzb24iCiAgICAgICAgICBoZWFkZXJzOgogICAgICAgICAgICAtIGtleTogIkFjY2VwdC1DaGFyc2V0IgogICAgICAgICAgICAgIHZhbHVlOiAidXRmLTgiCiAgICAtIHR5cGU6ICJzaC5rZXB0bi5ldmVudC5mYWlsZWR0YXNrLnRyaWdnZXJlZCIKICAgICAgc3Vic2NyaXB0aW9uSUQ6ICR7ZmFpbGVkdGFzay1zdWItaWR9CiAgICAgIHNlbmRGaW5pc2hlZDogdHJ1ZQogICAgICByZXF1ZXN0czoKICAgICAgICAtIHVybDogaHR0cDovL2xvY2FsOjgwODAge3suZGF0YS5wcm9qZWN0fX0ge3suZW52Lm15c2VjcmV0fX0KICAgICAgICAgIG1ldGhvZDogR0VUCiAgICAtIHR5cGU6ICJzaC5rZXB0bi5ldmVudC51bmFsbG93ZWR0YXNrLnRyaWdnZXJlZCIKICAgICAgc3Vic2NyaXB0aW9uSUQ6ICR7dW5hbGxvd2VkdGFzay1zdWItaWR9CiAgICAgIHNlbmRGaW5pc2hlZDogdHJ1ZQogICAgICBlbnZGcm9tOiAKICAgICAgICAtIG5hbWU6ICJzZWNyZXRLZXkiCiAgICAgICAgICBzZWNyZXRSZWY6CiAgICAgICAgICAgIG5hbWU6ICJteS13ZWJob29rLWs4cy1zZWNyZXQiCiAgICAgICAgICAgIGtleTogIm15LWtleSIKICAgICAgcmVxdWVzdHM6CiAgICAgICAgLSB1cmw6IGh0dHA6Ly9sb2NhbDo4MDgwIHt7LmRhdGEucHJvamVjdH19IHt7LmVudi5teXNlY3JldH19CiAgICAgICAgICBtZXRob2Q6IEdFVA=="

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

	for i := 1; i <= 1; i++ {
		fillProjectWithWebhooks(t, shipyardFilePath, webhookFilePath, projectName+fmt.Sprint(i), serviceName+fmt.Sprint(i))
		checkWebhooksHaveCorrectWersion(t, internalKeptnAPI, "keptn-"+projectName+fmt.Sprint(i), serviceName+fmt.Sprint(i), webhookConfigMigrationAlpha)
	}

	t.Logf("Executing dry-run for project %s", "keptn-"+projectName+"1")
	_, err = ExecuteCommandf("./../../cli/cli migrate-webhooks")
	require.Nil(t, err)

	for i := 1; i <= 1; i++ {
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
