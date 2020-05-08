package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	keptn "github.com/keptn/go-utils/pkg/lib"
	"github.com/keptn/keptn/configuration-service/common"
	"github.com/keptn/keptn/configuration-service/models"
	"github.com/keptn/keptn/configuration-service/restapi/operations/event"
	"github.com/mitchellh/mapstructure"
)

func HandleEventHandlerFunc(params event.HandleEventParams) middleware.Responder {
	if params.Body.Type != nil {
		keptnBase := &keptn.KeptnBase{}
		err := mapstructure.Decode(params.Body.Data, keptnBase)
		if err != nil {
			return event.NewHandleEventDefault(400).WithPayload(&models.Error{Message: swag.String("Could not parse event data: " + err.Error()), Code: 400})
		}

		if keptnBase.Project == "" {
			return event.NewHandleEventDefault(400).WithPayload(&models.Error{Message: swag.String("Project must not be empty"), Code: 400})
		}

		if keptnBase.Stage == "" {
			return event.NewHandleEventDefault(400).WithPayload(&models.Error{Message: swag.String("Stage must not be empty"), Code: 400})
		}

		if keptnBase.Service == "" {
			return event.NewHandleEventDefault(400).WithPayload(&models.Error{Message: swag.String("Service must not be empty"), Code: 400})
		}

		mv := common.GetProjectsMaterializedView()
		err = mv.UpdateEventOfService(params.Body.Data, *params.Body.Type, params.Body.Shkeptncontext, params.Body.ID)
		if err != nil {
			return event.NewHandleEventDefault(400).WithPayload(&models.Error{Message: swag.String("Could not update service: " + err.Error()), Code: 400})
		}

		return event.NewHandleEventOK()
	}

	return event.NewHandleEventDefault(400).WithPayload(&models.Error{Message: swag.String("Invalid event")})
}
