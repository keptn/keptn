package handler

import (
	"github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/remediation-service/pkg/sdk"
)

const remediationSpecVersion = "spec.keptn.sh/0.1.4"

type GetActionEventHandler struct {
	GetActionTriggeredData *v0_2_0.GetActionTriggeredEventData
}

func NewGetActionEventHandler() *GetActionEventHandler {
	return &GetActionEventHandler{GetActionTriggeredData: &v0_2_0.GetActionTriggeredEventData{}}
}

func (g *GetActionEventHandler) Execute(k sdk.IKeptn, ce interface{}, context sdk.Context) (sdk.Context, error) {
	data := ce.(*v0_2_0.GetActionTriggeredEventData)

	// get remediation.yaml resource
	resource, err := g.getRemediationResource(k, data)
	if err != nil {
		return context, err
	}

	// parse remediation.yaml resource
	remediation, err := ParseRemediationResource(resource)
	if err != nil {
		return context, err
	}

	// determine next action
	action, err := GetNextAction(remediation, data.ProblemDetails, data.ActionIndex)
	if err != nil {
		return context, err
	}

	// set finished data
	finishedEventData := v0_2_0.GetActionFinishedEventData{
		EventData:   data.EventData,
		Action:      *action,
		ActionIndex: data.ActionIndex + 1,
	}
	context.SetFinishedData(finishedEventData)

	return context, nil
}

func (g *GetActionEventHandler) GetData() interface{} {
	return g.GetActionTriggeredData
}

func (g *GetActionEventHandler) getRemediationResource(keptn sdk.IKeptn, eventData *v0_2_0.GetActionTriggeredEventData) (*models.Resource, error) {
	if eventData.Service == "" {
		return keptn.GetResourceHandler().GetStageResource(eventData.Project, eventData.Stage, "remediation.yaml")
	}

	return keptn.GetResourceHandler().GetServiceResource(eventData.Project, eventData.Stage, eventData.Service, "remediation.yaml")

}
