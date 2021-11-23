package sdk

import (
	"context"
	"fmt"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"
	"github.com/keptn/go-utils/pkg/api/models"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

const DefaultHTTPEventEndpoint = "http://localhost:8081/event"
const KeptnContextCEExtension = "shkeptncontext"
const TriggeredIDCEExtension = "triggeredid"
const ConfigurationServiceURL = "configuration-service:8080"

//go:generate moq  -out ./resourcehhandler_mock.go . ResourceHandler
type ResourceHandler interface {
	GetServiceResource(project string, stage string, service string, resourceURI string) (*models.Resource, error)
	GetStageResource(project string, stage string, resourceURI string) (*models.Resource, error)
	GetProjectResource(project string, resourceURI string) (*models.Resource, error)
}

//go:generate moq  -out ./eventsender_mock.go . EventSender
type EventSender interface {
	SendEvent(event cloudevents.Event) error
}

type EventReceiver interface {
	StartReceiver(ctx context.Context, fn interface{}) error
}

type IKeptn interface {
	// Start starts the internal event handling logic and needs to be called by the user
	// after creating value of IKeptn
	Start() error
	// GetResourceHandler returns a handler to fetch data from the configuration service
	GetResourceHandler() ResourceHandler
	// SendStartedEvent sends a started event for the given input event to the Keptn API
	SendStartedEvent(event KeptnEvent) error
	// SendFinishedEvent sends a finished event for the given input event to the Keptn API
	SendFinishedEvent(event KeptnEvent, result interface{}) error
}

//go:generate moq -out ./taskhandler_mock.go . TaskHandler
type TaskHandler interface {
	// Execute is called whenever the actual business-logic of the service shall be executed.
	// Thus, the core logic of the service shall be triggered/implemented in this method.
	//
	// Note, that the contract of the method is to return the payload of the .finished event to be sent out as well as a Error Pointer
	// or nil, if there was no error during execution.
	Execute(keptnHandle IKeptn, event KeptnEvent) (interface{}, *Error)
}

type KeptnEvent models.KeptnContextExtendedCE

// Opaque key type used for graceful shutdown context value
type gracefulShutdownKeyType struct{}

var gracefulShutdownKey = gracefulShutdownKeyType{}

type Error struct {
	StatusType keptnv2.StatusType
	ResultType keptnv2.ResultType
	Message    string
	Err        error
}

func (e Error) Error() string {
	return e.Message
}

// KeptnOption can be used to configure the keptn sdk
type KeptnOption func(*Keptn)

// WithTaskHandler registers a handler which is responsible for processing a .triggered event
func WithTaskHandler(eventType string, handler TaskHandler, filters ...func(keptnHandle IKeptn, event KeptnEvent) bool) KeptnOption {
	return func(k *Keptn) {
		k.taskRegistry.Add(eventType, TaskEntry{TaskHandler: handler, EventFilters: filters})
	}
}

// WithAutomaticResponse sets the option to instruct the sdk to automatically send a .started and .finished event.
// Per default this behavior is turned on and can be disabled with this function
func WithAutomaticResponse(autoResponse bool) KeptnOption {
	return func(k *Keptn) {
		k.automaticEventResponse = autoResponse
	}
}

// WithGracefulShutdown sets the option to ensure running tasks/handlers will finish in case of interrupt or forced termination
// Per default this behavior is turned on and can be disabled with this function
func WithGracefulShutdown(gracefulShutdown bool) KeptnOption {
	return func(k *Keptn) {
		k.gracefulShutdown = gracefulShutdown
	}
}

// Keptn is the default implementation of IKeptn
type Keptn struct {
	eventSender            EventSender
	eventReceiver          EventReceiver
	resourceHandler        ResourceHandler
	source                 string
	taskRegistry           *TaskRegistry
	syncProcessing         bool
	automaticEventResponse bool
	gracefulShutdown       bool
	recievingEvent         interface{}
}

// NewKeptn creates a new Keptn
func NewKeptn(source string, opts ...KeptnOption) *Keptn {
	client := NewHTTPClientFromEnv()
	resourceHandler := NewResourceHandlerFromEnv()
	taskRegistry := NewTasksMap()
	keptn := &Keptn{
		eventSender:            &keptnv2.HTTPEventSender{EventsEndpoint: DefaultHTTPEventEndpoint, Client: client},
		eventReceiver:          client,
		source:                 source,
		taskRegistry:           taskRegistry,
		resourceHandler:        resourceHandler,
		automaticEventResponse: true,
		gracefulShutdown:       true,
		syncProcessing:         false,
	}
	for _, opt := range opts {
		opt(keptn)
	}
	return keptn
}

func (k *Keptn) Start() error {
	ctx := getGracefulContext()
	err := k.eventReceiver.StartReceiver(ctx, k.gotEvent)
	if k.gracefulShutdown {
		val := ctx.Value(gracefulShutdownKey)
		if val != nil {

			if wg, ok := val.(*sync.WaitGroup); ok {
				wg.Wait()
			}
		}
	}
	return err
}

func (k *Keptn) GetResourceHandler() ResourceHandler {
	return k.resourceHandler
}

func (k *Keptn) SendStartedEvent(event KeptnEvent) error {
	inputCE := cloudevents.Event{}
	err := keptnv2.Decode(event, &inputCE)
	if err != nil {
		return err
	}
	startedEvent, err := k.createStartedEventForTriggeredEvent(inputCE)
	if err != nil {
		return err
	}
	return k.send(*startedEvent)
}

func (k *Keptn) SendFinishedEvent(event KeptnEvent, result interface{}) error {
	inputCE := cloudevents.Event{}
	err := keptnv2.Decode(event, &inputCE)
	if err != nil {
		return err
	}
	finishedEvent, err := k.createFinishedEventForReceivedEvent(inputCE, result)
	if err != nil {
		return err
	}
	return k.send(*finishedEvent)
}

func (k *Keptn) gotEvent(ctx context.Context, event cloudevents.Event) {
	if !keptnv2.IsTaskEventType(event.Type()) {
		log.Errorf("event with event type %s is no valid keptn task event type", event.Type())
		return
	}

	var val interface{} = nil
	if k.gracefulShutdown {
		val = ctx.Value(gracefulShutdownKey)
	}
	if val != nil {
		if wg, ok := val.(*sync.WaitGroup); ok {
			wg.Add(1)
		}
	}

	k.runEventTaskAction(func() {
		{
			defer func() {
				if val == nil {
					return
				}
				if wg, ok := val.(*sync.WaitGroup); ok {
					wg.Done()
				}
			}()

			if handler, ok := k.taskRegistry.Contains(event.Type()); ok {
				keptnEvent := &KeptnEvent{}
				if err := keptnv2.Decode(&event, keptnEvent); err != nil {
					errorLogEvent, err := k.createErrorLogEventForTriggeredEvent(event, nil, &Error{Err: err, StatusType: keptnv2.StatusErrored, ResultType: keptnv2.ResultFailed})
					if err != nil {
						log.Errorf("unable to create '.error.log' event from '.triggered' event: %v", err)
						return
					}
					// no started event sent yet, so it only makes sense to Send an error log event at this point
					if err := k.send(*errorLogEvent); err != nil {
						log.Errorf("unable to send '.finished' event: %v", err)
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
				if keptnv2.IsTaskEventType(event.Type()) && keptnv2.IsTriggeredEventType(event.Type()) && k.automaticEventResponse {
					startedEvent, err := k.createStartedEventForTriggeredEvent(event)
					if err != nil {
						log.Errorf("unable to create '.started' event from '.triggered' event: %v", err)
						return
					}
					if err := k.send(*startedEvent); err != nil {
						log.Errorf("unable to send '.started' event: %v", err)
						return
					}
				}

				result, err := handler.TaskHandler.Execute(k, *keptnEvent)
				if err != nil {
					log.Errorf("error during task execution %v", err.Err)
					if k.automaticEventResponse {
						errorEvent, err := k.createErrorEvent(event, result, err)
						if err != nil {
							log.Errorf("unable to create '.error' event: %v", err)
							return
						}
						if err := k.send(*errorEvent); err != nil {
							log.Errorf("unable to send '.error' event: %v", err)
							return
						}
					}
					return
				}
				if result == nil {
					log.Infof("no finished data set by task executor for event %s. Skipping sending finished event", event.Type())
				} else if keptnv2.IsTaskEventType(event.Type()) && keptnv2.IsTriggeredEventType(event.Type()) && k.automaticEventResponse {
					finishedEvent, err := k.createFinishedEventForReceivedEvent(event, result)
					if err != nil {
						log.Errorf("unable to create '.finished' event: %v", err)
						return
					}
					if err := k.send(*finishedEvent); err != nil {
						log.Errorf("unable to send '.finished' event: %v", err)
						return
					}
				}
			}
		}
	})
}

func (k *Keptn) runEventTaskAction(fn func()) {
	if k.syncProcessing {
		fn()
	} else {
		go fn()
	}
}

func (k *Keptn) send(event cloudevents.Event) error {
	log.Infof("Sending %s event", event.Type())
	if err := k.eventSender.SendEvent(event); err != nil {
		log.Println("Error sending .started event")
	}
	return nil
}

func (k *Keptn) createStartedEventForTriggeredEvent(triggeredEvent cloudevents.Event) (*cloudevents.Event, error) {
	startedEventType, err := keptnv2.ReplaceEventTypeKind(triggeredEvent.Type(), "started")
	if err != nil {
		return nil, fmt.Errorf("unable to create '.finished' event: %v from %s", err, triggeredEvent.Type())
	}
	keptnContext, err := triggeredEvent.Context.GetExtension(KeptnContextCEExtension)
	if err != nil {
		return nil, fmt.Errorf("unable to get keptn context from '.triggered' event: %v", err)
	}
	eventData := keptnv2.EventData{}
	triggeredEvent.DataAs(&eventData)
	c := cloudevents.NewEvent()
	c.SetID(uuid.New().String())
	c.SetType(startedEventType)
	c.SetDataContentType(cloudevents.ApplicationJSON)
	c.SetExtension(KeptnContextCEExtension, keptnContext)
	c.SetExtension(TriggeredIDCEExtension, triggeredEvent.ID())
	c.SetSource(k.source)
	c.SetData(cloudevents.ApplicationJSON, eventData)
	return &c, nil
}

func (k *Keptn) createFinishedEventForReceivedEvent(receivedEvent cloudevents.Event, eventData interface{}) (*cloudevents.Event, error) {
	var genericEvent map[string]interface{}
	keptnv2.Decode(eventData, &genericEvent)
	if genericEvent["status"] == nil || genericEvent["status"] == "" {
		genericEvent["status"] = "succeeded"
	}

	if genericEvent["result"] == nil || genericEvent["result"] == "" {
		genericEvent["result"] = "pass"
	}

	finishedEventType, err := keptnv2.ReplaceEventTypeKind(receivedEvent.Type(), "finished")
	if err != nil {
		return nil, fmt.Errorf("unable to create '.finished' event: %v from %s", err, receivedEvent.Type())
	}
	keptnContext, err := receivedEvent.Context.GetExtension(KeptnContextCEExtension)
	if err != nil {
		return nil, fmt.Errorf("unable to get keptn context from event %s: %v", receivedEvent.Type(), err)
	}
	c := cloudevents.NewEvent()
	c.SetID(uuid.New().String())
	c.SetType(finishedEventType)
	c.SetDataContentType(cloudevents.ApplicationJSON)
	c.SetExtension(KeptnContextCEExtension, keptnContext)
	c.SetExtension(TriggeredIDCEExtension, receivedEvent.ID())
	c.SetSource(k.source)
	c.SetData(cloudevents.ApplicationJSON, genericEvent)
	return &c, nil
}

func (k *Keptn) createErrorEvent(event cloudevents.Event, eventData interface{}, err *Error) (*cloudevents.Event, error) {
	if keptnv2.IsTaskEventType(event.Type()) && keptnv2.IsTriggeredEventType(event.Type()) {
		errorFinishedEvent, err2 := k.createErrorFinishedEventForTriggeredEvent(event, eventData, err)
		if err2 != nil {
			return nil, err2
		}
		return errorFinishedEvent, nil
	}
	errorLogEvent, err2 := k.createErrorLogEventForTriggeredEvent(event, eventData, err)
	if err2 != nil {
		return nil, err2
	}
	return errorLogEvent, nil
}

func (k *Keptn) createErrorLogEventForTriggeredEvent(triggeredEvent cloudevents.Event, eventData interface{}, err *Error) (*cloudevents.Event, error) {
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

	keptnContext, err2 := triggeredEvent.Context.GetExtension(KeptnContextCEExtension)
	if err2 != nil {
		return nil, fmt.Errorf("unable to get keptn context from '.triggered' event: %v", err2)
	}
	c := cloudevents.NewEvent()
	c.SetID(uuid.New().String())
	c.SetType(keptnv2.ErrorLogEventName)
	c.SetDataContentType(cloudevents.ApplicationJSON)
	c.SetExtension(KeptnContextCEExtension, keptnContext)
	c.SetExtension(TriggeredIDCEExtension, triggeredEvent.ID())
	c.SetSource(k.source)
	c.SetData(cloudevents.ApplicationJSON, errorEventData)
	return &c, nil
}

func (k *Keptn) createErrorFinishedEventForTriggeredEvent(event cloudevents.Event, eventData interface{}, err *Error) (*cloudevents.Event, error) {
	commonEventData := keptnv2.EventData{}
	if eventData == nil {
		event.DataAs(&commonEventData)
	}

	commonEventData.Result = err.ResultType
	commonEventData.Status = err.StatusType
	commonEventData.Message = err.Message

	finishedEventType, err2 := keptnv2.ReplaceEventTypeKind(event.Type(), "finished")
	if err2 != nil {
		return nil, fmt.Errorf("unable to create '.finished' event: %v from %s", err2, event.Type())
	}
	keptnContext, err2 := event.Context.GetExtension(KeptnContextCEExtension)
	if err2 != nil {
		return nil, fmt.Errorf("unable to get keptn context from '.triggered' event: %v", err2)
	}
	c := cloudevents.NewEvent()
	c.SetID(uuid.New().String())
	c.SetType(finishedEventType)
	c.SetDataContentType(cloudevents.ApplicationJSON)
	c.SetExtension(KeptnContextCEExtension, keptnContext)
	c.SetExtension(TriggeredIDCEExtension, event.ID())
	c.SetSource(k.source)
	c.SetData(cloudevents.ApplicationJSON, commonEventData)
	return &c, nil
}

// getGracefulContext returns a context with cancel and a wait group to sync before shutdown
func getGracefulContext() context.Context {

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	wg := &sync.WaitGroup{}
	ctx, cancel := context.WithCancel(cloudevents.WithEncodingStructured(context.WithValue(context.Background(), gracefulShutdownKey, wg)))

	go func() {
		<-ch
		cancel()
	}()

	return ctx
}
