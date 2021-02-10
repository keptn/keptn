package models

// TaskSequenceEvent godoc
type TaskSequenceEvent struct {
	TaskSequenceName string `json:"taskSequenceName" bson:"taskSequenceName"`
	TriggeredEventID string `json:"triggeredEventID" bson:"triggeredEventID"`
	Stage            string `json:"stage" bson:"stage"`
	KeptnContext     string `json:"keptnContext" bson:"keptnContext"`
}
