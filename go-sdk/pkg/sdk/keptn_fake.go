package sdk

import (
	"context"
	"errors"
	"fmt"
	api "github.com/keptn/go-utils/pkg/api/utils"
	"github.com/keptn/keptn/cp-connector/pkg/controlplane"
	"io/ioutil"
	"path/filepath"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/keptn/go-utils/pkg/api/models"
)

type FakeKeptn struct {
	TestResourceHandler ResourceHandler
	TestEventSource     *TestEventSource
	Keptn               *Keptn
}

func (f *FakeKeptn) Start() error {
	return f.Keptn.Start()
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
		MetaData:   controlplane.EventUpdateMetaData{Subject: "sh.keptn.event.faketask.triggered"},
	})
}

//func (f *FakeKeptn) NewEvent(event cloudevents.Event) {
//	testReceiver := f.Keptn.eventReceiver.(*TestReceiver)
//	testReceiver.NewEvent(context.WithValue(context.Background(), gracefulShutdownKey, &nopWG{}), event)
//}

func (f *FakeKeptn) GetEventSource() *TestEventSource {
	return f.TestEventSource
}

//func (f *FakeKeptn) GetEventSender() *TestSender {
//	return f.Keptn.eventSender.(*TestSender)
//}

func (f *FakeKeptn) SetAutomaticResponse(autoResponse bool) {
	f.Keptn.automaticEventResponse = autoResponse
}

func (f *FakeKeptn) SetResourceHandler(handler ResourceHandler) {
	f.TestResourceHandler = handler
	f.Keptn.resourceHandler = handler
}

func (f *FakeKeptn) AddTaskHandler(eventType string, handler TaskHandler, filters ...func(keptnHandle IKeptn, event KeptnEvent) bool) {
	f.Keptn.taskRegistry.Add(eventType, TaskEntry{TaskHandler: handler, EventFilters: filters})
}

func NewFakeKeptn(source string) *FakeKeptn {
	testSubscriptionSource := controlplane.NewFixedSubscriptionSource(controlplane.WithFixedSubscriptions(models.EventSubscription{Event: "sh.keptn.event.faketask.triggered"}))
	testEventSource := NewTestEventSource()
	cp := controlplane.New(testSubscriptionSource, testEventSource)
	resourceHandler := &TestResourceHandler{}
	logger := NewDefaultLogger()
	var fakeKeptn = &FakeKeptn{
		TestResourceHandler: resourceHandler,
		TestEventSource:     testEventSource,
		Keptn: &Keptn{
			controlPlane:           cp,
			resourceHandler:        resourceHandler,
			source:                 source,
			taskRegistry:           NewTasksMap(),
			syncProcessing:         true,
			automaticEventResponse: true,
			gracefulShutdown:       false,
			logger:                 logger,
		},
	}
	return fakeKeptn
}

// TestSender fakes the sending of CloudEvents
type TestSender struct {
	SentEvents []cloudevents.Event
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
	s.SentEvents = append(s.SentEvents, event)
	return nil
}

// AssertSentEventTypes checks if the given event types have been passed to the SendEvent function
func (s *TestSender) AssertSentEventTypes(eventTypes []string) error {
	if len(s.SentEvents) != len(eventTypes) {
		return fmt.Errorf("expected %d event, got %d", len(s.SentEvents), len(eventTypes))
	}
	for index, event := range s.SentEvents {
		if event.Type() != eventTypes[index] {
			return fmt.Errorf("received event type '%s' != %s", event.Type(), eventTypes[index])
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
