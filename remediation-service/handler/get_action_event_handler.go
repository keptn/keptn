package handler

import (
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/go-utils/pkg/lib/v0_1_4"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/go-sdk/pkg/sdk"
)

const remediationSpecVersion = "spec.keptn.sh/0.1.4"
const remediationResourceFileName = "remediation.yaml"

type GetActionEventHandler struct {
}

func NewGetActionEventHandler() *GetActionEventHandler {
	return &GetActionEventHandler{}
}

func (g *GetActionEventHandler) Execute(k sdk.IKeptn, data interface{}, _ string) (interface{}, *sdk.Error) {
	getActionTriggeredData := &keptnv2.GetActionTriggeredEventData{}

	if err := keptnv2.Decode(data, getActionTriggeredData); err != nil {
		return nil, &sdk.Error{Err: err, StatusType: keptnv2.StatusErrored, ResultType: keptnv2.ResultFailed, Message: "Could not decode input event data"}
	}

	// get remediation.yaml resource
	resource, err := g.getRemediationResource(k, getActionTriggeredData)
	if err != nil {
		return nil, &sdk.Error{Err: err, StatusType: keptnv2.StatusErrored, ResultType: keptnv2.ResultFailed, Message: "Could not get remediation.yaml file for services " + getActionTriggeredData.Service + " in stage " + getActionTriggeredData.Stage + "."}
	}

	// parse remediation.yaml resource
	remediation, err := ParseRemediationResource(resource)
	if err != nil {
		return nil, &sdk.Error{Err: err, StatusType: keptnv2.StatusErrored, ResultType: keptnv2.ResultFailed, Message: "Could not parse remediation.yaml file. Please validate it against the specification."}
	}

	// determine next action
	action, err := GetNextAction(remediation, getActionTriggeredData.Problem, getActionTriggeredData.ActionIndex)
	if err != nil {
		return nil, &sdk.Error{Err: err, StatusType: keptnv2.StatusSucceeded, ResultType: keptnv2.ResultFailed, Message: "No more actions defined for problem type " + getActionTriggeredData.Problem.RootCause + " in remediation.yaml file."}
	}

	finishedEventData := keptnv2.GetActionFinishedEventData{
		EventData:   getActionTriggeredData.EventData,
		Action:      *action,
		ActionIndex: getActionTriggeredData.ActionIndex + 1,
	}

	return finishedEventData, nil
}

func (g *GetActionEventHandler) getRemediationResource(keptn sdk.IKeptn, eventData *keptnv2.GetActionTriggeredEventData) (*models.Resource, error) {
	if eventData.Service == "" {
		return keptn.GetResourceHandler().GetStageResource(eventData.Project, eventData.Stage, remediationResourceFileName)
	}
	return keptn.GetResourceHandler().GetServiceResource(eventData.Project, eventData.Stage, eventData.Service, remediationResourceFileName)
}

// ParseRemediationResource returns the in-memory representation of a keptn resource.
// Note, that the spec version of the remediation.yaml file needs to match "spec.keptn.sh/0.1.4"
func ParseRemediationResource(resource *models.Resource) (*v0_1_4.Remediation, error) {
	remediationData := &v0_1_4.Remediation{}
	err := yaml.Unmarshal([]byte(resource.ResourceContent), remediationData)
	if err != nil {
		return nil, fmt.Errorf("could not parse remediation.yaml: %w", err)
	}

	if remediationData.ApiVersion != remediationSpecVersion {
		return nil, fmt.Errorf("remediation.yaml file does not conform to remediation spec %s", remediationSpecVersion)
	}
	return remediationData, nil
}

// GetNextAction contains the logic to determine, what will be the next remediation action according to the remediation.yaml file
// It searches for a problem type matching the root cause of the problem. If no problem type is found a problem type matching the problem title
// will be searched as a fallback. If still no problem type is found it will return an error.

// The actionIndex parameter specifies which action to take if a problem type was found.
func GetNextAction(remediation *v0_1_4.Remediation, problemDetails keptnv2.ProblemDetails, actionIndex int) (*keptnv2.ActionInfo, error) {
	rootCause := problemDetails.RootCause
	problemTitle := problemDetails.ProblemTitle

	var actions []v0_1_4.RemediationActionsOnOpen
	// search problem type matching root cause
	for _, r := range remediation.Spec.Remediations {
		if r.ProblemType == rootCause {
			actions = r.ActionsOnOpen
			break
		}
	}

	// fallback: search problem type matching problem title
	if actions == nil {
		for _, r := range remediation.Spec.Remediations {
			if r.ProblemType == problemTitle {
				actions = r.ActionsOnOpen
				break
			}
		}
	}

	// fallback: search problem type default
	if actions == nil {
		for _, r := range remediation.Spec.Remediations {
			if r.ProblemType == "default" {
				actions = r.ActionsOnOpen
				break
			}
		}
	}

	// we did not find an action
	if actions == nil {
		return nil, fmt.Errorf("unable to find actions for root cause %s", rootCause)
	}

	// the required action does not exist
	if actionIndex >= len(actions) {
		return nil, fmt.Errorf("failed to get action for root cause %s. There is no action with index %d", rootCause, actionIndex)
	}

	action := actions[actionIndex]
	return &keptnv2.ActionInfo{
		Name:        action.Name,
		Action:      action.Action,
		Description: action.Description,
		Value:       action.Value,
	}, nil

}
