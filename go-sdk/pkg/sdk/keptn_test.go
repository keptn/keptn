package sdk

import (
	"fmt"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func Test_WhenReceivingAnEvent_StartedEventAndFinishedEventsAreSent(t *testing.T) {
	taskHandler := &TaskHandlerMock{}
	taskHandler.ExecuteFunc = func(keptnHandle IKeptn, event KeptnEvent) (interface{}, *Error) {
		return FakeTaskData{}, nil
	}

	taskEntry := TaskEntry{
		TaskHandler: taskHandler,
	}

	taskEntries := map[string]TaskEntry{"sh.keptn.event.faketask.triggered": taskEntry}

	eventReceiver := &TestReceiver{}
	eventSender := &EventSenderMock{}

	eventSender.SendEventFunc = func(eventMoqParam event.Event) error {
		return nil
	}

	taskRegistry := &TaskRegistry{
		Entries: taskEntries,
	}

	keptn := Keptn{
		eventSender:            eventSender,
		eventReceiver:          eventReceiver,
		taskRegistry:           taskRegistry,
		automaticEventResponse: true,
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

func Test_WhenReceivingEvent_OnlyStartedEventIsSent(t *testing.T) {
	taskHandler := &TaskHandlerMock{}
	taskHandler.ExecuteFunc = func(keptnHandle IKeptn, event KeptnEvent) (interface{}, *Error) {
		return FakeTaskData{}, nil
	}

	taskEntry := TaskEntry{
		TaskHandler: taskHandler,
	}

	taskEntries := map[string]TaskEntry{"sh.keptn.event.faketask.triggered": taskEntry}

	eventReceiver := &TestReceiver{}
	eventSender := &EventSenderMock{}

	eventSender.SendEventFunc = func(eventMoqParam event.Event) error {
		return nil
	}

	taskRegistry := &TaskRegistry{
		Entries: taskEntries,
	}

	keptn := Keptn{
		eventSender:            eventSender,
		eventReceiver:          eventReceiver,
		taskRegistry:           taskRegistry,
		automaticEventResponse: false,
	}

	keptn.Start()
	eventReceiver.NewEvent(newTestTaskTriggeredEvent())

	require.Eventuallyf(t, func() bool {
		fmt.Println(len(eventSender.SendEventCalls()))
		return len(eventSender.SendEventCalls()) == 0
	}, time.Second, 10*time.Millisecond, "error message %s", "formatted")

}

func newTestTaskTriggeredEvent() cloudevents.Event {
	c := cloudevents.NewEvent()
	c.SetID(uuid.New().String())
	c.SetType("sh.keptn.event.faketask.triggered")
	c.SetDataContentType(cloudevents.ApplicationJSON)
	c.SetExtension(KeptnContextCEExtension, "keptncontext")
	c.SetExtension(TriggeredIDCEExtension, "ID")
	c.SetSource("unittest")
	c.SetData(cloudevents.ApplicationJSON, FakeTaskData{})
	return c
}

type FakeTaskData struct {
}
