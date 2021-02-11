package fake

import "github.com/keptn/keptn/remediation-service/models"

type RemediationRepo struct {
	// Remediations describes the current state of the remediation repo (incl. remediations that are assumed to be available when a unit test is started)
	Remediations []*models.Remediation
	// ReceivedRemediations acts as a kind of recorder to keep track of what remediations have been added/deleted during a test
	ReceivedRemediations []*models.Remediation
}

func (r *RemediationRepo) GetRemediations(keptnContext, project string) ([]*models.Remediation, error) {
	result := []*models.Remediation{}

	for _, remediation := range r.Remediations {
		if remediation.KeptnContext == keptnContext {
			result = append(result, remediation)
		}
	}
	return result, nil
}

func (r *RemediationRepo) CreateRemediation(project string, remediation *models.Remediation) error {
	if r.Remediations == nil {
		r.Remediations = []*models.Remediation{}
	}
	if r.ReceivedRemediations == nil {
		r.ReceivedRemediations = []*models.Remediation{}
	}
	r.Remediations = append(r.Remediations, remediation)
	r.ReceivedRemediations = append(r.ReceivedRemediations, remediation)
	return nil
}

func (r *RemediationRepo) DeleteRemediation(keptnContext, project string) error {
	newRemediations := []*models.Remediation{}
	for _, remediation := range r.Remediations {
		if remediation.KeptnContext != keptnContext {
			newRemediations = append(newRemediations, remediation)
		}
	}
	r.Remediations = newRemediations

	newReceivedRemediations := []*models.Remediation{}
	for _, remediation := range r.ReceivedRemediations {
		if remediation.KeptnContext != keptnContext {
			newReceivedRemediations = append(newReceivedRemediations, remediation)
		}
	}
	r.ReceivedRemediations = newReceivedRemediations
	return nil
}

func (r *RemediationRepo) GetReceivedRemediations() []*models.Remediation {
	if r.ReceivedRemediations == nil {
		r.ReceivedRemediations = []*models.Remediation{}
	}
	return r.ReceivedRemediations
}
