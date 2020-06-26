// Code generated by go-swagger; DO NOT EDIT.

package project

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	models "github.com/keptn/keptn/configuration-service/models"
)

// GetProjectOKCode is the HTTP code returned for type GetProjectOK
const GetProjectOKCode int = 200

/*GetProjectOK Success

swagger:response getProjectOK
*/
type GetProjectOK struct {

	/*
	  In: Body
	*/
	Payload *models.ExpandedProjects `json:"body,omitempty"`
}

// NewGetProjectOK creates GetProjectOK with default headers values
func NewGetProjectOK() *GetProjectOK {

	return &GetProjectOK{}
}

// WithPayload adds the payload to the get project o k response
func (o *GetProjectOK) WithPayload(payload *models.ExpandedProjects) *GetProjectOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get project o k response
func (o *GetProjectOK) SetPayload(payload *models.ExpandedProjects) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetProjectOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

/*GetProjectDefault Error

swagger:response getProjectDefault
*/
type GetProjectDefault struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetProjectDefault creates GetProjectDefault with default headers values
func NewGetProjectDefault(code int) *GetProjectDefault {
	if code <= 0 {
		code = 500
	}

	return &GetProjectDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the get project default response
func (o *GetProjectDefault) WithStatusCode(code int) *GetProjectDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the get project default response
func (o *GetProjectDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the get project default response
func (o *GetProjectDefault) WithPayload(payload *models.Error) *GetProjectDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get project default response
func (o *GetProjectDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetProjectDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
