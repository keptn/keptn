// Code generated by go-swagger; DO NOT EDIT.

package service_resource

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// NewGetProjectProjectNameStageStageNameServiceServiceNameResourceParams creates a new GetProjectProjectNameStageStageNameServiceServiceNameResourceParams object
// with the default values initialized.
func NewGetProjectProjectNameStageStageNameServiceServiceNameResourceParams() GetProjectProjectNameStageStageNameServiceServiceNameResourceParams {

	var (
		// initialize parameters with default values

		disableUpstreamSyncDefault = bool(false)

		pageSizeDefault = int64(20)
	)

	return GetProjectProjectNameStageStageNameServiceServiceNameResourceParams{
		DisableUpstreamSync: &disableUpstreamSyncDefault,

		PageSize: &pageSizeDefault,
	}
}

// GetProjectProjectNameStageStageNameServiceServiceNameResourceParams contains all the bound params for the get project project name stage stage name service service name resource operation
// typically these are obtained from a http.Request
//
// swagger:parameters GetProjectProjectNameStageStageNameServiceServiceNameResource
type GetProjectProjectNameStageStageNameServiceServiceNameResourceParams struct {

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
// To ensure default values, the struct must have been initialized with NewGetProjectProjectNameStageStageNameServiceServiceNameResourceParams() beforehand.
func (o *GetProjectProjectNameStageStageNameServiceServiceNameResourceParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	qs := runtime.Values(r.URL.Query())

	qDisableUpstreamSync, qhkDisableUpstreamSync, _ := qs.GetOK("disableUpstreamSync")
	if err := o.bindDisableUpstreamSync(qDisableUpstreamSync, qhkDisableUpstreamSync, route.Formats); err != nil {
		res = append(res, err)
	}

	qNextPageKey, qhkNextPageKey, _ := qs.GetOK("nextPageKey")
	if err := o.bindNextPageKey(qNextPageKey, qhkNextPageKey, route.Formats); err != nil {
		res = append(res, err)
	}

	qPageSize, qhkPageSize, _ := qs.GetOK("pageSize")
	if err := o.bindPageSize(qPageSize, qhkPageSize, route.Formats); err != nil {
		res = append(res, err)
	}

	rProjectName, rhkProjectName, _ := route.Params.GetOK("projectName")
	if err := o.bindProjectName(rProjectName, rhkProjectName, route.Formats); err != nil {
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

// bindDisableUpstreamSync binds and validates parameter DisableUpstreamSync from query.
func (o *GetProjectProjectNameStageStageNameServiceServiceNameResourceParams) bindDisableUpstreamSync(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false

	if raw == "" { // empty values pass all other validations
		// Default values have been previously initialized by NewGetProjectProjectNameStageStageNameServiceServiceNameResourceParams()
		return nil
	}

	value, err := swag.ConvertBool(raw)
	if err != nil {
		return errors.InvalidType("disableUpstreamSync", "query", "bool", raw)
	}
	o.DisableUpstreamSync = &value

	return nil
}

// bindNextPageKey binds and validates parameter NextPageKey from query.
func (o *GetProjectProjectNameStageStageNameServiceServiceNameResourceParams) bindNextPageKey(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false

	if raw == "" { // empty values pass all other validations
		return nil
	}
	o.NextPageKey = &raw

	return nil
}

// bindPageSize binds and validates parameter PageSize from query.
func (o *GetProjectProjectNameStageStageNameServiceServiceNameResourceParams) bindPageSize(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false

	if raw == "" { // empty values pass all other validations
		// Default values have been previously initialized by NewGetProjectProjectNameStageStageNameServiceServiceNameResourceParams()
		return nil
	}

	value, err := swag.ConvertInt64(raw)
	if err != nil {
		return errors.InvalidType("pageSize", "query", "int64", raw)
	}
	o.PageSize = &value

	if err := o.validatePageSize(formats); err != nil {
		return err
	}

	return nil
}

// validatePageSize carries on validations for parameter PageSize
func (o *GetProjectProjectNameStageStageNameServiceServiceNameResourceParams) validatePageSize(formats strfmt.Registry) error {

	if err := validate.MinimumInt("pageSize", "query", *o.PageSize, 1, false); err != nil {
		return err
	}

	if err := validate.MaximumInt("pageSize", "query", *o.PageSize, 50, false); err != nil {
		return err
	}

	return nil
}

// bindProjectName binds and validates parameter ProjectName from path.
func (o *GetProjectProjectNameStageStageNameServiceServiceNameResourceParams) bindProjectName(rawData []string, hasKey bool, formats strfmt.Registry) error {
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
func (o *GetProjectProjectNameStageStageNameServiceServiceNameResourceParams) bindServiceName(rawData []string, hasKey bool, formats strfmt.Registry) error {
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
func (o *GetProjectProjectNameStageStageNameServiceServiceNameResourceParams) bindStageName(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// Parameter is provided by construction from the route
	o.StageName = raw

	return nil
}
