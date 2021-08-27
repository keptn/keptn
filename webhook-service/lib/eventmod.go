package lib

import (
	"fmt"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/go-sdk/pkg/sdk"
)

type EventDataModifier struct {
	eventData    keptnv2.EventData
	eventDataMap map[string]interface{}
}

func NewEventDataModifier(event sdk.KeptnEvent) (*EventDataModifier, error) {
	eventData := keptnv2.EventData{}
	if err := keptnv2.Decode(event.Data, &eventData); err != nil {
		return nil, fmt.Errorf("could not decode incoming event payload: %w", err)
	}

	eventDataMap := map[string]interface{}{}
	if err := keptnv2.Decode(event, &eventDataMap); err != nil {
		return nil, fmt.Errorf("could not apply attributes from incoming event: %w", err)
	}
	return &EventDataModifier{eventData: eventData, eventDataMap: eventDataMap}, nil
}

func (e *EventDataModifier) Get() map[string]interface{} {
	return e.eventDataMap
}

func (e *EventDataModifier) Project() string {
	return e.eventData.Project
}

func (e *EventDataModifier) Stage() string {
	return e.eventData.Stage
}

func (e *EventDataModifier) Service() string {
	return e.eventData.Service
}

func (e *EventDataModifier) Add(key string, value interface{}) {
	e.eventDataMap[key] = value
}
