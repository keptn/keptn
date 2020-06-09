package handler

import (
	"fmt"
	cloudevents "github.com/cloudevents/sdk-go"
	keptn "github.com/keptn/go-utils/pkg/lib"
	"time"
)

const waitTimeInMinutes = 10

type ActionFinishedEventHandler struct {
	KeptnHandler *keptn.Keptn
	Logger       keptn.LoggerInterface
	Event        cloudevents.Event
	Remediation  *Remediation
}

func (eh *ActionFinishedEventHandler) HandleEvent() error {
	actionFinishedEvent := &keptn.ActionFinishedEventData{}

	err := eh.Event.DataAs(actionFinishedEvent)
	if err != nil {
		eh.Logger.Error("Could not parse incoming action.finished event: " + err.Error())
		return err
	}
	eh.Logger.Info(fmt.Sprintf("Received action.finished event for remediationStatus action. result = %v", actionFinishedEvent.Action.Result))
	eh.Logger.Info(fmt.Sprintf("Waiting for %d minutes for action to take effect", waitTimeInMinutes))
	<-time.After(waitTimeInMinutes * time.Minute)
	eh.Logger.Info("Wait time is over. Sending start-evaluation event.")

	err = eh.Remediation.sendStartEvaluationEvent()
	if err != nil {
		eh.Logger.Error("Could not send start-evaluation event: " + err.Error())
		eh.Remediation.sendRemediationFinishedEvent(keptn.RemediationStatusErrored, keptn.RemediationResultFailed, "could not send start-evaluation event")
		return err
	}
	return nil
}
