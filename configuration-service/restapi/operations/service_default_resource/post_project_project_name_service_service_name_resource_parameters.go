// Code generated by go-swagger; DO NOT EDIT.

package service_default_resource

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"

	strfmt "github.com/go-openapi/strfmt"
)

// NewPostProjectProjectNameServiceServiceNameResourceParams creates a new PostProjectProjectNameServiceServiceNameResourceParams object
// no default values defined in spec.
func NewPostProjectProjectNameServiceServiceNameResourceParams() PostProjectProjectNameServiceServiceNameResourceParams {

	return PostProjectProjectNameServiceServiceNameResourceParams{}
}

// PostProjectProjectNameServiceServiceNameResourceParams contains all the bound params for the post project project name service service name resource operation
// typically these are obtained from a http.Request
//
// swagger:parameters PostProjectProjectNameServiceServiceNameResource
type PostProjectProjectNameServiceServiceNameResourceParams struct {

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
	Resources PostProjectProjectNameServiceServiceNameResourceBody
	/*Name of the service
	  Required: true
	  In: path
	*/
	ServiceName string
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewPostProjectProjectNameServiceServiceNameResourceParams() beforehand.
func (o *PostProjectProjectNameServiceServiceNameResourceParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	rProjectName, rhkProjectName, _ := route.Params.GetOK("projectName")
	if err := o.bindProjectName(rProjectName, rhkProjectName, route.Formats); err != nil {
		res = append(res, err)
	}

	if runtime.HasBody(r) {
		defer r.Body.Close()
		var body PostProjectProjectNameServiceServiceNameResourceBody
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
	rServiceName, rhkServiceName, _ := route.Params.GetOK("serviceName")
	if err := o.bindServiceName(rServiceName, rhkServiceName, route.Formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindProjectName binds and validates parameter ProjectName from path.
func (o *PostProjectProjectNameServiceServiceNameResourceParams) bindProjectName(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// Parameter is provided by construction from the route

	o.ProjectName = raw

	return nil
}

// bindServiceName binds and validates parameter ServiceName from path.
func (o *PostProjectProjectNameServiceServiceNameResourceParams) bindServiceName(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// Parameter is provided by construction from the route

	o.ServiceName = raw

	return nil
}
