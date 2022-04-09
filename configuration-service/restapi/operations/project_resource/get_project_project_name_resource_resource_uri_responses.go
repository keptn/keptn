// Code generated by go-swagger; DO NOT EDIT.

package project_resource

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/keptn/keptn/configuration-service/models"
)

// GetProjectProjectNameResourceResourceURIOKCode is the HTTP code returned for type GetProjectProjectNameResourceResourceURIOK
const GetProjectProjectNameResourceResourceURIOKCode int = 200

/*GetProjectProjectNameResourceResourceURIOK Success

swagger:response getProjectProjectNameResourceResourceUriOK
*/
type GetProjectProjectNameResourceResourceURIOK struct {

	/*
	  In: Body
	*/
	Payload *models.Resource `json:"body,omitempty"`
}

// NewGetProjectProjectNameResourceResourceURIOK creates GetProjectProjectNameResourceResourceURIOK with default headers values
func NewGetProjectProjectNameResourceResourceURIOK() *GetProjectProjectNameResourceResourceURIOK {

	return &GetProjectProjectNameResourceResourceURIOK{}
}

// WithPayload adds the payload to the get project project name resource resource URI o k response
func (o *GetProjectProjectNameResourceResourceURIOK) WithPayload(payload *models.Resource) *GetProjectProjectNameResourceResourceURIOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get project project name resource resource URI o k response
func (o *GetProjectProjectNameResourceResourceURIOK) SetPayload(payload *models.Resource) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetProjectProjectNameResourceResourceURIOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetProjectProjectNameResourceResourceURINotFoundCode is the HTTP code returned for type GetProjectProjectNameResourceResourceURINotFound
const GetProjectProjectNameResourceResourceURINotFoundCode int = 404

/*GetProjectProjectNameResourceResourceURINotFound Failed. Project resource could not be found.

swagger:response getProjectProjectNameResourceResourceUriNotFound
*/
type GetProjectProjectNameResourceResourceURINotFound struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetProjectProjectNameResourceResourceURINotFound creates GetProjectProjectNameResourceResourceURINotFound with default headers values
func NewGetProjectProjectNameResourceResourceURINotFound() *GetProjectProjectNameResourceResourceURINotFound {

	return &GetProjectProjectNameResourceResourceURINotFound{}
}

// WithPayload adds the payload to the get project project name resource resource URI not found response
func (o *GetProjectProjectNameResourceResourceURINotFound) WithPayload(payload *models.Error) *GetProjectProjectNameResourceResourceURINotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get project project name resource resource URI not found response
func (o *GetProjectProjectNameResourceResourceURINotFound) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetProjectProjectNameResourceResourceURINotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

/*GetProjectProjectNameResourceResourceURIDefault Error

swagger:response getProjectProjectNameResourceResourceUriDefault
*/
type GetProjectProjectNameResourceResourceURIDefault struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetProjectProjectNameResourceResourceURIDefault creates GetProjectProjectNameResourceResourceURIDefault with default headers values
func NewGetProjectProjectNameResourceResourceURIDefault(code int) *GetProjectProjectNameResourceResourceURIDefault {
	if code <= 0 {
		code = 500
	}

	return &GetProjectProjectNameResourceResourceURIDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the get project project name resource resource URI default response
func (o *GetProjectProjectNameResourceResourceURIDefault) WithStatusCode(code int) *GetProjectProjectNameResourceResourceURIDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the get project project name resource resource URI default response
func (o *GetProjectProjectNameResourceResourceURIDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the get project project name resource resource URI default response
func (o *GetProjectProjectNameResourceResourceURIDefault) WithPayload(payload *models.Error) *GetProjectProjectNameResourceResourceURIDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get project project name resource resource URI default response
func (o *GetProjectProjectNameResourceResourceURIDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetProjectProjectNameResourceResourceURIDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
