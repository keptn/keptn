package lib

import (
	"context"
	"errors"
	"fmt"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	keptnapimodels "github.com/keptn/go-utils/pkg/api/models"
	keptn "github.com/keptn/go-utils/pkg/api/utils"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	logger "github.com/sirupsen/logrus"
	"strings"
)

type UniformLog interface {
	Start(ctx context.Context, eventChannel chan cloudevents.Event)
}

type EventUniformLog struct {
	IntegrationID string
	logHandler    keptn.ILogHandler
}

func NewEventUniformLog(integrationID string, logHandler keptn.ILogHandler) *EventUniformLog {
	return &EventUniformLog{
		IntegrationID: integrationID,
		logHandler:    logHandler,
	}
}

func (l *EventUniformLog) Start(ctx context.Context, eventChannel chan cloudevents.Event) {
	logger.Info("starting UniformLog")
	l.logHandler.Start(ctx)
	go func() {
		for {
			select {
			case event := <-eventChannel:
				logger.Infof("UniformLogger: received event: %s", event.Context.GetType())
				if err := l.OnEvent(event); err != nil {
					logger.Errorf("could not handle event: %s", err.Error())
				}
			case <-ctx.Done():
				logger.Info("closing UniformLogger")
			}

		}
	}()
}

func (l *EventUniformLog) OnEvent(event cloudevents.Event) error {
	keptnEvent, err := keptnv2.ToKeptnEvent(event)
	if err != nil {
		return fmt.Errorf("could not decode CloudEvent to Keptn event: %v", err.Error())
	}
	if keptnEvent.Type == nil {
		return errors.New("received keptnEvent has no type")
	}
	logger.Infof("Received Event: %v", keptnEvent)
	if strings.HasSuffix(*keptnEvent.Type, ".finished") {
		eventData := &keptnv2.EventData{}
		if err := keptnv2.EventDataAs(keptnEvent, eventData); err != nil {
			return fmt.Errorf("could not decode Keptn event data: %v", err.Error())
		}

		taskName, _, _ := keptnv2.ParseTaskEventType(*keptnEvent.Type)

		if eventData.Status == keptnv2.StatusErrored {
			logger.Info("UniformLogger: received .finished event with status errored. forwarding log message to log ingestion API")
			l.Log(keptnapimodels.LogEntry{
				IntegrationID: l.IntegrationID,
				Message:       eventData.Message,
				KeptnContext:  keptnEvent.Shkeptncontext,
				Task:          taskName,
				TriggeredID:   keptnEvent.Triggeredid,
			})
		}
		return nil
	} else if *keptnEvent.Type == keptnv2.ErrorLogEventName {
		logger.Info("received log.error event. forwarding log message to log ingestion API")

		eventData := &keptnv2.ErrorLogEvent{}
		if err := keptnv2.EventDataAs(keptnEvent, eventData); err != nil {
			return fmt.Errorf("could not decode Keptn event data: %v", err.Error())
		}
		l.Log(keptnapimodels.LogEntry{
			IntegrationID: l.IntegrationID,
			Message:       eventData.Message,
			KeptnContext:  keptnEvent.Shkeptncontext,
			Task:          eventData.Task,
			TriggeredID:   keptnEvent.Triggeredid,
		})
	}
	return nil
}

func (l *EventUniformLog) Log(entry keptnapimodels.LogEntry) {
	l.logHandler.Log([]keptnapimodels.LogEntry{entry})
}
