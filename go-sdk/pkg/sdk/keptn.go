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
	GetProjectResource(project string, resourceURI string) (*models.Resource, error)
}

type KeptnEvent models.KeptnContextExtendedCE

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
	Execute(keptnHandle IKeptn, event KeptnEvent) (interface{}, *Error)
}

type KeptnOption func(IKeptn)

// WithHandler registers a handler which is responsible for processing a .triggered event
func WithHandler(eventType string, handler TaskHandler, filters ...func(keptnHandle IKeptn, event KeptnEvent) bool) KeptnOption {
	return func(k IKeptn) {
		k.GetTaskRegistry().Add(eventType, TaskEntry{TaskHandler: handler, EventFilters: filters})
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
	SyncProcessing  bool
	ReceivingEvent  interface{}
}

// NewKeptn creates a new Keptn
func NewKeptn(source string, opts ...KeptnOption) *Keptn {
	client := NewHTTPClientFromEnv()
	resourceHandler := NewResourceHandlerFromEnv()
	taskRegistry := NewTasksMap()
	keptn := &Keptn{
		EventSender:     &keptnv2.HTTPEventSender{EventsEndpoint: DefaultHTTPEventEndpoint, Client: client},
		EventReceiver:   client,
		Source:          source,
		TaskRegistry:    taskRegistry,
		ResourceHandler: resourceHandler,
		SyncProcessing:  false,
	}
	for _, opt := range opts {
		opt(keptn)
	}
	return keptn
}

func (k *Keptn) Start() error {
	go api.RunHealthEndpoint("10998")
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
	if !keptnv2.IsTaskEventType(event.Type()) {
		log.Errorf("event with event type %s is no valid keptn task event type", event.Type())
		return
	}

	k.runEventTaskAction(func() {
		{
			if handler, ok := k.TaskRegistry.Contains(event.Type()); ok {
				keptnEvent := &KeptnEvent{}
				if err := keptnv2.Decode(&event, keptnEvent); err != nil {
					// no started event sent yet, so it only makes sense to send an error log event at this point
					if err := k.send(k.createErrorLogEventForTriggeredEvent(event, nil, &Error{Err: err, StatusType: keptnv2.StatusErrored, ResultType: keptnv2.ResultFailed})); err != nil {
						log.Errorf("unable to send .finished event: %v", err)
						return
					}
				}

				// execute the filtering functions of the task handler to determine whether the incoming event should be handled
				// only if all functions return true, the event will be handled
				for _, filterFn := range handler.EventFilters {
					if !filterFn(k, *keptnEvent) {
						log.Infof("Will not handle incoming %s event", event.Type())
						return
					}
				}

				// only respond with .started event if the incoming event is a task.triggered event
				if keptnv2.IsTaskEventType(event.Type()) && keptnv2.IsTriggeredEventType(event.Type()) {
					if err := k.send(k.createStartedEventForTriggeredEvent(event)); err != nil {
						log.Errorf("unable to send .started event: %v", err)
						return
					}
				}

				result, err := handler.TaskHandler.Execute(k, *keptnEvent)
				if err != nil {
					log.Errorf("error during task execution %v", err.Err)
					if err := k.send(k.createErrorEvent(event, result, err)); err != nil {
						log.Errorf("unable to send .finished event: %v", err)
						return
					}
					return
				}
				if result == nil {
					log.Errorf("no finished data set by task executor for event %s. Skipping sending finished event", event.Type())
				} else if keptnv2.IsTaskEventType(event.Type()) && keptnv2.IsTriggeredEventType(event.Type()) {
					if err := k.send(k.createFinishedEventForTriggeredEvent(event, result)); err != nil {
						log.Errorf("unable to send .finished event: %v", err)
					}
				}
			}
		}
	})
}

func (k *Keptn) runEventTaskAction(fn func()) {
	if k.SyncProcessing {
		fn()
	} else {
		go fn()
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
	startedEventType, _ := keptnv2.ReplaceEventTypeKind(triggeredEvent.Type(), "started")
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

func (k *Keptn) createErrorEvent(event cloudevents.Event, eventData interface{}, err *Error) cloudevents.Event {
	if keptnv2.IsTaskEventType(event.Type()) && keptnv2.IsTriggeredEventType(event.Type()) {
		return k.createErrorFinishedEventForTriggeredEvent(event, eventData, err)
	}
	return k.createErrorLogEventForTriggeredEvent(event, eventData, err)
}

func (k *Keptn) createErrorLogEventForTriggeredEvent(triggeredEvent cloudevents.Event, eventData interface{}, err *Error) cloudevents.Event {
	errorEventData := keptnv2.ErrorLogEvent{}
	if eventData == nil {
		triggeredEvent.DataAs(&errorEventData)
	}

	if keptnv2.IsTaskEventType(triggeredEvent.Type()) {
		taskName, _, err2 := keptnv2.ParseTaskEventType(triggeredEvent.Type())
		if err2 == nil && taskName != "" {
			errorEventData.Task = taskName
		}
	}

	errorEventData.Message = err.Message

	keptnContext, _ := triggeredEvent.Context.GetExtension(KeptnContextCEExtension)
	c := cloudevents.NewEvent()
	c.SetID(uuid.New().String())
	c.SetType(keptnv2.ErrorLogEventName)
	c.SetDataContentType(cloudevents.ApplicationJSON)
	c.SetExtension(KeptnContextCEExtension, keptnContext)
	c.SetExtension(TriggeredIDCEExtension, triggeredEvent.ID())
	c.SetSource(k.Source)
	c.SetData(cloudevents.ApplicationJSON, errorEventData)
	return c
}

func (k *Keptn) createErrorFinishedEventForTriggeredEvent(event cloudevents.Event, eventData interface{}, err *Error) cloudevents.Event {
	commonEventData := keptnv2.EventData{}
	if eventData == nil {
		event.DataAs(&commonEventData)
	}

	commonEventData.Result = err.ResultType
	commonEventData.Status = err.StatusType
	commonEventData.Message = err.Message

	finishedEventType, _ := keptnv2.ReplaceEventTypeKind(event.Type(), "finished")
	keptnContext, _ := event.Context.GetExtension(KeptnContextCEExtension)
	c := cloudevents.NewEvent()
	c.SetID(uuid.New().String())
	c.SetType(finishedEventType)
	c.SetDataContentType(cloudevents.ApplicationJSON)
	c.SetExtension(KeptnContextCEExtension, keptnContext)
	c.SetExtension(TriggeredIDCEExtension, event.ID())
	c.SetSource(k.Source)
	c.SetData(cloudevents.ApplicationJSON, commonEventData)
	return c
}
