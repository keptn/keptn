package log

import (
	"context"
	"fmt"
	"strings"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	apimodels "github.com/keptn/go-utils/pkg/api/models"
	keptn "github.com/keptn/go-utils/pkg/api/utils"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	logger "github.com/sirupsen/logrus"
)

type UniformLog interface {
	Start(ctx context.Context, eventChannel chan cloudevents.Event)
}

type EventUniformLog struct {
	IntegrationID string
	logHandler    keptn.ILogHandler
}

func New(integrationID string, logHandler keptn.LogsV1Interface) *EventUniformLog {
	return &EventUniformLog{
		IntegrationID: integrationID,
		logHandler:    logHandler,
	}
}

func (l *EventUniformLog) Start(ctx context.Context, eventChannel chan cloudevents.Event) {
	logger.Infof("Starting uniform log for Keptn service with integration ID %s", l.IntegrationID)
	l.logHandler.Start(ctx)
	go func() {
		for {
			select {
			case event := <-eventChannel:
				logger.Debugf("Received event: %s", event.Context.GetType())
				if err := l.onEvent(event); err != nil {
					logger.Errorf("Could not handle event: %v", err)
				}
			case <-ctx.Done():
				logger.Info("Closing UniformLogger")
				return
			}
		}
	}()
}

func (l *EventUniformLog) onEvent(event cloudevents.Event) error {
	keptnEvent, err := keptnv2.ToKeptnEvent(event)
	if err != nil {
		return fmt.Errorf("could not decode CloudEvent to Keptn event: %w", err)
	}
	if strings.HasSuffix(*keptnEvent.Type, ".finished") {
		eventData := &keptnv2.EventData{}
		if err := keptnv2.EventDataAs(keptnEvent, eventData); err != nil {
			return fmt.Errorf("could not decode Keptn event data: %w", err)
		}

		taskName, _, _ := keptnv2.ParseTaskEventType(*keptnEvent.Type)

		if eventData.Status == keptnv2.StatusErrored {
			logger.Info("Received '.finished' event with status 'errored'. Forwarding log message to log ingestion API")
			l.Log(apimodels.LogEntry{
				IntegrationID: l.IntegrationID,
				Message:       eventData.Message,
				KeptnContext:  keptnEvent.Shkeptncontext,
				Task:          taskName,
				TriggeredID:   keptnEvent.Triggeredid,
			})
		}
		return nil
	} else if *keptnEvent.Type == keptnv2.ErrorLogEventName {
		logger.Info("Received 'log.error' event. Forwarding log message to log ingestion API")

		eventData := &keptnv2.ErrorLogEvent{}
		if err := keptnv2.EventDataAs(keptnEvent, eventData); err != nil {
			return fmt.Errorf("unable decode Keptn event data: %w", err)
		}

		integrationID := l.IntegrationID
		if eventData.IntegrationID != "" {
			// overwrite default integrationID if it has been set in the event
			integrationID = eventData.IntegrationID
		}
		l.Log(apimodels.LogEntry{
			IntegrationID: integrationID,
			Message:       eventData.Message,
			KeptnContext:  keptnEvent.Shkeptncontext,
			Task:          eventData.Task,
			TriggeredID:   keptnEvent.Triggeredid,
		})
	}
	return nil
}

func (l *EventUniformLog) Log(entry apimodels.LogEntry) {
	l.logHandler.Log([]apimodels.LogEntry{entry})
}
