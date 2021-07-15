package models

import (
	"time"
)

// QueueItem is a type used to persist events that are queued for dispatching
type QueueItem struct {
	Scope     EventScope `json:"scope" bson:"scope"`
	EventID   string     `json:"eventID" bson:"eventID"`
	Timestamp time.Time  `json:"timestamp" bson:"timestamp"`
}
