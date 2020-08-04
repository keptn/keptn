package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/keptn/keptn/shipyard-controller/db"
	"github.com/keptn/keptn/shipyard-controller/models"
	"github.com/keptn/keptn/shipyard-controller/restapi/operations"
)

func GetTriggeredEvents(params operations.GetTriggeredEventsParams) middleware.Responder {
	em := getEventManagerInstance()

	var events []models.Event
	var err error

	if params.ProjectName != "" {
		events, err = em.GetTriggeredEventsOfProject(params.ProjectName)
	} else {
		events, err = em.GetAllTriggeredEvents()
	}

	return nil
}

func HandleEvent(params operations.HandleEventParams) middleware.Responder {
	return nil
}

var eventManagerInstance *eventManager

type eventManager struct {
	projectRepo        db.ProjectRepo
	triggeredEventRepo db.TriggeredEventRepo
}

func getEventManagerInstance() *eventManager {
	if eventManagerInstance == nil {
		eventManagerInstance = &eventManager{
			projectRepo:        &db.ProjectMongoDBRepo{},
			triggeredEventRepo: &db.MongoDBTriggeredEventsRepo{},
		}
	}
	return eventManagerInstance
}

func (em *eventManager) GetAllTriggeredEvents() ([]models.Event, error) {
	projects, err := em.projectRepo.GetProjects()

	if err != nil {
		return nil, err
	}

	allEvents := []models.Event{}
	for _, project := range projects {
		events, err := em.triggeredEventRepo.GetEvents(project)
		if err == nil {
			allEvents = append(allEvents, events...)
		}
	}
	return allEvents, nil
}

func (em *eventManager) GetTriggeredEventsOfProject(project string) ([]models.Event, error) {
	return em.triggeredEventRepo.GetEvents(project)
}
