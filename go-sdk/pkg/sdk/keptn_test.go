package sdk

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/go-utils/pkg/common/strutils"
	"github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/stretchr/testify/require"
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

func Test_WhenReceivingAnEvent_TaskHandlerFails(t *testing.T) {
	taskHandler := &TaskHandlerMock{}
	taskHandler.ExecuteFunc = func(keptnHandle IKeptn, event KeptnEvent) (interface{}, *Error) {
		return nil, &Error{
			StatusType: v0_2_0.StatusErrored,
			ResultType: v0_2_0.ResultFailed,
			Message:    "something went wrong",
			Err:        fmt.Errorf("something went wrong"),
		}
	}
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
	fakeKeptn.AssertSentEventStatus(t, 1, v0_2_0.StatusErrored)
	fakeKeptn.AssertSentEventResult(t, 1, v0_2_0.ResultFailed)
}

func Test_WhenReceivingBadEvent_NoEventIsSent(t *testing.T) {
	taskHandler := &TaskHandlerMock{}
	taskHandler.ExecuteFunc = func(keptnHandle IKeptn, event KeptnEvent) (interface{}, *Error) { return FakeTaskData{}, nil }
	fakeKeptn := NewFakeKeptn("fake")
	fakeKeptn.AddTaskHandler("sh.keptn.event.faketask.triggered", taskHandler)
	fakeKeptn.NewEvent(newTestTaskBadTriggeredEvent())
	fakeKeptn.AssertNumberOfEventSent(t, 0)
}

func Test_WhenReceivingAnEvent_AndNoFilterMatches_NoEventIsSent(t *testing.T) {
	taskHandler := &TaskHandlerMock{}
	taskHandler.ExecuteFunc = func(keptnHandle IKeptn, event KeptnEvent) (interface{}, *Error) { return FakeTaskData{}, nil }
	fakeKeptn := NewFakeKeptn("fake")
	fakeKeptn.AddTaskHandler("sh.keptn.event.faketask.triggered", taskHandler, func(keptnHandle IKeptn, event KeptnEvent) bool { return false })
	fakeKeptn.NewEvent(models.KeptnContextExtendedCE{
		Data:           v0_2_0.EventData{Project: "prj", Stage: "stg", Service: "svc"},
		ID:             "id",
		Shkeptncontext: "context",
		Source:         strutils.Stringp("source"),
		Type:           strutils.Stringp("sh.keptn.event.faketask.triggered"),
	})

	fakeKeptn.AssertNumberOfEventSent(t, 0)
}

func Test_NoFinishedEventDataProvided(t *testing.T) {
	taskHandler := &TaskHandlerMock{}
	taskHandler.ExecuteFunc = func(keptnHandle IKeptn, event KeptnEvent) (interface{}, *Error) {
		return nil, nil
	}
	fakeKeptn := NewFakeKeptn("fake")
	fakeKeptn.AddTaskHandler("sh.keptn.event.faketask.triggered", taskHandler)
	fakeKeptn.NewEvent(models.KeptnContextExtendedCE{
		Data:           v0_2_0.EventData{Project: "prj", Stage: "stg", Service: "svc"},
		ID:             "id",
		Shkeptncontext: "context",
		Source:         strutils.Stringp("source"),
		Type:           strutils.Stringp("sh.keptn.event.faketask.triggered"),
	})

	fakeKeptn.AssertNumberOfEventSent(t, 1)
	fakeKeptn.AssertSentEventType(t, 0, "sh.keptn.event.faketask.started")
}

func Test_InitialRegistrationData(t *testing.T) {
	keptn := Keptn{env: envConfig{
		PubSubTopic:       "sh.keptn.event.task1.triggered,sh.keptn.event.task2.triggered",
		Location:          "localhost",
		Version:           "v1",
		K8sDeploymentName: "k8s-deployment",
		K8sNamespace:      "k8s-namespace",
		K8sPodName:        "k8s-podname",
		K8sNodeName:       "k8s-nodename",
	}}

	regData := keptn.RegistrationData()
	require.Equal(t, "v1", regData.MetaData.IntegrationVersion)
	require.Equal(t, "localhost", regData.MetaData.Location)
	require.Equal(t, "k8s-deployment", regData.MetaData.KubernetesMetaData.DeploymentName)
	require.Equal(t, "k8s-namespace", regData.MetaData.KubernetesMetaData.Namespace)
	require.Equal(t, "k8s-podname", regData.MetaData.KubernetesMetaData.PodName)
	require.Equal(t, "k8s-nodename", regData.MetaData.Hostname)
	require.Equal(t, []models.EventSubscription{{Event: "sh.keptn.event.task1.triggered"}, {Event: "sh.keptn.event.task2.triggered"}}, regData.Subscriptions)
}

func Test_InitialRegistrationData_EmptyPubSubTopics(t *testing.T) {
	keptn := Keptn{env: envConfig{PubSubTopic: ""}}
	regData := keptn.RegistrationData()
	require.Equal(t, 0, len(regData.Subscriptions))
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
type TaskHandlerMock struct {
	// ExecuteFunc mocks the Execute method.
	ExecuteFunc func(keptnHandle IKeptn, event KeptnEvent) (interface{}, *Error)
}

func (mock *TaskHandlerMock) Execute(keptnHandle IKeptn, event KeptnEvent) (interface{}, *Error) {
	if mock.ExecuteFunc == nil {
		panic("TaskHandlerMock.ExecuteFunc: method is nil but taskHandler.Execute was just called")
	}
	return mock.ExecuteFunc(keptnHandle, event)
}
