package controller

import (
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/golang/mock/gomock"
	keptnevents "github.com/keptn/go-utils/pkg/lib"
	. "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/helm-service/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestHandleReleaseTriggeredEvent_WhenDeploymentStrategyDirect_ThenNoActionRequired(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockedBaseHandler := NewMockedHandler(createKeptn(), "")
	mockedMesh := mocks.NewMockMesh(ctrl)
	mockedChartGenerator := mocks.NewMockChartGenerator(ctrl)
	mockedConfigurationChanger := mocks.NewMockIConfigurationChanger(ctrl)

	releaseTriggeredEventData := ReleaseTriggeredEventData{
		EventData: EventData{},
		Deployment: DeploymentFinishedData{
			DeploymentStrategy: keptnevents.Direct.String(),
		},
	}

	instance := &ReleaseHandler{
		Handler:               mockedBaseHandler,
		mesh:                  mockedMesh,
		generatedChartHandler: mockedChartGenerator,
		configurationChanger:  mockedConfigurationChanger,
	}

	ce := cloudevents.NewEvent()
	_ = ce.SetData(cloudevents.ApplicationJSON, releaseTriggeredEventData)

	instance.HandleEvent(ce)

	expectedReleaseStartedEvent := cloudevents.NewEvent()
	expectedReleaseStartedEvent.SetType("sh.keptn.event.release.started")
	expectedReleaseStartedEvent.SetSource("helm-service")
	expectedReleaseStartedEvent.SetDataContentType(cloudevents.ApplicationJSON)
	expectedReleaseStartedEvent.SetExtension("triggeredid", "")
	expectedReleaseStartedEvent.SetExtension("shkeptncontext", "")
	expectedReleaseStartedEvent.SetData(cloudevents.ApplicationJSON, ReleaseStartedEventData{
		EventData: EventData{
			Status: StatusSucceeded,
		},
	})

	expectedReleaseFinishedEvent := cloudevents.NewEvent()
	expectedReleaseFinishedEvent.SetType("sh.keptn.event.release.finished")
	expectedReleaseFinishedEvent.SetSource("helm-service")
	expectedReleaseFinishedEvent.SetDataContentType(cloudevents.ApplicationJSON)
	expectedReleaseFinishedEvent.SetExtension("triggeredid", "")
	expectedReleaseFinishedEvent.SetExtension("shkeptncontext", "")
	expectedReleaseFinishedEvent.SetData(cloudevents.ApplicationJSON, ReleaseFinishedEventData{
		EventData: EventData{
			Status:  StatusSucceeded,
			Message: "Finished release",
		},
	})

	require.Equal(t, 2, len(mockedBaseHandler.sentCloudEvents))
	require.Equal(t, 0, len(mockedBaseHandler.handledErrorEvents))
	assert.Equal(t, expectedReleaseStartedEvent, mockedBaseHandler.sentCloudEvents[0])
	assert.Equal(t, expectedReleaseFinishedEvent, mockedBaseHandler.sentCloudEvents[1])
}

func TestHandleReleaseTriggeredEvent_WithInvalidDeploymentStrategy(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockedBaseHandler := NewMockedHandler(createKeptn(), "")
	mockedMesh := mocks.NewMockMesh(ctrl)
	mockedChartGenerator := mocks.NewMockChartGenerator(ctrl)
	mockedConfigurationChanger := mocks.NewMockIConfigurationChanger(ctrl)

	releaseTriggeredEventData := ReleaseTriggeredEventData{
		EventData: EventData{},
		Deployment: DeploymentFinishedData{
			DeploymentStrategy: "???", // <-- OHOH
		},
	}

	instance := &ReleaseHandler{
		Handler:               mockedBaseHandler,
		mesh:                  mockedMesh,
		generatedChartHandler: mockedChartGenerator,
		configurationChanger:  mockedConfigurationChanger,
	}

	ce := cloudevents.NewEvent()
	_ = ce.SetData(cloudevents.ApplicationJSON, releaseTriggeredEventData)

	instance.HandleEvent(ce)

	expectedReleaseStartedEvent := cloudevents.NewEvent()
	expectedReleaseStartedEvent.SetType("sh.keptn.event.release.started")
	expectedReleaseStartedEvent.SetSource("helm-service")
	expectedReleaseStartedEvent.SetDataContentType(cloudevents.ApplicationJSON)
	expectedReleaseStartedEvent.SetExtension("triggeredid", "")
	expectedReleaseStartedEvent.SetExtension("shkeptncontext", "")
	expectedReleaseStartedEvent.SetData(cloudevents.ApplicationJSON, ReleaseStartedEventData{
		EventData: EventData{
			Status: StatusSucceeded,
		},
	})

	expectedErrorData := ReleaseFinishedEventData{
		EventData: EventData{
			Labels:  nil,
			Status:  StatusErrored,
			Result:  ResultFailed,
			Message: "The deployment strategy ??? is invalid",
		},
	}

	require.Equal(t, 1, len(mockedBaseHandler.sentCloudEvents))
	require.Equal(t, 1, len(mockedBaseHandler.handledErrorEvents))
	assert.Equal(t, expectedReleaseStartedEvent, mockedBaseHandler.sentCloudEvents[0])
	assert.Equal(t, expectedErrorData, mockedBaseHandler.handledErrorEvents[0])

}
