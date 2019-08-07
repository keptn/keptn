// Code generated by go-swagger; DO NOT EDIT.

package service_resource

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	models "github.com/keptn/keptn/configuration-service/models"
)

// PutProjectProjectNameStageStageNameServiceServiceNameResourceResourceURICreatedCode is the HTTP code returned for type PutProjectProjectNameStageStageNameServiceServiceNameResourceResourceURICreated
const PutProjectProjectNameStageStageNameServiceServiceNameResourceResourceURICreatedCode int = 201

/*PutProjectProjectNameStageStageNameServiceServiceNameResourceResourceURICreated Success. Service resource has been updated. The version of the new configuration is returned.

swagger:response putProjectProjectNameStageStageNameServiceServiceNameResourceResourceUriCreated
*/
type PutProjectProjectNameStageStageNameServiceServiceNameResourceResourceURICreated struct {

	/*
	  In: Body
	*/
	Payload *models.Version `json:"body,omitempty"`
}

// NewPutProjectProjectNameStageStageNameServiceServiceNameResourceResourceURICreated creates PutProjectProjectNameStageStageNameServiceServiceNameResourceResourceURICreated with default headers values
func NewPutProjectProjectNameStageStageNameServiceServiceNameResourceResourceURICreated() *PutProjectProjectNameStageStageNameServiceServiceNameResourceResourceURICreated {

	return &PutProjectProjectNameStageStageNameServiceServiceNameResourceResourceURICreated{}
}

// WithPayload adds the payload to the put project project name stage stage name service service name resource resource Uri created response
func (o *PutProjectProjectNameStageStageNameServiceServiceNameResourceResourceURICreated) WithPayload(payload *models.Version) *PutProjectProjectNameStageStageNameServiceServiceNameResourceResourceURICreated {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the put project project name stage stage name service service name resource resource Uri created response
func (o *PutProjectProjectNameStageStageNameServiceServiceNameResourceResourceURICreated) SetPayload(payload *models.Version) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PutProjectProjectNameStageStageNameServiceServiceNameResourceResourceURICreated) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(201)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// PutProjectProjectNameStageStageNameServiceServiceNameResourceResourceURIBadRequestCode is the HTTP code returned for type PutProjectProjectNameStageStageNameServiceServiceNameResourceResourceURIBadRequest
const PutProjectProjectNameStageStageNameServiceServiceNameResourceResourceURIBadRequestCode int = 400

/*PutProjectProjectNameStageStageNameServiceServiceNameResourceResourceURIBadRequest Failed. Service resource could not be updated.

swagger:response putProjectProjectNameStageStageNameServiceServiceNameResourceResourceUriBadRequest
*/
type PutProjectProjectNameStageStageNameServiceServiceNameResourceResourceURIBadRequest struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewPutProjectProjectNameStageStageNameServiceServiceNameResourceResourceURIBadRequest creates PutProjectProjectNameStageStageNameServiceServiceNameResourceResourceURIBadRequest with default headers values
func NewPutProjectProjectNameStageStageNameServiceServiceNameResourceResourceURIBadRequest() *PutProjectProjectNameStageStageNameServiceServiceNameResourceResourceURIBadRequest {

	return &PutProjectProjectNameStageStageNameServiceServiceNameResourceResourceURIBadRequest{}
}

// WithPayload adds the payload to the put project project name stage stage name service service name resource resource Uri bad request response
func (o *PutProjectProjectNameStageStageNameServiceServiceNameResourceResourceURIBadRequest) WithPayload(payload *models.Error) *PutProjectProjectNameStageStageNameServiceServiceNameResourceResourceURIBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the put project project name stage stage name service service name resource resource Uri bad request response
func (o *PutProjectProjectNameStageStageNameServiceServiceNameResourceResourceURIBadRequest) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PutProjectProjectNameStageStageNameServiceServiceNameResourceResourceURIBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(400)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
