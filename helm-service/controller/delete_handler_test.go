package controller

import (
	"errors"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/golang/mock/gomock"
	"github.com/keptn/go-utils/pkg/api/models"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/helm-service/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCreateDeleteHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockedStagesHandler := mocks.NewMockIStagesHandler(ctrl)
	instance := NewDeleteHandler(createKeptn(), mockedStagesHandler, "")
	assert.NotNil(t, instance)
}

func TestHandleDeleteEvent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockedStagesHandler := mocks.NewMockIStagesHandler(ctrl)
	mockedHelmExecutor := mocks.NewMockHelmExecutor(ctrl)

	mockedBaseHandler := &MockedHandler{
		keptnHandler: createKeptn(),
		helmExecutor: mockedHelmExecutor,
	}

	instance := &DeleteHandler{
		Handler:       mockedBaseHandler,
		stagesHandler: mockedStagesHandler,
	}

	eventData := keptnv2.ServiceDeleteFinishedEventData{
		EventData: keptnv2.EventData{
			Project: "my-project",
			Service: "my-service",
		},
	}

	ce := cloudevents.NewEvent()
	_ = ce.SetData(cloudevents.ApplicationJSON, eventData)

	stages := []*models.Stage{{Services: []*models.Service{{}}, StageName: "dev"}}
	mockedStagesHandler.EXPECT().GetAllStages("my-project").Return(stages, nil)
	mockedHelmExecutor.EXPECT().UninstallRelease("my-project-dev-my-service", "my-project-dev").Return(nil)
	mockedHelmExecutor.EXPECT().UninstallRelease("my-project-dev-my-service-generated", "my-project-dev").Return(nil)

	expectedDeleteFinishedEvent := cloudevents.NewEvent()
	expectedDeleteFinishedEvent.SetType("sh.keptn.event.service.delete.finished")
	expectedDeleteFinishedEvent.SetSource("helm-service")
	expectedDeleteFinishedEvent.SetDataContentType(cloudevents.ApplicationJSON)
	expectedDeleteFinishedEvent.SetExtension("triggeredid", "")
	expectedDeleteFinishedEvent.SetExtension("shkeptncontext", "")
	expectedDeleteFinishedEvent.SetData(cloudevents.ApplicationJSON, keptnv2.ServiceDeleteFinishedEventData{
		EventData: keptnv2.EventData{
			Project: "my-project",
			Service: "my-service",
			Status:  keptnv2.StatusSucceeded,
			Result:  keptnv2.ResultPass,
			Message: "Finished uninstalling service my-service in project my-project",
		},
	})

	instance.HandleEvent(ce)

	require.Equal(t, 1, len(mockedBaseHandler.sentCloudEvents))
	assert.Equal(t, expectedDeleteFinishedEvent, mockedBaseHandler.sentCloudEvents[0])
}

func TestWhenReceivingUnparsableEvent_ThenErrorMessageIsSent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockedBaseHandler := NewMockedHandler(createKeptn(), "")
	mockedStagesHandler := mocks.NewMockIStagesHandler(ctrl)
	instance := DeleteHandler{
		Handler:       mockedBaseHandler,
		stagesHandler: mockedStagesHandler,
	}

	instance.HandleEvent(createUnparsableEvent())
	require.Equal(t, 1, len(mockedBaseHandler.handledErrorEvents))

}

func TestWhenGettingStagesFails_Then(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockedBaseHandler := NewMockedHandler(createKeptn(), "")
	mockedStagesHandler := mocks.NewMockIStagesHandler(ctrl)
	instance := DeleteHandler{
		Handler:       mockedBaseHandler,
		stagesHandler: mockedStagesHandler,
	}

	mockedStagesHandler.EXPECT().GetAllStages(gomock.Any()).Return(nil, errors.New("Failed to get stages"))

	eventData := keptnv2.ServiceDeleteFinishedEventData{}

	ce := cloudevents.NewEvent()
	_ = ce.SetData(cloudevents.ApplicationJSON, eventData)

	instance.HandleEvent(ce)
	require.Equal(t, 1, len(mockedBaseHandler.handledErrorEvents))
}

func TestWhenUninstallingReleaseFails_FinishedEventIsStillSent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockedStagesHandler := mocks.NewMockIStagesHandler(ctrl)
	mockedHelmExecutor := mocks.NewMockHelmExecutor(ctrl)

	mockedBaseHandler := &MockedHandler{
		keptnHandler: createKeptn(),
		helmExecutor: mockedHelmExecutor,
	}

	instance := &DeleteHandler{
		Handler:       mockedBaseHandler,
		stagesHandler: mockedStagesHandler,
	}

	eventData := keptnv2.ServiceDeleteFinishedEventData{
		EventData: keptnv2.EventData{
			Project: "my-project",
			Service: "my-service",
		},
	}

	ce := cloudevents.NewEvent()
	_ = ce.SetData(cloudevents.ApplicationJSON, eventData)

	stages := []*models.Stage{{Services: []*models.Service{{}}, StageName: "dev"}}
	mockedStagesHandler.EXPECT().GetAllStages("my-project").Return(stages, nil)
	mockedHelmExecutor.EXPECT().UninstallRelease("my-project-dev-my-service", "my-project-dev").Return(errors.New("failed to uninstall release"))
	mockedHelmExecutor.EXPECT().UninstallRelease("my-project-dev-my-service-generated", "my-project-dev").Return(nil)

	expectedDeleteFinishedEvent := cloudevents.NewEvent()
	expectedDeleteFinishedEvent.SetType("sh.keptn.event.service.delete.finished")
	expectedDeleteFinishedEvent.SetSource("helm-service")
	expectedDeleteFinishedEvent.SetDataContentType(cloudevents.ApplicationJSON)
	expectedDeleteFinishedEvent.SetExtension("triggeredid", "")
	expectedDeleteFinishedEvent.SetExtension("shkeptncontext", "")
	expectedDeleteFinishedEvent.SetData(cloudevents.ApplicationJSON, keptnv2.ServiceDeleteFinishedEventData{
		EventData: keptnv2.EventData{
			Project: "my-project",
			Service: "my-service",
			Status:  keptnv2.StatusSucceeded,
			Result:  keptnv2.ResultPass,
			Message: "Finished uninstalling service my-service in project my-project",
		},
	})

	instance.HandleEvent(ce)

	require.Equal(t, 1, len(mockedBaseHandler.sentCloudEvents))
	assert.Equal(t, expectedDeleteFinishedEvent, mockedBaseHandler.sentCloudEvents[0])
}
