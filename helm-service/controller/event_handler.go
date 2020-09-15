package controller

import (
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	cloudevents "github.com/cloudevents/sdk-go/v2"
)

type EventHandler interface {
	Handler
	HandleEvent(ce cloudevents.Event, closeLogger func(keptnHandler *keptnv2.Keptn))
}