package actions

import (
	keptnevents "github.com/keptn/go-utils/pkg/events"
	keptnmodels "github.com/keptn/go-utils/pkg/models"
)

type ActionExecutor interface {
	GetAction() string
	ExecuteAction(problem *keptnevents.ProblemEventData, shkeptncontext string, action *keptnmodels.RemediationAction) error
}
