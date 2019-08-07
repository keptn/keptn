// Code generated by go-swagger; DO NOT EDIT.

package project_resource

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	models "github.com/keptn/keptn/configuration-service/models"
)

// PutProjectProjectNameResourceResourceURICreatedCode is the HTTP code returned for type PutProjectProjectNameResourceResourceURICreated
const PutProjectProjectNameResourceResourceURICreatedCode int = 201

/*PutProjectProjectNameResourceResourceURICreated Success. Project resource has been updated. The version of the new configuration is returned.

swagger:response putProjectProjectNameResourceResourceUriCreated
*/
type PutProjectProjectNameResourceResourceURICreated struct {

	/*
	  In: Body
	*/
	Payload *models.Version `json:"body,omitempty"`
}

// NewPutProjectProjectNameResourceResourceURICreated creates PutProjectProjectNameResourceResourceURICreated with default headers values
func NewPutProjectProjectNameResourceResourceURICreated() *PutProjectProjectNameResourceResourceURICreated {

	return &PutProjectProjectNameResourceResourceURICreated{}
}

// WithPayload adds the payload to the put project project name resource resource Uri created response
func (o *PutProjectProjectNameResourceResourceURICreated) WithPayload(payload *models.Version) *PutProjectProjectNameResourceResourceURICreated {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the put project project name resource resource Uri created response
func (o *PutProjectProjectNameResourceResourceURICreated) SetPayload(payload *models.Version) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PutProjectProjectNameResourceResourceURICreated) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(201)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// PutProjectProjectNameResourceResourceURIBadRequestCode is the HTTP code returned for type PutProjectProjectNameResourceResourceURIBadRequest
const PutProjectProjectNameResourceResourceURIBadRequestCode int = 400

/*PutProjectProjectNameResourceResourceURIBadRequest Failed. Project resource could not be updated.

swagger:response putProjectProjectNameResourceResourceUriBadRequest
*/
type PutProjectProjectNameResourceResourceURIBadRequest struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewPutProjectProjectNameResourceResourceURIBadRequest creates PutProjectProjectNameResourceResourceURIBadRequest with default headers values
func NewPutProjectProjectNameResourceResourceURIBadRequest() *PutProjectProjectNameResourceResourceURIBadRequest {

	return &PutProjectProjectNameResourceResourceURIBadRequest{}
}

// WithPayload adds the payload to the put project project name resource resource Uri bad request response
func (o *PutProjectProjectNameResourceResourceURIBadRequest) WithPayload(payload *models.Error) *PutProjectProjectNameResourceResourceURIBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the put project project name resource resource Uri bad request response
func (o *PutProjectProjectNameResourceResourceURIBadRequest) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PutProjectProjectNameResourceResourceURIBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(400)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
