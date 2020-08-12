package models

// TaskSequence godoc
type TaskSequenceEvent struct {
	TaskSequenceName string `json:"name"`
	TriggeredEventID string `json:"triggeredEventID"`
}
