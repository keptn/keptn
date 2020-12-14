package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/configuration-service/common"
	"github.com/keptn/keptn/configuration-service/models"
	"github.com/keptn/keptn/configuration-service/restapi/operations/event"
)

func HandleEventHandlerFunc(params event.HandleEventParams) middleware.Responder {
	if params.Body.Type != nil {
		keptnBase := &keptnv2.EventData{}
		err := keptnv2.Decode(params.Body.Data, keptnBase)
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

		common.LockProject(keptnBase.Project)
		defer common.UnlockProject(keptnBase.Project)
		mv := common.GetProjectsMaterializedView()
		err = mv.UpdateEventOfService(params.Body.Data, *params.Body.Type, params.Body.Shkeptncontext, params.Body.ID, params.Body.Triggeredid)
		if err != nil {
			return event.NewHandleEventDefault(400).WithPayload(&models.Error{Message: swag.String("Could not update service: " + err.Error()), Code: 400})
		}

		return event.NewHandleEventOK()
	}

	return event.NewHandleEventDefault(400).WithPayload(&models.Error{Message: swag.String("Invalid event")})
}
