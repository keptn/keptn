package operations

import "github.com/keptn/keptn/shipyard-controller/common"

type SequenceControlCommand struct {
	State common.SequenceControlState `json:"state" binding:"required"`
	Stage string                      `json:"stage" binding:"required"`
}

type SequenceControlResponse struct {
}
