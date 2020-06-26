package handlers

import (
	"fmt"
	"net/url"
	"time"

	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/types"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"github.com/google/uuid"
	keptnevents "github.com/keptn/go-utils/pkg/lib"
	keptnutils "github.com/keptn/go-utils/pkg/lib"
	"github.com/keptn/keptn/api/models"
	"github.com/keptn/keptn/api/restapi/operations/project"
	"github.com/keptn/keptn/api/utils"
	"github.com/keptn/keptn/api/ws"
)

type projectCreateEventData struct {
	keptnevents.ProjectCreateEventData `json:",inline"`
	EventContext                       models.EventContext `json:"eventContext"`
}

type projectDeleteEventData struct {
	keptnevents.ProjectDeleteEventData `json:",inline"`
	EventContext                       models.EventContext `json:"eventContext"`
}

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

	eventContext := models.EventContext{KeptnContext: &keptnContext, Token: &token}
	source, _ := url.Parse("https://github.com/keptn/keptn/api")

	prjData := keptnevents.ProjectCreateEventData{
		Project:      *params.Project.Name,
		Shipyard:     *params.Project.Shipyard,
		GitUser:      params.Project.GitUser,
		GitToken:     params.Project.GitToken,
		GitRemoteURL: params.Project.GitRemoteURL,
	}
	forwardData := projectCreateEventData{ProjectCreateEventData: prjData, EventContext: eventContext}

	contentType := "application/json"
	event := cloudevents.Event{
		Context: cloudevents.EventContextV02{
			ID:          uuid.New().String(),
			Time:        &types.Timestamp{Time: time.Now()},
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

	return project.NewPostProjectOK().WithPayload(&eventContext)
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

	eventContext := models.EventContext{KeptnContext: &keptnContext, Token: &token}

	source, _ := url.Parse("https://github.com/keptn/keptn/api")

	prjData := keptnevents.ProjectDeleteEventData{
		Project: params.ProjectName,
	}
	forwardData := projectDeleteEventData{ProjectDeleteEventData: prjData, EventContext: eventContext}

	contentType := "application/json"
	event := cloudevents.Event{
		Context: cloudevents.EventContextV02{
			ID:          uuid.New().String(),
			Time:        &types.Timestamp{Time: time.Now()},
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

	return project.NewDeleteProjectProjectNameOK().WithPayload(&eventContext)
}

func getProjectPostInternalError(err error) *project.PostProjectDefault {
	return project.NewPostProjectDefault(500).WithPayload(
		&models.Error{Code: 500, Message: swag.String(err.Error())})
}

func getProjectDeleteInternalError(err error) *project.DeleteProjectProjectNameDefault {
	return project.NewDeleteProjectProjectNameDefault(500).WithPayload(
		&models.Error{Code: 500, Message: swag.String(err.Error())})
}
