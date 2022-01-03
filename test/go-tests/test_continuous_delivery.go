package go_tests

import (
	"path"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const onboardServiceShipyard = `apiVersion: "spec.keptn.sh/0.2.3"
kind: "Shipyard"
metadata:
  name: "shipyard-sockshop"
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

    - name: "staging"
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
            - event: "staging.delivery.finished"
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

    - name: "prod-a"
      sequences:
        - name: "delivery"
          triggeredOn:
            - event: "staging.delivery.finished"
          tasks:
            - name: "deployment"
              properties:
                deploymentstrategy: "blue_green_service"
            - name: "release"
        - name: "rollback"
          triggeredOn:
            - event: "prod-a.delivery.finished"
              selector:
                match:
                  result: "fail"
          tasks:
            - name: "rollback"
        - name: "delivery-direct"
          triggeredOn:
            - event: "staging.delivery-direct.finished"
          tasks:
            - name: "deployment"
              properties:
                deploymentstrategy: "direct"
            - name: "release"

    - name: "prod-b"
      sequences:
        - name: "delivery"
          triggeredOn:
            - event: "staging.delivery.finished"
          tasks:
            - name: "deployment"
              properties:
                deploymentstrategy: "blue_green_service"
            - name: "release"
        - name: "rollback"
          triggeredOn:
            - event: "prod-b.delivery.finished"
              selector:
                match:
                  result: "fail"
          tasks:
            - name: "rollback"
        - name: "delivery-direct"
          triggeredOn:
            - event: "staging.delivery-direct.finished"
          tasks:
            - name: "deployment"
              properties:
                deploymentstrategy: "direct"
            - name: "release"
`

func TestContinuousDelivery(t *testing.T) {

	repoLocalDir := "../assets/podtato-head"
	keptnProjectName := "podtato-head"
	serviceName := "helloservice"
	serviceChartLocalDir := path.Join(repoLocalDir, "helm-charts", "helloservice.tgz")
	serviceJmeterDir := path.Join(repoLocalDir, "jmeter")
	serviceHealthCheckEndpoint := "/metrics"
	const artifactTagInitial = "v0.1.0"
	const artifactTagUpdated = "v0.1.1"
	const artifactImage = "ghcr.io/podtato-head/podtatoserver"
	const prodStageA = "prod-a"
	const prodStageB = "prod-b"

	t.Logf("Creating a new project %s without a GIT Upstream", keptnProjectName)
	shipyardFilePath, err := CreateTmpShipyardFile(onboardServiceShipyard)
	require.Nil(t, err)
	err = CreateProject(keptnProjectName, shipyardFilePath, true)
	require.Nil(t, err)

	t.Logf("Creating service %s in project %s", serviceName, keptnProjectName)
	_, err = ExecuteCommandf("keptn create service %s --project %s", serviceName, keptnProjectName)
	require.Nil(t, err)

	t.Logf("Adding resource for service %s in project %s", serviceName, keptnProjectName)
	_, err = ExecuteCommandf("keptn add-resource --project %s --service=%s --all-stages --resource=%s --resourceUri=%s", keptnProjectName, serviceName, serviceChartLocalDir, "helm/helloservice.tgz")
	require.Nil(t, err)

	t.Log("Adding jmeter config in staging")
	_, err = ExecuteCommandf("keptn add-resource --project=%s --service=%s --stage=%s --resource=%s --resourceUri=%s", keptnProjectName, serviceName, "staging", serviceJmeterDir+"/jmeter.conf.yaml", "jmeter/jmeter.conf.yaml")
	require.Nil(t, err)

	t.Log("Adding load test resources for jmeter in staging")
	_, err = ExecuteCommandf("keptn add-resource --project=%s --service=%s --stage=%s --resource=%s --resourceUri=%s", keptnProjectName, serviceName, "staging", serviceJmeterDir+"/load.jmx", "jmeter/load.jmx")
	require.Nil(t, err)

	///////////////////////////////////////
	// Deploy v0.1.0
	///////////////////////////////////////

	t.Logf("Trigger delivery of helloservice:v0.1.0")
	_, err = ExecuteCommandf("keptn trigger delivery --project=%s --service=%s --image=%s --tag=%s --sequence=%s", keptnProjectName, serviceName, artifactImage, artifactTagInitial, "delivery")
	require.Nil(t, err)

	t.Logf("Sleeping for 60s...")
	time.Sleep(60 * time.Second)
	t.Logf("Continue to work...")

	t.Log("Verify Direct delivery of helloservice in stage dev")
	err = VerifyDirectDeployment(serviceName, keptnProjectName, "dev", artifactImage, artifactTagInitial)
	require.Nil(t, err)

	t.Log("Verify network access to public URI of helloservice in stage dev")
	cartPubURL, err := GetPublicURLOfService(serviceName, keptnProjectName, "dev")
	require.Nil(t, err)
	err = WaitForURL(cartPubURL+serviceHealthCheckEndpoint, time.Minute)
	require.Nil(t, err)

	t.Log("Verify delivery of helloservice:v0.1.0 in stage staging")
	err = VerifyBlueGreenDeployment(serviceName, keptnProjectName, "staging", artifactImage, artifactTagInitial)
	require.Nil(t, err)

	t.Log("Verify network access to public URI of helloservice in stage staging")
	cartPubURL, err = GetPublicURLOfService(serviceName, keptnProjectName, "staging")
	require.Nil(t, err)
	err = WaitForURL(cartPubURL+serviceHealthCheckEndpoint, time.Minute)
	require.Nil(t, err)

	t.Log("Verify delivery of helloservice:v0.1.0 in stage prod-a")
	err = VerifyBlueGreenDeployment(serviceName, keptnProjectName, prodStageA, artifactImage, artifactTagInitial)
	require.Nil(t, err)

	t.Log("Verify network access to public URI of helloservice in stage prod-a")
	cartPubURL, err = GetPublicURLOfService(serviceName, keptnProjectName, prodStageA)
	require.Nil(t, err)
	err = WaitForURL(cartPubURL+serviceHealthCheckEndpoint, time.Minute)
	require.Nil(t, err)

	t.Log("Verify delivery of helloservice:v0.1.0 in stage prod-b")
	err = VerifyBlueGreenDeployment(serviceName, keptnProjectName, prodStageB, artifactImage, artifactTagInitial)
	require.Nil(t, err)

	t.Log("Verify network access to public URI of helloservice in stage prod-b")
	cartPubURL, err = GetPublicURLOfService(serviceName, keptnProjectName, prodStageB)
	require.Nil(t, err)
	err = WaitForURL(cartPubURL+serviceHealthCheckEndpoint, time.Minute)
	require.Nil(t, err)

	///////////////////////////////////////
	// Deploy v0.1.1
	///////////////////////////////////////

	t.Logf("Trigger delivery of helloservice:v0.1.1")
	_, err = ExecuteCommandf("keptn trigger delivery --project=%s --service=%s --image=%s --tag=%s --sequence=%s", keptnProjectName, serviceName, artifactImage, artifactTagUpdated, "delivery")
	require.Nil(t, err)

	t.Logf("Sleeping for 60s...")
	time.Sleep(60 * time.Second)
	t.Logf("Continue to work...")

	t.Log("Verify Direct delivery of helloservice in stage dev")
	err = VerifyDirectDeployment(serviceName, keptnProjectName, "dev", artifactImage, artifactTagUpdated)
	require.Nil(t, err)

	t.Log("Verify network access to public URI of helloservice in stage dev")
	cartPubURL, err = GetPublicURLOfService(serviceName, keptnProjectName, "dev")
	require.Nil(t, err)
	err = WaitForURL(cartPubURL+serviceHealthCheckEndpoint, time.Minute)
	require.Nil(t, err)

	t.Log("Verify delivery of helloservice:v0.1.1 in stage staging")
	err = VerifyBlueGreenDeployment(serviceName, keptnProjectName, "staging", artifactImage, artifactTagUpdated)
	require.Nil(t, err)

	t.Log("Verify network access to public URI of helloservice in stage staging")
	cartPubURL, err = GetPublicURLOfService(serviceName, keptnProjectName, "staging")
	require.Nil(t, err)
	err = WaitForURL(cartPubURL+serviceHealthCheckEndpoint, time.Minute)
	require.Nil(t, err)

	t.Log("Verify delivery of helloservice:v0.1.1 in stage prod-a")
	err = VerifyBlueGreenDeployment(serviceName, keptnProjectName, prodStageA, artifactImage, artifactTagUpdated)
	require.Nil(t, err)

	t.Log("Verify network access to public URI of helloservice in stage prod-a")
	cartPubURL, err = GetPublicURLOfService(serviceName, keptnProjectName, prodStageA)
	require.Nil(t, err)
	err = WaitForURL(cartPubURL+serviceHealthCheckEndpoint, time.Minute)
	require.Nil(t, err)

	t.Log("Verify delivery of helloservice:v0.1.1 in stage prod-b")
	err = VerifyBlueGreenDeployment(serviceName, keptnProjectName, prodStageB, artifactImage, artifactTagUpdated)
	require.Nil(t, err)

	t.Log("Verify network access to public URI of helloservice in stage prod-b")
	cartPubURL, err = GetPublicURLOfService(serviceName, keptnProjectName, prodStageB)
	require.Nil(t, err)
	err = WaitForURL(cartPubURL+serviceHealthCheckEndpoint, time.Minute)
	require.Nil(t, err)
}
