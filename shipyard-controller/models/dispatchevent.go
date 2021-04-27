package models

import (
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"time"
)

// DispatcherEvent defines the type of a event which will be dispatches by the EventDispatcher
// It wraps the event to be dispatched and adds a timestamp for deciding when the EventDispatcher will
// dispatch the event
type DispatcherEvent struct {
	Event     cloudevents.Event
	TimeStamp time.Time
}
