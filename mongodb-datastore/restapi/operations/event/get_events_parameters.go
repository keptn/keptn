// Code generated by go-swagger; DO NOT EDIT.

package event

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"

	strfmt "github.com/go-openapi/strfmt"
)

// NewGetEventsParams creates a new GetEventsParams object
// with the default values initialized.
func NewGetEventsParams() GetEventsParams {

	var (
		// initialize parameters with default values

		pageSizeDefault = int64(20)
	)

	return GetEventsParams{
		PageSize: &pageSizeDefault,
	}
}

// GetEventsParams contains all the bound params for the get events operation
// typically these are obtained from a http.Request
//
// swagger:parameters getEvents
type GetEventsParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*From time to fetch keptn cloud events
	  In: query
	*/
	FromTime *string
	/*keptnContext of the events to get
	  In: query
	*/
	KeptnContext *string
	/*Key of the page to be returned
	  In: query
	*/
	NextPageKey *string
	/*Page size to be returned
	  Maximum: 100
	  Minimum: 1
	  In: query
	  Default: 20
	*/
	PageSize *int64
	/*Name of the project
	  In: query
	*/
	Project *string
	/*Set to load only root events
	  In: query
	*/
	Root *string
	/*Name of the service
	  In: query
	*/
	Service *string
	/*Name of the event source
	  In: query
	*/
	Source *string
	/*Name of the stage
	  In: query
	*/
	Stage *string
	/*Type of the keptn cloud event
	  In: query
	*/
	Type *string
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewGetEventsParams() beforehand.
func (o *GetEventsParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	qs := runtime.Values(r.URL.Query())

	qFromTime, qhkFromTime, _ := qs.GetOK("fromTime")
	if err := o.bindFromTime(qFromTime, qhkFromTime, route.Formats); err != nil {
		res = append(res, err)
	}

	qKeptnContext, qhkKeptnContext, _ := qs.GetOK("keptnContext")
	if err := o.bindKeptnContext(qKeptnContext, qhkKeptnContext, route.Formats); err != nil {
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

	qProject, qhkProject, _ := qs.GetOK("project")
	if err := o.bindProject(qProject, qhkProject, route.Formats); err != nil {
		res = append(res, err)
	}

	qRoot, qhkRoot, _ := qs.GetOK("root")
	if err := o.bindRoot(qRoot, qhkRoot, route.Formats); err != nil {
		res = append(res, err)
	}

	qService, qhkService, _ := qs.GetOK("service")
	if err := o.bindService(qService, qhkService, route.Formats); err != nil {
		res = append(res, err)
	}

	qSource, qhkSource, _ := qs.GetOK("source")
	if err := o.bindSource(qSource, qhkSource, route.Formats); err != nil {
		res = append(res, err)
	}

	qStage, qhkStage, _ := qs.GetOK("stage")
	if err := o.bindStage(qStage, qhkStage, route.Formats); err != nil {
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

// bindFromTime binds and validates parameter FromTime from query.
func (o *GetEventsParams) bindFromTime(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false
	if raw == "" { // empty values pass all other validations
		return nil
	}

	o.FromTime = &raw

	return nil
}

// bindKeptnContext binds and validates parameter KeptnContext from query.
func (o *GetEventsParams) bindKeptnContext(rawData []string, hasKey bool, formats strfmt.Registry) error {
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

// bindNextPageKey binds and validates parameter NextPageKey from query.
func (o *GetEventsParams) bindNextPageKey(rawData []string, hasKey bool, formats strfmt.Registry) error {
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
func (o *GetEventsParams) bindPageSize(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false
	if raw == "" { // empty values pass all other validations
		// Default values have been previously initialized by NewGetEventsParams()
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
func (o *GetEventsParams) validatePageSize(formats strfmt.Registry) error {

	if err := validate.MinimumInt("pageSize", "query", int64(*o.PageSize), 1, false); err != nil {
		return err
	}

	if err := validate.MaximumInt("pageSize", "query", int64(*o.PageSize), 100, false); err != nil {
		return err
	}

	return nil
}

// bindProject binds and validates parameter Project from query.
func (o *GetEventsParams) bindProject(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false
	if raw == "" { // empty values pass all other validations
		return nil
	}

	o.Project = &raw

	return nil
}

// bindRoot binds and validates parameter Root from query.
func (o *GetEventsParams) bindRoot(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false
	if raw == "" { // empty values pass all other validations
		return nil
	}

	o.Root = &raw

	return nil
}

// bindService binds and validates parameter Service from query.
func (o *GetEventsParams) bindService(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false
	if raw == "" { // empty values pass all other validations
		return nil
	}

	o.Service = &raw

	return nil
}

// bindSource binds and validates parameter Source from query.
func (o *GetEventsParams) bindSource(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false
	if raw == "" { // empty values pass all other validations
		return nil
	}

	o.Source = &raw

	return nil
}

// bindStage binds and validates parameter Stage from query.
func (o *GetEventsParams) bindStage(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: false
	// AllowEmptyValue: false
	if raw == "" { // empty values pass all other validations
		return nil
	}

	o.Stage = &raw

	return nil
}

// bindType binds and validates parameter Type from query.
func (o *GetEventsParams) bindType(rawData []string, hasKey bool, formats strfmt.Registry) error {
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
