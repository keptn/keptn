package db

import "github.com/keptn/keptn/shipyard-controller/models"

//go:generate moq --skip-ensure -pkg db_mock -out ./mock/staterepo_mock.go . StateRepo
type StateRepo interface {
	CreateState(state models.SequenceState) error
	FindStates(filter models.StateFilter) (*models.SequenceStates, error)
	UpdateState(state models.SequenceState) error
	DeleteStates(filter models.StateFilter) error
}
