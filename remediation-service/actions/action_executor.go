package actions

import (
	keptnevents "github.com/keptn/go-utils/pkg/events"
	keptnmodels "github.com/keptn/go-utils/pkg/models"
	"sigs.k8s.io/yaml"
)

const envConfigSvcURL = "CONFIGURATION_SERVICE"

type ActionExecutor interface {
	GetAction() string
	ExecuteAction(problem *keptnevents.ProblemEventData, shkeptncontext string, action *keptnmodels.RemediationAction) error
}

type ProblemDetails struct {
	ClientIP string `json:"ClientIP"`
}

func GetIP(problem *keptnevents.ProblemEventData) (string, error) {

	details := ProblemDetails{}
	err := yaml.Unmarshal(problem.ProblemDetails, &details)
	if err != nil {
		return "", err
	}

	return details.ClientIP, nil
}
