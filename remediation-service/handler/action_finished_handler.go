package handler

import (
	"fmt"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"os"
	"time"
)

const waitTimeInMinutes = 10

// ActionFinishedEventHandler handles action.finished events
type ActionFinishedEventHandler struct {
	KeptnHandler *keptnv2.Keptn
	Event        cloudevents.Event
	Remediation  *Remediation
	WaitFunction waitFunction
}

type waitFunction func()

// HandleEvent handles the incoming event
func (eh *ActionFinishedEventHandler) HandleEvent() error {
	actionFinishedEvent := &keptnv2.ActionFinishedEventData{}

	err := eh.Event.DataAs(actionFinishedEvent)
	if err != nil {
		eh.KeptnHandler.Logger.Error("Could not parse incoming action.finished event: " + err.Error())
		return err
	}
	eh.KeptnHandler.Logger.Info(fmt.Sprintf("Received action.finished event for remediationStatus action. result = %v, status = %v", actionFinishedEvent.Result, actionFinishedEvent.Status))

	if eh.WaitFunction == nil {
		eh.WaitFunction = func() {

			waitTime := getWaitTime()
			eh.KeptnHandler.Logger.Info(fmt.Sprintf("Waiting for %s for action to take effect", waitTime.String()))
			<-time.After(waitTime)
		}
	}
	eh.WaitFunction()
	eh.KeptnHandler.Logger.Info("Wait time is over. Sending start-evaluation event.")

	err = eh.Remediation.sendEvaluationTriggeredEvent()
	if err != nil {
		eh.KeptnHandler.Logger.Error("Could not send start-evaluation event: " + err.Error())
		eh.Remediation.sendRemediationFinishedEvent(keptnv2.StatusErrored, keptnv2.ResultFailed, "could not send start-evaluation event")
		return err
	}
	return nil
}

func getWaitTime() time.Duration {
	waitTime, err := time.ParseDuration(os.Getenv("WAIT_TIME_MINUTES"))
	if err != nil {
		waitTime = waitTimeInMinutes * time.Minute
	}
	return waitTime
}
