package common

import "github.com/keptn/keptn/shipyard-controller/models"

type SequenceTimeout struct {
	KeptnContext string
	LastEvent    models.Event
}

type SequenceControlState string

const (
	PauseSequence  SequenceControlState = "pause"
	ResumeSequence SequenceControlState = "resume"
	AbortSequence  SequenceControlState = "abort"
)

type SequenceControl struct {
	State        SequenceControlState
	KeptnContext string
	Stage        string
	Project      string
}
