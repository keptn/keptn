// Code generated by go-swagger; DO NOT EDIT.

package service

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/validate"

	"github.com/keptn/keptn/configuration-service/models"
)

// NewPostProjectProjectNameStageStageNameServiceParams creates a new PostProjectProjectNameStageStageNameServiceParams object
//
// There are no default values defined in the spec.
func NewPostProjectProjectNameStageStageNameServiceParams() PostProjectProjectNameStageStageNameServiceParams {

	return PostProjectProjectNameStageStageNameServiceParams{}
}

// PostProjectProjectNameStageStageNameServiceParams contains all the bound params for the post project project name stage stage name service operation
// typically these are obtained from a http.Request
//
// swagger:parameters PostProjectProjectNameStageStageNameService
type PostProjectProjectNameStageStageNameServiceParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*Name of the project
	  Required: true
	  In: path
	*/
	ProjectName string
	/*Service entity
	  In: body
	*/
	Service *models.Service
	/*Name of the stage
	  Required: true
	  In: path
	*/
	StageName string
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewPostProjectProjectNameStageStageNameServiceParams() beforehand.
func (o *PostProjectProjectNameStageStageNameServiceParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	rProjectName, rhkProjectName, _ := route.Params.GetOK("projectName")
	if err := o.bindProjectName(rProjectName, rhkProjectName, route.Formats); err != nil {
		res = append(res, err)
	}

	if runtime.HasBody(r) {
		defer r.Body.Close()
		var body models.Service
		if err := route.Consumer.Consume(r.Body, &body); err != nil {
			res = append(res, errors.NewParseError("service", "body", "", err))
		} else {
			// validate body object
			if err := body.Validate(route.Formats); err != nil {
				res = append(res, err)
			}

			ctx := validate.WithOperationRequest(context.Background())
			if err := body.ContextValidate(ctx, route.Formats); err != nil {
				res = append(res, err)
			}

			if len(res) == 0 {
				o.Service = &body
			}
		}
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
func (o *PostProjectProjectNameStageStageNameServiceParams) bindProjectName(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// Parameter is provided by construction from the route
	o.ProjectName = raw

	return nil
}

// bindStageName binds and validates parameter StageName from path.
func (o *PostProjectProjectNameStageStageNameServiceParams) bindStageName(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// Parameter is provided by construction from the route
	o.StageName = raw

	return nil
}
