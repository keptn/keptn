// Code generated by go-swagger; DO NOT EDIT.

package event

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	models "github.com/keptn/keptn/api/models"
)

// SendEventOKCode is the HTTP code returned for type SendEventOK
const SendEventOKCode int = 200

/*SendEventOK Forwarded

swagger:response sendEventOK
*/
type SendEventOK struct {

	/*
	  In: Body
	*/
	Payload *models.EventContext `json:"body,omitempty"`
}

// NewSendEventOK creates SendEventOK with default headers values
func NewSendEventOK() *SendEventOK {

	return &SendEventOK{}
}

// WithPayload adds the payload to the send event o k response
func (o *SendEventOK) WithPayload(payload *models.EventContext) *SendEventOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the send event o k response
func (o *SendEventOK) SetPayload(payload *models.EventContext) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *SendEventOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

/*SendEventDefault Error

swagger:response sendEventDefault
*/
type SendEventDefault struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewSendEventDefault creates SendEventDefault with default headers values
func NewSendEventDefault(code int) *SendEventDefault {
	if code <= 0 {
		code = 500
	}

	return &SendEventDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the send event default response
func (o *SendEventDefault) WithStatusCode(code int) *SendEventDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the send event default response
func (o *SendEventDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the send event default response
func (o *SendEventDefault) WithPayload(payload *models.Error) *SendEventDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the send event default response
func (o *SendEventDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *SendEventDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
