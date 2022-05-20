package sdk

import (
	"context"
	"github.com/kelseyhightower/envconfig"
	"github.com/keptn/go-utils/pkg/api/models"
	api "github.com/keptn/go-utils/pkg/api/utils"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	api2 "github.com/keptn/keptn/cp-connector/pkg/api"
	"github.com/keptn/keptn/cp-connector/pkg/controlplane"
	"github.com/keptn/keptn/cp-connector/pkg/nats"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
)

const (
	shkeptnspecversion = "0.2.4"
	cloudeventsversion = "1.0"
)

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

type TaskHandler interface {
	// Execute is called whenever the actual business-logic of the service shall be executed.
	// Thus, the core logic of the service shall be triggered/implemented in this method.
	//
	// Note, that the contract of the method is to return the payload of the .finished event to be sent out as well as a Error Pointer
	// or nil, if there was no error during execution.
	Execute(keptnHandle IKeptn, event KeptnEvent) (interface{}, *Error)
}

type KeptnEvent models.KeptnContextExtendedCE

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

type ResourceHandler interface {
	GetResource(scope api.ResourceScope, options ...api.URIOption) (*models.Resource, error)
}

type healthEndpointRunner func(port string, cp *controlplane.ControlPlane)

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

// WithTaskHandler registers a handler which is responsible for processing a .triggered event
func WithTaskHandler(eventType string, handler TaskHandler, filters ...func(keptnHandle IKeptn, event KeptnEvent) bool) KeptnOption {
	return func(k *Keptn) {
		k.taskRegistry.Add(eventType, taskEntry{taskHandler: handler, eventFilters: filters})
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

// Keptn is the default implementation of IKeptn
type Keptn struct {
	controlPlane           *controlplane.ControlPlane
	eventSender            controlplane.EventSender
	resourceHandler        ResourceHandler
	source                 string
	taskRegistry           *taskRegistry
	syncProcessing         bool
	automaticEventResponse bool
	gracefulShutdown       bool
	logger                 Logger
	env                    envConfig
	healthEndpointRunner   healthEndpointRunner
}

// NewKeptn creates a new Keptn
func NewKeptn(source string, opts ...KeptnOption) *Keptn {
	env := newEnvConfig()
	controlPlane, eventSender := newControlPlaneFromEnv()
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
		healthEndpointRunner:   newHealthEndpointRunner,
	}
	for _, opt := range opts {
		opt(keptn)
	}
	return keptn
}

func (k *Keptn) OnEvent(ctx context.Context, event models.KeptnContextExtendedCE) error {
	eventSender, ok := ctx.Value(controlplane.EventSenderKey).(controlplane.EventSender)
	if !ok {
		k.logger.Errorf("Unable to get event sender. Skip processing of event %s", event.ID)
		return nil
	}

	if event.Type == nil {
		k.logger.Errorf("Unable to get event type. Skip processing of event %s", event.ID)
		return nil
	}

	if !keptnv2.IsTaskEventType(*event.Type) {
		k.logger.Errorf("Event type %s does not match format for task events. Skip Processing of event %s", *event.Type, event.ID)
		return nil
	}
	wg, ok := ctx.Value(gracefulShutdownKey).(wgInterface)
	if !ok {
		k.logger.Errorf("Unable to get graceful shutdown wait group. Skip processing of event %s", event.ID)
		return nil
	}
	wg.Add(1)
	k.runEventTaskAction(func() {
		{
			defer wg.Done()
			if handler, ok := k.taskRegistry.Contains(*event.Type); ok {
				keptnEvent := &KeptnEvent{}
				if err := keptnv2.Decode(&event, keptnEvent); err != nil {
					errorLogEvent, err := createErrorLogEvent(k.source, event, nil, &Error{Err: err, StatusType: keptnv2.StatusErrored, ResultType: keptnv2.ResultFailed})
					if err != nil {
						k.logger.Errorf("Unable to create '.error.log' event from '.triggered' event: %v", err)
						return
					}
					// no started event sent yet, so it only makes sense to Send an error log event at this point
					if err := eventSender(*errorLogEvent); err != nil {
						k.logger.Errorf("Unable to send '.finished' event: %v", err)
						return
					}
				}

				// execute the filtering functions of the task handler to determine whether the incoming event should be handled
				// only if all functions return true, the event will be handled
				for _, filterFn := range handler.eventFilters {
					if !filterFn(k, *keptnEvent) {
						k.logger.Infof("Will not handle incoming %s event", *event.Type)
						return
					}
				}

				// only respond with .started event if the incoming event is a task.triggered event
				if keptnv2.IsTaskEventType(*event.Type) && keptnv2.IsTriggeredEventType(*event.Type) && k.automaticEventResponse {
					startedEvent, err := createStartedEvent(k.source, event)
					if err != nil {
						k.logger.Errorf("Unable to create '.started' event from '.triggered' event: %v", err)
						return
					}
					if err := eventSender(*startedEvent); err != nil {
						k.logger.Errorf("Unable to send '.started' event: %v", err)
						return
					}
				}

				result, err := handler.taskHandler.Execute(k, *keptnEvent)
				if err != nil {
					k.logger.Errorf("Error during task execution %v", err.Err)
					if k.automaticEventResponse {
						errorEvent, err := createErrorEvent(k.source, event, result, err)
						if err != nil {
							k.logger.Errorf("Unable to create '.error' event: %v", err)
							return
						}
						if err := eventSender(*errorEvent); err != nil {
							k.logger.Errorf("Unable to send '.error' event: %v", err)
							return
						}
					}
					return
				}
				if result == nil {
					k.logger.Infof("no finished data set by task executor for event %s. Skipping sending finished event", *event.Type)
				} else if keptnv2.IsTaskEventType(*event.Type) && keptnv2.IsTriggeredEventType(*event.Type) && k.automaticEventResponse {
					finishedEvent, err := createFinishedEvent(k.source, event, result)
					if err != nil {
						k.logger.Errorf("Unable to create '.finished' event: %v", err)
						return
					}
					if err := eventSender(*finishedEvent); err != nil {
						k.logger.Errorf("Unable to send '.finished' event: %v", err)
						return
					}
				}
			}
		}
	})
	return nil
}

func (k *Keptn) RegistrationData() controlplane.RegistrationData {
	subscriptions := []models.EventSubscription{}
	subjects := []string{}
	if k.env.PubSubTopic != "" {
		subjects = strings.Split(k.env.PubSubTopic, ",")
	}

	for _, s := range subjects {
		subscriptions = append(subscriptions, models.EventSubscription{Event: s})
	}
	return controlplane.RegistrationData{
		Name: k.source,
		MetaData: models.MetaData{
			Hostname:           k.env.K8sNodeName,
			IntegrationVersion: k.env.Version,
			Location:           k.env.Location,
			DistributorVersion: "0.15.0", // note: to be deleted when bridge stops requiring this info
			KubernetesMetaData: models.KubernetesMetaData{
				Namespace:      k.env.K8sNamespace,
				PodName:        k.env.K8sPodName,
				DeploymentName: k.env.K8sDeploymentName,
			},
		},
		Subscriptions: subscriptions,
	}
}

func (k *Keptn) Start() error {
	if k.env.HealthEndpointEnabled {
		k.healthEndpointRunner(k.env.HealthEndpointPort, k.controlPlane)
	}
	ctx, wg := k.getContext(k.gracefulShutdown)
	err := k.controlPlane.Register(ctx, k)
	wg.Wait()
	return err
}

func (k *Keptn) GetResourceHandler() ResourceHandler {
	return k.resourceHandler
}

func (k *Keptn) SendStartedEvent(event KeptnEvent) error {
	finishedEvent, err := createStartedEvent(k.source, models.KeptnContextExtendedCE(event))
	if err != nil {
		return err
	}
	return k.eventSender(*finishedEvent)
}

func (k *Keptn) SendFinishedEvent(event KeptnEvent, result interface{}) error {
	finishedEvent, err := createFinishedEvent(k.source, models.KeptnContextExtendedCE(event), result)
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

func (k *Keptn) getContext(graceful bool) (context.Context, wgInterface) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	var wg wgInterface
	if graceful {
		wg = &sync.WaitGroup{}
	} else {
		wg = &nopWG{}
	}
	ctx, cancel := context.WithCancel(context.WithValue(context.Background(), gracefulShutdownKey, wg))
	go func() {
		<-ch
		cancel()
	}()
	return ctx, wg
}

func noOpHealthEndpointRunner(port string, cp *controlplane.ControlPlane) {}

func newHealthEndpointRunner(port string, cp *controlplane.ControlPlane) {
	go func() {
		api.RunHealthEndpoint(port, api.WithReadinessConditionFunc(func() bool {
			return cp.IsRegistered()
		}))
	}()
}

func newResourceHandlerFromEnv() *api.ResourceHandler {
	var env envConfig
	if err := envconfig.Process("", &env); err != nil {
		log.Fatalf("failed to process env var: %s", err)
	}
	return api.NewResourceHandler(env.ConfigurationServiceURL)
}

func newControlPlaneFromEnv() (*controlplane.ControlPlane, controlplane.EventSender) {
	var env envConfig
	if err := envconfig.Process("", &env); err != nil {
		log.Fatalf("failed to process env var: %s", err)
	}
	apiSet, err := api2.NewInternal(nil)
	if err != nil {
		log.Fatal(err)
	}
	natsConnector, err := nats.Connect(env.EventBrokerURL)
	if err != nil {
		log.Fatal(err)
	}
	eventSource := controlplane.NewNATSEventSource(natsConnector)
	eventSender := eventSource.Sender()
	subscriptionSource := controlplane.NewUniformSubscriptionSource(apiSet.UniformV1())
	logForwarder := controlplane.NewLogForwarder(apiSet.LogsV1())
	controlPlane := controlplane.New(subscriptionSource, eventSource, logForwarder)
	return controlPlane, eventSender
}
