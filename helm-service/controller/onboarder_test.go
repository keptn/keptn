package controller

import (
	"encoding/json"
	"testing"

	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	"github.com/google/uuid"
	"github.com/keptn/keptn/helm-service/controller/helm"
	"github.com/keptn/keptn/helm-service/controller/mesh"

	keptnevents "github.com/keptn/go-utils/pkg/events"
	keptnmodels "github.com/keptn/go-utils/pkg/models"
	keptnutils "github.com/keptn/go-utils/pkg/utils"
	"github.com/stretchr/testify/assert"
)

const configBaseURL = "localhost:8080"
const projectName = "sockshop"
const serviceName = "carts"
const stage1 = "dev"
const stage2 = "prod"

func createTestProjet(t *testing.T) {

	prjHandler := keptnutils.NewProjectHandler(configBaseURL)
	prj := keptnmodels.Project{ProjectName: projectName}
	respErr, err := prjHandler.CreateProject(prj)
	check(err, t)
	assert.Nil(t, respErr, "Creating a project failed")
}

func createTestStages(t *testing.T) {

	stageHandler := keptnutils.NewStageHandler(configBaseURL)
	respErr, err := stageHandler.CreateStage(projectName, stage1)
	check(err, t)
	assert.Nil(t, respErr, "Creating a stage failed")
	respErr, err = stageHandler.CreateStage(projectName, stage2)
	check(err, t)
	assert.Nil(t, respErr, "Creating a stage failed")
}

func TestDoOnboard(t *testing.T) {

	createTestProjet(t)
	createTestStages(t)

	data := helm.CreateHelmChartData(t)
	ce := cloudevents.New("0.2")
	dataBytes, err := json.Marshal(keptnevents.ServiceCreateEventData{Project: projectName, Service: serviceName, HelmChart: data})
	check(err, t)
	ce.Data = dataBytes

	id := uuid.New().String()
	err = DoOnboard(ce, mesh.NewIstioMesh(), keptnutils.NewLogger(id, "service.create", "helm-service"), id, configBaseURL)

	check(err, t)
}

func check(e error, t *testing.T) {
	if e != nil {
		t.Error(e)
	}
}
