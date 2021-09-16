package sdk_test

import (
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/google/uuid"
	"github.com/keptn/keptn/go-sdk/pkg/sdk"
	"github.com/keptn/keptn/go-sdk/pkg/sdk/fake"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func Test_WhenReceivingAnEvent_StartedEventAndFinishedEventsAreSent(t *testing.T) {
	taskHandler := &fake.TaskHandlerMock{}
	taskHandler.ExecuteFunc = func(keptnHandle sdk.IKeptn, ce interface{}, eventType string) (interface{}, *sdk.Error) {
		return FakeTaskData{}, nil
	}

	taskEntry := sdk.TaskEntry{
		TaskHandler:    taskHandler,
		ReceivingEvent: &FakeTaskData{},
	}

	taskEntries := map[string]sdk.TaskEntry{"sh.keptn.event.faketask.triggered": taskEntry}

	eventReceiver := &fake.TestReceiver{}
	eventSender := &fake.EventSenderMock{}

	eventSender.SendEventFunc = func(eventMoqParam event.Event) error {
		return nil
	}

	taskRegistry := &sdk.TaskRegistry{
		Entries: taskEntries,
	}

	keptn := sdk.Keptn{
		EventSender:   eventSender,
		EventReceiver: eventReceiver,
		TaskRegistry:  taskRegistry,
	}

	keptn.Start()
	eventReceiver.NewEvent(newTestTaskTriggeredEvent())

	require.Eventuallyf(t, func() bool {
		return len(eventSender.SendEventCalls()) == 2
	}, time.Second, 10*time.Millisecond, "error message %s", "formatted")

	require.Eventuallyf(t, func() bool {
		return eventSender.SendEventCalls()[0].EventMoqParam.Type() == "sh.keptn.event.faketask.started"
	}, time.Second, 10*time.Millisecond, "error message %s", "formatted")

	require.Eventuallyf(t, func() bool {
		return eventSender.SendEventCalls()[1].EventMoqParam.Type() == "sh.keptn.event.faketask.finished"
	}, time.Second, 10*time.Millisecond, "error message %s", "formatted")

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
