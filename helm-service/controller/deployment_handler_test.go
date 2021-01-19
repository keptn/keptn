package controller

import (
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/golang/mock/gomock"
	keptn "github.com/keptn/go-utils/pkg/lib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/helm-service/mocks"
	"testing"
)

func TestHandleEventWithNoConfigurationChangeAndDirectDeploymentStrategy(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockedBaseHandler := NewMockedHandler(createKeptn(), "")
	mockedOnboarder := mocks.NewMockOnboarder(ctrl)
	mockedChartGenerator := mocks.NewMockChartGenerator(ctrl)

	deploymentHandler := DeploymentHandler{
		Handler:               mockedBaseHandler,
		mesh:                  mocks.NewMockMesh(ctrl),
		generatedChartHandler: mockedChartGenerator,
		onboarder:             mockedOnboarder,
	}

	deploymentTriggeredEventData := keptnv2.DeploymentTriggeredEventData{
		EventData: keptnv2.EventData{
			Project: "my-project",
			Stage:   "my-stage",
			Service: "my-service",
		},
		ConfigurationChange: keptnv2.ConfigurationChange{},
		Deployment: keptnv2.DeploymentWithStrategy{
			DeploymentStrategy: keptn.Direct.String(),
		},
	}

	ce := cloudevents.NewEvent()
	_ = ce.SetData(cloudevents.ApplicationJSON, deploymentTriggeredEventData)
	deploymentHandler.HandleEvent(ce)

	expectedActionFinishedEvent := cloudevents.NewEvent()
	expectedActionFinishedEvent.SetType("sh.keptn.event.deployment.finished")
	expectedActionFinishedEvent.SetSource("helm-service")
	expectedActionFinishedEvent.SetDataContentType(cloudevents.ApplicationJSON)
	expectedActionFinishedEvent.SetExtension("triggeredid", "")
	expectedActionFinishedEvent.SetExtension("shkeptncontext", "")
	expectedActionFinishedEvent.SetData(cloudevents.ApplicationJSON, keptnv2.DeploymentFinishedEventData{
		EventData: keptnv2.EventData{
			Project: "my-project",
			Stage:   "my-stage",
			Service: "my-service",
			Status:  keptnv2.StatusSucceeded,
			Result:  keptnv2.ResultPass,
			Message: "Successfully deployed",
		},
		Deployment: keptnv2.DeploymentData{
			DeploymentStrategy:   "direct",
			DeploymentURIsLocal:  []string{"http://my-service.my-project-my-stage:80"},
			DeploymentURIsPublic: []string{"http://my-service.my-project-my-stage.svc.cluster.local:80"},
			DeploymentNames:      []string{"direct"},
			GitCommit:            "USER_CHART_GIT_ID",
		},
	})

	require.Equal(t, 2, len(mockedBaseHandler.sentCloudEvents))
	assert.Equal(t, expectedActionFinishedEvent, mockedBaseHandler.sentCloudEvents[1])
	require.Equal(t, 2, len(mockedBaseHandler.upgradeChartInvocations))
	assert.Equal(t, "carts", mockedBaseHandler.upgradeChartInvocations[0].ch.Metadata.Name)
	assert.Equal(t, deploymentTriggeredEventData.EventData, mockedBaseHandler.upgradeChartInvocations[0].event)
	assert.Equal(t, keptn.Direct, mockedBaseHandler.upgradeChartInvocations[0].strategy)
	assert.Equal(t, "carts-generated", mockedBaseHandler.upgradeChartInvocations[1].ch.Metadata.Name)
	assert.Equal(t, deploymentTriggeredEventData.EventData, mockedBaseHandler.upgradeChartInvocations[1].event)
	assert.Equal(t, keptn.Direct, mockedBaseHandler.upgradeChartInvocations[1].strategy)
}
