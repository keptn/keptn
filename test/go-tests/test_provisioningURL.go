package go_tests

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/go-utils/pkg/common/retry"
	"github.com/stretchr/testify/require"
)

const provisioningShipyard = `YXBpVmVyc2lvbjogInNwZWMua2VwdG4uc2gvMC4yLjMiCmtpbmQ6ICJTaGlweWFyZCIKbWV0YWRhdGE6CiAgbmFtZTogInNoaXB5YXJkLXBvZHRhdG8tb2hlYWQiCnNwZWM6CiAgc3RhZ2VzOgogICAgLSBuYW1lOiAiZGV2IgogICAgICBzZXF1ZW5jZXM6CiAgICAgICAgLSBuYW1lOiAiZGVsaXZlcnkiCiAgICAgICAgICB0YXNrczoKICAgICAgICAgICAgLSBuYW1lOiAiZGVwbG95bWVudCIKICAgICAgICAgICAgICBwcm9wZXJ0aWVzOgogICAgICAgICAgICAgICAgZGVwbG95bWVudHN0cmF0ZWd5OiAiZGlyZWN0IgogICAgICAgICAgICAtIG5hbWU6ICJ0ZXN0IgogICAgICAgICAgICAgIHByb3BlcnRpZXM6CiAgICAgICAgICAgICAgICB0ZXN0c3RyYXRlZ3k6ICJmdW5jdGlvbmFsIgogICAgICAgICAgICAtIG5hbWU6ICJldmFsdWF0aW9uIgogICAgICAgICAgICAtIG5hbWU6ICJyZWxlYXNlIgogICAgICAgIC0gbmFtZTogImRlbGl2ZXJ5LWRpcmVjdCIKICAgICAgICAgIHRhc2tzOgogICAgICAgICAgICAtIG5hbWU6ICJkZXBsb3ltZW50IgogICAgICAgICAgICAgIHByb3BlcnRpZXM6CiAgICAgICAgICAgICAgICBkZXBsb3ltZW50c3RyYXRlZ3k6ICJkaXJlY3QiCiAgICAgICAgICAgIC0gbmFtZTogInJlbGVhc2UiCgogICAgLSBuYW1lOiAicHJvZCIKICAgICAgc2VxdWVuY2VzOgogICAgICAgIC0gbmFtZTogImRlbGl2ZXJ5IgogICAgICAgICAgdHJpZ2dlcmVkT246CiAgICAgICAgICAgIC0gZXZlbnQ6ICJkZXYuZGVsaXZlcnkuZmluaXNoZWQiCiAgICAgICAgICB0YXNrczoKICAgICAgICAgICAgLSBuYW1lOiAiZGVwbG95bWVudCIKICAgICAgICAgICAgICBwcm9wZXJ0aWVzOgogICAgICAgICAgICAgICAgZGVwbG95bWVudHN0cmF0ZWd5OiAiYmx1ZV9ncmVlbl9zZXJ2aWNlIgogICAgICAgICAgICAtIG5hbWU6ICJ0ZXN0IgogICAgICAgICAgICAgIHByb3BlcnRpZXM6CiAgICAgICAgICAgICAgICB0ZXN0c3RyYXRlZ3k6ICJwZXJmb3JtYW5jZSIKICAgICAgICAgICAgLSBuYW1lOiAiZXZhbHVhdGlvbiIKICAgICAgICAgICAgLSBuYW1lOiAicmVsZWFzZSIKICAgICAgICAtIG5hbWU6ICJyb2xsYmFjayIKICAgICAgICAgIHRyaWdnZXJlZE9uOgogICAgICAgICAgICAtIGV2ZW50OiAicHJvZC5kZWxpdmVyeS5maW5pc2hlZCIKICAgICAgICAgICAgICBzZWxlY3RvcjoKICAgICAgICAgICAgICAgIG1hdGNoOgogICAgICAgICAgICAgICAgICByZXN1bHQ6ICJmYWlsIgogICAgICAgICAgdGFza3M6CiAgICAgICAgICAgIC0gbmFtZTogInJvbGxiYWNrIgoKICAgICAgICAtIG5hbWU6ICJkZWxpdmVyeS1kaXJlY3QiCiAgICAgICAgICB0cmlnZ2VyZWRPbjoKICAgICAgICAgICAgLSBldmVudDogImRldi5kZWxpdmVyeS1kaXJlY3QuZmluaXNoZWQiCiAgICAgICAgICB0YXNrczoKICAgICAgICAgICAgLSBuYW1lOiAiZGVwbG95bWVudCIKICAgICAgICAgICAgICBwcm9wZXJ0aWVzOgogICAgICAgICAgICAgICAgZGVwbG95bWVudHN0cmF0ZWd5OiAiZGlyZWN0IgogICAgICAgICAgICAtIG5hbWU6ICJyZWxlYXNlIg==`
const baseProjectPath = "/controlPlane/v1/project"

var provisioningConfigMapTemplate = `apiVersion: v1
kind: ConfigMap
metadata:
  name: mockserver-config
data:
  initializerJson.json: |-
    [
      {
        "httpRequest": {
          "path": "/repository",
          "method": "POST"
        },
        "httpResponse": {
          "body": {
            "gitRemoteURL": "http://gitea-http:3000/%s/url-provisioning",
            "gitToken":     "%s",
            "gitUser":      "%s"
          },
          "statusCode": 201
        }
      },
      {
        "httpRequest": {
          "path": "/repository",
          "method": "DELETE"
        },
        "httpResponse": {
          "body": "Repository deleted",
          "statusCode": 200
        }
      },
    ]
  mockserver.properties: |-
    ###############################
    # MockServer & Proxy Settings #
    ###############################
    # Socket & Port Settings
    # socket timeout in milliseconds (default 120000)
    mockserver.maxSocketTimeout=120000
    # Certificate Generation
    # dynamically generated CA key pair (if they don't already exist in
    specified directory)
    mockserver.dynamicallyCreateCertificateAuthorityCertificate=true
    # save dynamically generated CA key pair in working directory
    mockserver.directoryToSaveDynamicSSLCertificate=.
    # certificate domain name (default "localhost")
    mockserver.sslCertificateDomainName=localhost
    # comma separated list of ip addresses for Subject Alternative Name domain
    names (default empty list)
    mockserver.sslSubjectAlternativeNameDomains=www.example.com,www.another.com
    # comma separated list of ip addresses for Subject Alternative Name ips
    (default empty list)
    mockserver.sslSubjectAlternativeNameIps=127.0.0.1
    # CORS
    # enable CORS for MockServer REST API
    mockserver.enableCORSForAPI=true
    # enable CORS for all responses
    mockserver.enableCORSForAllResponses=true
    # Json Initialization
    mockserver.initializationJsonPath=/config/initializerJson.json
`

/* @Description: Test_ProvisioningURL tests the behaviour of the system, when no git credentials are submitted by user as
 * they should be provisioned by a 3rd partly REST API when AUTOMATIC_PROVISIONING_URL env variable is set. We are using
 * a simple mockserver to mock the responses from the 3rd partly REST API, which responses with a valid git credentials.
 * @Outcome Creating and deleting the project is tested using the provisoned git credentials
 * @Issue 7149
 */
func Test_ProvisioningURL(t *testing.T) {
	projectName := "url-provisioning"
	mockserverConfigFileName := "mockserver-config.yaml"
	keptnNamespace := GetKeptnNameSpaceFromEnv()
	mockServerIP := "http://mockserver:1080"
	user := GetGiteaUser()
	url := fmt.Sprintf("http://gitea-http:3000/%s/%s", user, projectName)
	token, _ := GetGiteaToken()
	provisioningConfigMap := fmt.Sprintf(provisioningConfigMapTemplate, user, token, user)
	shipyardPod := "shipyard-controller"

	mockserverconfigFilePath, err := CreateTmpFile(mockserverConfigFileName, provisioningConfigMap)
	defer func() {
		if err := os.Remove(mockserverconfigFilePath); err != nil {
			t.Logf("Could not delete file: %v", err)
		}
	}()

	t.Logf("Create mock server ConfigMap")
	_, err = ExecuteCommandf("kubectl apply -f %s -n %s", mockserverconfigFilePath, keptnNamespace)
	require.Nil(t, err)

	t.Logf("Restarting mockserver deployment")
	_, err = ExecuteCommandf("kubectl rollout restart deployment mockserver -n %s", keptnNamespace)
	require.Nil(t, err)

	t.Logf("Setting up AUTOMATIC_PROVISIONING_URL env variable")
	t.Logf("External mockserver IP address: %s", mockServerIP)
	_, err = ExecuteCommandf("kubectl set env deployment/shipyard-controller AUTOMATIC_PROVISIONING_URL=%s -n %s", mockServerIP, keptnNamespace)
	require.Nil(t, err)

	t.Logf("Sleeping for 120s to make sure shipyard pod is deleted...")
	// due to termination grace period this takes quite some time now...
	time.Sleep(120 * time.Second)

	t.Logf("Checking if mockserver is running")
	err = WaitForPodOfDeployment("mockserver")
	require.Nil(t, err)

	//kubectl set kills shipyard pod, instead of sleeping we wait for the deployment to be ready
	t.Logf("Waiting for Shipyard-controller to be running")
	err = WaitForPodOfDeployment(shipyardPod)
	require.Nil(t, err)

	t.Logf("Creating a new upstream repository for project %s", projectName)
	err = retry.Retry(func() error {
		if err := RecreateProjectUpstream(projectName); err != nil {
			return err
		}
		return nil
	}, retry.NumberOfRetries(10), retry.DelayBetweenRetries(10*time.Second))
	require.Nil(t, err)

	t.Logf("Creating a new project %s with a provisioned Gitea Upstream", projectName)
	projectParams := map[string]string{
		"name":     projectName,
		"shipyard": provisioningShipyard,
	}
	createProjectRequestData, err := json.Marshal(projectParams)
	require.Nil(t, err)

	resp, err := ApiPOSTRequest(baseProjectPath, createProjectRequestData, 3)
	require.Nil(t, err)
	if resp.Response().StatusCode != 201 {
		t.Logf("%+v %s", resp.Response().Body, resp.Response().Body)
	}
	require.Equal(t, 201, resp.Response().StatusCode)

	t.Logf("Getting project %s with a provisioned Gitea Upstream", projectName)
	resp, err = ApiGETRequest(baseProjectPath+"/"+projectName, 3)
	require.Nil(t, err)
	require.Equal(t, 200, resp.Response().StatusCode)

	t.Logf("Checking if upstream was provisioned")
	project := models.ExpandedProject{}
	err = resp.ToJSON(&project)
	require.Nil(t, err)
	require.Equal(t, user, project.GitCredentials.User)
	require.Equal(t, url, project.GitCredentials.RemoteURL)

	t.Logf("Deleting project %s with a provisioned Gitea Upstream", projectName)
	resp, err = ApiDELETERequest(baseProjectPath+"/"+projectName, 3)
	require.Nil(t, err)
	require.Equal(t, 200, resp.Response().StatusCode)

	t.Logf("Unsetting AUTOMATIC_PROVISIONING_URL env variable")
	_, err = ExecuteCommandf("kubectl set env deployment/shipyard-controller AUTOMATIC_PROVISIONING_URL=%s -n %s", "", keptnNamespace)
	require.Nil(t, err)

	t.Logf("Sleeping for 30s to make sure shipyard pod is deleted...")
	time.Sleep(30 * time.Second)
	t.Logf("Waiting for Shipyard-controller to be running")
	err = WaitForPodOfDeployment(shipyardPod)
	require.Nil(t, err)

	t.Log("Deleting mockserver ConfigMap")
	_, err = ExecuteCommandf("kubectl delete configmap mockserver-config -n %s", keptnNamespace)
	require.Nil(t, err)
}

func Test_ProvisioningURL_hiddenURL(t *testing.T) {
	projectName := "url-provisioning"
	mockserverConfigFileName := "mockserver-config.yaml"
	keptnNamespace := GetKeptnNameSpaceFromEnv()
	mockServerIP := "http://mockserver:1080"
	user := GetGiteaUser()
	token, _ := GetGiteaToken()
	provisioningConfigMap := fmt.Sprintf(provisioningConfigMapTemplate, user, token, user)
	shipyardPod := "shipyard-controller"

	mockserverconfigFilePath, err := CreateTmpFile(mockserverConfigFileName, provisioningConfigMap)
	defer func() {
		if err := os.Remove(mockserverconfigFilePath); err != nil {
			t.Logf("Could not delete file: %v", err)
		}
	}()

	t.Logf("Create mock server ConfigMap")
	_, err = ExecuteCommandf("kubectl apply -f %s -n %s", mockserverconfigFilePath, keptnNamespace)
	require.Nil(t, err)

	t.Logf("Restarting mockserver deployment")
	_, err = ExecuteCommandf("kubectl rollout restart deployment mockserver -n %s", keptnNamespace)
	require.Nil(t, err)

	t.Logf("Setting up AUTOMATIC_PROVISIONING_URL env variable")
	t.Logf("External mockserver IP address: %s", mockServerIP)
	_, err = ExecuteCommandf("kubectl set env deployment/shipyard-controller AUTOMATIC_PROVISIONING_URL=%s -n %s", mockServerIP, keptnNamespace)
	require.Nil(t, err)

	t.Logf("Setting up HIDE_AUTOMATIC_PROVISIONED_URL env variable to true")
	t.Logf("External mockserver IP address: %s", mockServerIP)
	_, err = ExecuteCommandf("kubectl set env deployment/shipyard-controller HIDE_AUTOMATIC_PROVISIONED_URL=%s -n %s", "true", keptnNamespace)
	require.Nil(t, err)

	t.Logf("Sleeping for 120s to make sure shipyard pod is deleted...")
	// due to termination grace period this takes quite some time now...
	time.Sleep(120 * time.Second)

	t.Logf("Checking if mockserver is running")
	err = WaitForPodOfDeployment("mockserver")
	require.Nil(t, err)

	//kubectl set kills shipyard pod, instead of sleeping we wait for the deployment to be ready
	t.Logf("Waiting for Shipyard-controller to be running")
	err = WaitForPodOfDeployment(shipyardPod)
	require.Nil(t, err)

	t.Logf("Creating a new upstream repository for project %s", projectName)
	err = retry.Retry(func() error {
		if err := RecreateProjectUpstream(projectName); err != nil {
			return err
		}
		return nil
	}, retry.NumberOfRetries(10), retry.DelayBetweenRetries(10*time.Second))
	require.Nil(t, err)

	t.Logf("Creating a new project %s with a provisioned Gitea Upstream", projectName)
	projectParams := map[string]string{
		"name":     projectName,
		"shipyard": provisioningShipyard,
	}
	createProjectRequestData, err := json.Marshal(projectParams)
	require.Nil(t, err)

	resp, err := ApiPOSTRequest(baseProjectPath, createProjectRequestData, 3)
	require.Nil(t, err)
	require.Equal(t, 201, resp.Response().StatusCode)

	t.Logf("Getting project %s with a provisioned Gitea Upstream", projectName)
	resp, err = ApiGETRequest(baseProjectPath+"/"+projectName, 3)
	require.Nil(t, err)
	require.Equal(t, 200, resp.Response().StatusCode)

	t.Logf("Checking if provisioned upstream was hidden")
	project := models.ExpandedProject{}
	err = resp.ToJSON(&project)
	require.Nil(t, err)
	require.Equal(t, "", project.GitCredentials.User)
	require.Equal(t, "", project.GitCredentials.RemoteURL)

	t.Logf("Deleting project %s with a provisioned Gitea Upstream", projectName)
	resp, err = ApiDELETERequest(baseProjectPath+"/"+projectName, 3)
	require.Nil(t, err)
	require.Equal(t, 200, resp.Response().StatusCode)

	t.Logf("Unsetting AUTOMATIC_PROVISIONING_URL env variable")
	_, err = ExecuteCommandf("kubectl set env deployment/shipyard-controller AUTOMATIC_PROVISIONING_URL=%s -n %s", "", keptnNamespace)
	require.Nil(t, err)

	t.Logf("Unsetting HIDE_AUTOMATIC_PROVISIONED_URL env variable to false")
	_, err = ExecuteCommandf("kubectl set env deployment/shipyard-controller HIDE_AUTOMATIC_PROVISIONED_URL=%s -n %s", "false", keptnNamespace)
	require.Nil(t, err)

	t.Logf("Sleeping for 30s to make sure shipyard pod is deleted...")
	time.Sleep(30 * time.Second)
	t.Logf("Waiting for Shipyard-controller to be running")
	err = WaitForPodOfDeployment(shipyardPod)
	require.Nil(t, err)

	t.Log("Deleting mockserver ConfigMap")
	_, err = ExecuteCommandf("kubectl delete configmap mockserver-config -n %s", keptnNamespace)
	require.Nil(t, err)
}
