package handler

import (
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/go-utils/pkg/lib/v0_1_4"
	"github.com/keptn/go-utils/pkg/lib/v0_2_0"
)

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

func GetNextAction(remediation *v0_1_4.Remediation, problemDetails v0_2_0.ProblemDetails, actionIndex int) (*v0_2_0.ActionInfo, error) {
	rootCause := problemDetails.RootCause
	problemTitle := problemDetails.ProblemTitle

	var actions []v0_1_4.RemediationActionsOnOpen
	for _, r := range remediation.Spec.Remediations {
		if r.ProblemType == rootCause || r.ProblemType == problemTitle {
			actions = r.ActionsOnOpen
			break
		}
	}

	if actions == nil {
		return nil, fmt.Errorf("unable to find actions for root cause %s", rootCause)
	}

	if actionIndex >= len(actions) {
		return nil, fmt.Errorf("failed to get action for root cause %s. There is no action with index %d", rootCause, actionIndex)

	}

	action := actions[actionIndex]
	return &v0_2_0.ActionInfo{
		Name:        action.Name,
		Action:      "",
		Description: action.Description,
		Value:       action.Value,
	}, nil

}
