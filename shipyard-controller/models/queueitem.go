package models

import "time"

type QueueItem struct {
	Scope     EventScope `json:"scope" bson:"scope"`
	EventID   string     `json:"eventID" bson:"eventID"`
	Timestamp time.Time  `json:"timestamp" bson:"timestamp"`
}
