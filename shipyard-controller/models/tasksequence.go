package models

import keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"

// TaskExecution godoc
type TaskExecution struct {
	TaskSequenceName string `json:"taskSequenceName" bson:"taskSequenceName"`
	TriggeredEventID string `json:"triggeredEventID" bson:"triggeredEventID"`
	Task             Task   `json:"task" bson:"task"`
	Stage            string `json:"stage" bson:"stage"`
	Service          string `json:"service" bson:"service"`
	KeptnContext     string `json:"keptnContext" bson:"keptnContext"`
}

type Task struct {
	keptnv2.Task
	TaskIndex int
}
