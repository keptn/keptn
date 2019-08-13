// Code generated by go-swagger; DO NOT EDIT.

package service_resource

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"

	strfmt "github.com/go-openapi/strfmt"

	models "github.com/keptn/keptn/configuration-service/models"
)

// NewPutProjectProjectNameStageStageNameServiceServiceNameResourceResourceURIParams creates a new PutProjectProjectNameStageStageNameServiceServiceNameResourceResourceURIParams object
// no default values defined in spec.
func NewPutProjectProjectNameStageStageNameServiceServiceNameResourceResourceURIParams() PutProjectProjectNameStageStageNameServiceServiceNameResourceResourceURIParams {

	return PutProjectProjectNameStageStageNameServiceServiceNameResourceResourceURIParams{}
}

// PutProjectProjectNameStageStageNameServiceServiceNameResourceResourceURIParams contains all the bound params for the put project project name stage stage name service service name resource resource URI operation
// typically these are obtained from a http.Request
//
// swagger:parameters PutProjectProjectNameStageStageNameServiceServiceNameResourceResourceURI
type PutProjectProjectNameStageStageNameServiceServiceNameResourceResourceURIParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*Name of the project
	  Required: true
	  In: path
	*/
	ProjectName string
	/*Resource
	  In: body
	*/
	Resource *models.Resource
	/*Resource URI
	  Required: true
	  In: path
	*/
	ResourceURI string
	/*Name of the service
	  Required: true
	  In: path
	*/
	ServiceName string
	/*Name of the stage
	  Required: true
	  In: path
	*/
	StageName string
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewPutProjectProjectNameStageStageNameServiceServiceNameResourceResourceURIParams() beforehand.
func (o *PutProjectProjectNameStageStageNameServiceServiceNameResourceResourceURIParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	rProjectName, rhkProjectName, _ := route.Params.GetOK("projectName")
	if err := o.bindProjectName(rProjectName, rhkProjectName, route.Formats); err != nil {
		res = append(res, err)
	}

	if runtime.HasBody(r) {
		defer r.Body.Close()
		var body models.Resource
		if err := route.Consumer.Consume(r.Body, &body); err != nil {
			res = append(res, errors.NewParseError("resource", "body", "", err))
		} else {
			// validate body object
			if err := body.Validate(route.Formats); err != nil {
				res = append(res, err)
			}

			if len(res) == 0 {
				o.Resource = &body
			}
		}
	}
	rResourceURI, rhkResourceURI, _ := route.Params.GetOK("resourceURI")
	if err := o.bindResourceURI(rResourceURI, rhkResourceURI, route.Formats); err != nil {
		res = append(res, err)
	}

	rServiceName, rhkServiceName, _ := route.Params.GetOK("serviceName")
	if err := o.bindServiceName(rServiceName, rhkServiceName, route.Formats); err != nil {
		res = append(res, err)
	}

	rStageName, rhkStageName, _ := route.Params.GetOK("stageName")
	if err := o.bindStageName(rStageName, rhkStageName, route.Formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindProjectName binds and validates parameter ProjectName from path.
func (o *PutProjectProjectNameStageStageNameServiceServiceNameResourceResourceURIParams) bindProjectName(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// Parameter is provided by construction from the route

	o.ProjectName = raw

	return nil
}

// bindResourceURI binds and validates parameter ResourceURI from path.
func (o *PutProjectProjectNameStageStageNameServiceServiceNameResourceResourceURIParams) bindResourceURI(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// Parameter is provided by construction from the route

	o.ResourceURI = raw

	return nil
}

// bindServiceName binds and validates parameter ServiceName from path.
func (o *PutProjectProjectNameStageStageNameServiceServiceNameResourceResourceURIParams) bindServiceName(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// Parameter is provided by construction from the route

	o.ServiceName = raw

	return nil
}

// bindStageName binds and validates parameter StageName from path.
func (o *PutProjectProjectNameStageStageNameServiceServiceNameResourceResourceURIParams) bindStageName(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// Parameter is provided by construction from the route

	o.StageName = raw

	return nil
}
