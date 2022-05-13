package sdk

import (
	"context"
	"errors"
	"fmt"
	api "github.com/keptn/go-utils/pkg/api/utils"
	"github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/cp-connector/pkg/controlplane"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"log"
	"path/filepath"
	"sync"
	"testing"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/keptn/go-utils/pkg/api/models"
)

type FakeKeptn struct {
	TestResourceHandler    ResourceHandler
	TestEventSource        *TestEventSource
	TestSubscriptionSource *TestSubscriptionSource
	Keptn                  *Keptn
}

func (f *FakeKeptn) StartAsync() error {
	return f.Keptn.Start()
}

func (f *FakeKeptn) Start() error {
	go func() {
		log.Fatal(f.Keptn.Start())
	}()

	select {
	case <-f.TestEventSource.Started:
		return nil
	case <-time.After(5 * time.Second):
		log.Fatal("Timed out waiting for FakeKeptn to be started")
	}
	return nil
}

func (f *FakeKeptn) SendStartedEvent(event KeptnEvent) error {
	return f.Keptn.SendStartedEvent(event)
}

func (f *FakeKeptn) SendFinishedEvent(event KeptnEvent, result interface{}) error {
	return f.Keptn.SendFinishedEvent(event, result)
}

func (f *FakeKeptn) Logger() Logger {
	return f.Keptn.Logger()
}

func (f *FakeKeptn) GetResourceHandler() ResourceHandler {
	if f.TestResourceHandler == nil {
		return &TestResourceHandler{}
	}
	return f.TestResourceHandler
}

func (f *FakeKeptn) NewEvent(event models.KeptnContextExtendedCE) {
	f.TestEventSource.NewEvent(controlplane.EventUpdate{
		KeptnEvent: event,
		MetaData:   controlplane.EventUpdateMetaData{Subject: *event.Type},
	})
}

func (f *FakeKeptn) AssertNumberOfEventSent(t *testing.T, numOfEvents int) {
	require.Eventuallyf(t, func() bool {
		return f.TestEventSource.GetNumberOfSetEvents() == numOfEvents
	}, time.Second, 10*time.Millisecond, "number of events expected: %d got: %d", numOfEvents, f.TestEventSource.GetNumberOfSetEvents())
}

func (f *FakeKeptn) AssertSentEvent(t *testing.T, eventIndex int, assertFn func(ce models.KeptnContextExtendedCE) bool) {
	if eventIndex >= f.TestEventSource.GetNumberOfSetEvents() {
		t.Fatalf("unable to assert sent event with index %d: too less events sent", eventIndex)
	}

	require.Eventually(t, func() bool {
		return assertFn(f.TestEventSource.SentEvents[eventIndex])
	}, time.Second, 10*time.Millisecond)
}

func (f *FakeKeptn) AssertSentEventType(t *testing.T, eventIndex int, eventType string) {
	if eventIndex >= f.TestEventSource.GetNumberOfSetEvents() {
		t.Fatalf("unable to assert sent event with index %d: too less events sent", eventIndex)
	}
	require.Equalf(t, eventType, *f.TestEventSource.SentEvents[eventIndex].Type, "event type expected: %s got %s", eventType, *f.TestEventSource.SentEvents[eventIndex].Type)
}

func (f *FakeKeptn) AssertSentEventStatus(t *testing.T, eventIndex int, status v0_2_0.StatusType) {
	if eventIndex >= f.TestEventSource.GetNumberOfSetEvents() {
		t.Fatalf("unable to assert sent event with index %d: too less events sent", eventIndex)
	}
	eventData := v0_2_0.EventData{}
	v0_2_0.EventDataAs(f.TestEventSource.SentEvents[eventIndex], &eventData)
	require.Equal(t, status, eventData.Status)
}

func (f *FakeKeptn) AssertSentEventResult(t *testing.T, eventIndex int, result v0_2_0.ResultType) {
	if eventIndex >= f.TestEventSource.GetNumberOfSetEvents() {
		t.Fatalf("unable to assert sent event with index %d: too less events sent", eventIndex)
	}
	eventData := v0_2_0.EventData{}
	v0_2_0.EventDataAs(f.TestEventSource.SentEvents[eventIndex], &eventData)
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
	f.TestSubscriptionSource.AddSubscription(models.EventSubscription{ID: subscriptionID, Event: eventType})
	f.Keptn.taskRegistry.Add(eventType, TaskEntry{TaskHandler: handler, EventFilters: filters})
}

func NewFakeKeptn(source string) *FakeKeptn {
	testSubscriptionSource := NewTestSubscriptionSource()
	testEventSource := NewTestEventSource()
	cp := controlplane.New(testSubscriptionSource, testEventSource, nil)
	resourceHandler := &TestResourceHandler{}
	logger := newDefaultLogger()
	var fakeKeptn = &FakeKeptn{
		TestResourceHandler:    resourceHandler,
		TestEventSource:        testEventSource,
		TestSubscriptionSource: testSubscriptionSource,
		Keptn: &Keptn{
			controlPlane:           cp,
			eventSender:            testEventSource.Sender(),
			resourceHandler:        resourceHandler,
			source:                 source,
			taskRegistry:           newTaskMap(),
			syncProcessing:         true,
			automaticEventResponse: true,
			gracefulShutdown:       false,
			logger:                 logger,
			healthEndpointRunner:   noOpHealthEndpointRunner,
		},
	}
	return fakeKeptn
}

// TestSender fakes the sending of CloudEvents
type TestSender struct {
	SentEvents []cloudevents.Event
	mutex      sync.Mutex
	Reactors   map[string]func(event cloudevents.Event) error
}

// SendEvent fakes the sending of CloudEvents
func (s *TestSender) SendEvent(event cloudevents.Event) error {
	if s.Reactors != nil {
		for eventTypeSelector, reactor := range s.Reactors {
			if eventTypeSelector == "*" || eventTypeSelector == event.Type() {
				if err := reactor(event); err != nil {
					return err
				}
			}
		}
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.SentEvents = append(s.SentEvents, event)
	return nil
}

// AssertSentEventTypes checks if the given event types have been passed to the SendEvent function
func (s *TestSender) AssertSentEventTypes(eventTypes []string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	sentTot := len(s.SentEvents)
	typesTot := len(eventTypes)
	if sentTot != typesTot {
		return fmt.Errorf("expected %d event, got %d", sentTot, typesTot)
	}
	for index, sentEvent := range s.SentEvents {
		if sentEvent.Type() != eventTypes[index] {
			return fmt.Errorf("received event type '%s' != %s", sentEvent.Type(), eventTypes[index])
		}
	}
	return nil
}

// AddReactor adds custom logic that should be applied when SendEvent is called for the given event type
func (s *TestSender) AddReactor(eventTypeSelector string, reactor func(event cloudevents.Event) error) {
	if s.Reactors == nil {
		s.Reactors = map[string]func(event cloudevents.Event) error{}
	}
	s.Reactors[eventTypeSelector] = reactor
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
