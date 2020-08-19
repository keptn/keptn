package db

import "github.com/keptn/keptn/shipyard-controller/models"

// TaskSequenceRepo godoc
type TaskSequenceRepo interface {
	// GetTaskSequence godoc
	GetTaskSequence(project, triggeredID string) (*models.TaskSequenceEvent, error)
	// CreateTaskSequenceMapping godoc
	CreateTaskSequenceMapping(project string, taskSequenceEvent models.TaskSequenceEvent) error
	// DeleteTaskSequenceMapping godoc
	DeleteTaskSequenceMapping(keptnContext, project, stage, taskSequenceName string) error
}
