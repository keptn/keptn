package sdk

import (
	"context"
	"errors"
	"fmt"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/keptn/go-utils/pkg/api/models"
	api "github.com/keptn/go-utils/pkg/api/utils"
	"github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/cp-connector/pkg/controlplane"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"path/filepath"
	"testing"
)

type FakeKeptn struct {
	TestResourceHandler ResourceHandler
	SentEvents          []models.KeptnContextExtendedCE
	Keptn               *Keptn
}

func (f *FakeKeptn) GetResourceHandler() ResourceHandler {
	if f.TestResourceHandler == nil {
		return &TestResourceHandler{}
	}
	return f.TestResourceHandler
}

func (f *FakeKeptn) NewEvent(event models.KeptnContextExtendedCE) {
	ctx := context.WithValue(context.TODO(), controlplane.EventSenderKey, controlplane.EventSender(f.fakeSender))
	ctx = context.WithValue(ctx, gracefulShutdownKey, &nopWG{})
	f.Keptn.OnEvent(ctx, event)
}

func (f *FakeKeptn) AssertNumberOfEventSent(t *testing.T, numOfEvents int) {
	require.Equalf(t, numOfEvents, len(f.SentEvents), "number of events expected: %d got: %d", numOfEvents, len(f.SentEvents))
}

func (f *FakeKeptn) AssertSentEvent(t *testing.T, eventIndex int, assertFn func(ce models.KeptnContextExtendedCE) bool) {
	if eventIndex >= len(f.SentEvents) {
		t.Fatalf("unable to assert sent event with index %d: too less events sent", eventIndex)
	}
	require.True(t, assertFn(f.SentEvents[eventIndex]))
}

func (f *FakeKeptn) AssertSentEventType(t *testing.T, eventIndex int, eventType string) {
	if eventIndex >= len(f.SentEvents) {
		t.Fatalf("unable to assert sent event with index %d: too less events sent", eventIndex)
	}
	require.Equalf(t, eventType, *f.SentEvents[eventIndex].Type, "event type expected: %s got %s", eventType, *f.SentEvents[eventIndex].Type)
}

func (f *FakeKeptn) AssertSentEventStatus(t *testing.T, eventIndex int, status v0_2_0.StatusType) {
	if eventIndex >= len(f.SentEvents) {
		t.Fatalf("unable to assert sent event with index %d: too less events sent", eventIndex)
	}
	eventData := v0_2_0.EventData{}
	v0_2_0.EventDataAs(f.SentEvents[eventIndex], &eventData)
	require.Equal(t, status, eventData.Status)
}

func (f *FakeKeptn) AssertSentEventResult(t *testing.T, eventIndex int, result v0_2_0.ResultType) {
	if eventIndex >= len(f.SentEvents) {
		t.Fatalf("unable to assert sent event with index %d: too less events sent", eventIndex)
	}
	eventData := v0_2_0.EventData{}
	v0_2_0.EventDataAs(f.SentEvents[eventIndex], &eventData)
	require.Equal(t, result, eventData.Result)
}

func (f *FakeKeptn) SetAutomaticResponse(autoResponse bool) {
	f.Keptn.automaticEventResponse = autoResponse
}

func (f *FakeKeptn) SetResourceHandler(handler ResourceHandler) {
	f.TestResourceHandler = handler
	f.Keptn.resourceHandler = handler
}

func (f *FakeKeptn) AddTaskHandler(eventType string, handler TaskHandler, filters ...func(keptnHandle IKeptn, event KeptnEvent) bool) {
	f.AddTaskHandlerWithSubscriptionID(eventType, handler, "", filters...)
}

func (f *FakeKeptn) AddTaskHandlerWithSubscriptionID(eventType string, handler TaskHandler, subscriptionID string, filters ...func(keptnHandle IKeptn, event KeptnEvent) bool) {
	f.Keptn.taskRegistry.Add(eventType, TaskEntry{TaskHandler: handler, EventFilters: filters})
}

func (f *FakeKeptn) fakeSender(ce models.KeptnContextExtendedCE) error {
	f.SentEvents = append(f.SentEvents, ce)
	return nil
}

func NewFakeKeptn(source string) *FakeKeptn {
	resourceHandler := &TestResourceHandler{}
	var fakeKeptn = &FakeKeptn{
		TestResourceHandler: resourceHandler,
		Keptn: &Keptn{
			resourceHandler:        resourceHandler,
			source:                 source,
			taskRegistry:           newTaskMap(),
			syncProcessing:         true,
			automaticEventResponse: true,
			gracefulShutdown:       false,
			logger:                 newDefaultLogger(),
			healthEndpointRunner:   noOpHealthEndpointRunner,
		},
	}
	fakeKeptn.Keptn.eventSender = fakeKeptn.fakeSender
	return fakeKeptn
}

type TestReceiver struct {
	receiverFn interface{}
}

func (t *TestReceiver) StartReceiver(ctx context.Context, fn interface{}) error {
	t.receiverFn = fn
	return nil
}

func (t *TestReceiver) NewEvent(ctx context.Context, e cloudevents.Event) {
	if ctx.Value(gracefulShutdownKey) == nil {
		ctx = context.WithValue(ctx, gracefulShutdownKey, &nopWG{})
	}
	t.receiverFn.(func(context.Context, event.Event))(ctx, e)
}

type TestResourceHandler struct {
	Resource models.Resource
}

func (t TestResourceHandler) GetResource(scope api.ResourceScope, options ...api.URIOption) (*models.Resource, error) {
	return newResourceFromFile(fmt.Sprintf("test/keptn/resources/%s%s%s%s", scope.GetProjectPath(), scope.GetStagePath(), scope.GetServicePath(), scope.GetResourcePath())), nil
}

func newResourceFromFile(filename string) *models.Resource {
	filename = filepath.Clean(filename)
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil
	}

	return &models.Resource{
		Metadata:        nil,
		ResourceContent: string(content),
		ResourceURI:     nil,
	}
}

type StringResourceHandler struct {
	ResourceContent string
}

func (s StringResourceHandler) GetResource(scope api.ResourceScope, options ...api.URIOption) (*models.Resource, error) {
	return &models.Resource{
		Metadata:        &models.Version{Version: "CommitID"},
		ResourceContent: s.ResourceContent,
		ResourceURI:     nil,
	}, nil
}

type FailingResourceHandler struct {
}

func (f FailingResourceHandler) GetResource(scope api.ResourceScope, options ...api.URIOption) (*models.Resource, error) {
	return nil, errors.New("unable to get resource")
}
