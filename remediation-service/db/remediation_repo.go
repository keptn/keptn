package db

import "github.com/keptn/keptn/remediation-service/models"

// RemediationRepo godoc
type RemediationRepo interface {
	// GetTaskSequence godoc
	GetRemediation(project, keptnContext string) (*models.Remediation, error)
	// CreateTaskSequenceMapping godoc
	CreateRemediation(project string, remediation *models.Remediation) error
	// DeleteTaskSequenceMapping godoc
	DeleteRemediation(keptnContext, project string) error
	// DeleteTaskSequenceCollection godoc
	DeleteTaskSequenceCollection(project string) error
}
