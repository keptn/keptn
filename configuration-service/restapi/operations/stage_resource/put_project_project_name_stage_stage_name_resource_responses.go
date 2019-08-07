// Code generated by go-swagger; DO NOT EDIT.

package stage_resource

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	models "github.com/keptn/keptn/configuration-service/models"
)

// PutProjectProjectNameStageStageNameResourceCreatedCode is the HTTP code returned for type PutProjectProjectNameStageStageNameResourceCreated
const PutProjectProjectNameStageStageNameResourceCreatedCode int = 201

/*PutProjectProjectNameStageStageNameResourceCreated Success. Stage resources have been updated. The version of the new configuration is returned.

swagger:response putProjectProjectNameStageStageNameResourceCreated
*/
type PutProjectProjectNameStageStageNameResourceCreated struct {

	/*
	  In: Body
	*/
	Payload *models.Version `json:"body,omitempty"`
}

// NewPutProjectProjectNameStageStageNameResourceCreated creates PutProjectProjectNameStageStageNameResourceCreated with default headers values
func NewPutProjectProjectNameStageStageNameResourceCreated() *PutProjectProjectNameStageStageNameResourceCreated {

	return &PutProjectProjectNameStageStageNameResourceCreated{}
}

// WithPayload adds the payload to the put project project name stage stage name resource created response
func (o *PutProjectProjectNameStageStageNameResourceCreated) WithPayload(payload *models.Version) *PutProjectProjectNameStageStageNameResourceCreated {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the put project project name stage stage name resource created response
func (o *PutProjectProjectNameStageStageNameResourceCreated) SetPayload(payload *models.Version) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PutProjectProjectNameStageStageNameResourceCreated) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(201)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// PutProjectProjectNameStageStageNameResourceBadRequestCode is the HTTP code returned for type PutProjectProjectNameStageStageNameResourceBadRequest
const PutProjectProjectNameStageStageNameResourceBadRequestCode int = 400

/*PutProjectProjectNameStageStageNameResourceBadRequest Failed. Stage resources could not be updated.

swagger:response putProjectProjectNameStageStageNameResourceBadRequest
*/
type PutProjectProjectNameStageStageNameResourceBadRequest struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewPutProjectProjectNameStageStageNameResourceBadRequest creates PutProjectProjectNameStageStageNameResourceBadRequest with default headers values
func NewPutProjectProjectNameStageStageNameResourceBadRequest() *PutProjectProjectNameStageStageNameResourceBadRequest {

	return &PutProjectProjectNameStageStageNameResourceBadRequest{}
}

// WithPayload adds the payload to the put project project name stage stage name resource bad request response
func (o *PutProjectProjectNameStageStageNameResourceBadRequest) WithPayload(payload *models.Error) *PutProjectProjectNameStageStageNameResourceBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the put project project name stage stage name resource bad request response
func (o *PutProjectProjectNameStageStageNameResourceBadRequest) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PutProjectProjectNameStageStageNameResourceBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(400)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
