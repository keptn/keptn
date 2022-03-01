package handler

import (
	"github.com/gin-gonic/gin"
	keptnmodels "github.com/keptn/go-utils/pkg/api/models"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	logger "github.com/sirupsen/logrus"
	"net/http"
	"os"
)

type IEventHandler interface {
	HandleEvent(context *gin.Context)
}

type EventHandler struct {
	pm *ProjectManager
}

func (eh *EventHandler) HandleEvent(c *gin.Context) {
	logger.Debug("Received an event")
	keptnEvent := &keptnmodels.KeptnContextExtendedCE{}
	if err := c.ShouldBindJSON(keptnEvent); err != nil {
		OnAPIError(c, err)
		return
	}

	logger.Debug("Handling event", keptnEvent)

	event := &keptnv2.ProjectDeleteFinishedEventData{}

	if err := keptnv2.Decode(keptnEvent.Data, event); err != nil {
		OnAPIError(c, err)
		return
	}
	logger.Debug("Decoded event", event)

	self := os.Getenv(EnvKubernetesPodName)
	// if shipyard-controller or any other replica managed to remove a project
	if event.Status == keptnv2.StatusSucceeded && *keptnEvent.Source != self {
		logger.Debug("Deleting project", event.Project)
		err := eh.pm.DeleteProject(event.Project)

		if err != nil {
			OnAPIError(c, err)
			return
		}
	}
	c.Status(http.StatusOK)

}

// NewEventHandler creates a new EventHandler
func NewEventHandler(pm *ProjectManager) IEventHandler {
	return &EventHandler{
		pm: pm,
	}
}
