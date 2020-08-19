package operations

// NewGetTriggeredEventsParams creates a new GetTriggeredEventsParams object
// with the default values initialized.
func NewGetTriggeredEventsParams() GetTriggeredEventsParams {

	var (
		// initialize parameters with default values

		pageSizeDefault = int64(20)
	)

	return GetTriggeredEventsParams{
		PageSize: &pageSizeDefault,
	}
}

// GetTriggeredEventsParams contains all the bound params for the get triggered events operation
// typically these are obtained from a http.Request
//
// swagger:parameters get triggered events
type GetTriggeredEventsParams struct {

	/*Stage name
	  In: query
	*/
	EventID *string `form:"eventID" json:"eventID"`
	/*Event type
	  Required: true
	  In: path
	*/
	EventType string `form:"eventType" json:"eventType"`
	/*Pointer to the next set of items
	  In: query
	*/
	NextPageKey *string `form:"nextPageKey" json:"nextPageKey"`
	/*The number of items to return
	  Maximum: 50
	  Minimum: 1
	  In: query
	  Default: 20
	*/
	PageSize *int64 `form:"pageSize" json:"pageSize"`
	/*Project name
	  In: query
	*/
	Project *string `form:"project" json:"project"`
	/*Service name
	  In: query
	*/
	Service *string `form:"service" json:"service"`
	/*Stage name
	  In: query
	*/
	Stage *string `form:"stage" json:"stage"`
}
