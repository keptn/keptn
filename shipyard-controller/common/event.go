package common

import (
	"github.com/cloudevents/sdk-go/v2/event"
)

// EventStatus indicates the status type of an event, i.e. 'triggered', 'started', or 'finished'
type EventStatus string

type EventSender interface {
	SendEvent(eventMoqParam event.Event) error
}

const (
	// TriggeredEvent describes a 'triggered' event
	TriggeredEvent EventStatus = "triggered"
	// StartedEvent describes a 'started' event
	StartedEvent EventStatus = "started"
	// FinishedEvent describes a 'finished' event
	FinishedEvent EventStatus = "finished"
	// RootEvent indicates that an event triggered a task sequence execution
	RootEvent EventStatus = "root"
)

// EventFilter allows to pass filters
type EventFilter struct {
	Type         string
	Stage        *string
	Service      *string
	ID           *string
	TriggeredID  *string
	Source       *string
	KeptnContext *string
}
