package operations

import (
	models "github.com/keptn/keptn/shipyard-controller/models"
)

// NewHandleEventParams creates a new HandleEventParams object
// no default values defined in spec.
func NewHandleEventParams() HandleEventParams {

	return HandleEventParams{}
}

// HandleEventParams contains all the bound params for the handle event operation
// typically these are obtained from a http.Request
//
// swagger:parameters handle event
type HandleEventParams struct {
	models.Event
}
