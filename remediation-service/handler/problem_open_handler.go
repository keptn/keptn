package handler

import (
	"fmt"
	cloudevents "github.com/cloudevents/sdk-go"
	keptn "github.com/keptn/go-utils/pkg/lib"
)

// ProblemOpenEventHandler handles incoming problem.open events
type ProblemOpenEventHandler struct {
	KeptnHandler *keptn.Keptn
	Event        cloudevents.Event
	Remediation  *Remediation
}

// HandleEvent handles the event
func (eh *ProblemOpenEventHandler) HandleEvent() error {
	var problemEvent *keptn.ProblemEventData
	if eh.Event.Type() == keptn.ProblemOpenEventType {
		eh.KeptnHandler.Logger.Debug("Received problem notification")
		problemEvent = &keptn.ProblemEventData{}
		if err := eh.Event.DataAs(problemEvent); err != nil {
			return err
		}
	}

	eh.KeptnHandler.Logger.Debug("Received problem event with state " + problemEvent.State)

	// check if remediation should be performed
	autoRemediate, err := eh.isRemediationEnabled()
	if err != nil {
		eh.KeptnHandler.Logger.Error(fmt.Sprintf("Failed to check if remediation is enabled: %s", err.Error()))
		return err
	}

	if autoRemediate {
		eh.KeptnHandler.Logger.Info(fmt.Sprintf("Remediation enabled for project %s in stage %s", problemEvent.Project, problemEvent.Stage))
	} else {
		eh.KeptnHandler.Logger.Info(fmt.Sprintf("Remediation disabled for project %s in stage %s", problemEvent.Project, problemEvent.Stage))
		return nil
	}

	// get remediation.yaml
	resource, err := eh.Remediation.getRemediationFile()
	if err != nil {
		return err
	}

	// get remediation action from remediation.yaml
	remediationData, err := eh.Remediation.getRemediation(resource)
	if err != nil {
		return err
	}

	err = eh.Remediation.sendRemediationTriggeredEvent(problemEvent)
	if err != nil {
		msg := "could not send remediation.triggered event"
		eh.KeptnHandler.Logger.Error(msg + ": " + err.Error())
		_ = eh.Remediation.sendRemediationFinishedEvent(keptn.RemediationStatusErrored, keptn.RemediationResultFailed, msg)
		return err
	}

	problemType := problemEvent.ProblemTitle

	actionIndex := 0
	action := eh.Remediation.getActionForProblemType(*remediationData, problemType, actionIndex)
	if action == nil {
		action = eh.Remediation.getActionForProblemType(*remediationData, "*", actionIndex)
	}

	if action != nil {
		err = eh.Remediation.triggerAction(action, actionIndex, keptn.ProblemDetails{
			State:          problemEvent.State,
			ProblemID:      problemEvent.ProblemID,
			ProblemTitle:   problemEvent.ProblemTitle,
			ProblemDetails: problemEvent.ProblemDetails,
			PID:            problemEvent.PID,
			ProblemURL:     problemEvent.ProblemURL,
			ImpactedEntity: problemEvent.ImpactedEntity,
			Tags:           problemEvent.Tags,
		})
		if err != nil {
			return err
		}
	} else {
		msg := "No remediation configured for problem type " + problemType
		eh.KeptnHandler.Logger.Info(msg)
		return eh.Remediation.sendRemediationFinishedEvent(keptn.RemediationStatusSucceeded, keptn.RemediationResultPass, "triggered all actions")
	}

	return nil
}

func (eh *ProblemOpenEventHandler) isRemediationEnabled() (bool, error) {
	shipyard, err := eh.KeptnHandler.GetShipyard()
	if err != nil {
		return false, err
	}
	for _, s := range shipyard.Stages {
		if s.Name == eh.KeptnHandler.KeptnBase.Stage && s.RemediationStrategy == "automated" {
			return true, nil
		}
	}

	return false, nil
}
