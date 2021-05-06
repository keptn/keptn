package sdk_test

import (
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/google/uuid"
	"github.com/keptn/keptn/remediation-service/internal/sdk"
	"github.com/keptn/keptn/remediation-service/internal/sdk/fake"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_WhenReceivingAnEvent_StartedEventAndFinishedEventsAreSent(t *testing.T) {

	taskHandler := &fake.TaskHandlerMock{}
	taskHandler.ExecuteFunc = func(keptnHandle sdk.IKeptn, ce interface{}) (interface{}, *sdk.Error) {
		return FakeTaskData{}, nil
	}
	taskHandler.GetTriggeredDataFunc = func() interface{} {
		return FakeTaskData{}
	}

	taskEntry := sdk.TaskEntry{
		TaskHandler: taskHandler,
	}

	taskEntries := map[string]sdk.TaskEntry{"sh.keptn.event.faketask.triggered": taskEntry}

	eventReceiver := &fake.TestReceiver{}
	eventSender := &fake.EventSenderMock{}

	eventSender.SendEventFunc = func(eventMoqParam event.Event) error {
		return nil
	}

	taskRegistry := sdk.TaskRegistry{
		Entries: taskEntries,
	}

	keptn := sdk.Keptn{
		EventSender:   eventSender,
		EventReceiver: eventReceiver,
		TaskRegistry:  taskRegistry,
	}

	keptn.Start()
	eventReceiver.NewEvent(newTestTaskTriggeredEvent())

	require.Equal(t, 2, len(eventSender.SendEventCalls()))
	require.Equal(t, "sh.keptn.event.faketask.started", eventSender.SendEventCalls()[0].EventMoqParam.Type())
	require.Equal(t, "sh.keptn.event.faketask.finished", eventSender.SendEventCalls()[1].EventMoqParam.Type())
}

func newTestTaskTriggeredEvent() cloudevents.Event {
	c := cloudevents.NewEvent()
	c.SetID(uuid.New().String())
	c.SetType("sh.keptn.event.faketask.triggered")
	c.SetDataContentType(cloudevents.ApplicationJSON)
	c.SetExtension(sdk.KeptnContextCEExtension, "keptncontext")
	c.SetExtension(sdk.TriggeredIDCEExtension, "ID")
	c.SetSource("unittest")
	c.SetData(cloudevents.ApplicationJSON, FakeTaskData{})
	return c
}

type FakeTaskData struct {
}
