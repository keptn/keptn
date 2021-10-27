package models

type SequenceControlCommand struct {
	State SequenceControlState `json:"state" binding:"required"`
	Stage string               `json:"stage"`
}

type SequenceControlResponse struct {
}
