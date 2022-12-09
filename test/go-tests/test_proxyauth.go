package go_tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/go-utils/pkg/common/osutils"
	"github.com/stretchr/testify/require"
)

const testingProxyShipyard = `apiVersion: "spec.keptn.sh/0.2.3"
kind: "Shipyard"
metadata:
  name: "shipyard-podtato-ohead"
spec:
  stages:
    - name: "dev"
      sequences:
        - name: "delivery"
          tasks:
            - name: "deployment"
              properties:
                deploymentstrategy: "direct"
            - name: "test"
              properties:
                teststrategy: "functional"
            - name: "evaluation"
            - name: "release"
        - name: "delivery-direct"
          tasks:
            - name: "deployment"
              properties:
                deploymentstrategy: "direct"
            - name: "release"

    - name: "prod"
      sequences:
        - name: "delivery"
          triggeredOn:
            - event: "dev.delivery.finished"
          tasks:
            - name: "deployment"
              properties:
                deploymentstrategy: "blue_green_service"
            - name: "test"
              properties:
                teststrategy: "performance"
            - name: "evaluation"
            - name: "release"
        - name: "rollback"
          triggeredOn:
            - event: "prod.delivery.finished"
              selector:
                match:
                  result: "fail"
          tasks:
            - name: "rollback"

        - name: "delivery-direct"
          triggeredOn:
            - event: "dev.delivery-direct.finished"
          tasks:
            - name: "deployment"
              properties:
                deploymentstrategy: "direct"
            - name: "release"
`
const baseProxyProjectPath = "/api/controlPlane/v1/project"

type Payload struct {
	HTTPRequest                  HTTPRequest                  `json:"httpRequest,omitempty"`
	HTTPOverrideForwardedRequest HTTPOverrideForwardedRequest `json:"httpOverrideForwardedRequest,omitempty"`
}

type Headers struct {
	Name   string   `json:"name,omitempty"`
	Values []string `json:"values,omitempty"`
}

type HTTPRequest struct {
	Method  string    `json:"method,omitempty"`
	Path    string    `json:"path,omitempty"`
	Headers []Headers `json:"headers,omitempty"`
}

type RequestOverride struct {
	Method  string    `json:"method,omitempty"`
	Path    string    `json:"path,omitempty"`
	Headers []Headers `json:"headers,omitempty"`
}

type HTTPOverrideForwardedRequest struct {
	RequestOverride RequestOverride `json:"requestOverride,omitempty"`
}

// UpdateMockserverConfig execute a call to the mockserver API to add a new configuration as proxy, it is equivalent to a curl command

func UpdateMockserverConfig(t *testing.T, project string) {
	namespace := osutils.GetOSEnvOrDefault(KeptnNamespaceEnvVar, DefaultKeptnNamespace)
	pods, err := GetPodNamesOfDeployment("app=mockserver")
	require.Nil(t, err)
	require.NotZero(t, len(pods))
	err = KubeCtlPortForwardSvc(context.Background(), pods[0], "1080", "1080", namespace)
	require.Nil(t, err)
	token, baseURL, err := GetApiCredentials()
	require.Nil(t, err)
	apiurl, err := url.Parse(baseURL)
	require.Nil(t, err)
	data := Payload{
		HTTPRequest: HTTPRequest{
			Path: baseProxyProjectPath + "/" + project,
		},
		HTTPOverrideForwardedRequest: HTTPOverrideForwardedRequest{
			RequestOverride: RequestOverride{
				Path: baseProxyProjectPath + "/" + project,
				Headers: []Headers{
					{Name: "Host",
						Values: []string{apiurl.Host},
					},
					{Name: "X-Token",
						Values: []string{token},
					},
				},
			},
		},
	}
	payloadBytes, err := json.Marshal(data)
	require.Nil(t, err)
	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("PUT", "http://localhost:1080/mockserver/expectation", body)
	require.Nil(t, err)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	require.Nil(t, err)
	defer resp.Body.Close()
}

func Test_ProxyAuth(t *testing.T) {
	projectName := "proxy-auth"
	serviceName := "helloservice"
	secondServiceName := "helloservice2"

	mockServerIP := "localhost:1080"

	defer func() {
		logs, err := PrintLogsWithDeploymentName("app.kubernetes.io/name=resource-service")
		require.Nil(t, err)
		t.Log("logs from RService: ")
		t.Log(logs)

		logsShippy, err := PrintLogsWithDeploymentName("app.kubernetes.io/name=shipyard-controller")
		require.Nil(t, err)
		t.Log("logs from Shippy: ")
		t.Log(logsShippy)
	}()

	t.Logf("Creating a new project %s with Gitea Upstream", projectName)
	shipyardFilePath, err := CreateTmpShipyardFile(testingProxyShipyard)
	require.Nil(t, err)
	projectName, err = CreateProjectWithProxy(projectName, shipyardFilePath, mockServerIP)
	require.Nil(t, err)

	t.Log("Adding new config to mockserver")
	UpdateMockserverConfig(t, projectName)

	t.Logf("Creating service %s in project %s", serviceName, projectName)
	_, err = ExecuteCommandf("keptn create service %s --project %s", serviceName, projectName)
	require.Nil(t, err)

	t.Logf("Getting project %s with a proxy", projectName)

	ApiCaller, err := NewAPICallerWithBaseURL("http://" + mockServerIP)
	resp, err := ApiCaller.Get(baseProxyProjectPath+"/"+projectName, 3)
	require.Nil(t, err)
	require.Equal(t, 200, resp.Response().StatusCode)

	t.Logf("Checking if upstream was provisioned")
	project := models.ExpandedProject{}
	err = resp.ToJSON(&project)
	require.Nil(t, err)
	require.NotNil(t, project.GitCredentials)
	require.NotNil(t, project.GitCredentials.HttpsAuth)
	require.NotNil(t, project.GitCredentials.HttpsAuth.Proxy)
	require.Equal(t, mockServerIP, project.GitCredentials.HttpsAuth.Proxy.URL)
	require.Equal(t, "http", project.GitCredentials.HttpsAuth.Proxy.Scheme)
	require.Equal(t, "", project.GitCredentials.HttpsAuth.Proxy.User)
	require.Equal(t, true, project.GitCredentials.HttpsAuth.InsecureSkipTLS)
	require.Equal(t, projectName, project.ProjectName)

	t.Logf("Updating project credentials")
	user := GetGiteaUser()
	token, err := GetGiteaToken()
	require.Nil(t, err)

	// apply the k8s job for creating the git upstream
	_, err = ExecuteCommand(fmt.Sprintf("keptn update project %s --git-remote-url=http://gitea-http:3000/%s/%s --git-user=%s --git-token=%s --git-proxy-url=%s --git-proxy-scheme=http --insecure-skip-tls", projectName, user, projectName, user, token, mockServerIP))
	require.Nil(t, err)

	t.Logf("Creating service %s in project %s", secondServiceName, projectName)
	_, err = ExecuteCommandf("keptn create service %s --project %s", secondServiceName, projectName)
	require.Nil(t, err)

	//Modify the proxy settings to be certain that no other project use the proxy
	_, err = ExecuteCommand(fmt.Sprintf("keptn update project %s --git-remote-url=http://gitea-http:3000/%s/%s --git-user=%s --git-token=%s --git-proxy-url=%s --git-proxy-scheme=http --insecure-skip-tls", projectName, user, projectName, user, token, mockServerIP))
	require.Nil(t, err)
}
