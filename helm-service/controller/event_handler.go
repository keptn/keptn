package controller

import (
	cloudevents "github.com/cloudevents/sdk-go/v2"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
)

// EventHandler provides an interface for handling CloudEvents
type EventHandler interface {
	HandleEvent(ce cloudevents.Event, closeLogger func(keptnHandler *keptnv2.Keptn))
}
