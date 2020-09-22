// Code generated by go-swagger; DO NOT EDIT.

package project_resource

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/keptn/keptn/configuration-service/models"
)

// DeleteProjectProjectNameResourceResourceURINoContentCode is the HTTP code returned for type DeleteProjectProjectNameResourceResourceURINoContent
const DeleteProjectProjectNameResourceResourceURINoContentCode int = 204

/*DeleteProjectProjectNameResourceResourceURINoContent Success. Project resource has been deleted.

swagger:response deleteProjectProjectNameResourceResourceUriNoContent
*/
type DeleteProjectProjectNameResourceResourceURINoContent struct {

	/*
	  In: Body
	*/
	Payload *models.Version `json:"body,omitempty"`
}

// NewDeleteProjectProjectNameResourceResourceURINoContent creates DeleteProjectProjectNameResourceResourceURINoContent with default headers values
func NewDeleteProjectProjectNameResourceResourceURINoContent() *DeleteProjectProjectNameResourceResourceURINoContent {

	return &DeleteProjectProjectNameResourceResourceURINoContent{}
}

// WithPayload adds the payload to the delete project project name resource resource Uri no content response
func (o *DeleteProjectProjectNameResourceResourceURINoContent) WithPayload(payload *models.Version) *DeleteProjectProjectNameResourceResourceURINoContent {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the delete project project name resource resource Uri no content response
func (o *DeleteProjectProjectNameResourceResourceURINoContent) SetPayload(payload *models.Version) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *DeleteProjectProjectNameResourceResourceURINoContent) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(204)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// DeleteProjectProjectNameResourceResourceURIBadRequestCode is the HTTP code returned for type DeleteProjectProjectNameResourceResourceURIBadRequest
const DeleteProjectProjectNameResourceResourceURIBadRequestCode int = 400

/*DeleteProjectProjectNameResourceResourceURIBadRequest Failed. Project resource could not be deleted.

swagger:response deleteProjectProjectNameResourceResourceUriBadRequest
*/
type DeleteProjectProjectNameResourceResourceURIBadRequest struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewDeleteProjectProjectNameResourceResourceURIBadRequest creates DeleteProjectProjectNameResourceResourceURIBadRequest with default headers values
func NewDeleteProjectProjectNameResourceResourceURIBadRequest() *DeleteProjectProjectNameResourceResourceURIBadRequest {

	return &DeleteProjectProjectNameResourceResourceURIBadRequest{}
}

// WithPayload adds the payload to the delete project project name resource resource Uri bad request response
func (o *DeleteProjectProjectNameResourceResourceURIBadRequest) WithPayload(payload *models.Error) *DeleteProjectProjectNameResourceResourceURIBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the delete project project name resource resource Uri bad request response
func (o *DeleteProjectProjectNameResourceResourceURIBadRequest) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *DeleteProjectProjectNameResourceResourceURIBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(400)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

/*DeleteProjectProjectNameResourceResourceURIDefault Error

swagger:response deleteProjectProjectNameResourceResourceUriDefault
*/
type DeleteProjectProjectNameResourceResourceURIDefault struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewDeleteProjectProjectNameResourceResourceURIDefault creates DeleteProjectProjectNameResourceResourceURIDefault with default headers values
func NewDeleteProjectProjectNameResourceResourceURIDefault(code int) *DeleteProjectProjectNameResourceResourceURIDefault {
	if code <= 0 {
		code = 500
	}

	return &DeleteProjectProjectNameResourceResourceURIDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the delete project project name resource resource URI default response
func (o *DeleteProjectProjectNameResourceResourceURIDefault) WithStatusCode(code int) *DeleteProjectProjectNameResourceResourceURIDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the delete project project name resource resource URI default response
func (o *DeleteProjectProjectNameResourceResourceURIDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the delete project project name resource resource URI default response
func (o *DeleteProjectProjectNameResourceResourceURIDefault) WithPayload(payload *models.Error) *DeleteProjectProjectNameResourceResourceURIDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the delete project project name resource resource URI default response
func (o *DeleteProjectProjectNameResourceResourceURIDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *DeleteProjectProjectNameResourceResourceURIDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
