package nats

import (
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/resource-service/handler"
	"github.com/keptn/keptn/resource-service/models"
	logger "github.com/sirupsen/logrus"
)

const shipyardController = "shipyard-controller"

type EventMsgHandler struct {
	pm handler.IProjectManager
}

func EventHandler(projectManager handler.IProjectManager) *EventMsgHandler {
	return &EventMsgHandler{
		pm: projectManager,
	}
}
func (eh *EventMsgHandler) Process(event models.Event) error {
	e := &keptnv2.ProjectDeleteFinishedEventData{}
	if err := keptnv2.Decode(event.Data, e); err != nil {
		return err
	}

	if e.Status == keptnv2.StatusSucceeded && *event.Source == shipyardController {
		logger.Infof("Deleting project %s", e.Project)
		err := eh.pm.DeleteProject(e.Project)

		if err != nil {
			return err
		}
	}
	return nil
}
