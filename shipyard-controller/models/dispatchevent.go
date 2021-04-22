package models

import (
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"time"
)

type DispatcherEvent struct {
	Event     cloudevents.Event
	TimeStamp time.Time
}
