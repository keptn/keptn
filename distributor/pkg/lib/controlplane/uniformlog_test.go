package controlplane

import (
	"context"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/keptn/go-utils/pkg/api/models"
	fakeapi "github.com/keptn/go-utils/pkg/api/utils/fake"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestEventUniformLog_LogEvent(t *testing.T) {
	fakeLogHandler, uniformLog := getTestUniformLogger()

	// send a new event with unhandled type -> should not be logged
	newEvent := event.New()
	newEvent.SetType("random-event")

	err := uniformLog.OnEvent(newEvent)

	require.Nil(t, err)
	require.Empty(t, fakeLogHandler.LogCalls())

	// send sh.keptn.log.error event
	newEvent.SetType(keptnv2.ErrorLogEventName)

	logEventData := &keptnv2.ErrorLogEvent{
		Message: "my-message",
		Task:    "my-task",
	}
	newEvent.SetExtension("shkeptncontext", "my-context")
	newEvent.SetExtension("triggeredid", "my-triggered-id")
	newEvent.SetData(event.ApplicationJSON, logEventData)

	err = uniformLog.OnEvent(newEvent)

	require.Nil(t, err)
	require.NotEmpty(t, fakeLogHandler.LogCalls())
	require.Equal(t, models.LogEntry{
		IntegrationID: "my-id",
		Message:       "my-message",
		KeptnContext:  "my-context",
		Task:          "my-task",
		TriggeredID:   "my-triggered-id",
	}, fakeLogHandler.LogCalls()[0].Logs[0])
}

func TestEventUniformLog_LogEventInvalidPayload(t *testing.T) {
	fakeLogHandler, uniformLog := getTestUniformLogger()

	newEvent := event.New()
	newEvent.SetType(keptnv2.ErrorLogEventName)

	// send sh.keptn.log.error event with invalid payload -> should result in an error
	newEvent.SetData(event.TextJSON, "invalid")

	err := uniformLog.OnEvent(newEvent)

	require.NotNil(t, err)
	require.Empty(t, fakeLogHandler.LogCalls())
}

func TestEventUniformLog_TaskFinishedWithErroredStatus(t *testing.T) {
	fakeLogHandler, uniformLog := getTestUniformLogger()

	newEvent := event.New()
	newEvent.SetType(keptnv2.GetFinishedEventType(keptnv2.DeploymentTaskName))
	newEvent.SetExtension("shkeptncontext", "my-context")
	newEvent.SetExtension("triggeredid", "my-triggered-id")

	finishedData := &keptnv2.EventData{
		Status:  keptnv2.StatusErrored,
		Message: "my-message",
	}
	newEvent.SetData(event.ApplicationJSON, finishedData)

	err := uniformLog.OnEvent(newEvent)

	require.Nil(t, err)
	require.Len(t, fakeLogHandler.LogCalls(), 1)
	require.Equal(t, models.LogEntry{
		IntegrationID: "my-id",
		Message:       "my-message",
		KeptnContext:  "my-context",
		Task:          keptnv2.DeploymentTaskName,
		TriggeredID:   "my-triggered-id",
	}, fakeLogHandler.LogCalls()[0].Logs[0])
}

func TestEventUniformLog_TaskFinishedWithSucceedStatus(t *testing.T) {
	fakeLogHandler, uniformLog := getTestUniformLogger()

	newEvent := event.New()
	newEvent.SetType(keptnv2.GetFinishedEventType(keptnv2.DeploymentTaskName))

	finishedData := &keptnv2.EventData{
		Status:  keptnv2.StatusSucceeded,
		Message: "my-message",
	}
	newEvent.SetData(event.ApplicationJSON, finishedData)

	err := uniformLog.OnEvent(newEvent)

	require.Nil(t, err)
	require.Empty(t, fakeLogHandler.LogCalls())
}

func TestEventUniformLog_TaskFinishedWithSucceedStatusAndResultFail(t *testing.T) {
	fakeLogHandler, uniformLog := getTestUniformLogger()

	newEvent := event.New()
	newEvent.SetType(keptnv2.GetFinishedEventType(keptnv2.DeploymentTaskName))

	finishedData := &keptnv2.EventData{
		Status:  keptnv2.StatusSucceeded,
		Result:  keptnv2.ResultFailed,
		Message: "my-message",
	}
	newEvent.SetData(event.ApplicationJSON, finishedData)

	err := uniformLog.OnEvent(newEvent)

	require.Nil(t, err)
	require.Empty(t, fakeLogHandler.LogCalls())
}

func TestEventUniformLog_TaskFinishedWithInvalidData(t *testing.T) {
	fakeLogHandler, uniformLog := getTestUniformLogger()

	newEvent := event.New()
	newEvent.SetType(keptnv2.GetFinishedEventType(keptnv2.DeploymentTaskName))

	newEvent.SetData(event.TextJSON, "invalid")

	err := uniformLog.OnEvent(newEvent)

	require.NotNil(t, err)
	require.Empty(t, fakeLogHandler.LogCalls())
}

func getTestUniformLogger() (*fakeapi.ILogHandlerMock, *EventUniformLog) {
	fakeLogHandler := &fakeapi.ILogHandlerMock{
		FlushFunc: func() error {
			return nil
		},
		LogFunc: func(logs []models.LogEntry) {

		},
		StartFunc: func(ctx context.Context) {

		},
	}
	uniformLog := NewEventUniformLog(
		"my-id",
		fakeLogHandler,
	)

	eventsChannel := make(chan event.Event)

	uniformLog.Start(context.TODO(), eventsChannel)
	return fakeLogHandler, uniformLog
}
