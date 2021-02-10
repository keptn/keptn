package controller

import (
	"encoding/base64"
	"errors"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/golang/mock/gomock"
	"github.com/keptn/keptn/helm-service/mocks"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/keptn/keptn/helm-service/pkg/helm"

	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"

	"github.com/keptn/go-utils/pkg/api/models"
)

//mockColleciton holds all the mocks
type mocksCollection struct {
	mockedBaseHandler      *MockedHandler
	mockedMesh             *mocks.MockMesh
	mockedProjectHandler   *mocks.MockIProjectHandler
	mockedNamespaceManager *mocks.MockINamespaceManager
	mockedStagesHandler    *mocks.MockIStagesHandler
	mockedServiceHandler   *mocks.MockIServiceHandler
	mockedChartStorer      *mocks.MockIChartStorer
	mockedChartGenerator   *mocks.MockChartGenerator
	mockedChartPackager    *mocks.MockIChartPackager
	mockedOnBoarder        *mocks.MockOnboarder
}

// testOnboarderCreator is an oboarder which has only mocked dependencies
type testOnboarderCreator struct {
}

// newTestOnboardHandlerCreator creates an instance of testOnboarderCreator which uses only mocks
func newTestOnboardHandlerCreator(t *testing.T, mockedBaseHandlerOptions ...MockedHandlerOption) (*gomock.Controller, OnboardHandler, mocksCollection) {

	ctrl := gomock.NewController(t)
	mockedBaseHandler := NewMockedHandler(createKeptn(), "", mockedBaseHandlerOptions...)
	mockedMesh := mocks.NewMockMesh(ctrl)
	mockedProjectHandler := mocks.NewMockIProjectHandler(ctrl)
	mockedStagesHandler := mocks.NewMockIStagesHandler(ctrl)
	mockedServiceHandler := mocks.NewMockIServiceHandler(ctrl)
	mockedOnBoarder := mocks.NewMockOnboarder(ctrl)

	onboardHandler := OnboardHandler{
		Handler:        mockedBaseHandler,
		projectHandler: mockedProjectHandler,
		stagesHandler:  mockedStagesHandler,
		onboarder:      mockedOnBoarder,
	}

	mocksCol := mocksCollection{
		mockedBaseHandler:    mockedBaseHandler,
		mockedMesh:           mockedMesh,
		mockedProjectHandler: mockedProjectHandler,
		mockedStagesHandler:  mockedStagesHandler,
		mockedServiceHandler: mockedServiceHandler,
		mockedOnBoarder:      mockedOnBoarder,
	}
	return ctrl, onboardHandler, mocksCol

}

func TestCreateOnboarderHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockedBaseHandler := createKeptnBaseHandlerMock()
	mockedProjectHandler := mocks.NewMockIProjectHandler(ctrl)
	mockedStagesHandler := mocks.NewMockIStagesHandler(ctrl)
	mockedOnboarder := mocks.NewMockOnboarder(ctrl)

	handler := NewOnboardHandler(
		mockedBaseHandler,
		mockedProjectHandler,
		mockedStagesHandler,
		mockedOnboarder)

	assert.NotNil(t, handler)

}

func TestHandleEvent_WhenPassingUnparsableEvent_ThenHandleErrorIsCalled(t *testing.T) {

	ctrl, instance, _ := newTestOnboardHandlerCreator(t)
	defer ctrl.Finish()
	instance.HandleEvent(createUnparsableEvent())
}

func TestHandleEvent_WhenHelmChartMissing_ThenNothingHappens(t *testing.T) {
	ctrl, instance, moqs := newTestOnboardHandlerCreator(t)
	defer ctrl.Finish()

	instance.HandleEvent(cloudevents.NewEvent())
	assert.Empty(t, moqs.mockedBaseHandler.sentCloudEvents)
	assert.Empty(t, moqs.mockedBaseHandler.handledErrorEvents)
}

func TestHandleEvent_WhenNoProjectExists_ThenHandleErrorIsCalled(t *testing.T) {
	ctrl, instance, moqs := newTestOnboardHandlerCreator(t)
	defer ctrl.Finish()

	moqs.mockedProjectHandler.EXPECT().GetProject(gomock.Any()).Return(nil, &models.Error{Message: stringp("")})

	instance.HandleEvent(createEvent(t, "EVENT_ID"))
	assert.Equal(t, 1, len(moqs.mockedBaseHandler.handledErrorEvents))
}

func TestHandleEvent_WhenNoStagesDefined_ThenHandleErrorIsCalled(t *testing.T) {
	ctrl, instance, moqs := newTestOnboardHandlerCreator(t)
	defer ctrl.Finish()

	moqs.mockedProjectHandler.EXPECT().GetProject(gomock.Any()).Return(&models.Project{}, nil)
	moqs.mockedStagesHandler.EXPECT().GetAllStages(gomock.Any()).Return([]*models.Stage{}, nil)

	instance.HandleEvent(createEvent(t, "EVENT_ID"))
	assert.Equal(t, 1, len(moqs.mockedBaseHandler.handledErrorEvents))
}

func TestHandleEvent_WhenStagesCannotBeFetched_ThenHandleErrorIsCalled(t *testing.T) {
	ctrl, instance, moqs := newTestOnboardHandlerCreator(t)
	defer ctrl.Finish()

	stages := []*models.Stage{{Services: []*models.Service{{}}, StageName: "dev"}}

	moqs.mockedProjectHandler.EXPECT().GetProject(gomock.Any()).Return(&models.Project{Stages: stages}, nil)
	moqs.mockedStagesHandler.EXPECT().GetAllStages(gomock.Any()).Return(nil, errors.New("unable to fetch stages"))

	instance.HandleEvent(createEvent(t, "EVENT_ID"))
	assert.Equal(t, 1, len(moqs.mockedBaseHandler.handledErrorEvents))
}

func stringp(s string) *string {
	return &s
}

func createEventWith(t *testing.T, id string, eventData keptnv2.EventData) event.Event {
	testEvent := cloudevents.NewEvent()
	testEvent.SetID(id)
	_ = testEvent.SetData(cloudevents.ApplicationJSON, keptnv2.ServiceCreateFinishedEventData{
		EventData: eventData,
		Helm: keptnv2.Helm{
			Chart: base64.StdEncoding.EncodeToString(helm.CreateTestHelmChartData(t)),
		},
	})
	return testEvent
}

func createEvent(t *testing.T, id string) event.Event {
	testEvent := cloudevents.NewEvent()
	testEvent.SetID(id)
	_ = testEvent.SetData(cloudevents.ApplicationJSON, keptnv2.ServiceCreateFinishedEventData{
		EventData: keptnv2.EventData{
			Project: "my-project",
			Stage:   "dev",
			Service: "carts",
		},
		Helm: keptnv2.Helm{
			Chart: base64.StdEncoding.EncodeToString(helm.CreateTestHelmChartData(t)),
		},
	})
	return testEvent
}
