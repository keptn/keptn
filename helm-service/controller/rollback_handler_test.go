package controller

import (
	"context"
	"errors"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/golang/mock/gomock"
	keptn "github.com/keptn/go-utils/pkg/lib"
	"github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/helm-service/mocks"
	"github.com/keptn/keptn/helm-service/pkg/helm"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestCreateRollbackHandler(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockedBaseHandler := NewMockedHandler(createKeptn(), "")
	mockedMesh := mocks.NewMockMesh(ctrl)
	mockedConfigurationChanger := mocks.NewMockIConfigurationChanger(ctrl)
	instance := NewRollbackHandler(mockedBaseHandler, mockedMesh, mockedConfigurationChanger)
	assert.NotNil(t, instance)
	assert.NotNil(t, instance.Handler)
	assert.NotNil(t, instance.mesh)
	assert.NotNil(t, instance.configurationChanger)
}

func TestHandleRollbackEvent(t *testing.T) {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	ctx, _ := context.WithCancel(cloudevents.WithEncodingStructured(context.WithValue(context.Background(), "Wg", wg)))
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockedBaseHandler := NewMockedHandler(createKeptn(), "")
	mockedMesh := mocks.NewMockMesh(ctrl)
	mockedConfigurationChanger := mocks.NewMockIConfigurationChanger(ctrl)
	testGenChart := helm.GetTestGeneratedChart()
	testUserChart := helm.GetTestUserChart()

	instance := &RollbackHandler{
		Handler:              mockedBaseHandler,
		mesh:                 mockedMesh,
		configurationChanger: mockedConfigurationChanger,
	}

	eventData := v0_2_0.EventData{
		Project: "my-project",
		Stage:   "my-stage",
		Service: "my-service",
	}
	rollbackTriggeredEventData := v0_2_0.RollbackTriggeredEventData{
		EventData: eventData,
	}

	ce := cloudevents.NewEvent()
	_ = ce.SetData(cloudevents.ApplicationJSON, rollbackTriggeredEventData)

	mockedConfigurationChanger.EXPECT().UpdateChart(gomock.Any(), gomock.Any(), gomock.Any()).Return(&testGenChart, "version", nil)
	instance.HandleEvent(ctx, ce)

	assert.Equal(t, 1, len(mockedBaseHandler.upgradeChartInvocations))
	assert.Equal(t, keptn.Duplicate, mockedBaseHandler.upgradeChartInvocations[0].strategy)
	assert.Equal(t, &testGenChart, mockedBaseHandler.upgradeChartInvocations[0].ch)
	assert.Equal(t, eventData, mockedBaseHandler.upgradeChartInvocations[0].event)

	assert.Equal(t, 1, len(mockedBaseHandler.upgradeChartWithReplicasInvocations))
	assert.Equal(t, keptn.Duplicate, mockedBaseHandler.upgradeChartWithReplicasInvocations[0].strategy)
	assert.Equal(t, &testUserChart, mockedBaseHandler.upgradeChartWithReplicasInvocations[0].ch)
	assert.Equal(t, eventData, mockedBaseHandler.upgradeChartWithReplicasInvocations[0].event)
	assert.Equal(t, 0, mockedBaseHandler.upgradeChartWithReplicasInvocations[0].replicas)

	assert.Equal(t, 0, len(mockedBaseHandler.handledErrorEvents))

	assert.Equal(t, "sh.keptn.event.rollback.started", mockedBaseHandler.sentCloudEvents[0].Type())
	assert.Equal(t, "sh.keptn.event.rollback.finished", mockedBaseHandler.sentCloudEvents[1].Type())
}

func TestHandleRollbackEvent_UpdatingChartFails(t *testing.T) {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	ctx, _ := context.WithCancel(cloudevents.WithEncodingStructured(context.WithValue(context.Background(), "Wg", wg)))
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockedBaseHandler := NewMockedHandler(createKeptn(), "")
	mockedMesh := mocks.NewMockMesh(ctrl)
	mockedConfigurationChanger := mocks.NewMockIConfigurationChanger(ctrl)

	expectedEventData := v0_2_0.RollbackFinishedEventData{
		EventData: v0_2_0.EventData{
			Project: "my-project",
			Stage:   "my-stage",
			Service: "my-service",
			Status:  v0_2_0.StatusErrored,
			Result:  v0_2_0.ResultFailed,
			Message: "Whoops...",
		},
	}

	instance := &RollbackHandler{
		Handler:              mockedBaseHandler,
		mesh:                 mockedMesh,
		configurationChanger: mockedConfigurationChanger,
	}

	eventData := v0_2_0.EventData{
		Project: "my-project",
		Stage:   "my-stage",
		Service: "my-service",
	}
	rollbackTriggeredEventData := v0_2_0.RollbackTriggeredEventData{
		EventData: eventData,
	}

	ce := cloudevents.NewEvent()
	_ = ce.SetData(cloudevents.ApplicationJSON, rollbackTriggeredEventData)

	mockedConfigurationChanger.EXPECT().UpdateChart(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, "", errors.New("Whoops..."))

	instance.HandleEvent(ctx, ce)

	assert.Equal(t, 1, len(mockedBaseHandler.handledErrorEvents))
	assert.Equal(t, 2, len(mockedBaseHandler.sentCloudEvents))
	assert.Equal(t, "sh.keptn.event.rollback.started", mockedBaseHandler.sentCloudEvents[0].Type())
	assert.Equal(t, "sh.keptn.event.rollback.finished", mockedBaseHandler.sentCloudEvents[1].Type())

	assert.Equal(t, expectedEventData, mockedBaseHandler.handledErrorEvents[0])
	assert.Equal(t, expectedEventData, mockedBaseHandler.handledErrorEvents[0])

}
