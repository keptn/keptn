package common

import "github.com/keptn/keptn/shipyard-controller/models"

type SequenceCancellationReason int

const (
	// there will be more reasons added later
	Timeout SequenceCancellationReason = iota
)

type SequenceCancellation struct {
	KeptnContext string
	Reason       SequenceCancellationReason
	LastEvent    models.Event
}

type SequencePauseRequest struct {
	KeptnContext string
	Stage        string
	Project      string
}
