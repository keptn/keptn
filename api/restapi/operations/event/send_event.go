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

	models "github.com/keptn/api/models"
)

// SendEventHandlerFunc turns a function with the right signature into a send event handler
type SendEventHandlerFunc func(SendEventParams, *models.Principal) middleware.Responder

// Handle executing the request and returning a response
func (fn SendEventHandlerFunc) Handle(params SendEventParams, principal *models.Principal) middleware.Responder {
	return fn(params, principal)
}

// SendEventHandler interface for that can handle valid send event params
type SendEventHandler interface {
	Handle(SendEventParams, *models.Principal) middleware.Responder
}

// NewSendEvent creates a new http.Handler for the send event operation
func NewSendEvent(ctx *middleware.Context, handler SendEventHandler) *SendEvent {
	return &SendEvent{Context: ctx, Handler: handler}
}

/*SendEvent swagger:route POST /event event sendEvent

Forwards a new event to the eventbroker

*/
type SendEvent struct {
	Context *middleware.Context
	Handler SendEventHandler
}

func (o *SendEvent) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewSendEventParams()

	uprinc, aCtx, err := o.Context.Authorize(r, route)
	if err != nil {
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}
	if aCtx != nil {
		r = aCtx
	}
	var principal *models.Principal
	if uprinc != nil {
		principal = uprinc.(*models.Principal) // this is really a models.Principal, I promise
	}

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params, principal) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}

// SendEventBody send event body
// swagger:model SendEventBody
type SendEventBody struct {

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
	// Required: true
	Shkeptncontext *string `json:"shkeptncontext"`
}

// UnmarshalJSON unmarshals this object from a JSON structure
func (o *SendEventBody) UnmarshalJSON(raw []byte) error {
	// SendEventParamsBodyAO0
	var dataSendEventParamsBodyAO0 struct {
		Contenttype string `json:"contenttype,omitempty"`

		Data interface{} `json:"data,omitempty"`

		Extensions interface{} `json:"extensions,omitempty"`

		ID *string `json:"id"`

		Source *string `json:"source"`

		Specversion *string `json:"specversion"`

		Time strfmt.DateTime `json:"time,omitempty"`

		Type *string `json:"type"`
	}
	if err := swag.ReadJSON(raw, &dataSendEventParamsBodyAO0); err != nil {
		return err
	}

	o.Contenttype = dataSendEventParamsBodyAO0.Contenttype

	o.Data = dataSendEventParamsBodyAO0.Data

	o.Extensions = dataSendEventParamsBodyAO0.Extensions

	o.ID = dataSendEventParamsBodyAO0.ID

	o.Source = dataSendEventParamsBodyAO0.Source

	o.Specversion = dataSendEventParamsBodyAO0.Specversion

	o.Time = dataSendEventParamsBodyAO0.Time

	o.Type = dataSendEventParamsBodyAO0.Type

	// SendEventParamsBodyAO1
	var dataSendEventParamsBodyAO1 struct {
		Shkeptncontext *string `json:"shkeptncontext"`
	}
	if err := swag.ReadJSON(raw, &dataSendEventParamsBodyAO1); err != nil {
		return err
	}

	o.Shkeptncontext = dataSendEventParamsBodyAO1.Shkeptncontext

	return nil
}

// MarshalJSON marshals this object to a JSON structure
func (o SendEventBody) MarshalJSON() ([]byte, error) {
	_parts := make([][]byte, 0, 2)

	var dataSendEventParamsBodyAO0 struct {
		Contenttype string `json:"contenttype,omitempty"`

		Data interface{} `json:"data,omitempty"`

		Extensions interface{} `json:"extensions,omitempty"`

		ID *string `json:"id"`

		Source *string `json:"source"`

		Specversion *string `json:"specversion"`

		Time strfmt.DateTime `json:"time,omitempty"`

		Type *string `json:"type"`
	}

	dataSendEventParamsBodyAO0.Contenttype = o.Contenttype

	dataSendEventParamsBodyAO0.Data = o.Data

	dataSendEventParamsBodyAO0.Extensions = o.Extensions

	dataSendEventParamsBodyAO0.ID = o.ID

	dataSendEventParamsBodyAO0.Source = o.Source

	dataSendEventParamsBodyAO0.Specversion = o.Specversion

	dataSendEventParamsBodyAO0.Time = o.Time

	dataSendEventParamsBodyAO0.Type = o.Type

	jsonDataSendEventParamsBodyAO0, errSendEventParamsBodyAO0 := swag.WriteJSON(dataSendEventParamsBodyAO0)
	if errSendEventParamsBodyAO0 != nil {
		return nil, errSendEventParamsBodyAO0
	}
	_parts = append(_parts, jsonDataSendEventParamsBodyAO0)

	var dataSendEventParamsBodyAO1 struct {
		Shkeptncontext *string `json:"shkeptncontext"`
	}

	dataSendEventParamsBodyAO1.Shkeptncontext = o.Shkeptncontext

	jsonDataSendEventParamsBodyAO1, errSendEventParamsBodyAO1 := swag.WriteJSON(dataSendEventParamsBodyAO1)
	if errSendEventParamsBodyAO1 != nil {
		return nil, errSendEventParamsBodyAO1
	}
	_parts = append(_parts, jsonDataSendEventParamsBodyAO1)

	return swag.ConcatJSON(_parts...), nil
}

// Validate validates this send event body
func (o *SendEventBody) Validate(formats strfmt.Registry) error {
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

	if err := o.validateShkeptncontext(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *SendEventBody) validateID(formats strfmt.Registry) error {

	if err := validate.Required("body"+"."+"id", "body", o.ID); err != nil {
		return err
	}

	return nil
}

func (o *SendEventBody) validateSource(formats strfmt.Registry) error {

	if err := validate.Required("body"+"."+"source", "body", o.Source); err != nil {
		return err
	}

	return nil
}

func (o *SendEventBody) validateSpecversion(formats strfmt.Registry) error {

	if err := validate.Required("body"+"."+"specversion", "body", o.Specversion); err != nil {
		return err
	}

	return nil
}

func (o *SendEventBody) validateTime(formats strfmt.Registry) error {

	if swag.IsZero(o.Time) { // not required
		return nil
	}

	if err := validate.FormatOf("body"+"."+"time", "body", "date-time", o.Time.String(), formats); err != nil {
		return err
	}

	return nil
}

func (o *SendEventBody) validateType(formats strfmt.Registry) error {

	if err := validate.Required("body"+"."+"type", "body", o.Type); err != nil {
		return err
	}

	return nil
}

func (o *SendEventBody) validateShkeptncontext(formats strfmt.Registry) error {

	if err := validate.Required("body"+"."+"shkeptncontext", "body", o.Shkeptncontext); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (o *SendEventBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *SendEventBody) UnmarshalBinary(b []byte) error {
	var res SendEventBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

// SendEventDefaultBody send event default body
// swagger:model SendEventDefaultBody
type SendEventDefaultBody struct {

	// code
	Code int64 `json:"code,omitempty"`

	// fields
	Fields string `json:"fields,omitempty"`

	// message
	// Required: true
	Message *string `json:"message"`
}

// Validate validates this send event default body
func (o *SendEventDefaultBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateMessage(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *SendEventDefaultBody) validateMessage(formats strfmt.Registry) error {

	if err := validate.Required("sendEvent default"+"."+"message", "body", o.Message); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (o *SendEventDefaultBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *SendEventDefaultBody) UnmarshalBinary(b []byte) error {
	var res SendEventDefaultBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
