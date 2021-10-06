package go_tests

import (
	"github.com/stretchr/testify/require"
	"path"
	"testing"
)

const onboardServiceShipyard = `apiVersion: "spec.keptn.sh/0.2.0"
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

func Test_Continuous_Delivery(t *testing.T) {

	gitExamplesRepositoryURL := "https://github.com/keptn/examples"
	gitExamplesBranchName := "master"
	gitExampleRepositoryLocalDir, _ := CreateTmpDir()
	keptnProjectName := "sockshop"
	cartsServiceName := "carts"
	cartsChartLocalDir := path.Join(gitExampleRepositoryLocalDir, "onboarding-carts", "carts")
	cartsJmeterDir := path.Join(gitExampleRepositoryLocalDir, "onboarding-carts", "jmeter")
	cartsDBServiceName := "carts-db"
	cartsDBChartLocalDir := path.Join(gitExampleRepositoryLocalDir, "onboarding-carts", "carts-db")

	t.Logf("Cloning Keptn examples GIT repository %s from branch %s", gitExamplesRepositoryURL, gitExamplesBranchName)
	_, err := ExecuteCommandf("git clone --branch %s %s --single-branch %s", gitExamplesBranchName, gitExamplesRepositoryURL, gitExampleRepositoryLocalDir)
	require.Nil(t, err)

	t.Logf("Creating a new project %s without a GIT Upstream", keptnProjectName)
	shipyardFilePath, err := CreateTmpShipyardFile(onboardServiceShipyard)
	require.Nil(t, err)
	err = CreateProject(keptnProjectName, shipyardFilePath, true)
	require.Nil(t, err)

	t.Logf("Onboarding service %s in project %s with chart %s", cartsServiceName, keptnProjectName, cartsChartLocalDir)
	_, err = ExecuteCommandf("keptn onboard service %s --project %s --chart=%s", cartsServiceName, keptnProjectName, cartsChartLocalDir)
	require.Nil(t, err)

	t.Log("Adding functional test resources for jmeter")
	_, err = ExecuteCommandf("keptn add-resource --project=%s --service=%s --stage=%s --resource=%s --resourceUri=%s", keptnProjectName, cartsServiceName, "dev", cartsJmeterDir+"/basiccheck.jmx", "jmeter/basiccheck.jmx")
	require.Nil(t, err)

	t.Log("Adding performance test resources for jmeter")
	// Note: in order to speed up the tests we use basiccheck also for performance test
	_, err = ExecuteCommandf("keptn add-resource --project=%s --service=%s --stage=%s --resource=%s --resourceUri=%s", keptnProjectName, cartsServiceName, "staging", cartsJmeterDir+"/basiccheck.jmx", "jmeter/basiccheck.jmx")
	require.Nil(t, err)

	t.Logf("Onboarding service %s in project %s with chart %s", cartsDBServiceName, keptnProjectName, cartsDBChartLocalDir)
	_, err = ExecuteCommandf("keptn onboard service %s --project %s --chart=%s", cartsDBServiceName, keptnProjectName, cartsDBChartLocalDir)
	require.Nil(t, err)

	t.Logf("Trigger delivery of carts db")
	_, err = ExecuteCommandf("keptn trigger delivery --project=%s --service=%s --image=mongo:latest --sequence=delivery-direct", keptnProjectName, cartsDBServiceName)
	require.Nil(t, err)

	t.Logf("Trigger delivery of carts")
	_, err = ExecuteCommandf("keptn trigger delivery --project=%s --service=%s --image=%s --tag=%s --sequence=%s", keptnProjectName, cartsServiceName, "docker.io/keptnexamples/carts", "0.10.1", "delivery")
	require.Nil(t, err)

	t.Log("Verify direct delivery of carts db in stage dev")
	err = VerifyDirectDeployment(cartsDBServiceName, keptnProjectName, "dev", "mongo", "latest")
	require.Nil(t, err)

	t.Log("Verify Direct delivery of carts in stage dev")
	err = VerifyDirectDeployment(cartsServiceName, keptnProjectName, "dev", "docker.io/keptnexamples/carts", "0.10.1")
	require.Nil(t, err)

	t.Log("Verify direct delivery of carts db in stage staging")
	err = VerifyDirectDeployment(cartsDBServiceName, keptnProjectName, "staging", "mongo", "latest")
	require.Nil(t, err)

	t.Log("Verify delivery of carts in stage staging")
	err = VerifyBlueGreenDeployment(cartsServiceName, keptnProjectName, "staging", "docker.io/keptnexamples/carts", "0.10.1")
	require.Nil(t, err)

	t.Log("Verify direct delivery of carts db in stage prod-a")
	err = VerifyDirectDeployment(cartsDBServiceName, keptnProjectName, "prod-a", "mongo", "latest")
	require.Nil(t, err)

	t.Log("Verify delivery of carts in stage prod-a")
	err = VerifyBlueGreenDeployment(cartsServiceName, keptnProjectName, "prod-a", "docker.io/keptnexamples/carts", "0.10.1")
	require.Nil(t, err)

	t.Log("Verify direct delivery of carts db in stage prod-b")
	err = VerifyDirectDeployment(cartsDBServiceName, keptnProjectName, "prod-b", "mongo", "latest")
	require.Nil(t, err)

	t.Log("Verify delivery of carts in stage prod-a")
	err = VerifyBlueGreenDeployment(cartsServiceName, keptnProjectName, "prod-b", "docker.io/keptnexamples/carts", "0.10.1")
	require.Nil(t, err)

}
