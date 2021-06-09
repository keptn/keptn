package lib

import (
	"context"
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
	GetChannel() chan cloudevents.Event
}

type EventUniformLog struct {
	Integration  keptnapimodels.Integration
	logHandler   keptn.LogHandler
	eventChannel chan cloudevents.Event
}

func NewEventUniformLog(integration keptnapimodels.Integration) *EventUniformLog {
	return &EventUniformLog{
		Integration: integration,
		logHandler:  keptn.LogHandler{},
	}
}

func (l *EventUniformLog) GetChannel() chan cloudevents.Event {
	return l.eventChannel
}

func (l *EventUniformLog) Start(ctx context.Context, eventChannel chan cloudevents.Event) {
	l.logHandler.Start(ctx)
	go func() {
		for {
			event := <-eventChannel
			if err := l.OnEvent(event); err != nil {
				logger.Errorf("could not handle event: %s", err.Error())
			}
			return
		}
	}()
}

func (l *EventUniformLog) OnEvent(event cloudevents.Event) error {
	if !strings.HasSuffix(event.Context.GetType(), ".finished") {
		return nil
	}
	// TODO: also check for log event (needs to be defined)
	keptnEvent, err := keptnv2.ToKeptnEvent(event)
	if err != nil {
		return fmt.Errorf("could not decode CloudEvent to Keptn event: %v", err.Error())
	}

	eventData := &keptnv2.EventData{}
	if err := keptnv2.EventDataAs(keptnEvent, eventData); err != nil {
		return fmt.Errorf("could not decode Keptn event data: %v", err.Error())
	}

	if eventData.Status == keptnv2.StatusErrored {
		l.Log(keptnapimodels.LogEntry{
			IntegrationID: l.Integration.ID,
			Message:       eventData.Message,
		})
	}
	return nil
}

func (l *EventUniformLog) Log(entry keptnapimodels.LogEntry) {
	l.logHandler.Log([]keptnapimodels.LogEntry{entry})
}
