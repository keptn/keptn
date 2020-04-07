package actions

import (
	"github.com/keptn/go-utils/pkg/lib"
)

const envConfigSvcURL = "CONFIGURATION_SERVICE"

type ActionExecutor interface {
	GetAction() string
	ExecuteAction(problem *keptn.ProblemEventData, keptnHandler *keptn.Keptn, action *keptn.RemediationAction) error
	ResolveAction(problem *keptn.ProblemEventData, keptnHandler *keptn.Keptn, action *keptn.RemediationAction) error
}
