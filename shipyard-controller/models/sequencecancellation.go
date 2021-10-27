package models

type SequenceTimeout struct {
	KeptnContext string
	LastEvent    Event
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
