package models

// TaskSequenceEvent godoc
type TaskSequenceEvent struct {
	TaskSequenceName string `json:"taskSequenceName"`
	TriggeredEventID string `json:"triggeredEventID"`
	Stage            string `json:"stage"`
	KeptnContext     string `json:"keptnContext"`
}
