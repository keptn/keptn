package handlers

import (
	"fmt"
	"net/url"

	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/types"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"github.com/google/uuid"
	keptnevents "github.com/keptn/go-utils/pkg/events"
	keptnutils "github.com/keptn/go-utils/pkg/utils"
	"github.com/keptn/keptn/api/models"
	"github.com/keptn/keptn/api/restapi/operations/project"
	"github.com/keptn/keptn/api/utils"
	"github.com/keptn/keptn/api/ws"
)

// PostProjectHandlerFunc creates a new project
func PostProjectHandlerFunc(params project.PostProjectParams, p *models.Principal) middleware.Responder {

	keptnContext := uuid.New().String()
	l := keptnutils.NewLogger(keptnContext, "", "api")
	l.Info("API received create for project")

	token, err := ws.CreateChannelInfo(keptnContext)
	if err != nil {
		l.Error(fmt.Sprintf("Error creating channel info %s", err.Error()))
		return getProjectPostInternalError(err)
	}

	channelInfo := models.ChannelInfo{ChannelID: &keptnContext, Token: &token}

	source, _ := url.Parse("https://github.com/keptn/keptn/api")

	prjData := keptnevents.ProjectCreateEventData{
		Project:      params.Project.Name,
		Shipyard:     params.Project.Shipyard,
		GitUser:      params.Project.GitUser,
		GitToken:     params.Project.GitToken,
		GitRemoteURL: params.Project.GitRemoteURL,
	}
	forwardData := EnrichedCEData{Data: prjData, ChannelInfo: channelInfo}

	contentType := "application/json"
	event := cloudevents.Event{
		Context: cloudevents.EventContextV02{
			ID:          uuid.New().String(),
			Type:        keptnevents.InternalProjectCreateEventType,
			Source:      types.URLRef{URL: *source},
			ContentType: &contentType,
			Extensions:  map[string]interface{}{"shkeptncontext": keptnContext},
		}.AsV02(),
		Data: forwardData,
	}

	_, err = utils.PostToEventBroker(event)
	if err != nil {
		l.Error(fmt.Sprintf("Error sending CloudEvent %s", err.Error()))
		return getProjectPostInternalError(err)
	}

	return project.NewPostProjectOK().WithPayload(&channelInfo)
}

// DeleteProjectProjectNameHandlerFunc deletes a project
func DeleteProjectProjectNameHandlerFunc(params project.DeleteProjectProjectNameParams, p *models.Principal) middleware.Responder {

	keptnContext := uuid.New().String()
	l := keptnutils.NewLogger(keptnContext, "", "api")
	l.Info("API received delete for project")

	token, err := ws.CreateChannelInfo(keptnContext)
	if err != nil {
		l.Error(fmt.Sprintf("Error creating channel info %s", err.Error()))
		return getProjectDeleteInternalError(err)
	}

	channelInfo := models.ChannelInfo{ChannelID: &keptnContext, Token: &token}

	source, _ := url.Parse("https://github.com/keptn/keptn/api")

	prjData := keptnevents.ProjectDeleteEventData{
		Project: params.ProjectName,
	}
	forwardData := EnrichedCEData{Data: prjData, ChannelInfo: channelInfo}

	contentType := "application/json"
	event := cloudevents.Event{
		Context: cloudevents.EventContextV02{
			ID:          uuid.New().String(),
			Type:        keptnevents.InternalProjectDeleteEventType,
			Source:      types.URLRef{URL: *source},
			ContentType: &contentType,
			Extensions:  map[string]interface{}{"shkeptncontext": keptnContext},
		}.AsV02(),
		Data: forwardData,
	}
	event.Extensions()["shkeptncontext"] = keptnContext

	_, err = utils.PostToEventBroker(event)
	if err != nil {
		l.Error(fmt.Sprintf("Error sending CloudEvent %s", err.Error()))
		return getProjectDeleteInternalError(err)
	}

	return project.NewDeleteProjectProjectNameOK().WithPayload(&channelInfo)
}

func getProjectPostInternalError(err error) *project.PostProjectDefault {
	return project.NewPostProjectDefault(500).WithPayload(
		&models.Error{Code: 500, Message: swag.String(err.Error())})
}

func getProjectDeleteInternalError(err error) *project.DeleteProjectProjectNameDefault {
	return project.NewDeleteProjectProjectNameDefault(500).WithPayload(
		&models.Error{Code: 500, Message: swag.String(err.Error())})
}
