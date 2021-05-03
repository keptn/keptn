package sdk

import (
	"context"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"
	"github.com/keptn/go-utils/pkg/api/models"
	api "github.com/keptn/go-utils/pkg/api/utils"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"log"
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

type Context struct {
	FinishedData interface{}
}

func (c *Context) SetFinishedData(data interface{}) {
	c.FinishedData = data
}

type KeptnEventData struct {
	Project string
	Stage   string
	Service string
}

//go:generate moq  -pkg fake -out ./fake/task_handler_mock.go . TaskHandler
type TaskHandler interface {
	Execute(keptnHandle IKeptn, ce interface{}, context Context) (Context, error)
	GetData() interface{}
}

type KeptnOption func(IKeptn)

func WithHandler(handler TaskHandler, eventType string) KeptnOption {
	return func(k IKeptn) {
		k.GetTaskRegistry().Add(eventType, TaskEntry{TaskHandler: handler})
	}
}

type IKeptn interface {
	Start()
	GetResourceHandler() ResourceHandler
	GetTaskRegistry() *TaskRegistry
}

type Keptn struct {
	EventSender     EventSender
	EventReceiver   EventReceiver
	ResourceHandler ResourceHandler
	Source          string
	TaskRegistry    TaskRegistry
}

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

func (k *Keptn) Start() {
	ctx := context.Background()
	ctx = cloudevents.WithEncodingStructured(ctx)
	err := k.EventReceiver.StartReceiver(ctx, k.gotEvent)
	_ = err
}

func (k *Keptn) GetResourceHandler() ResourceHandler {
	return k.ResourceHandler
}

func (k *Keptn) GetTaskRegistry() *TaskRegistry {
	return &k.TaskRegistry
}

func (k *Keptn) gotEvent(event cloudevents.Event) {
	if handler, ok := k.TaskRegistry.Contains(event.Type()); ok {
		data := handler.TaskHandler.GetData()
		if err := event.DataAs(&data); err != nil {
			k.send(k.createErrorFinishedEventForTriggeredEvent(event, nil))
		}
		k.send(k.createStartedEventForTriggeredEvent(event))

		newContext, err := handler.TaskHandler.Execute(k, data, handler.Context)
		if err != nil {
			k.send(k.createErrorFinishedEventForTriggeredEvent(event, newContext.FinishedData))
			return
		}
		k.send(k.createFinishedEventForTriggeredEvent(event, newContext.FinishedData))
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

	finishedEventType := strings.TrimSuffix(triggeredEvent.Type(), ".triggered") + ".finished"
	keptnContext, _ := triggeredEvent.Context.GetExtension(KeptnContextCEExtension)
	c := cloudevents.NewEvent()
	c.SetID(uuid.New().String())
	c.SetType(finishedEventType)
	c.SetDataContentType(cloudevents.ApplicationJSON)
	c.SetExtension(KeptnContextCEExtension, keptnContext)
	c.SetExtension(TriggeredIDCEExtension, triggeredEvent.ID())
	c.SetSource(k.Source)
	c.SetData(cloudevents.ApplicationJSON, eventData)
	return c

}

func (k *Keptn) createErrorFinishedEventForTriggeredEvent(triggeredEvent cloudevents.Event, eventData interface{}) cloudevents.Event {
	commonEventData := keptnv2.EventData{}
	if eventData != nil {
		keptnv2.Decode(eventData, &commonEventData)
		commonEventData.Status = keptnv2.StatusErrored
		commonEventData.Result = keptnv2.ResultFailed
		commonEventData.Message = "MSG_FAIL"
	} else {
		commonEventData.Status = keptnv2.StatusErrored
		commonEventData.Result = keptnv2.ResultFailed
		commonEventData.Message = "no valid event data"
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
