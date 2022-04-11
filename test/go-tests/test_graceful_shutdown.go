package go_tests

import (
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/mholt/archiver/v3"
	"github.com/stretchr/testify/require"
	"os"
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

type Setup struct {
	Project        string
	Service        string
	Chart          string
	Jmeter         string
	HealthEndpoint string
}

func newSetup(t *testing.T) *Setup {
	repoLocalDir, err := filepath.Abs("../assets/podtato-head")
	require.Nil(t, err)
	chartSourceDir := repoLocalDir + "/helm-charts/helloservice"
	chartArchiveDir := repoLocalDir + "/helm-charts/helloservice.tgz"

	err = archiver.Archive([]string{chartSourceDir}, chartArchiveDir)
	require.Nil(t, err)

	return &Setup{
		Project:        "tinypodtato",
		Service:        "helloservice",
		Chart:          chartArchiveDir,
		Jmeter:         repoLocalDir + "/jmeter",
		HealthEndpoint: "/metrics",
	}
}

func Test_GracefulShutdown(t *testing.T) {
	shipyardPod := "shipyard-controller"
	setup := newSetup(t)

	// Delete chart archive at the end of the test
	defer func(path string) {
		err := os.RemoveAll(path)
		require.Nil(t, err)
	}(setup.Chart)

	_ = startDelivery(t, setup)

	waitAndKill(t, shipyardPod, 35)

	t.Logf("Sleeping for 60s...")
	time.Sleep(60 * time.Second)
	checkSuccessfulDeployment(t, shipyardPod, setup)
}

func checkSuccessfulDeployment(t *testing.T, shipyardPod string, setup *Setup) {
	t.Logf("Continue to work...")
	err := WaitForPodOfDeployment(shipyardPod)
	require.Nil(t, err)

	t.Log("Verify Direct delivery of helloservice in stage dev")
	err = VerifyDirectDeployment(setup.Service, setup.Project, "dev", "ghcr.io/podtato-head/podtatoserver", "v0.1.0")
	logError(err, t, shipyardPod)

	t.Log("Verify network access to public URI of helloservice in stage dev")
	cartPubURL, err := GetPublicURLOfService(setup.Service, setup.Project, "dev")
	logError(err, t, shipyardPod)

	err = WaitForURL(cartPubURL+setup.HealthEndpoint, time.Minute)
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

func Test_GracefulLeader(t *testing.T) {

	shipyardPod := "shipyard-controller"
	setup := newSetup(t)
	setup.Project = "leader_election"
	keptnContext := startDelivery(t, setup)

	require.Eventually(t, func() bool {
		t.Log("checking if evaluation.finished event is available")
		event, err := GetLatestEventOfType(keptnContext, setup.Project, "dev", keptnv2.GetStartedEventType(keptnv2.DeploymentTaskName))
		if err != nil || event == nil {
			return false
		}
		waitAndKill(t, shipyardPod, 0)
		return true
	}, 1*time.Minute, 10*time.Second)

	checkSuccessfulDeployment(t, shipyardPod, setup)

}

func startDelivery(t *testing.T, setup *Setup) string {
	t.Logf("Creating a new project %s", setup.Project)
	shipyardFilePath, err := CreateTmpShipyardFile(tinyShipyard)
	require.Nil(t, err)
	setup.Project, err = CreateProject(setup.Project, shipyardFilePath)
	require.Nil(t, err)

	t.Logf("Creating service %s in project %s", setup.Service, setup.Project)
	_, err = ExecuteCommandf("keptn create service %s --project %s", setup.Service, setup.Project)
	require.Nil(t, err)

	t.Logf("Adding resource for service %s in project %s", setup.Service, setup.Project)
	_, err = ExecuteCommandf("keptn add-resource --project %s --service=%s --stage=%s --resource=%s --resourceUri=%s", setup.Project, setup.Service, "dev", setup.Chart, "helm/helloservice.tgz")
	require.Nil(t, err)

	t.Log("Adding jmeter config in staging")
	_, err = ExecuteCommandf("keptn add-resource --project=%s --service=%s --stage=%s --resource=%s --resourceUri=%s", setup.Project, setup.Service, "dev", setup.Jmeter+"/jmeter.conf.yaml", "jmeter/jmeter.conf.yaml")
	require.Nil(t, err)

	///////////////////////////////////////
	// Deploy v0.1.0
	///////////////////////////////////////

	t.Logf("Trigger delivery of helloservice:v0.1.0")
	keptnContext, err := ExecuteCommandf("keptn trigger delivery --project=%s --service=%s --image=%s:%s --sequence=%s", setup.Project, setup.Service, "ghcr.io/podtato-head/podtatoserver", "v0.1.0", "delivery")
	require.Nil(t, err)
	return keptnContext
}

func logError(err error, t *testing.T, service string) {
	if err != nil {
		t.Log(GetDiagnostics(service, ""))
	}
	require.Nil(t, err)
}
