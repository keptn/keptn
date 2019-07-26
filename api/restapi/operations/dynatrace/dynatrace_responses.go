// Code generated by go-swagger; DO NOT EDIT.

package dynatrace

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	models "github.com/keptn/keptn/api/models"
)

// DynatraceCreatedCode is the HTTP code returned for type DynatraceCreated
const DynatraceCreatedCode int = 201

/*DynatraceCreated event forwarded

swagger:response dynatraceCreated
*/
type DynatraceCreated struct {
}

// NewDynatraceCreated creates DynatraceCreated with default headers values
func NewDynatraceCreated() *DynatraceCreated {

	return &DynatraceCreated{}
}

// WriteResponse to the client
func (o *DynatraceCreated) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(201)
}

/*DynatraceDefault error

swagger:response dynatraceDefault
*/
type DynatraceDefault struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewDynatraceDefault creates DynatraceDefault with default headers values
func NewDynatraceDefault(code int) *DynatraceDefault {
	if code <= 0 {
		code = 500
	}

	return &DynatraceDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the dynatrace default response
func (o *DynatraceDefault) WithStatusCode(code int) *DynatraceDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the dynatrace default response
func (o *DynatraceDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the dynatrace default response
func (o *DynatraceDefault) WithPayload(payload *models.Error) *DynatraceDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the dynatrace default response
func (o *DynatraceDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *DynatraceDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
