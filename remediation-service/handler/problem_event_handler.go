package handler

import (
	cloudevents "github.com/cloudevents/sdk-go"
	keptn "github.com/keptn/go-utils/pkg/lib"
)

// ProblemOpenEventHandler handles incoming problem.open events
type ProblemEventHandler struct {
	KeptnHandler *keptn.Keptn
	Logger       keptn.LoggerInterface
	Event        cloudevents.Event
	Remediation  *Remediation
}

// HandleEvent handles the event
func (eh *ProblemEventHandler) HandleEvent() error {
	var problemEvent *keptn.ProblemEventData
	if eh.Event.Type() == keptn.ProblemEventType {
		eh.Logger.Debug("Received problem notification")
		problemEvent = &keptn.ProblemEventData{}
		if err := eh.Event.DataAs(problemEvent); err != nil {
			return err
		}
	}

	eh.Logger.Debug("Received problem event with state " + problemEvent.State)

	// this service should only react to events with STATE=CLOSED. Opened problems are handled by the ProblemOpenEventHandler
	if problemEvent.State == "CLOSED" {
		msg := "Problem " + problemEvent.PID + " of type " + problemEvent.ProblemTitle + " has been closed."
		eh.Logger.Info(msg)
		_ = eh.Remediation.sendRemediationFinishedEvent(keptn.RemediationStatusSucceeded, keptn.RemediationResultFailed, msg)
		err := deleteRemediation(eh.KeptnHandler.KeptnContext, *eh.KeptnHandler.KeptnBase)
		if err != nil {
			eh.Logger.Error("Could not close remediation: " + err.Error())
			return err
		}
		return nil
	}
	return nil
}
