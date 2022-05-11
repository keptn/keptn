package sdk

import (
	"context"
	"errors"
	"fmt"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/keptn/go-utils/pkg/api/models"
	api "github.com/keptn/go-utils/pkg/api/utils"
	"github.com/keptn/go-utils/pkg/common/strutils"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/cp-connector/pkg/controlplane"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

//go:generate moq  -out ./resourcehandler_mock.go . ResourceHandler
type ResourceHandler interface {
	GetResource(scope api.ResourceScope, options ...api.URIOption) (*models.Resource, error)
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
	// Logger returns the logger used by the sdk
	// Per default DefaultLogger is used which internally just uses the go logging package
	// Another logger can be configured using the sdk.WithLogger function
	Logger() Logger
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

type wgInterface interface {
	Add(delta int)
	Done()
	Wait()
}

type nopWG struct {
	// --
}

func (w *nopWG) Add(delta int) {
	// --
}
func (w *nopWG) Done() {
	// --
}
func (w *nopWG) Wait() {
	// --
}

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

// WithLogger configures keptn to use another logger
func WithLogger(logger Logger) KeptnOption {
	return func(k *Keptn) {
		k.logger = logger
	}
}

func WithControlPlaneFromEnv(env envConfig) KeptnOption {
	return func(k *Keptn) {
		if env.Location == "remote-execution-plane" {
			panic("go-sdk remote-execution-plane mode not yet supported...")
		}
		k.controlPlane, k.eventSender = newControlPlaneFromEnv(env)
	}
}

// Keptn is the default implementation of IKeptn
type Keptn struct {
	controlPlane           *controlplane.ControlPlane
	eventSender            controlplane.EventSender
	resourceHandler        ResourceHandler
	source                 string
	taskRegistry           *TaskRegistry
	syncProcessing         bool
	automaticEventResponse bool
	gracefulShutdown       bool
	logger                 Logger
	env                    envConfig
}

// NewKeptn creates a new Keptn
func NewKeptn(source string, opts ...KeptnOption) *Keptn {
	env := newEnvConfig()
	controlPlane, eventSender := newControlPlane()
	resourceHandler := newResourceHandlerFromEnv()
	taskRegistry := newTaskMap()
	logger := newDefaultLogger()
	keptn := &Keptn{
		controlPlane:           controlPlane,
		eventSender:            eventSender,
		source:                 source,
		taskRegistry:           taskRegistry,
		resourceHandler:        resourceHandler,
		automaticEventResponse: true,
		gracefulShutdown:       true,
		syncProcessing:         false,
		logger:                 logger,
		env:                    env,
	}
	for _, opt := range opts {
		opt(keptn)
	}
	return keptn
}

func (k *Keptn) OnEvent(ctx context.Context, event models.KeptnContextExtendedCE) error {
	eventSender := ctx.Value(controlplane.EventSenderKey).(controlplane.EventSender)
	if !keptnv2.IsTaskEventType(*event.Type) {
		k.logger.Errorf("event with event type %s is no valid keptn task event type", event.Type)
		return nil
	}
	ctx.Value(gracefulShutdownKey).(wgInterface).Add(1)
	k.runEventTaskAction(func() {
		{
			defer ctx.Value(gracefulShutdownKey).(wgInterface).Done()
			if handler, ok := k.taskRegistry.Contains(*event.Type); ok {
				keptnEvent := &KeptnEvent{}
				if err := keptnv2.Decode(&event, keptnEvent); err != nil {
					errorLogEvent, err := k.createErrorLogEventForTriggeredEvent(event, nil, &Error{Err: err, StatusType: keptnv2.StatusErrored, ResultType: keptnv2.ResultFailed})
					if err != nil {
						k.logger.Errorf("unable to create '.error.log' event from '.triggered' event: %v", err)
						return
					}
					// no started event sent yet, so it only makes sense to Send an error log event at this point
					if err := eventSender(*errorLogEvent); err != nil {
						k.logger.Errorf("unable to send '.finished' event: %v", err)
						return
					}
				}

				// execute the filtering functions of the task handler to determine whether the incoming event should be handled
				// only if all functions return true, the event will be handled
				for _, filterFn := range handler.EventFilters {
					if !filterFn(k, *keptnEvent) {
						k.logger.Infof("Will not handle incoming %s event", *event.Type)
						return
					}
				}

				// only respond with .started event if the incoming event is a task.triggered event
				if keptnv2.IsTaskEventType(*event.Type) && keptnv2.IsTriggeredEventType(*event.Type) && k.automaticEventResponse {
					startedEvent, err := k.createStartedEventForTriggeredEvent(event)
					if err != nil {
						k.logger.Errorf("unable to create '.started' event from '.triggered' event: %v", err)
						return
					}
					if err := eventSender(*startedEvent); err != nil {
						k.logger.Errorf("unable to send '.started' event: %v", err)
						return
					}
				}

				result, err := handler.TaskHandler.Execute(k, *keptnEvent)
				if err != nil {
					k.logger.Errorf("error during task execution %v", err.Err)
					if k.automaticEventResponse {
						errorEvent, err := k.createErrorEvent(event, result, err)
						if err != nil {
							k.logger.Errorf("unable to create '.error' event: %v", err)
							return
						}
						if err := eventSender(*errorEvent); err != nil {
							k.logger.Errorf("unable to send '.error' event: %v", err)
							return
						}
					}
					return
				}
				if result == nil {
					k.logger.Infof("no finished data set by task executor for event %s. Skipping sending finished event", *event.Type)
				} else if keptnv2.IsTaskEventType(*event.Type) && keptnv2.IsTriggeredEventType(*event.Type) && k.automaticEventResponse {
					finishedEvent, err := k.createFinishedEventForReceivedEvent(event, result)
					if err != nil {
						k.logger.Errorf("unable to create '.finished' event: %v", err)
						return
					}
					if err := eventSender(*finishedEvent); err != nil {
						k.logger.Errorf("unable to send '.finished' event: %v", err)
						return
					}
				}
			}
		}
	})
	return nil
}

func (k *Keptn) RegistrationData() controlplane.RegistrationData {
	return controlplane.RegistrationData{
		Name: "local-service",
		MetaData: models.MetaData{
			Hostname:           "localhost",
			IntegrationVersion: k.env.Version,
			Location:           k.env.Location,
			KubernetesMetaData: models.KubernetesMetaData{
				Namespace:      k.env.K8sNamespace,
				PodName:        k.env.K8sPodName,
				DeploymentName: k.env.K8sDeploymentName,
			},
		},
		Subscriptions: []models.EventSubscription{},
	}
}

func (k *Keptn) Start() error {
	ctx := getContext(k.gracefulShutdown)
	err := k.controlPlane.Register(ctx, k)
	ctx.Value(gracefulShutdownKey).(wgInterface).Wait()
	return err
}

func (k *Keptn) GetResourceHandler() ResourceHandler {
	return k.resourceHandler
}

func (k *Keptn) SendStartedEvent(event KeptnEvent) error {
	finishedEvent, err := k.createStartedEventForTriggeredEvent(models.KeptnContextExtendedCE(event))
	if err != nil {
		return err
	}
	return k.eventSender(*finishedEvent)
}

func (k *Keptn) SendFinishedEvent(event KeptnEvent, result interface{}) error {
	finishedEvent, err := k.createFinishedEventForReceivedEvent(models.KeptnContextExtendedCE(event), result)
	if err != nil {
		return err
	}
	return k.eventSender(*finishedEvent)
}

func (k *Keptn) Logger() Logger {
	return k.logger
}

func (k *Keptn) runEventTaskAction(fn func()) {
	if k.syncProcessing {
		fn()
	} else {
		go fn()
	}
}

func (k *Keptn) createStartedEventForTriggeredEvent(triggeredEvent models.KeptnContextExtendedCE) (*models.KeptnContextExtendedCE, error) {
	startedEventType, err := keptnv2.ReplaceEventTypeKind(*triggeredEvent.Type, "started")
	if err != nil {
		return nil, fmt.Errorf("unable to create '.finished' event: %v from %s", err, *triggeredEvent.Type)
	}
	if triggeredEvent.Shkeptncontext == "" {
		return nil, errors.New("unable to get keptn context from '.triggered' event")
	}
	eventData := keptnv2.EventData{}
	triggeredEvent.DataAs(&eventData)
	event, err := keptnv2.KeptnEvent(startedEventType, k.source, eventData).WithKeptnContext(triggeredEvent.Shkeptncontext).WithTriggeredID(triggeredEvent.ID).Build()
	if err != nil {
		return nil, fmt.Errorf("could not create error event: %v", err)
	}
	return &event, nil
}

func (k *Keptn) createFinishedEventForReceivedEvent(receivedEvent models.KeptnContextExtendedCE, eventData interface{}) (*models.KeptnContextExtendedCE, error) {
	var genericEvent map[string]interface{}
	keptnv2.Decode(eventData, &genericEvent)
	if genericEvent["status"] == nil || genericEvent["status"] == "" {
		genericEvent["status"] = "succeeded"
	}

	if genericEvent["result"] == nil || genericEvent["result"] == "" {
		genericEvent["result"] = "pass"
	}

	finishedEventType, err := keptnv2.ReplaceEventTypeKind(*receivedEvent.Type, "finished")
	if err != nil {
		return nil, fmt.Errorf("unable to create '.finished' event: %v from %s", err, *receivedEvent.Type)
	}
	if receivedEvent.Shkeptncontext == "" {
		return nil, fmt.Errorf("unable to get keptn context from event %s", *receivedEvent.Type)
	}

	event := models.KeptnContextExtendedCE{
		Contenttype:        cloudevents.ApplicationJSON,
		Data:               genericEvent,
		Source:             strutils.Stringp(k.source),
		Shkeptnspecversion: "1.0",
		Specversion:        "0.2.4",
		Type:               strutils.Stringp(finishedEventType),
		Triggeredid:        receivedEvent.ID,
		Shkeptncontext:     receivedEvent.Shkeptncontext,
	}

	return &event, nil
}

func (k *Keptn) createErrorEvent(event models.KeptnContextExtendedCE, eventData interface{}, err *Error) (*models.KeptnContextExtendedCE, error) {
	if keptnv2.IsTaskEventType(*event.Type) && keptnv2.IsTriggeredEventType(*event.Type) {
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

func (k *Keptn) createErrorLogEventForTriggeredEvent(triggeredEvent models.KeptnContextExtendedCE, eventData interface{}, errVal *Error) (*models.KeptnContextExtendedCE, error) {
	errorEventData := keptnv2.ErrorLogEvent{}
	if eventData == nil {
		triggeredEvent.DataAs(&errorEventData)
	}

	if keptnv2.IsTaskEventType(*triggeredEvent.Type) {
		taskName, _, err := keptnv2.ParseTaskEventType(*triggeredEvent.Type)
		if err == nil && taskName != "" {
			errorEventData.Task = taskName
		}
	}

	errorEventData.Message = errVal.Message

	if triggeredEvent.Shkeptncontext == "" {
		return nil, errors.New("unable to get keptn context from '.triggered' event")
	}

	errEvent, err := keptnv2.KeptnEvent(keptnv2.ErrorLogEventName, k.source, errorEventData).WithKeptnContext(triggeredEvent.Shkeptncontext).WithTriggeredID(triggeredEvent.ID).Build()
	if err != nil {
		return nil, fmt.Errorf("could not create error event: %v", err)
	}

	return &errEvent, nil
}

func (k *Keptn) createErrorFinishedEventForTriggeredEvent(triggeredEvent models.KeptnContextExtendedCE, eventData interface{}, errVal *Error) (*models.KeptnContextExtendedCE, error) {
	commonEventData := keptnv2.EventData{}
	if eventData == nil {
		triggeredEvent.DataAs(&commonEventData)
	}

	commonEventData.Result = errVal.ResultType
	commonEventData.Status = errVal.StatusType
	commonEventData.Message = errVal.Message

	finishedEventType, err := keptnv2.ReplaceEventTypeKind(*triggeredEvent.Type, "finished")
	if err != nil {
		return nil, fmt.Errorf("unable to create '.finished' event: %v from %s", err, *triggeredEvent.Type)
	}
	if triggeredEvent.Shkeptncontext == "" {
		return nil, fmt.Errorf("unable to get keptn context from '.triggered' event")
	}

	finishedEvent, errCreateEvent := keptnv2.KeptnEvent(finishedEventType, k.source, commonEventData).WithTriggeredID(triggeredEvent.ID).WithKeptnContext(triggeredEvent.Shkeptncontext).Build()
	if errCreateEvent != nil {
		return nil, errCreateEvent
	}
	return &finishedEvent, nil
}

func getContext(graceful bool) context.Context {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	var wg wgInterface
	if graceful {
		wg = &sync.WaitGroup{}
	} else {
		wg = &nopWG{}
	}
	ctx, cancel := context.WithCancel(cloudevents.WithEncodingStructured(context.WithValue(context.Background(), gracefulShutdownKey, wg)))

	go func() {
		<-ch
		cancel()
	}()

	return ctx
}
