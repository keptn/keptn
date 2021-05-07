package db

import "github.com/keptn/keptn/shipyard-controller/models"

//go:generate moq --skip-ensure -pkg db_mock -out ./mock/sequencestaterepo_mock.go . SequenceStateRepo
type SequenceStateRepo interface {
	CreateSequenceState(state models.SequenceState) error
	FindSequenceStates(filter models.StateFilter) (*models.SequenceStates, error)
	UpdateSequenceState(state models.SequenceState) error
	DeleteSequenceStates(filter models.StateFilter) error
}
