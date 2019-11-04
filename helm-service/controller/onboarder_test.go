package controller

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	"github.com/google/uuid"
	"github.com/keptn/keptn/helm-service/controller/helm"
	"github.com/keptn/keptn/helm-service/controller/mesh"
	"github.com/keptn/keptn/helm-service/pkg/helmtest"

	configmodels "github.com/keptn/go-utils/pkg/configuration-service/models"
	configutils "github.com/keptn/go-utils/pkg/configuration-service/utils"
	keptnevents "github.com/keptn/go-utils/pkg/events"
	keptnutils "github.com/keptn/go-utils/pkg/utils"
	"github.com/stretchr/testify/assert"
)

const configBaseURL = "localhost:6060"
const projectName = "sockshop"
const serviceName = "carts"
const stage1 = "dev"
const stage2 = "staging"
const stage3 = "production"

const shipyard = `project: sockshop
stages:
    - {deployment_strategy: direct, name: dev, test_strategy: functional}
    - {deployment_strategy: blue_green_service, name: staging, test_strategy: performance}
    - {deployment_strategy: blue_green_service, name: production}`

func createTestProjet(t *testing.T) {

	prjHandler := configutils.NewProjectHandler(configBaseURL)
	prj := configmodels.Project{ProjectName: projectName}
	respErr, err := prjHandler.CreateProject(prj)
	check(err, t)
	assert.Nil(t, respErr, "Creating a project failed")

	// Send shipyard
	rHandler := configutils.NewResourceHandler(configBaseURL)
	shipyardURI := "shipyard.yaml"
	shipyardResource := configmodels.Resource{ResourceURI: &shipyardURI, ResourceContent: shipyard}
	resources := []*configmodels.Resource{&shipyardResource}
	_, err = rHandler.CreateProjectResources(projectName, resources)
	check(err, t)

	// Create stages
	stageHandler := configutils.NewStageHandler(configBaseURL)
	for _, stage := range []string{stage1, stage2, stage3} {

		respErr, err := stageHandler.CreateStage(projectName, stage)
		check(err, t)
		assert.Nil(t, respErr, "Creating a stage failed")
	}
}

// TestDoOnboard tests the onboarding of a new chart. Therefore, this test requires the configuration-service
// on localhost:8080
func TestDoOnboard(t *testing.T) {

	_, err := http.Get("http://" + configBaseURL)
	if err != nil {
		t.Skip("Skipping test; no local configuration-service available")
	}
	os.Setenv("CONFIGURATION_SERVICE", configBaseURL)

	createTestProjet(t)

	data := helmtest.CreateHelmChartData(t)
	encodedChart := base64.StdEncoding.EncodeToString(data)
	fmt.Println(encodedChart)
	ce := cloudevents.New("0.2")
	dataBytes, err := json.Marshal(keptnevents.ServiceCreateEventData{Project: projectName, Service: serviceName, HelmChart: encodedChart})
	check(err, t)
	ce.Data = dataBytes

	id := uuid.New().String()
	onboarder := NewOnboarder(mesh.NewIstioMesh(), helm.NewCanaryOnNamespaceGenerator(),
		keptnutils.NewLogger(id, "service.create", "helm-service"), "test.keptn.sh")
	loggingDone := make(chan bool)
	err = onboarder.DoOnboard(ce, loggingDone)

	check(err, t)
}

func check(e error, t *testing.T) {
	if e != nil {
		t.Error(e)
	}
}
