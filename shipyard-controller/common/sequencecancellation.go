package common

import "github.com/keptn/keptn/shipyard-controller/models"

type SequenceCancellationReason int

const (
	Cancelled SequenceCancellationReason = iota
	Timeout
)

type SequenceCancellation struct {
	KeptnContext string
	Reason       SequenceCancellationReason
	LastEvent    models.Event
}
