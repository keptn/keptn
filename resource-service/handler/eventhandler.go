package handler

import (
	"github.com/gin-gonic/gin"
	keptnmodels "github.com/keptn/go-utils/pkg/api/models"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/resource-service/models"
	logger "github.com/sirupsen/logrus"
	"net/http"
)

type IEventHandler interface {
	HandleEvent(context *gin.Context)
	Process(event models.Event, sync bool) error
}

type EventHandler struct {
	pm   *ProjectManager
	Self string
}

func (eh *EventHandler) HandleEvent(c *gin.Context) {

	event := &models.Event{}
	if err := c.ShouldBindJSON(event); err != nil {
		OnAPIError(c, err)
		return
	}
	keptnEvent := &keptnmodels.KeptnContextExtendedCE{}
	if err := keptnv2.Decode(event, keptnEvent); err != nil {
		OnAPIError(c, err)
		return
	}
	if err := keptnEvent.Validate(); err != nil {
		OnAPIError(c, err)
		return
	}

	err := eh.Process(*event, false)

	if err != nil {
		OnAPIError(c, err)
		return
	}
	c.Status(http.StatusOK)

}

// NewEventHandler creates a new EventHandler
func NewEventHandler(pm *ProjectManager, name string) IEventHandler {
	return &EventHandler{
		pm:   pm,
		Self: name,
	}
}

func (eh *EventHandler) Process(event models.Event, sync bool) error {
	keptnEvent := &keptnv2.ProjectDeleteFinishedEventData{}

	if err := keptnv2.Decode(event.Data, keptnEvent); err != nil {
		return err
	}

	// if shipyard-controller or any other replica managed to remove a project
	if keptnEvent.Status == keptnv2.StatusSucceeded && *event.Source != eh.Self {
		logger.Debug("Deleting project", keptnEvent.Project)
		err := eh.pm.DeleteProject(keptnEvent.Project)

		if err != nil {
			return err
		}
	}
	return nil
}
