package controller

import (
	"encoding/base64"
	"errors"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/golang/mock/gomock"
	keptncommon "github.com/keptn/go-utils/pkg/lib/keptn"
	"github.com/keptn/keptn/helm-service/mocks"

	"testing"

	"github.com/keptn/keptn/helm-service/pkg/helm"

	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"

	"github.com/keptn/go-utils/pkg/api/models"
	configmodels "github.com/keptn/go-utils/pkg/api/models"
	configutils "github.com/keptn/go-utils/pkg/api/utils"
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

//go:ignore
func TestHandleEvent(t *testing.T) {
	//
	//
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ce := cloudevents.NewEvent()
	keptn, _ := keptnv2.NewKeptn(&ce, keptncommon.KeptnOpts{})

	//mockedHandler := mocks.NewMockHandler(ctrl)
	mockedMesh := mocks.NewMockMesh(ctrl)
	mockedProjectHandler := mocks.NewMockProjectOperator(ctrl)
	mockedNamespaceManager := mocks.NewMockINamespaceManager(ctrl)
	mockedStagesHandler := mocks.NewMockIStagesHandler(ctrl)

	onboarder := NewOnboarder(
		keptn,
		mockedMesh,
		mockedProjectHandler,
		mockedNamespaceManager,
		mockedStagesHandler,
		"",
	)

	eventData := keptnv2.EventData{
		Project: "my-project",
		Stage:   "dev",
		Service: "carts",
		Labels:  nil,
		Status:  "some-status",
		Result:  "some-result",
		Message: "MESSAGE",
	}

	data := helm.CreateTestHelmChartData(t)

	serviceCreateFinishedEventData := keptnv2.ServiceCreateFinishedEventData{
		EventData: eventData,
		Helm: keptnv2.Helm{
			Chart: base64.StdEncoding.EncodeToString(data),
		},
	}

	event := cloudevents.NewEvent()
	event.SetType("test-type")
	event.SetSource("test-source")
	event.SetData("", serviceCreateFinishedEventData)

	mockedProjectHandler.EXPECT().GetProject(gomock.Any()).Return(nil, nil)
	mockedStagesHandler.EXPECT().GetAllStages(gomock.Any()).Return([]*models.Stage{&models.Stage{
		Services: []*models.Service{&models.Service{
			CreationDate:   "",
			DeployedImage:  "",
			LastEventTypes: nil,
			OpenApprovals:  nil,
			ServiceName:    "",
		}},
		StageName: "dev",
	}}, nil)

	onboarder.HandleEvent(event, nilCloser)

}

func nilCloser(keptnHandler *keptnv2.Keptn) {

}

func createTestProject(t *testing.T) {

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
	_, err2 := rHandler.CreateProjectResources(projectName, resources)
	if err2 != nil {
		t.Error(err)
	}

	// Create stages
	stageHandler := configutils.NewStageHandler(configBaseURL)
	for _, stage := range []string{stage1, stage2, stage3} {

		respErr, err := stageHandler.CreateStage(projectName, stage)
		check(err, t)
		assert.Nil(t, respErr, "Creating a stage failed")
	}
}

func TestCheckAndSetServiceName(t *testing.T) {

	mockHandler := &HandlerBase{
		keptnHandler:     nil,
		helmExecutor:     nil,
		configServiceURL: configBaseURL,
	}

	o := Onboarder{
		Handler: mockHandler,
		mesh:    nil,
	}
	data := helm.CreateTestHelmChartData(t)

	testCases := []struct {
		name        string
		event       *keptnv2.ServiceCreateFinishedEventData
		error       error
		serviceName string
	}{
		{"Mismatch", &keptnv2.ServiceCreateFinishedEventData{EventData: keptnv2.EventData{Service: "carts-1"},
			Helm: keptnv2.Helm{Chart: base64.StdEncoding.EncodeToString(data)}},
			errors.New("Provided Keptn service name \"carts-1\" does not match Kubernetes service name \"carts\""), "carts-1"},
		{"Match", &keptnv2.ServiceCreateFinishedEventData{EventData: keptnv2.EventData{Service: "carts"},
			Helm: keptnv2.Helm{Chart: base64.StdEncoding.EncodeToString(data)}},
			nil, "carts"},
		{"Set", &keptnv2.ServiceCreateFinishedEventData{EventData: keptnv2.EventData{Service: ""},
			Helm: keptnv2.Helm{Chart: base64.StdEncoding.EncodeToString(data)}},
			nil, "carts"},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			res := o.checkAndSetServiceName(tt.event)
			if res == nil && res != tt.error {
				t.Errorf("got nil, want %s", tt.error.Error())
			} else if res != nil && tt.error != nil && res.Error() != tt.error.Error() {
				t.Errorf("got %s, want %s", res.Error(), tt.error.Error())
			} else if res != nil && tt.error == nil {
				t.Errorf("got %s, want nil", res.Error())
			}

			if tt.event.Service != tt.serviceName {
				t.Errorf("got %s, want %s", tt.event.Service, tt.serviceName)
			}
		})
	}
}

func stringp(s string) *string {
	return &s
}

func check(e *configmodels.Error, t *testing.T) {
	if e != nil {
		t.Error(e.Message)
	}
}
