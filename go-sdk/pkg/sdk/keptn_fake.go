package sdk

import (
	"context"
	"errors"
	"fmt"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/keptn/go-utils/pkg/api/models"
	"io/ioutil"
)

type FakeKeptn struct {
	TestResourceHandler ResourceHandler
	Keptn               *Keptn
}

func (f *FakeKeptn) GetResourceHandler() ResourceHandler {
	if f.TestResourceHandler == nil {
		return &TestResourceHandler{}
	}
	return f.TestResourceHandler
}

func (f *FakeKeptn) GetTaskRegistry() *TaskRegistry {
	return f.Keptn.GetTaskRegistry()
}

func (f *FakeKeptn) SetConfigurationServiceURL(configurationServiceURL string) {
	panic("implement me")
}

func (f *FakeKeptn) NewEvent(event cloudevents.Event) {
	testReceiver := f.Keptn.eventReceiver.(*TestReceiver)
	testReceiver.NewEvent(event)
}

func (f *FakeKeptn) GetEventSender() *TestSender {
	return f.Keptn.eventSender.(*TestSender)
}

func (f *FakeKeptn) SendStartedEvent(event KeptnEvent) error {
	return f.Keptn.SendStartedEvent(event)
}

func (f *FakeKeptn) SendFinishedEvent(event KeptnEvent, result interface{}) error {
	return f.Keptn.SendFinishedEvent(event, result)
}

func (f *FakeKeptn) SetAutomaticResponse(autoResponse bool) {
	f.Keptn.automaticEventResponse = autoResponse
}

func (f *FakeKeptn) SetResourceHandler(handler ResourceHandler) {
	f.TestResourceHandler = handler
	f.Keptn.resourceHandler = handler
}

func (f *FakeKeptn) AddHandler(eventType string, handler TaskHandler, filters ...func(keptnHandle IKeptn, event KeptnEvent) bool) {
	f.Keptn.taskRegistry.Add(eventType, TaskEntry{TaskHandler: handler, EventFilters: filters})
}

func NewFakeKeptn(source string) *FakeKeptn {
	eventReceiver := &TestReceiver{}
	eventSender := &TestSender{}
	resourceHandler := &TestResourceHandler{}

	var fakeKeptn = &FakeKeptn{
		TestResourceHandler: resourceHandler,
		Keptn: &Keptn{
			eventSender:            eventSender,
			eventReceiver:          eventReceiver,
			resourceHandler:        resourceHandler,
			source:                 source,
			taskRegistry:           NewTasksMap(),
			syncProcessing:         true,
			automaticEventResponse: true,
		},
	}
	return fakeKeptn
}

func (f *FakeKeptn) Start() error {
	return f.Keptn.Start()
}

//---

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

// ---

type TestReceiver struct {
	receiverFn interface{}
}

func (t *TestReceiver) StartReceiver(ctx context.Context, fn interface{}) error {
	t.receiverFn = fn
	return nil
}

func (t *TestReceiver) NewEvent(e cloudevents.Event) {
	t.receiverFn.(func(event.Event))(e)
}

// ---

type TestResourceHandler struct {
	Resource models.Resource
}

func (t TestResourceHandler) GetServiceResource(project string, stage string, service string, resourceURI string) (*models.Resource, error) {
	return newResourceFromFile(fmt.Sprintf("test/keptn/resources/%s/%s/%s/%s", project, stage, service, resourceURI)), nil
}

func (t TestResourceHandler) GetStageResource(project string, stage string, resourceURI string) (*models.Resource, error) {
	return newResourceFromFile(fmt.Sprintf("test/keptn/resources/%s/%s/%s", project, stage, resourceURI)), nil
}

func (t TestResourceHandler) GetProjectResource(project string, resourceURI string) (*models.Resource, error) {
	return newResourceFromFile(fmt.Sprintf("test/keptn/resources/%s/%s", project, resourceURI)), nil
}

func newResourceFromFile(filename string) *models.Resource {
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

func (s StringResourceHandler) GetServiceResource(project string, stage string, service string, resourceURI string) (*models.Resource, error) {
	return &models.Resource{
		Metadata:        nil,
		ResourceContent: string(s.ResourceContent),
		ResourceURI:     nil,
	}, nil
}

func (s StringResourceHandler) GetStageResource(project string, stage string, resourceURI string) (*models.Resource, error) {
	return &models.Resource{
		Metadata:        nil,
		ResourceContent: string(s.ResourceContent),
		ResourceURI:     nil,
	}, nil
}

func (s StringResourceHandler) GetProjectResource(project string, resourceURI string) (*models.Resource, error) {
	return &models.Resource{
		Metadata:        nil,
		ResourceContent: string(s.ResourceContent),
		ResourceURI:     nil,
	}, nil
}

type FailingResourceHandler struct {
}

func (f FailingResourceHandler) GetServiceResource(project string, stage string, service string, resourceURI string) (*models.Resource, error) {
	return nil, errors.New("unable to get resource")
}

func (f FailingResourceHandler) GetStageResource(project string, stage string, resourceURI string) (*models.Resource, error) {
	return nil, errors.New("unable to get resource")
}

func (f FailingResourceHandler) GetProjectResource(project string, resourceURI string) (*models.Resource, error) {
	return nil, errors.New("unable to get resource")
}
