package fake

import (
	"github.com/keptn/keptn/shipyard-controller/common"
	"github.com/keptn/keptn/shipyard-controller/models"
)

type ShipyardController struct {
	GetAllTriggeredEventsFunc       func(filter common.EventFilter) ([]models.Event, error)
	GetTriggeredEventsOfProjectFunc func(project string, filter common.EventFilter) ([]models.Event, error)
	HandleIncomingEventFunc         func(event models.Event) error
}

func (s *ShipyardController) GetAllTriggeredEvents(filter common.EventFilter) ([]models.Event, error) {
	return s.GetAllTriggeredEventsFunc(filter)
}

func (s *ShipyardController) GetTriggeredEventsOfProject(project string, filter common.EventFilter) ([]models.Event, error) {
	return s.GetTriggeredEventsOfProjectFunc(project, filter)
}

func (s *ShipyardController) HandleIncomingEvent(event models.Event) error {
	return s.HandleIncomingEventFunc(event)
}
