// Code generated by go-swagger; DO NOT EDIT.

package event

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"
)

// GetEventsOKCode is the HTTP code returned for type GetEventsOK
const GetEventsOKCode int = 200

/*GetEventsOK ok

swagger:response getEventsOK
*/
type GetEventsOK struct {

	/*
	  In: Body
	*/
	Payload *GetEventsOKBody `json:"body,omitempty"`
}

// NewGetEventsOK creates GetEventsOK with default headers values
func NewGetEventsOK() *GetEventsOK {

	return &GetEventsOK{}
}

// WithPayload adds the payload to the get events o k response
func (o *GetEventsOK) WithPayload(payload *GetEventsOKBody) *GetEventsOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get events o k response
func (o *GetEventsOK) SetPayload(payload *GetEventsOKBody) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetEventsOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

/*GetEventsDefault error

swagger:response getEventsDefault
*/
type GetEventsDefault struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload *GetEventsDefaultBody `json:"body,omitempty"`
}

// NewGetEventsDefault creates GetEventsDefault with default headers values
func NewGetEventsDefault(code int) *GetEventsDefault {
	if code <= 0 {
		code = 500
	}

	return &GetEventsDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the get events default response
func (o *GetEventsDefault) WithStatusCode(code int) *GetEventsDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the get events default response
func (o *GetEventsDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the get events default response
func (o *GetEventsDefault) WithPayload(payload *GetEventsDefaultBody) *GetEventsDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get events default response
func (o *GetEventsDefault) SetPayload(payload *GetEventsDefaultBody) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetEventsDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
