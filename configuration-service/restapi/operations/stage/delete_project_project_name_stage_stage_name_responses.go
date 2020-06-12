// Code generated by go-swagger; DO NOT EDIT.

package stage

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	models "github.com/keptn/keptn/configuration-service/models"
)

// DeleteProjectProjectNameStageStageNameNoContentCode is the HTTP code returned for type DeleteProjectProjectNameStageStageNameNoContent
const DeleteProjectProjectNameStageStageNameNoContentCode int = 204

/*DeleteProjectProjectNameStageStageNameNoContent Success. Stage has been deleted. Response does not have a body.

swagger:response deleteProjectProjectNameStageStageNameNoContent
*/
type DeleteProjectProjectNameStageStageNameNoContent struct {
}

// NewDeleteProjectProjectNameStageStageNameNoContent creates DeleteProjectProjectNameStageStageNameNoContent with default headers values
func NewDeleteProjectProjectNameStageStageNameNoContent() *DeleteProjectProjectNameStageStageNameNoContent {

	return &DeleteProjectProjectNameStageStageNameNoContent{}
}

// WriteResponse to the client
func (o *DeleteProjectProjectNameStageStageNameNoContent) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(204)
}

// DeleteProjectProjectNameStageStageNameBadRequestCode is the HTTP code returned for type DeleteProjectProjectNameStageStageNameBadRequest
const DeleteProjectProjectNameStageStageNameBadRequestCode int = 400

/*DeleteProjectProjectNameStageStageNameBadRequest Failed. Stage could not be deleted.

swagger:response deleteProjectProjectNameStageStageNameBadRequest
*/
type DeleteProjectProjectNameStageStageNameBadRequest struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewDeleteProjectProjectNameStageStageNameBadRequest creates DeleteProjectProjectNameStageStageNameBadRequest with default headers values
func NewDeleteProjectProjectNameStageStageNameBadRequest() *DeleteProjectProjectNameStageStageNameBadRequest {

	return &DeleteProjectProjectNameStageStageNameBadRequest{}
}

// WithPayload adds the payload to the delete project project name stage stage name bad request response
func (o *DeleteProjectProjectNameStageStageNameBadRequest) WithPayload(payload *models.Error) *DeleteProjectProjectNameStageStageNameBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the delete project project name stage stage name bad request response
func (o *DeleteProjectProjectNameStageStageNameBadRequest) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *DeleteProjectProjectNameStageStageNameBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(400)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

/*DeleteProjectProjectNameStageStageNameDefault Error

swagger:response deleteProjectProjectNameStageStageNameDefault
*/
type DeleteProjectProjectNameStageStageNameDefault struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewDeleteProjectProjectNameStageStageNameDefault creates DeleteProjectProjectNameStageStageNameDefault with default headers values
func NewDeleteProjectProjectNameStageStageNameDefault(code int) *DeleteProjectProjectNameStageStageNameDefault {
	if code <= 0 {
		code = 500
	}

	return &DeleteProjectProjectNameStageStageNameDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the delete project project name stage stage name default response
func (o *DeleteProjectProjectNameStageStageNameDefault) WithStatusCode(code int) *DeleteProjectProjectNameStageStageNameDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the delete project project name stage stage name default response
func (o *DeleteProjectProjectNameStageStageNameDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the delete project project name stage stage name default response
func (o *DeleteProjectProjectNameStageStageNameDefault) WithPayload(payload *models.Error) *DeleteProjectProjectNameStageStageNameDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the delete project project name stage stage name default response
func (o *DeleteProjectProjectNameStageStageNameDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *DeleteProjectProjectNameStageStageNameDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
