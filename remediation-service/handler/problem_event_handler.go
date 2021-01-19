package handler

import (
	cloudevents "github.com/cloudevents/sdk-go"
	keptn "github.com/keptn/go-utils/pkg/lib"
	keptncommon "github.com/keptn/go-utils/pkg/lib/keptn"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
)

// ProblemEventHandler handles incoming problem events
type ProblemEventHandler struct {
	KeptnHandler *keptnv2.Keptn
	Logger       keptncommon.LoggerInterface
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
		return eh.Remediation.sendRemediationFinishedEvent(keptnv2.StatusSucceeded, keptnv2.ResultPass, msg)
	}
	return nil
}
