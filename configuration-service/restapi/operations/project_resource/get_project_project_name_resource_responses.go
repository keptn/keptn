// Code generated by go-swagger; DO NOT EDIT.

package project_resource

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/keptn/keptn/configuration-service/models"
)

// GetProjectProjectNameResourceOKCode is the HTTP code returned for type GetProjectProjectNameResourceOK
const GetProjectProjectNameResourceOKCode int = 200

/*GetProjectProjectNameResourceOK Success

swagger:response getProjectProjectNameResourceOK
*/
type GetProjectProjectNameResourceOK struct {

	/*
	  In: Body
	*/
	Payload *models.Resources `json:"body,omitempty"`
}

// NewGetProjectProjectNameResourceOK creates GetProjectProjectNameResourceOK with default headers values
func NewGetProjectProjectNameResourceOK() *GetProjectProjectNameResourceOK {

	return &GetProjectProjectNameResourceOK{}
}

// WithPayload adds the payload to the get project project name resource o k response
func (o *GetProjectProjectNameResourceOK) WithPayload(payload *models.Resources) *GetProjectProjectNameResourceOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get project project name resource o k response
func (o *GetProjectProjectNameResourceOK) SetPayload(payload *models.Resources) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetProjectProjectNameResourceOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetProjectProjectNameResourceNotFoundCode is the HTTP code returned for type GetProjectProjectNameResourceNotFound
const GetProjectProjectNameResourceNotFoundCode int = 404

/*GetProjectProjectNameResourceNotFound Failed. Containing project could not be found.

swagger:response getProjectProjectNameResourceNotFound
*/
type GetProjectProjectNameResourceNotFound struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetProjectProjectNameResourceNotFound creates GetProjectProjectNameResourceNotFound with default headers values
func NewGetProjectProjectNameResourceNotFound() *GetProjectProjectNameResourceNotFound {

	return &GetProjectProjectNameResourceNotFound{}
}

// WithPayload adds the payload to the get project project name resource not found response
func (o *GetProjectProjectNameResourceNotFound) WithPayload(payload *models.Error) *GetProjectProjectNameResourceNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get project project name resource not found response
func (o *GetProjectProjectNameResourceNotFound) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetProjectProjectNameResourceNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

/*GetProjectProjectNameResourceDefault Error

swagger:response getProjectProjectNameResourceDefault
*/
type GetProjectProjectNameResourceDefault struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetProjectProjectNameResourceDefault creates GetProjectProjectNameResourceDefault with default headers values
func NewGetProjectProjectNameResourceDefault(code int) *GetProjectProjectNameResourceDefault {
	if code <= 0 {
		code = 500
	}

	return &GetProjectProjectNameResourceDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the get project project name resource default response
func (o *GetProjectProjectNameResourceDefault) WithStatusCode(code int) *GetProjectProjectNameResourceDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the get project project name resource default response
func (o *GetProjectProjectNameResourceDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the get project project name resource default response
func (o *GetProjectProjectNameResourceDefault) WithPayload(payload *models.Error) *GetProjectProjectNameResourceDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get project project name resource default response
func (o *GetProjectProjectNameResourceDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetProjectProjectNameResourceDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
