package controller

import (
	"encoding/base64"
	"errors"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/golang/mock/gomock"
	"github.com/keptn/keptn/helm-service/mocks"
	"testing"

	"github.com/keptn/keptn/helm-service/pkg/helm"

	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"

	"github.com/keptn/go-utils/pkg/api/models"
	configmodels "github.com/keptn/go-utils/pkg/api/models"
)

//MOCKS
var mockedBaseHandler *MockHandler
var mockedMesh *mocks.MockMesh
var mockedProjectHandler *mocks.MockProjectOperator
var mockedNamespaceManager *mocks.MockINamespaceManager
var mockedStagesHandler *mocks.MockIStagesHandler
var mockedServiceHandler *mocks.MockIServiceHandler
var mockedChartStorer *mocks.MockChartStorer

func createMocks(t *testing.T) *gomock.Controller {

	ctrl := gomock.NewController(t)
	mockedBaseHandler = NewMockHandler(ctrl)
	mockedMesh = mocks.NewMockMesh(ctrl)
	mockedProjectHandler = mocks.NewMockProjectOperator(ctrl)
	mockedNamespaceManager = mocks.NewMockINamespaceManager(ctrl)
	mockedStagesHandler = mocks.NewMockIStagesHandler(ctrl)
	mockedServiceHandler = mocks.NewMockIServiceHandler(ctrl)
	mockedChartStorer = mocks.NewMockChartStorer(ctrl)
	return ctrl

}

func TestHandleEvent_WhenPassingUnparsableEvent_ThenHandleErrorIsCalled(t *testing.T) {
	//PREPARE MOCKS
	ctrl := createMocks(t)
	defer ctrl.Finish()

	mockedBaseHandler.EXPECT().handleError(gomock.Eq("EVENT_ID"), gomock.Any(), gomock.Eq("service.create"), gomock.Any())

	instance := Onboarder{
		Handler:          mockedBaseHandler,
		mesh:             mockedMesh,
		projectHandler:   mockedProjectHandler,
		namespaceManager: mockedNamespaceManager,
		stagesHandler:    mockedStagesHandler,
	}

	instance.HandleEvent(createUnparsableEvent(), nilCloser)
}

func TestHandleEvent_WhenHelmChartMissing_ThenNothingHappens(t *testing.T) {
	//PREPARE MOCKS
	ctrl := createMocks(t)
	defer ctrl.Finish()

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
	ctrl := createMocks(t)
	defer ctrl.Finish()

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

	instance.HandleEvent(createEvent(t, "EVENT_ID"), nilCloser)
}

func TestHandleEvent_WhenNoStagesDefined_ThenHandleErrorIsCalled(t *testing.T) {
	//PREPARE MOCKS
	ctrl := createMocks(t)
	defer ctrl.Finish()

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

	instance.HandleEvent(createEvent(t, "EVENT_ID"), nilCloser)

}

func TestHandleEvent_OnboardsService(t *testing.T) {
	ctrl := createMocks(t)
	defer ctrl.Finish()

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

	instance.HandleEvent(createEvent(t, "EVENT_ID"), nilCloser)

}

func TestCheckAndSetServiceName(t *testing.T) {

	o := Onboarder{
		Handler: &HandlerBase{},
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

func createEvent(t *testing.T, id string) event.Event {
	event := cloudevents.NewEvent()
	event.SetID(id)
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
	return event
}

func createUnparsableEvent() event.Event {
	event := cloudevents.NewEvent()
	event.SetData(cloudevents.ApplicationJSON, "WEIRD_JSON_CONTENT")
	event.SetID("EVENT_ID")
	return event
}

func nilCloser(keptnHandler *keptnv2.Keptn) {
	//No-op
}
