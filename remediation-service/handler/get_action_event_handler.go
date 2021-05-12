package handler

import (
	"github.com/keptn/go-utils/pkg/api/models"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/remediation-service/internal/sdk"
)

const remediationSpecVersion = "spec.keptn.sh/0.1.4"
const remediationResourceFileName = "remediation.yaml"

type GetActionEventHandler struct {
	GetActionTriggeredData *keptnv2.GetActionTriggeredEventData
}

func NewGetActionEventHandler() *GetActionEventHandler {
	return &GetActionEventHandler{GetActionTriggeredData: &keptnv2.GetActionTriggeredEventData{}}
}

func (g *GetActionEventHandler) Execute(k sdk.IKeptn, ce interface{}) (interface{}, *sdk.Error) {
	data := ce.(*keptnv2.GetActionTriggeredEventData)

	// get remediation.yaml resource
	resource, err := g.getRemediationResource(k, data)
	if err != nil {
		return nil, &sdk.Error{Err: err, StatusType: keptnv2.StatusErrored, ResultType: keptnv2.ResultFailed, Message: "unable to get remediation.yaml"}
	}

	// parse remediation.yaml resource
	remediation, err := ParseRemediationResource(resource)
	if err != nil {
		return nil, &sdk.Error{Err: err, StatusType: keptnv2.StatusErrored, ResultType: keptnv2.ResultFailed, Message: "unable to parse remediation.yaml"}
	}

	// determine next action
	action, err := GetNextAction(remediation, data.Problem, data.ActionIndex)
	if err != nil {
		return nil, &sdk.Error{Err: err, StatusType: keptnv2.StatusSucceeded, ResultType: keptnv2.ResultFailed, Message: "unable to get next action from remediation.yaml"}
	}

	finishedEventData := keptnv2.GetActionFinishedEventData{
		EventData:   data.EventData,
		Action:      *action,
		ActionIndex: data.ActionIndex + 1,
	}

	return finishedEventData, nil
}

func (g *GetActionEventHandler) GetTriggeredData() interface{} {
	return g.GetActionTriggeredData
}

func (g *GetActionEventHandler) getRemediationResource(keptn sdk.IKeptn, eventData *keptnv2.GetActionTriggeredEventData) (*models.Resource, error) {
	if eventData.Service == "" {
		return keptn.GetResourceHandler().GetStageResource(eventData.Project, eventData.Stage, remediationResourceFileName)
	}

	return keptn.GetResourceHandler().GetServiceResource(eventData.Project, eventData.Stage, eventData.Service, remediationResourceFileName)

}
