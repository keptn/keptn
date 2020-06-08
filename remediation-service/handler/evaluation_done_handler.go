package handler

import (
	cloudevents "github.com/cloudevents/sdk-go"
	keptn "github.com/keptn/go-utils/pkg/lib"
)

type EvaluationDoneEventHandler struct {
	KeptnHandler   *keptn.Keptn
	Logger         keptn.LoggerInterface
	Event          cloudevents.Event
	RemediationLog *RemediationLog
}

func (eh *EvaluationDoneEventHandler) HandleEvent() error {

}
