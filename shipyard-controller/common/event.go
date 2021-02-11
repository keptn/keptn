package common

// EventStatus indicates the status type of an event, i.e. 'triggered', 'started', or 'finished'
type EventStatus string

const (
	// TriggeredEvent describes a 'triggered' event
	TriggeredEvent EventStatus = "triggered"
	// StartedEvent describes a 'started' event
	StartedEvent EventStatus = "started"
	// FinishedEvent describes a 'finished' event
	FinishedEvent EventStatus = "finished"
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
