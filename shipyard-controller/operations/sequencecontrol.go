package operations

import "github.com/keptn/keptn/shipyard-controller/common"

type SequenceControlCommand struct {
	State common.SequenceControlState
	Stage string
}

type SequenceControlResponse struct {
}
