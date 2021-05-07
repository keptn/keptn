package sdk

import (
	"context"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"
	"github.com/keptn/go-utils/pkg/api/models"
	api "github.com/keptn/go-utils/pkg/api/utils"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	log "github.com/sirupsen/logrus"
	"strings"
)

const KeptnContextCEExtension = "shkeptncontext"
const TriggeredIDCEExtension = "triggeredid"
const ConfigurationServiceURL = "configuration-service:8080"

//go:generate moq  -pkg fake -out ./fake/resource_handler_mock.go . ResourceHandler
type ResourceHandler interface {
	GetServiceResource(project string, stage string, service string, resourceURI string) (*models.Resource, error)
	GetStageResource(project string, stage string, resourceURI string) (*models.Resource, error)
}

type Error struct {
	StatusType keptnv2.StatusType
	ResultType keptnv2.ResultType
	Message    string
	Err        error
}

//go:generate moq  -pkg fake -out ./fake/task_handler_mock.go . TaskHandler
type TaskHandler interface {
	// Execute is called whenever the actual business-logic of the service shall be executed.
	// Thus, the core logic of the service shall be triggered/implemented in this method.
	//
	// Note, that the contract of the method is to return the payload of the .finished event to be sent out as well as a Error Pointer
	// or nil, if there was no error during execution.
	Execute(keptnHandle IKeptn, ce interface{}) (interface{}, *Error)

	// GetTriggeredData is called when a new event was received. It is expected that this method returns a pointer to
	// a struct value of the .triggered event data the service is supposed to process.
	GetTriggeredData() interface{}
}

type KeptnOption func(IKeptn)

// WithHandler registers a handler which is responsible for processing a .triggered event
func WithHandler(handler TaskHandler, eventType string) KeptnOption {
	return func(k IKeptn) {
		k.GetTaskRegistry().Add(eventType, TaskEntry{TaskHandler: handler})
	}
}

type IKeptn interface {
	// Start starts the internal event handling logic and needs to be called by the user
	// after creating value of IKeptn
	Start() error
	// GetResourceHandler returns a handler to fetch data from the configuration service
	GetResourceHandler() ResourceHandler
	// GetTaskRegistry provides access to the internal data structure used for organizing task executors
	GetTaskRegistry() *TaskRegistry
}

// Keptn is the default implementation of IKeptn
type Keptn struct {
	EventSender     EventSender
	EventReceiver   EventReceiver
	ResourceHandler ResourceHandler
	Source          string
	TaskRegistry    *TaskRegistry
}

// NewKeptn creates a new Keptn
func NewKeptn(ceClient cloudevents.Client, source string, opts ...KeptnOption) *Keptn {

	keptn := &Keptn{
		EventSender:     &keptnv2.HTTPEventSender{EventsEndpoint: DefaultHTTPEventEndpoint, Client: ceClient},
		EventReceiver:   ceClient,
		Source:          source,
		TaskRegistry:    NewTasksMap(),
		ResourceHandler: api.NewResourceHandler(ConfigurationServiceURL),
	}
	for _, opt := range opts {
		opt(keptn)
	}
	return keptn
}

func (k *Keptn) Start() error {
	go api.RunHealthEndpoint("10999")
	ctx := context.Background()
	ctx = cloudevents.WithEncodingStructured(ctx)
	return k.EventReceiver.StartReceiver(ctx, k.gotEvent)
}

func (k *Keptn) GetResourceHandler() ResourceHandler {
	return k.ResourceHandler
}

func (k *Keptn) GetTaskRegistry() *TaskRegistry {
	return k.TaskRegistry
}

func (k *Keptn) gotEvent(event cloudevents.Event) {
	if handler, ok := k.TaskRegistry.Contains(event.Type()); ok {
		data := handler.TaskHandler.GetTriggeredData()
		if err := event.DataAs(&data); err != nil {
			log.Errorf("error during decoding of .triggered event: %v", err)
			if err := k.send(k.createErrorFinishedEventForTriggeredEvent(event, nil, &Error{Err: err, StatusType: keptnv2.StatusErrored, ResultType: keptnv2.ResultFailed})); err != nil {
				log.Errorf("unable to send .finished event: %v", err)
				return
			}
		}
		if err := k.send(k.createStartedEventForTriggeredEvent(event)); err != nil {
			log.Errorf("unable to send .started event: %v", err)
			return
		}

		result, err := handler.TaskHandler.Execute(k, data)
		if err != nil {
			log.Errorf("error during task execution %v", err.Err)
			if err := k.send(k.createErrorFinishedEventForTriggeredEvent(event, result, err)); err != nil {
				log.Errorf("unable to send .finished event: %v", err)
				return
			}
			return
		}
		if result == nil {
			log.Errorf("no finished data set by task executor for event %s. Skipping sending finished event", event.Type())
		} else if err := k.send(k.createFinishedEventForTriggeredEvent(event, result)); err != nil {
			log.Errorf("unable to send .finished event: %v", err)
		}
	}
}

func (k *Keptn) send(event cloudevents.Event) error {
	log.Infof("Sending %s event", event.Type())
	if err := k.EventSender.SendEvent(event); err != nil {
		log.Println("Error sending .started event")
	}
	return nil
}

func (k *Keptn) createStartedEventForTriggeredEvent(triggeredEvent cloudevents.Event) cloudevents.Event {

	startedEventType := strings.TrimSuffix(triggeredEvent.Type(), ".triggered") + ".started"
	keptnContext, _ := triggeredEvent.Context.GetExtension(KeptnContextCEExtension)
	eventData := keptnv2.EventData{}
	triggeredEvent.DataAs(&eventData)
	c := cloudevents.NewEvent()
	c.SetID(uuid.New().String())
	c.SetType(startedEventType)
	c.SetDataContentType(cloudevents.ApplicationJSON)
	c.SetExtension(KeptnContextCEExtension, keptnContext)
	c.SetExtension(TriggeredIDCEExtension, triggeredEvent.ID())
	c.SetSource(k.Source)
	c.SetData(cloudevents.ApplicationJSON, eventData)
	return c
}

func (k *Keptn) createFinishedEventForTriggeredEvent(triggeredEvent cloudevents.Event, eventData interface{}) cloudevents.Event {
	var genericEvent map[string]interface{}
	keptnv2.Decode(eventData, &genericEvent)
	if genericEvent["status"] == nil || genericEvent["status"] == "" {
		genericEvent["status"] = "succeeded"
	}

	if genericEvent["result"] == nil || genericEvent["result"] == "" {
		genericEvent["result"] = "pass"
	}

	finishedEventType := strings.TrimSuffix(triggeredEvent.Type(), ".triggered") + ".finished"
	keptnContext, _ := triggeredEvent.Context.GetExtension(KeptnContextCEExtension)
	c := cloudevents.NewEvent()
	c.SetID(uuid.New().String())
	c.SetType(finishedEventType)
	c.SetDataContentType(cloudevents.ApplicationJSON)
	c.SetExtension(KeptnContextCEExtension, keptnContext)
	c.SetExtension(TriggeredIDCEExtension, triggeredEvent.ID())
	c.SetSource(k.Source)
	c.SetData(cloudevents.ApplicationJSON, genericEvent)
	return c

}

func (k *Keptn) createErrorFinishedEventForTriggeredEvent(triggeredEvent cloudevents.Event, eventData interface{}, err *Error) cloudevents.Event {

	commonEventData := keptnv2.EventData{}
	if eventData == nil {
		triggeredEvent.DataAs(&commonEventData)
	}

	commonEventData.Result = err.ResultType
	commonEventData.Status = err.StatusType
	commonEventData.Message = err.Message

	finishedEventType := strings.TrimSuffix(triggeredEvent.Type(), ".triggered") + ".finished"
	keptnContext, _ := triggeredEvent.Context.GetExtension(KeptnContextCEExtension)
	c := cloudevents.NewEvent()
	c.SetID(uuid.New().String())
	c.SetType(finishedEventType)
	c.SetDataContentType(cloudevents.ApplicationJSON)
	c.SetExtension(KeptnContextCEExtension, keptnContext)
	c.SetExtension(TriggeredIDCEExtension, triggeredEvent.ID())
	c.SetSource(k.Source)
	c.SetData(cloudevents.ApplicationJSON, commonEventData)
	return c
}
