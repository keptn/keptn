// Code generated by go-swagger; DO NOT EDIT.

package project_resource

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"

	strfmt "github.com/go-openapi/strfmt"
)

// NewPutProjectProjectNameResourceParams creates a new PutProjectProjectNameResourceParams object
// no default values defined in spec.
func NewPutProjectProjectNameResourceParams() PutProjectProjectNameResourceParams {

	return PutProjectProjectNameResourceParams{}
}

// PutProjectProjectNameResourceParams contains all the bound params for the put project project name resource operation
// typically these are obtained from a http.Request
//
// swagger:parameters PutProjectProjectNameResource
type PutProjectProjectNameResourceParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*Name of the project
	  Required: true
	  In: path
	*/
	ProjectName string
	/*List of resources
	  In: body
	*/
	Resources PutProjectProjectNameResourceBody
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewPutProjectProjectNameResourceParams() beforehand.
func (o *PutProjectProjectNameResourceParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	rProjectName, rhkProjectName, _ := route.Params.GetOK("projectName")
	if err := o.bindProjectName(rProjectName, rhkProjectName, route.Formats); err != nil {
		res = append(res, err)
	}

	if runtime.HasBody(r) {
		defer r.Body.Close()
		var body PutProjectProjectNameResourceBody
		if err := route.Consumer.Consume(r.Body, &body); err != nil {
			res = append(res, errors.NewParseError("resources", "body", "", err))
		} else {
			// validate body object
			if err := body.Validate(route.Formats); err != nil {
				res = append(res, err)
			}

			if len(res) == 0 {
				o.Resources = body
			}
		}
	}
	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindProjectName binds and validates parameter ProjectName from path.
func (o *PutProjectProjectNameResourceParams) bindProjectName(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// Parameter is provided by construction from the route

	o.ProjectName = raw

	return nil
}
