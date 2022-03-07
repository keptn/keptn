package nats

import (
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/resource-service/handler"
	"github.com/keptn/keptn/resource-service/models"
	logger "github.com/sirupsen/logrus"
)

type EventMsgHandler struct {
	pm   handler.IProjectManager
	Self string
}

func EventHandler(projectManager handler.IProjectManager, name string) *EventMsgHandler {
	return &EventMsgHandler{
		pm:   projectManager,
		Self: name,
	}
}
func (eh *EventMsgHandler) Process(event models.Event, sync bool) error {
	e := &keptnv2.ProjectDeleteFinishedEventData{}
	if err := keptnv2.Decode(event.Data, e); err != nil {
		return err
	}

	// if shipyard-controller or any other replica managed to remove a project
	if e.Status == keptnv2.StatusSucceeded && *event.Source != eh.Self {
		logger.Debug("Deleting project", e.Project)
		err := eh.pm.DeleteProject(e.Project)
		if err != nil {
			return err
		}
	}
	return nil
}
