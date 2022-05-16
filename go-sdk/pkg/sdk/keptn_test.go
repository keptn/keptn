package sdk

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/go-utils/pkg/common/strutils"
	"github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/cp-connector/pkg/controlplane"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func Test_WhenReceivingAnEvent_StartedEventAndFinishedEventsAreSent(t *testing.T) {
	taskHandler := &TaskHandlerMock{}
	taskHandler.ExecuteFunc = func(keptnHandle IKeptn, event KeptnEvent) (interface{}, *Error) { return FakeTaskData{}, nil }
	taskEntry := TaskEntry{TaskHandler: taskHandler}
	taskEntries := map[string]TaskEntry{"sh.keptn.event.faketask.triggered": taskEntry}
	taskRegistry := &TaskRegistry{Entries: taskEntries}

	testSubscriptionSource := controlplane.NewFixedSubscriptionSource(controlplane.WithFixedSubscriptions(models.EventSubscription{Event: "sh.keptn.event.faketask.triggered"}))
	testEventSource := NewTestEventSource()
	cp := controlplane.New(testSubscriptionSource, testEventSource, nil)

	keptn := Keptn{
		controlPlane:           cp,
		taskRegistry:           taskRegistry,
		automaticEventResponse: true,
		logger:                 newDefaultLogger(),
		healthEndpointRunner:   noOpHealthEndpointRunner,
	}

	go keptn.Start()
	<-testEventSource.Started

	testEventSource.NewEvent(controlplane.EventUpdate{
		KeptnEvent: models.KeptnContextExtendedCE{
			Data:           v0_2_0.EventData{Project: "prj", Stage: "stg", Service: "svc"},
			ID:             "id",
			Shkeptncontext: "context",
			Source:         strutils.Stringp("source"),
			Type:           strutils.Stringp("sh.keptn.event.faketask.triggered"),
		},
		MetaData: controlplane.EventUpdateMetaData{Subject: "sh.keptn.event.faketask.triggered"},
	})

	require.Eventuallyf(t, func() bool {
		fmt.Println(testEventSource.GetNumberOfSetEvents())
		return testEventSource.GetNumberOfSetEvents() == 2
	}, time.Second*10, time.Second, "error message %s", "formatted")

	require.Eventuallyf(t, func() bool {
		return *testEventSource.GetSentEvents()[0].Type == "sh.keptn.event.faketask.started"
	}, time.Second*10, time.Second, "error message %s", "formatted")

	require.Eventuallyf(t, func() bool {
		return *testEventSource.GetSentEvents()[1].Type == "sh.keptn.event.faketask.finished"
	}, time.Second*10, time.Second, "error message %s", "formatted")

}

func Test_WhenReceivingEvent_OnlyStartedEventIsSent(t *testing.T) {
	taskHandler := &TaskHandlerMock{}
	taskHandler.ExecuteFunc = func(keptnHandle IKeptn, event KeptnEvent) (interface{}, *Error) { return FakeTaskData{}, nil }
	taskEntry := TaskEntry{TaskHandler: taskHandler}
	taskEntries := map[string]TaskEntry{"sh.keptn.event.faketask.triggered": taskEntry}
	taskRegistry := &TaskRegistry{Entries: taskEntries}
	testSubscriptionSource := controlplane.NewFixedSubscriptionSource(controlplane.WithFixedSubscriptions(models.EventSubscription{Event: "sh.keptn.event.faketask.triggered"}))
	testEventSource := NewTestEventSource()
	cp := controlplane.New(testSubscriptionSource, testEventSource, nil)

	keptn := Keptn{
		controlPlane:           cp,
		taskRegistry:           taskRegistry,
		automaticEventResponse: true,
		logger:                 newDefaultLogger(),
		healthEndpointRunner:   noOpHealthEndpointRunner,
	}

	go keptn.Start()
	<-testEventSource.Started

	testEventSource.NewEvent(controlplane.EventUpdate{
		KeptnEvent: newTestTaskTriggeredEvent(),
		MetaData:   controlplane.EventUpdateMetaData{Subject: "sh.keptn.event.faketask.triggered"},
	})

	require.Eventuallyf(t, func() bool {
		return len(testEventSource.GetSentEvents()) == 0
	}, time.Second*10, time.Second, "error message %s", "formatted")

}

func Test_WhenReceivingBadEvent_NoEventIsSent(t *testing.T) {
	taskHandler := &TaskHandlerMock{}
	taskHandler.ExecuteFunc = func(keptnHandle IKeptn, event KeptnEvent) (interface{}, *Error) { return FakeTaskData{}, nil }
	taskEntry := TaskEntry{TaskHandler: taskHandler}
	taskEntries := map[string]TaskEntry{"sh.keptn.event.faketask.triggered": taskEntry}
	taskRegistry := &TaskRegistry{Entries: taskEntries}
	testSubscriptionSource := controlplane.NewFixedSubscriptionSource(controlplane.WithFixedSubscriptions(models.EventSubscription{Event: "sh.keptn.event.faketask.triggered"}))
	testEventSource := NewTestEventSource()
	cp := controlplane.New(testSubscriptionSource, testEventSource, nil)

	keptn := Keptn{
		controlPlane:           cp,
		taskRegistry:           taskRegistry,
		automaticEventResponse: true,
		logger:                 newDefaultLogger(),
		healthEndpointRunner:   noOpHealthEndpointRunner,
	}

	go keptn.Start()
	<-testEventSource.Started

	testEventSource.NewEvent(controlplane.EventUpdate{
		KeptnEvent: newTestTaskBadTriggeredEvent(),
		MetaData:   controlplane.EventUpdateMetaData{Subject: "sh.keptn.event.faketask.finished.triggered"},
	})

	require.Eventuallyf(t, func() bool {
		return len(testEventSource.GetSentEvents()) == 0
	}, time.Second*10, time.Second, "error message %s", "formatted")
}

func newTestTaskTriggeredEvent() models.KeptnContextExtendedCE {
	return models.KeptnContextExtendedCE{
		Contenttype:    "application/json",
		Data:           FakeTaskData{},
		ID:             uuid.New().String(),
		Shkeptncontext: "keptncontext",
		Triggeredid:    "ID",
		GitCommitID:    "mycommitid",
		Source:         strutils.Stringp("unittest"),
		Type:           strutils.Stringp("sh.keptn.event.faketask.triggered"),
	}
}

func newTestTaskBadTriggeredEvent() models.KeptnContextExtendedCE {
	return models.KeptnContextExtendedCE{
		Contenttype:    "application/json",
		Data:           FakeTaskData{},
		ID:             uuid.New().String(),
		Shkeptncontext: "keptncontext",
		Triggeredid:    "ID",
		GitCommitID:    "mycommitid",
		Source:         strutils.Stringp("unittest"),
		Type:           strutils.Stringp("sh.keptn.event.faketask.finished.triggered"),
	}
}

type FakeTaskData struct {
}
