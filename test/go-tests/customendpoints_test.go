package go_tests

import (
	"context"
	"fmt"
	"github.com/imroc/req"
	"github.com/keptn/go-utils/pkg/api/models"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	keptnkubeutils "github.com/keptn/kubernetes-utils/pkg"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"
)

const customEndpointShipyard = `apiVersion: "spec.keptn.sh/0.2.0"
kind: "Shipyard"
metadata:
  name: "shipyard"
spec:
  stages:
    - name: "dev"
      sequences:
        - name: "delivery"
          tasks:
            - name: "deployment"
              properties:
                deploymentstrategy: "user_managed"`

const customEndpoints = `deploymentURIsLocal:
  - "http://my-local-url:80"
deploymentURIsPublic:
  - "${PUBLIC_URL}"`

func Test_CustomUserManagedEndpointsTest(t *testing.T) {
	projectName := "user-managed"
	serviceName := "nginx"
	stageName := "dev"
	sequenceName := "delivery"
	serviceChartPath := "https://charts.bitnami.com/bitnami/nginx-8.9.0.tgz"

	shipyardFilePath, err := CreateTmpShipyardFile(customEndpointShipyard)

	require.Nil(t, err)
	defer os.Remove(shipyardFilePath)

	_, err = ExecuteCommand(fmt.Sprintf("wget %s -O chart.tgz", serviceChartPath))
	require.Nil(t, err)

	defer os.Remove("chart.tgz")

	// make sure the namespace from a previous test run has been deleted properly
	exists, err := keptnkubeutils.ExistsNamespace(false, projectName+"-dev")
	if exists {
		t.Logf("Deleting namespace %s-dev from previous test execution", projectName)
		clientset, err := keptnkubeutils.GetClientset(false)
		require.Nil(t, err)
		err = clientset.CoreV1().Namespaces().Delete(context.TODO(), projectName+"-dev", v1.DeleteOptions{})
		require.Nil(t, err)
	}

	require.Eventually(t, func() bool {
		t.Logf("Checking if namespace %s-dev is still there", projectName)
		exists, err := keptnkubeutils.ExistsNamespace(false, projectName+"-dev")
		if err != nil || exists {
			t.Logf("Namespace %s-dev is still there", projectName)
			return false
		}
		t.Logf("Namespace %s-dev has been removed - proceeding with test execution", projectName)
		return true
	}, 60*time.Second, 5*time.Second)

	// check if the project is already available - if not, delete it before creating it again
	err = CreateProject(projectName, shipyardFilePath, true)
	require.Nil(t, err)

	// create the service
	t.Logf("Creating service %s in project %s", serviceName, projectName)
	output, err := ExecuteCommand(fmt.Sprintf("keptn create service %s --project=%s", serviceName, projectName))
	require.Nil(t, err)
	require.Contains(t, output, "created successfully")

	// upload the service's helm chart
	t.Logf("Uploading the helm chart of service %s in project %s", serviceName, projectName)
	_, err = ExecuteCommand(fmt.Sprintf("keptn add-resource --service=%s --project=%s --all-stages --resource=./chart.tgz --resourceUri=helm/%s.tgz", serviceName, projectName, serviceName))
	require.Nil(t, err)

	// trigger the sequence without defining custom endpoints first
	t.Logf("Triggering the first delivery sequence without providing custom endpoints")
	keptnContextID, err := TriggerSequence(projectName, serviceName, stageName, sequenceName, nil)
	require.Nil(t, err)
	require.NotEmpty(t, keptnContextID)

	// wait until we get a deployment.finished event
	var deploymentFinishedEvent *models.KeptnContextExtendedCE
	t.Log("Waiting for deployment to complete")
	require.Eventually(t, func() bool {
		deploymentFinishedEvent, err = GetLatestEventOfType(keptnContextID, projectName, stageName, keptnv2.GetFinishedEventType(keptnv2.DeploymentTaskName))
		if err != nil || deploymentFinishedEvent == nil {
			t.Log("Deployment has not been completed yet... Waiting a couple of seconds before checking again")
			return false
		}
		return true
	}, 60*time.Second, 5*time.Second)
	t.Log("Deployment has been completed")

	deploymentFinishedEventData := &keptnv2.DeploymentFinishedEventData{}
	err = keptnv2.EventDataAs(*deploymentFinishedEvent, deploymentFinishedEventData)

	// if no custom deployment URIs have been defined, they should also be nil in the deployment.finished event
	require.Nil(t, err)
	require.Nil(t, deploymentFinishedEventData.Deployment.DeploymentURIsPublic)
	require.Nil(t, deploymentFinishedEventData.Deployment.DeploymentURIsLocal)

	// get the LoadBalancer endpoint of the deployed service so we can define its URL in the next delivery
	serviceEndpoint, err := keptnkubeutils.GetKeptnEndpointFromService(false, projectName+"-dev", projectName+"-dev-"+serviceName)
	require.Nil(t, err)
	require.NotEmpty(t, serviceEndpoint)

	publicURL := fmt.Sprintf("http://%s", serviceEndpoint)
	// create endpoints.yaml containing the IP of the service
	endpointsFileContent := strings.Replace(customEndpoints, "${PUBLIC_URL}", publicURL, 1)
	endpointsFilePath, err := CreateTmpFile("endpoints-*.yaml", endpointsFileContent)

	require.Nil(t, err)
	defer os.Remove(endpointsFilePath)

	// now, let's add an endpoints.yaml file to our service
	t.Log("Adding endpoints.yaml to our service")
	_, err = ExecuteCommand(fmt.Sprintf("keptn add-resource --project=%s --service=%s --resource=%s --resourceUri=helm/endpoints.yaml --all-stages", projectName, serviceName, endpointsFilePath))

	require.Nil(t, err)

	// trigger the sequence again
	t.Logf("Triggering the delivery sequence again - this time with custom endpoints")
	keptnContextID, err = TriggerSequence(projectName, serviceName, stageName, sequenceName, nil)
	require.Nil(t, err)
	require.NotEmpty(t, keptnContextID)

	// wait until we get a deployment.finished event
	t.Log("Waiting for deployment to complete")
	require.Eventually(t, func() bool {
		deploymentFinishedEvent, err = GetLatestEventOfType(keptnContextID, projectName, stageName, keptnv2.GetFinishedEventType(keptnv2.DeploymentTaskName))
		if err != nil || deploymentFinishedEvent == nil {
			t.Log("Deployment has not been completed yet... Waiting a couple of seconds before checking again")
			return false
		}
		return true
	}, 60*time.Second, 5*time.Second)
	t.Log("Deployment has been completed")

	err = keptnv2.EventDataAs(*deploymentFinishedEvent, deploymentFinishedEventData)

	t.Log("Verifying if deploymentURIsLocal and deploymentURIsPublic have been set properly")
	// if no custom deployment URIs have been defined, they should also be nil in the deployment.finished event
	require.Nil(t, err)
	require.Equal(t, []string{publicURL}, deploymentFinishedEventData.Deployment.DeploymentURIsPublic)
	require.Equal(t, []string{"http://my-local-url:80"}, deploymentFinishedEventData.Deployment.DeploymentURIsLocal)

	t.Log("deploymentURIsLocal and deploymentURIsPublic have been set properly")

	// check if the service has been deployed (i.e. is reachable)
	t.Log("Checking reachability of service")
	resp, err := req.Get(publicURL)
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.Response().StatusCode)
	t.Log("Service is reachable")
}
