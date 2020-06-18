// Code generated by go-swagger; DO NOT EDIT.

package project

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	models "github.com/keptn/keptn/api/models"
)

// DeleteProjectProjectNameOKCode is the HTTP code returned for type DeleteProjectProjectNameOK
const DeleteProjectProjectNameOKCode int = 200

/*DeleteProjectProjectNameOK Deleting of project triggered

swagger:response deleteProjectProjectNameOK
*/
type DeleteProjectProjectNameOK struct {

	/*
	  In: Body
	*/
	Payload *models.EventContext `json:"body,omitempty"`
}

// NewDeleteProjectProjectNameOK creates DeleteProjectProjectNameOK with default headers values
func NewDeleteProjectProjectNameOK() *DeleteProjectProjectNameOK {

	return &DeleteProjectProjectNameOK{}
}

// WithPayload adds the payload to the delete project project name o k response
func (o *DeleteProjectProjectNameOK) WithPayload(payload *models.EventContext) *DeleteProjectProjectNameOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the delete project project name o k response
func (o *DeleteProjectProjectNameOK) SetPayload(payload *models.EventContext) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *DeleteProjectProjectNameOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// DeleteProjectProjectNameBadRequestCode is the HTTP code returned for type DeleteProjectProjectNameBadRequest
const DeleteProjectProjectNameBadRequestCode int = 400

/*DeleteProjectProjectNameBadRequest Failed. Project could not be deleted

swagger:response deleteProjectProjectNameBadRequest
*/
type DeleteProjectProjectNameBadRequest struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewDeleteProjectProjectNameBadRequest creates DeleteProjectProjectNameBadRequest with default headers values
func NewDeleteProjectProjectNameBadRequest() *DeleteProjectProjectNameBadRequest {

	return &DeleteProjectProjectNameBadRequest{}
}

// WithPayload adds the payload to the delete project project name bad request response
func (o *DeleteProjectProjectNameBadRequest) WithPayload(payload *models.Error) *DeleteProjectProjectNameBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the delete project project name bad request response
func (o *DeleteProjectProjectNameBadRequest) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *DeleteProjectProjectNameBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(400)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

/*DeleteProjectProjectNameDefault Error

swagger:response deleteProjectProjectNameDefault
*/
type DeleteProjectProjectNameDefault struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewDeleteProjectProjectNameDefault creates DeleteProjectProjectNameDefault with default headers values
func NewDeleteProjectProjectNameDefault(code int) *DeleteProjectProjectNameDefault {
	if code <= 0 {
		code = 500
	}

	return &DeleteProjectProjectNameDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the delete project project name default response
func (o *DeleteProjectProjectNameDefault) WithStatusCode(code int) *DeleteProjectProjectNameDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the delete project project name default response
func (o *DeleteProjectProjectNameDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the delete project project name default response
func (o *DeleteProjectProjectNameDefault) WithPayload(payload *models.Error) *DeleteProjectProjectNameDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the delete project project name default response
func (o *DeleteProjectProjectNameDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *DeleteProjectProjectNameDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
