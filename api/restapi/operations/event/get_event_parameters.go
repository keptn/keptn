// Code generated by go-swagger; DO NOT EDIT.

package event

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"

	strfmt "github.com/go-openapi/strfmt"
)

// NewGetEventParams creates a new GetEventParams object
// no default values defined in spec.
func NewGetEventParams() GetEventParams {

	return GetEventParams{}
}

// GetEventParams contains all the bound params for the get event operation
// typically these are obtained from a http.Request
//
// swagger:parameters GetEvent
type GetEventParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*KeptnContext of the events to get
	  In: query
	*/
	KeptnContext *string
	/*Type of the Keptn cloud event
	  In: query
	*/
	Type *string
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewGetEventParams() beforehand.
func (o *GetEventParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	qs := runtime.Values(r.URL.Query())

	qKeptnContext, qhkKeptnContext, _ := qs.GetOK("keptnContext")
	if err := o.bindKeptnContext(qKeptnContext, qhkKeptnContext, route.Formats); err != nil {
		res = append(res, err)
	}

	qType, qhkType, _ := qs.GetOK("type")
	if err := o.bindType(qType, qhkType, route.Formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindKeptnContext binds and validates parameter KeptnContext from query.
func (o *GetEventParams) bindKeptnContext(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false
	if raw == "" { // empty values pass all other validations
		return nil
	}

	o.KeptnContext = &raw

	return nil
}

// bindType binds and validates parameter Type from query.
func (o *GetEventParams) bindType(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false
	if raw == "" { // empty values pass all other validations
		return nil
	}

	o.Type = &raw

	return nil
}
