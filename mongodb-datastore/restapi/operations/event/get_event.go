// Code generated by go-swagger; DO NOT EDIT.

package event

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	errors "github.com/go-openapi/errors"
	middleware "github.com/go-openapi/runtime/middleware"
	strfmt "github.com/go-openapi/strfmt"
	swag "github.com/go-openapi/swag"
	validate "github.com/go-openapi/validate"
)

// GetEventHandlerFunc turns a function with the right signature into a get event handler
type GetEventHandlerFunc func(GetEventParams) middleware.Responder

// Handle executing the request and returning a response
func (fn GetEventHandlerFunc) Handle(params GetEventParams) middleware.Responder {
	return fn(params)
}

// GetEventHandler interface for that can handle valid get event params
type GetEventHandler interface {
	Handle(GetEventParams) middleware.Responder
}

// NewGetEvent creates a new http.Handler for the get event operation
func NewGetEvent(ctx *middleware.Context, handler GetEventHandler) *GetEvent {
	return &GetEvent{Context: ctx, Handler: handler}
}

/*GetEvent swagger:route GET /events/id/{id} event getEvent

Gets events from the data store

*/
type GetEvent struct {
	Context *middleware.Context
	Handler GetEventHandler
}

func (o *GetEvent) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewGetEventParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}

// GetEventDefaultBody get event default body
// swagger:model GetEventDefaultBody
type GetEventDefaultBody struct {

	// code
	Code int64 `json:"code,omitempty"`

	// fields
	Fields string `json:"fields,omitempty"`

	// message
	// Required: true
	Message *string `json:"message"`
}

// Validate validates this get event default body
func (o *GetEventDefaultBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateMessage(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *GetEventDefaultBody) validateMessage(formats strfmt.Registry) error {

	if err := validate.Required("getEvent default"+"."+"message", "body", o.Message); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (o *GetEventDefaultBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *GetEventDefaultBody) UnmarshalBinary(b []byte) error {
	var res GetEventDefaultBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

// GetEventOKBodyItems0 get event o k body items0
// swagger:model GetEventOKBodyItems0
type GetEventOKBodyItems0 struct {

	// contenttype
	Contenttype string `json:"contenttype,omitempty"`

	// data
	Data interface{} `json:"data,omitempty"`

	// extensions
	Extensions interface{} `json:"extensions,omitempty"`

	// id
	// Required: true
	ID *string `json:"id"`

	// source
	// Required: true
	Source *string `json:"source"`

	// specversion
	// Required: true
	Specversion *string `json:"specversion"`

	// time
	// Format: date-time
	Time strfmt.DateTime `json:"time,omitempty"`

	// type
	// Required: true
	Type *string `json:"type"`

	// shkeptncontext
	Shkeptncontext string `json:"shkeptncontext,omitempty"`
}

// UnmarshalJSON unmarshals this object from a JSON structure
func (o *GetEventOKBodyItems0) UnmarshalJSON(raw []byte) error {
	// AO0
	var dataAO0 struct {
		Contenttype string `json:"contenttype,omitempty"`

		Data interface{} `json:"data,omitempty"`

		Extensions interface{} `json:"extensions,omitempty"`

		ID *string `json:"id"`

		Source *string `json:"source"`

		Specversion *string `json:"specversion"`

		Time strfmt.DateTime `json:"time,omitempty"`

		Type *string `json:"type"`
	}
	if err := swag.ReadJSON(raw, &dataAO0); err != nil {
		return err
	}

	o.Contenttype = dataAO0.Contenttype

	o.Data = dataAO0.Data

	o.Extensions = dataAO0.Extensions

	o.ID = dataAO0.ID

	o.Source = dataAO0.Source

	o.Specversion = dataAO0.Specversion

	o.Time = dataAO0.Time

	o.Type = dataAO0.Type

	// AO1
	var dataAO1 struct {
		Shkeptncontext string `json:"shkeptncontext,omitempty"`
	}
	if err := swag.ReadJSON(raw, &dataAO1); err != nil {
		return err
	}

	o.Shkeptncontext = dataAO1.Shkeptncontext

	return nil
}

// MarshalJSON marshals this object to a JSON structure
func (o GetEventOKBodyItems0) MarshalJSON() ([]byte, error) {
	_parts := make([][]byte, 0, 2)

	var dataAO0 struct {
		Contenttype string `json:"contenttype,omitempty"`

		Data interface{} `json:"data,omitempty"`

		Extensions interface{} `json:"extensions,omitempty"`

		ID *string `json:"id"`

		Source *string `json:"source"`

		Specversion *string `json:"specversion"`

		Time strfmt.DateTime `json:"time,omitempty"`

		Type *string `json:"type"`
	}

	dataAO0.Contenttype = o.Contenttype

	dataAO0.Data = o.Data

	dataAO0.Extensions = o.Extensions

	dataAO0.ID = o.ID

	dataAO0.Source = o.Source

	dataAO0.Specversion = o.Specversion

	dataAO0.Time = o.Time

	dataAO0.Type = o.Type

	jsonDataAO0, errAO0 := swag.WriteJSON(dataAO0)
	if errAO0 != nil {
		return nil, errAO0
	}
	_parts = append(_parts, jsonDataAO0)

	var dataAO1 struct {
		Shkeptncontext string `json:"shkeptncontext,omitempty"`
	}

	dataAO1.Shkeptncontext = o.Shkeptncontext

	jsonDataAO1, errAO1 := swag.WriteJSON(dataAO1)
	if errAO1 != nil {
		return nil, errAO1
	}
	_parts = append(_parts, jsonDataAO1)

	return swag.ConcatJSON(_parts...), nil
}

// Validate validates this get event o k body items0
func (o *GetEventOKBodyItems0) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateID(formats); err != nil {
		res = append(res, err)
	}

	if err := o.validateSource(formats); err != nil {
		res = append(res, err)
	}

	if err := o.validateSpecversion(formats); err != nil {
		res = append(res, err)
	}

	if err := o.validateTime(formats); err != nil {
		res = append(res, err)
	}

	if err := o.validateType(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *GetEventOKBodyItems0) validateID(formats strfmt.Registry) error {

	if err := validate.Required("id", "body", o.ID); err != nil {
		return err
	}

	return nil
}

func (o *GetEventOKBodyItems0) validateSource(formats strfmt.Registry) error {

	if err := validate.Required("source", "body", o.Source); err != nil {
		return err
	}

	return nil
}

func (o *GetEventOKBodyItems0) validateSpecversion(formats strfmt.Registry) error {

	if err := validate.Required("specversion", "body", o.Specversion); err != nil {
		return err
	}

	return nil
}

func (o *GetEventOKBodyItems0) validateTime(formats strfmt.Registry) error {

	if swag.IsZero(o.Time) { // not required
		return nil
	}

	if err := validate.FormatOf("time", "body", "date-time", o.Time.String(), formats); err != nil {
		return err
	}

	return nil
}

func (o *GetEventOKBodyItems0) validateType(formats strfmt.Registry) error {

	if err := validate.Required("type", "body", o.Type); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (o *GetEventOKBodyItems0) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *GetEventOKBodyItems0) UnmarshalBinary(b []byte) error {
	var res GetEventOKBodyItems0
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
