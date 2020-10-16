package api

import (
	"errors"
	"fmt"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	keptnapi "github.com/keptn/go-utils/pkg/api/utils"
	keptncommon "github.com/keptn/go-utils/pkg/lib/keptn"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/shipyard-controller/common"
	"github.com/keptn/keptn/shipyard-controller/models"
	"github.com/keptn/keptn/shipyard-controller/operations"
	"net/http"
	"net/url"
)

// CreateService godoc
// @Summary Create a new service
// @Description Create a new service
// @Tags Services
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param   project     body    operations.CreateServiceParams     true        "Project"
// @Success 200 {object} operations.CreateServiceResponse	"ok"
// @Failure 400 {object} models.Error "Invalid payload"
// @Failure 500 {object} models.Error "Internal error"
// @Router /project/:project/service [post]
func CreateService(c *gin.Context) {
	// validate the input
	createServiceParams := &operations.CreateServiceParams{}
	if err := c.ShouldBindJSON(createServiceParams); err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Code:    400,
			Message: stringp("Invalid request format: " + err.Error()),
		})
		return
	}
	if err := validateCreateServiceParams(createServiceParams); err != nil {
		c.JSON(http.StatusBadRequest, models.Error{
			Code:    400,
			Message: stringp("Could not validate payload: " + err.Error()),
		})
		return
	}

	pm, err := newServiceManager()
	if err != nil {

		c.JSON(http.StatusInternalServerError, models.Error{
			Code:    500,
			Message: stringp("Could not process request: " + err.Error()),
		})
		return
	}

	if err := pm.createService(createServiceParams); err != nil {
		if err == errServiceAlreadyExists {
			c.JSON(http.StatusConflict, models.Error{
				Code:    http.StatusConflict,
				Message: stringp(err.Error()),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.Error{
			Code:    http.StatusInternalServerError,
			Message: stringp(err.Error()),
		})
		return
	}
}

func validateCreateServiceParams(params *operations.CreateServiceParams) error {
	return nil
}

var errServiceAlreadyExists = errors.New("project already exists")

type serviceManager struct {
	projectAPI  *keptnapi.ProjectHandler
	stagesAPI   *keptnapi.StageHandler
	resourceAPI *keptnapi.ResourceHandler
	logger      keptncommon.LoggerInterface
}

func newServiceManager() (*serviceManager, error) {
	csEndpoint, err := keptncommon.GetServiceEndpoint("CONFIGURATION_SERVICE")
	if err != nil {
		return nil, fmt.Errorf("could not get configuration-service URL: %s", err.Error())
	}
	if err != nil {
		return nil, fmt.Errorf("could not initilize secret store: " + err.Error())
	}
	return &serviceManager{
		projectAPI:  keptnapi.NewProjectHandler(csEndpoint.String()),
		stagesAPI:   keptnapi.NewStageHandler(csEndpoint.String()),
		resourceAPI: keptnapi.NewResourceHandler(csEndpoint.String()),
		logger:      keptncommon.NewLogger("", "", "shipyard-controller"),
	}, nil
}

func (sm *serviceManager) createService(params *operations.CreateServiceParams) error {
	return nil
}

func (pm *projectManager) sendServiceCreateStartedEvent(projectName string, params *operations.CreateServiceParams) error {
	eventPayload := keptnv2.ProjectCreateStartedEventData{
		EventData: keptnv2.EventData{
			Project: projectName,
			Service: *params.Name,
		},
	}
	source, _ := url.Parse("shipyard-controller")
	eventType := keptnv2.GetFinishedEventType(keptnv2.ProjectCreateTaskName)
	event := cloudevents.NewEvent()
	event.SetType(eventType)
	event.SetSource(source.String())
	event.SetDataContentType(cloudevents.ApplicationJSON)
	event.SetExtension("shkeptncontext", uuid.New().String())
	event.SetData(cloudevents.ApplicationJSON, eventPayload)

	if err := common.SendEvent(event); err != nil {
		return errors.New("could not send create.project.started event: " + err.Error())
	}
	return nil
}

func (pm *projectManager) sendServiceCreateSuccessFinishedEvent(params *operations.CreateServiceParams) error {
	finishedEventPayload := keptnv2.ServiceCreateFinishedEventData{
		EventData: keptnv2.EventData{
			Project: *params.Name,
			Status:  keptnv2.StatusSucceeded,
			Result:  keptnv2.ResultPass,
		},
		Helm: keptnv2.Helm{
			Chart: "", // TODO
		},
	}
	source, _ := url.Parse("shipyard-controller")
	eventType := keptnv2.GetFinishedEventType(keptnv2.ProjectCreateTaskName)
	event := cloudevents.NewEvent()
	event.SetType(eventType)
	event.SetSource(source.String())
	event.SetDataContentType(cloudevents.ApplicationJSON)
	event.SetExtension("shkeptncontext", uuid.New().String())
	event.SetData(cloudevents.ApplicationJSON, finishedEventPayload)

	if err := common.SendEvent(event); err != nil {
		return errors.New("could not send create.project.finished event: " + err.Error())
	}
	return nil
}
