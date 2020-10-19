package controller

import (
	"encoding/base64"
	"errors"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/golang/mock/gomock"
	configutils "github.com/keptn/go-utils/pkg/api/utils"
	"github.com/keptn/keptn/helm-service/mocks"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/keptn/keptn/helm-service/pkg/helm"

	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"

	"github.com/keptn/go-utils/pkg/api/models"
	configmodels "github.com/keptn/go-utils/pkg/api/models"
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

//MOCKS
var mockedBaseHandler *MockHandler
var mockedMesh *mocks.MockMesh
var mockedProjectHandler *mocks.MockProjectOperator
var mockedNamespaceManager *mocks.MockINamespaceManager
var mockedStagesHandler *mocks.MockIStagesHandler
var mockedServiceHandler *mocks.MockIServiceHandler
var mockedChartStorer *mocks.MockChartStorer

func createMocks(ctrl *gomock.Controller) {
	mockedBaseHandler = NewMockHandler(ctrl)
	mockedMesh = mocks.NewMockMesh(ctrl)
	mockedProjectHandler = mocks.NewMockProjectOperator(ctrl)
	mockedNamespaceManager = mocks.NewMockINamespaceManager(ctrl)
	mockedStagesHandler = mocks.NewMockIStagesHandler(ctrl)
	mockedServiceHandler = mocks.NewMockIServiceHandler(ctrl)
	mockedChartStorer = mocks.NewMockChartStorer(ctrl)

}

func TestHandleEvent_WhenPassingUnparsableEvent_ThenHandleErrorIsCalled(t *testing.T) {
	//PREPARE MOCKS
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	createMocks(ctrl)

	mockedBaseHandler.EXPECT().handleError(gomock.Eq("EVENT_ID"), gomock.Any(), gomock.Eq("service.create"), gomock.Any())

	instance := Onboarder{
		Handler:          mockedBaseHandler,
		mesh:             mockedMesh,
		projectHandler:   mockedProjectHandler,
		namespaceManager: mockedNamespaceManager,
		stagesHandler:    mockedStagesHandler,
	}

	event := cloudevents.NewEvent()
	event.SetData(cloudevents.ApplicationJSON, "WEIRD_JSON_CONTENT")
	event.SetID("EVENT_ID")
	instance.HandleEvent(event, nilCloser)
}

func TestHandleEvent_WhenHelmChartMissing_ThenNothingHappens(t *testing.T) {
	//PREPARE MOCKS
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	createMocks(ctrl)

	instance := Onboarder{
		Handler:          mockedBaseHandler,
		mesh:             mockedMesh,
		projectHandler:   mockedProjectHandler,
		namespaceManager: mockedNamespaceManager,
		stagesHandler:    mockedStagesHandler,
	}

	event := cloudevents.NewEvent()
	instance.HandleEvent(event, nilCloser)
}

func TestHandleEvent_WhenNoProjectExists_ThenHandleErrorIsCalled(t *testing.T) {
	//PREPARE MOCKS
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	createMocks(ctrl)

	mockedBaseHandler.EXPECT().getKeptnHandler().Return(nil)
	mockedProjectHandler.EXPECT().GetProject(gomock.Any()).Return(nil, &models.Error{Message: stringp("")})
	mockedBaseHandler.EXPECT().handleError(gomock.Eq("EVENT_ID"), gomock.Any(), gomock.Eq("service.create"), gomock.Any())

	instance := Onboarder{
		Handler:          mockedBaseHandler,
		mesh:             mockedMesh,
		projectHandler:   mockedProjectHandler,
		namespaceManager: mockedNamespaceManager,
		stagesHandler:    mockedStagesHandler,
	}

	event := cloudevents.NewEvent()
	event.SetID("EVENT_ID")
	event.SetData(cloudevents.ApplicationJSON, keptnv2.ServiceCreateFinishedEventData{
		EventData: keptnv2.EventData{
			Project: "my-project",
			Stage:   "dev",
			Service: "carts",
		},
		Helm: keptnv2.Helm{
			Chart: base64.StdEncoding.EncodeToString(helm.CreateTestHelmChartData(t)),
		},
	})

	instance.HandleEvent(event, nilCloser)
}

func TestHandleEvent_WhenNoStagesDefined_ThenHandleErrorIsCalled(t *testing.T) {
	//PREPARE MOCKS
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	createMocks(ctrl)

	mockedBaseHandler.EXPECT().getKeptnHandler().Return(nil)
	mockedProjectHandler.EXPECT().GetProject(gomock.Any()).Return(&models.Project{}, nil)
	mockedStagesHandler.EXPECT().GetAllStages(gomock.Any()).Return([]*models.Stage{}, nil)
	mockedBaseHandler.EXPECT().handleError(gomock.Eq("EVENT_ID"), gomock.Any(), gomock.Eq("service.create"), gomock.Any())

	instance := Onboarder{
		Handler:          mockedBaseHandler,
		mesh:             mockedMesh,
		projectHandler:   mockedProjectHandler,
		namespaceManager: mockedNamespaceManager,
		stagesHandler:    mockedStagesHandler,
	}

	event := cloudevents.NewEvent()
	event.SetID("EVENT_ID")
	event.SetData(cloudevents.ApplicationJSON, keptnv2.ServiceCreateFinishedEventData{
		EventData: keptnv2.EventData{
			Project: "my-project",
			Stage:   "dev",
			Service: "carts",
		},
		Helm: keptnv2.Helm{
			Chart: base64.StdEncoding.EncodeToString(helm.CreateTestHelmChartData(t)),
		},
	})

	instance.HandleEvent(event, nilCloser)

}

func TestHandleEvent_OnboardsService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	createMocks(ctrl)

	event := cloudevents.NewEvent()
	event.SetID("EVENT_ID")
	event.SetData(cloudevents.ApplicationJSON, keptnv2.ServiceCreateFinishedEventData{
		EventData: keptnv2.EventData{
			Project: "my-project",
			Stage:   "dev",
			Service: "carts",
		},
		Helm: keptnv2.Helm{
			Chart: base64.StdEncoding.EncodeToString(helm.CreateTestHelmChartData(t)),
		},
	})

	stages := []*models.Stage{&models.Stage{Services: []*models.Service{&models.Service{}}, StageName: "dev"}}

	mockedBaseHandler.EXPECT().getKeptnHandler().Return(nil)
	mockedProjectHandler.EXPECT().GetProject(gomock.Any()).Return(&models.Project{Stages: stages}, nil)
	mockedStagesHandler.EXPECT().GetAllStages(gomock.Any()).Return(stages, nil)
	mockedNamespaceManager.EXPECT().InitNamespaces(gomock.Eq("my-project"), gomock.Eq([]string{"dev"}))
	mockedServiceHandler.EXPECT().GetService(gomock.Eq("my-project"), gomock.Eq("dev"), gomock.Eq("carts")).Return(nil, nil)
	mockedBaseHandler.EXPECT().getConfigServiceURL().Return("")
	mockedChartStorer.EXPECT().StoreChart(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return("", nil)
	mockedBaseHandler.EXPECT().sendEvent(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

	instance := Onboarder{
		Handler:          mockedBaseHandler,
		mesh:             mockedMesh,
		projectHandler:   mockedProjectHandler,
		namespaceManager: mockedNamespaceManager,
		stagesHandler:    mockedStagesHandler,
		serviceHandler:   mockedServiceHandler,
		chartStorer:      mockedChartStorer,
	}

	instance.HandleEvent(event, nilCloser)

}

func nilCloser(keptnHandler *keptnv2.Keptn) {
	//No-op
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
