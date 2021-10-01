package lib

import (
	"errors"
	"fmt"
	keptnmodels "github.com/keptn/go-utils/pkg/api/models"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/go-sdk/pkg/sdk"
)

type DistributorData struct {
	SubscriptionID string `json:"subscriptionID"`
}

type TemporaryData struct {
	TemporaryData struct {
		Distributor DistributorData `json:"distributor"`
	} `json:"temporaryData"`
}

type EventDataAdapter struct {
	event        keptnmodels.KeptnContextExtendedCE
	eventData    keptnv2.EventData
	eventDataMap map[string]interface{}
}

func NewEventDataAdapter(event sdk.KeptnEvent) (*EventDataAdapter, error) {
	eventData := keptnv2.EventData{}
	if err := keptnv2.Decode(event.Data, &eventData); err != nil {
		return nil, fmt.Errorf("could not decode incoming event payload: %w", err)
	}

	if eventData.Project == "" || eventData.Stage == "" || eventData.Service == "" {
		return nil, fmt.Errorf("project, stage and service must be present in the event data")
	}
	eventDataMap := map[string]interface{}{}
	if err := keptnv2.Decode(event, &eventDataMap); err != nil {
		return nil, fmt.Errorf("could not apply attributes from incoming event: %w", err)
	}
	keptnEvent := keptnmodels.KeptnContextExtendedCE{}
	if err := keptnv2.Decode(event, &keptnEvent); err != nil {
		return nil, fmt.Errorf("could not decode incoming event payload: %w", err)
	}

	return &EventDataAdapter{event: keptnEvent, eventData: eventData, eventDataMap: eventDataMap}, nil
}

func (e *EventDataAdapter) Get() map[string]interface{} {
	return e.eventDataMap
}

func (e *EventDataAdapter) Project() string {
	return e.eventData.Project
}

func (e *EventDataAdapter) Stage() string {
	return e.eventData.Stage
}

func (e *EventDataAdapter) Service() string {
	return e.eventData.Service
}

func (e *EventDataAdapter) SubscriptionID() (string, error) {
	distributorData := &DistributorData{}
	if err := e.event.GetTemporaryData("distributor", distributorData); err != nil {
		return "", err
	}
	if distributorData.SubscriptionID == "" {
		return "", errors.New("no subscription ID found in event")
	}
	return distributorData.SubscriptionID, nil
}

func (e *EventDataAdapter) Labels() interface{} {
	return e.eventData.Labels
}

func (e *EventDataAdapter) Add(key string, value interface{}) {
	e.eventDataMap[key] = value
}

func (e *EventDataAdapter) Remove(key string) {
	delete(e.eventDataMap, key)
}
