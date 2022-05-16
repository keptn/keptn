package sdk

import (
	"github.com/google/uuid"
	"github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/go-utils/pkg/common/strutils"
	"github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"testing"
)

func Test_WhenReceivingAnEvent_StartedEventAndFinishedEventsAreSent(t *testing.T) {
	taskHandler := &TaskHandlerMock{}
	taskHandler.ExecuteFunc = func(keptnHandle IKeptn, event KeptnEvent) (interface{}, *Error) { return FakeTaskData{}, nil }
	fakeKeptn := NewFakeKeptn("fake")
	fakeKeptn.AddTaskHandler("sh.keptn.event.faketask.triggered", taskHandler)
	fakeKeptn.NewEvent(models.KeptnContextExtendedCE{
		Data:           v0_2_0.EventData{Project: "prj", Stage: "stg", Service: "svc"},
		ID:             "id",
		Shkeptncontext: "context",
		Source:         strutils.Stringp("source"),
		Type:           strutils.Stringp("sh.keptn.event.faketask.triggered"),
	})

	fakeKeptn.AssertNumberOfEventSent(t, 2)
	fakeKeptn.AssertSentEventType(t, 0, "sh.keptn.event.faketask.started")
	fakeKeptn.AssertSentEventType(t, 1, "sh.keptn.event.faketask.finished")

}

func Test_WhenReceivingEvent_OnlyStartedEventIsSent(t *testing.T) {
	taskHandler := &TaskHandlerMock{}
	taskHandler.ExecuteFunc = func(keptnHandle IKeptn, event KeptnEvent) (interface{}, *Error) { return FakeTaskData{}, nil }
	fakeKeptn := NewFakeKeptn("fake")
	fakeKeptn.AddTaskHandler("sh.keptn.event.faketask.triggered", taskHandler)
	fakeKeptn.NewEvent(newTestTaskTriggeredEvent())
	fakeKeptn.AssertNumberOfEventSent(t, 0)
}

func Test_WhenReceivingBadEvent_NoEventIsSent(t *testing.T) {
	taskHandler := &TaskHandlerMock{}
	taskHandler.ExecuteFunc = func(keptnHandle IKeptn, event KeptnEvent) (interface{}, *Error) { return FakeTaskData{}, nil }
	fakeKeptn := NewFakeKeptn("fake")
	fakeKeptn.AddTaskHandler("sh.keptn.event.faketask.triggered", taskHandler)
	fakeKeptn.NewEvent(newTestTaskBadTriggeredEvent())
	fakeKeptn.AssertNumberOfEventSent(t, 0)
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
