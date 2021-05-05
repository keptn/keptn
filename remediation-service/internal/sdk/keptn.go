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

type Context struct {
	FinishedData interface{}
}

func (c *Context) SetFinishedData(data interface{}) {
	c.FinishedData = data
}

//go:generate moq  -pkg fake -out ./fake/task_handler_mock.go . TaskHandler
type TaskHandler interface {
	// Execute is called whenever the actual business-logic of the service shall be executed.
	// Thus, the core logic of the service shall be triggered/implemented in this method.
	//
	// Note, that the contract of the method is to return a valid Context as well as a Error Pointer
	// or nil, if there was no error during execution.
	//
	// During or at the end of execution the implementation is expected to call Context.SetFinishedData(data interface{})
	// to set the data of the .finished event which will eventually be sent out
	//
	Execute(keptnHandle IKeptn, ce interface{}, context Context) (Context, *Error)

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
	// GetTaskRegistry provides access to the internal data structure used for organizing task execuctors
	GetTaskRegistry() *TaskRegistry
}

// Keptn is the default implementation of IKeptn
type Keptn struct {
	EventSender     EventSender
	EventReceiver   EventReceiver
	ResourceHandler ResourceHandler
	Source          string
	TaskRegistry    TaskRegistry
}

// NewKeptn creates a new Keptn
func NewKeptn(ceClient cloudevents.Client, source string, opts ...KeptnOption) *Keptn {

	keptn := &Keptn{
		EventSender:     NewHTTPEventSender(ceClient),
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
	return &k.TaskRegistry
}

func (k *Keptn) gotEvent(event cloudevents.Event) {
	if handler, ok := k.TaskRegistry.Contains(event.Type()); ok {
		data := handler.TaskHandler.GetTriggeredData()
		if err := event.DataAs(&data); err != nil {
			k.send(k.createErrorFinishedEventForTriggeredEvent(event, nil, &Error{Err: err, StatusType: keptnv2.StatusErrored, ResultType: keptnv2.ResultFailed}))
		}
		k.send(k.createStartedEventForTriggeredEvent(event))

		newContext, err := handler.TaskHandler.Execute(k, data, handler.Context)
		if err != nil {
			log.Errorf("error during task execution %v", err.Err)
			k.send(k.createErrorFinishedEventForTriggeredEvent(event, newContext.FinishedData, err))
			return
		}
		if newContext.FinishedData == nil {
			log.Errorf("no finished data set by task executor for event %s. Skipping sending finished event", event.Type())
		} else {
			k.send(k.createFinishedEventForTriggeredEvent(event, newContext.FinishedData))
		}
	}
}

func (k *Keptn) send(event cloudevents.Event) error {
	if err := k.EventSender.SendEvent(event); err != nil {
		log.Println("Error sending .started event")
	}
	return nil
}

func (k *Keptn) createStartedEventForTriggeredEvent(triggeredEvent cloudevents.Event) cloudevents.Event {
	startedEventType := strings.TrimSuffix(triggeredEvent.Type(), ".triggered") + ".started"
	keptnContext, _ := triggeredEvent.Context.GetExtension(KeptnContextCEExtension)
	c := cloudevents.NewEvent()
	c.SetID(uuid.New().String())
	c.SetType(startedEventType)
	c.SetDataContentType(cloudevents.ApplicationJSON)
	c.SetExtension(KeptnContextCEExtension, keptnContext)
	c.SetExtension(TriggeredIDCEExtension, triggeredEvent.ID())
	c.SetSource(k.Source)
	c.SetData(cloudevents.ApplicationJSON, keptnv2.EventData{})
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
	if eventData != nil {
		keptnv2.Decode(eventData, &commonEventData)
		commonEventData.Status = err.StatusType
		commonEventData.Result = err.ResultType
		commonEventData.Message = err.Message
	} else {
		commonEventData.Status = err.StatusType
		commonEventData.Result = err.ResultType
		commonEventData.Message = err.Message
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
	c.SetData(cloudevents.ApplicationJSON, commonEventData)
	return c
}
