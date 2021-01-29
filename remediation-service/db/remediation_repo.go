package db

import "github.com/keptn/keptn/remediation-service/models"

// IRemediationRepo godoc
type IRemediationRepo interface {
	GetRemediations(keptnContext, project string) ([]*models.Remediation, error)
	CreateRemediation(project string, remediation *models.Remediation) error
	DeleteRemediation(keptnContext, project string) error
}
