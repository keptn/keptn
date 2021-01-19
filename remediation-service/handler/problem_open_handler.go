package handler

import (
	"fmt"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	keptn "github.com/keptn/go-utils/pkg/lib"
)

// ProblemOpenEventHandler handles incoming problem.open events
type ProblemOpenEventHandler struct {
	KeptnHandler *keptnv2.Keptn
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
		eh.KeptnHandler.Logger.Error(err.Error())
		_ = eh.Remediation.sendRemediationFinishedEvent(keptnv2.StatusErrored, keptnv2.ResultFailed, err.Error())
		return err
	}

	if autoRemediate {
		eh.KeptnHandler.Logger.Info(fmt.Sprintf("Remediation enabled for project %s in stage %s", problemEvent.Project, problemEvent.Stage))
	} else {
		msg := fmt.Sprintf("Remediation disabled for service %s in project %s in stage %s", problemEvent.Service, problemEvent.Project, problemEvent.Stage)
		eh.KeptnHandler.Logger.Info(msg)
		_ = eh.Remediation.sendRemediationFinishedEvent(keptnv2.StatusErrored, keptnv2.ResultFailed, msg)
		return nil
	}

	// get remediation.yaml
	resource, err := eh.Remediation.getRemediationFile()
	if err != nil {
		eh.KeptnHandler.Logger.Info(err.Error())
		_ = eh.Remediation.sendRemediationFinishedEvent(keptnv2.StatusErrored, keptnv2.ResultFailed, err.Error())
		return err
	}

	// get remediation action from remediation.yaml
	remediationData, err := eh.Remediation.getRemediation(resource)
	if err != nil {
		eh.KeptnHandler.Logger.Error(err.Error())
		_ = eh.Remediation.sendRemediationFinishedEvent(keptnv2.StatusErrored, keptnv2.ResultFailed, err.Error())
		return err
	}

	err = eh.Remediation.sendRemediationTriggeredEvent(problemEvent)
	if err != nil {
		msg := "could not send remediation.triggered event"
		eh.KeptnHandler.Logger.Error(msg + ": " + err.Error())
		_ = eh.Remediation.sendRemediationFinishedEvent(keptnv2.StatusErrored, keptnv2.ResultFailed, msg)
		return err
	}

	problemType := problemEvent.ProblemTitle

	actionIndex := 0
	action := eh.Remediation.getActionForProblemType(*remediationData, problemType, actionIndex)
	if action == nil {
		action = eh.Remediation.getActionForProblemType(*remediationData, "default", actionIndex)
	}

	if action != nil {
		err = eh.Remediation.triggerAction(action, actionIndex, keptnv2.ProblemDetails{
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
			eh.KeptnHandler.Logger.Error(err.Error())
			_ = eh.Remediation.sendRemediationFinishedEvent(keptnv2.StatusErrored, keptnv2.ResultFailed, err.Error())
			return err
		}
	} else {
		msg := "No remediation configured for problem type " + problemType
		eh.KeptnHandler.Logger.Info(msg)
		return eh.Remediation.sendRemediationFinishedEvent(keptnv2.StatusSucceeded, keptnv2.ResultPass, "triggered all actions")
	}

	return nil
}

func (eh *ProblemOpenEventHandler) isRemediationEnabled() (bool, error) {
	remediationFile, err := eh.Remediation.getRemediationFile()
	if err != nil {
		if err == errNoRemediationFileAvailable {
			return false, nil
		}
		return false, fmt.Errorf("Failed to check if remediation is enabled: %s", err.Error())
	} else if remediationFile == nil {
		return false, nil
	}
	eh.KeptnHandler.Logger.Debug("remediation.yaml for service found")
	return true, nil
}
