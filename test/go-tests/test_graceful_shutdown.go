package go_tests

import (
	"github.com/stretchr/testify/require"
	"path/filepath"
	"testing"
	"time"
)

const tinyShipyard = `apiVersion: "spec.keptn.sh/0.2.3"
kind: "Shipyard"
metadata:
  name: "shipyard-podtato"
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
            - name: "release"`

func Test_GracefulShutdown(t *testing.T) {
	repoLocalDir, err := filepath.Abs("../assets/podtato-head")
	require.Nil(t, err)
	t.Log("Current local dir is : ", repoLocalDir)

	projectName := "tinypodtato"
	serviceName := "helloservice"
	serviceChartLocalDir := repoLocalDir + "/helm-charts/helloservice.tgz"
	serviceJmeterDir := repoLocalDir + "/jmeter"
	serviceHealthCheckEndpoint := "/metrics"
	shipyardPod := "shipyard-controller"

	t.Logf("Creating a new project %s", keptnProjectName)
	shipyardFilePath, err := CreateTmpShipyardFile(tinyShipyard)
	require.Nil(t, err)
	projectName, err = CreateProject(projectName, shipyardFilePath, true)
	require.Nil(t, err)

	t.Logf("Creating service %s in project %s", serviceName, keptnProjectName)
	_, err = ExecuteCommandf("keptn create service %s --project %s", serviceName, keptnProjectName)
	require.Nil(t, err)

	t.Logf("Adding resource for service %s in project %s", serviceName, keptnProjectName)
	_, err = ExecuteCommandf("keptn add-resource --project %s --service=%s --stage=%s --resource=%s --resourceUri=%s", keptnProjectName, serviceName, "dev", serviceChartLocalDir, "helm/helloservice.tgz")
	require.Nil(t, err)

	t.Log("Adding jmeter config in staging")
	_, err = ExecuteCommandf("keptn add-resource --project=%s --service=%s --stage=%s --resource=%s --resourceUri=%s", keptnProjectName, serviceName, "dev", serviceJmeterDir+"/jmeter.conf.yaml", "jmeter/jmeter.conf.yaml")
	require.Nil(t, err)

	///////////////////////////////////////
	// Deploy v0.1.0
	///////////////////////////////////////

	t.Logf("Trigger delivery of helloservice:v0.1.0")
	_, err = ExecuteCommandf("keptn trigger delivery --project=%s --service=%s --image=%s --tag=%s --sequence=%s", keptnProjectName, serviceName, "ghcr.io/podtato-head/podtatoserver", "v0.1.0", "delivery")
	require.Nil(t, err)

	waitAndKill(t, shipyardPod, 35)

	t.Logf("Sleeping for 60s...")
	time.Sleep(60 * time.Second)
	t.Logf("Continue to work...")
	err = WaitForPodOfDeployment(shipyardPod)
	require.Nil(t, err)

	t.Log("Verify Direct delivery of helloservice in stage dev")
	err = VerifyDirectDeployment(serviceName, keptnProjectName, "dev", "ghcr.io/podtato-head/podtatoserver", "v0.1.0")
	logError(err, t, shipyardPod)

	t.Log("Verify network access to public URI of helloservice in stage dev")
	cartPubURL, err := GetPublicURLOfService(serviceName, keptnProjectName, "dev")
	logError(err, t, shipyardPod)

	err = WaitForURL(cartPubURL+serviceHealthCheckEndpoint, time.Minute)
	logError(err, t, shipyardPod)

}

func waitAndKill(t *testing.T, keptnServiceName string, waitFor int) {
	t.Logf("Sleeping %d seconds...\n", waitFor)
	time.Sleep(time.Duration(waitFor) * time.Second)

	t.Logf("this is the state of %s before shutdown", keptnServiceName)
	GetDiagnostics(keptnServiceName, "")

	t.Logf("Killing %s Pod", keptnServiceName)
	err := RestartPod(keptnServiceName)
	logError(err, t, keptnServiceName)
}

func logError(err error, t *testing.T, service string) {
	if err != nil {
		t.Log(GetDiagnostics(service, ""))
	}
	require.Nil(t, err)
}
