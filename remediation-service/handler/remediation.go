package handler

import (
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/go-utils/pkg/lib/v0_1_4"
	"github.com/keptn/go-utils/pkg/lib/v0_2_0"
)

// ParseRemediationResource returns the in-memory representation of a keptn resource.
// Note, that the spec version of the remediation.yaml file needs to match "spec.keptn.sh/0.1.4"
func ParseRemediationResource(resource *models.Resource) (*v0_1_4.Remediation, error) {
	remediationData := &v0_1_4.Remediation{}
	err := yaml.Unmarshal([]byte(resource.ResourceContent), remediationData)
	if err != nil {
		return nil, fmt.Errorf("could not parse remediation.yaml: %w", err)
	}

	// TODO: it is probably no good idea to hard-code this here
	if remediationData.ApiVersion != remediationSpecVersion {
		return nil, fmt.Errorf("remediation.yaml file does not conform to remediation spec %s", remediationSpecVersion)
	}
	return remediationData, nil
}

// GetNextAction contains the logic to determine, what will be the next remediation action according to the remediation.yaml file
// It searches for a problem type matching the root cause of the problem. If no problem type is found a problemtype matching the problem title
// will be searched as a fallback. If still no problem type is found it will return an error.

// The actionIndex parameter specifies which action to take if a problem type was found.
func GetNextAction(remediation *v0_1_4.Remediation, problemDetails v0_2_0.ProblemDetails, actionIndex int) (*v0_2_0.ActionInfo, error) {
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

	// we did not find an action
	if actions == nil {
		return nil, fmt.Errorf("unable to find actions for root cause %s", rootCause)
	}

	// the required action does not exist
	if actionIndex >= len(actions) {
		return nil, fmt.Errorf("failed to get action for root cause %s. There is no action with index %d", rootCause, actionIndex)
	}

	action := actions[actionIndex]
	return &v0_2_0.ActionInfo{
		Name:        action.Name,
		Action:      action.Action,
		Description: action.Description,
		Value:       action.Value,
	}, nil

}
