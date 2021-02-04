package db

import "github.com/keptn/keptn/shipyard-controller/models"

// TaskSequenceRepo godoc
//go:generate moq --skip-ensure -pkg db_mock -out ./mock/task_sequence_repo_moq.go . TaskSequenceRepo
type TaskSequenceRepo interface {
	// GetTaskSequence godoc
	GetTaskSequence(project, triggeredID string) (*models.TaskSequenceEvent, error)
	// CreateTaskSequenceMapping godoc
	CreateTaskSequenceMapping(project string, taskSequenceEvent models.TaskSequenceEvent) error
	// DeleteTaskSequenceMapping godoc
	DeleteTaskSequenceMapping(keptnContext, project, stage, taskSequenceName string) error
	// DeleteTaskSequenceCollection godoc
	DeleteTaskSequenceCollection(project string) error
}
