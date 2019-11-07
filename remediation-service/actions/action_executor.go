package actions

import (
	keptnevents "github.com/keptn/go-utils/pkg/events"
	keptnmodels "github.com/keptn/go-utils/pkg/models"
)

const envConfigSvcURL = "CONFIGURATION_SERVICE"

type ActionExecutor interface {
	GetAction() string
	ExecuteAction(problem *keptnevents.ProblemEventData, shkeptncontext string, action *keptnmodels.RemediationAction) error
	ResolveAction(problem *keptnevents.ProblemEventData, shkeptncontext string, action *keptnmodels.RemediationAction) error
}
