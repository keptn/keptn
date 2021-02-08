package operations

import "net/http"

type GetStagesParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*Disable sync of upstream repo before reading content
	  In: query
	  Default: false
	*/
	DisableUpstreamSync *bool
	/*Pointer to the next set of items
	  In: query
	*/
	NextPageKey *string
	/*The number of items to return
	  Maximum: 50
	  Minimum: 1
	  In: query
	  Default: 20
	*/
	PageSize *int64
	/*Name of the project
	  Required: true
	  In: path
	*/
	ProjectName string
}

type GetStageParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*Disable sync of upstream repo before reading content
	  In: query
	  Default: false
	*/
	DisableUpstreamSync *bool
	/*Name of the project
	  Required: true
	  In: path
	*/
	ProjectName string
	/*Name of the stage
	  Required: true
	  In: path
	*/
	StageName string
}
