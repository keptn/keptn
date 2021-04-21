package handler

import (
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"time"
)

type DispatcherEvent struct {
	event     cloudevents.Event
	timestamp time.Time
}

//go:generate moq -pkg fake -skip-ensure -out ./fake/eventdispatcher.go . IEventDispatcher
type IEventDispatcher interface {
	Add(event DispatcherEvent)
}
