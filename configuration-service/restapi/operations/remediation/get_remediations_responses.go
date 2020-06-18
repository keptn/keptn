// Code generated by go-swagger; DO NOT EDIT.

package remediation

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	models "github.com/keptn/keptn/configuration-service/models"
)

// GetRemediationsOKCode is the HTTP code returned for type GetRemediationsOK
const GetRemediationsOKCode int = 200

/*GetRemediationsOK Success

swagger:response getRemediationsOK
*/
type GetRemediationsOK struct {

	/*
	  In: Body
	*/
	Payload *models.Remediations `json:"body,omitempty"`
}

// NewGetRemediationsOK creates GetRemediationsOK with default headers values
func NewGetRemediationsOK() *GetRemediationsOK {

	return &GetRemediationsOK{}
}

// WithPayload adds the payload to the get remediations o k response
func (o *GetRemediationsOK) WithPayload(payload *models.Remediations) *GetRemediationsOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get remediations o k response
func (o *GetRemediationsOK) SetPayload(payload *models.Remediations) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetRemediationsOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetRemediationsNotFoundCode is the HTTP code returned for type GetRemediationsNotFound
const GetRemediationsNotFoundCode int = 404

/*GetRemediationsNotFound Failed. Service could not be found.

swagger:response getRemediationsNotFound
*/
type GetRemediationsNotFound struct {

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetRemediationsNotFound creates GetRemediationsNotFound with default headers values
func NewGetRemediationsNotFound() *GetRemediationsNotFound {

	return &GetRemediationsNotFound{}
}

// WithPayload adds the payload to the get remediations not found response
func (o *GetRemediationsNotFound) WithPayload(payload *models.Error) *GetRemediationsNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get remediations not found response
func (o *GetRemediationsNotFound) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetRemediationsNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

/*GetRemediationsDefault Error

swagger:response getRemediationsDefault
*/
type GetRemediationsDefault struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetRemediationsDefault creates GetRemediationsDefault with default headers values
func NewGetRemediationsDefault(code int) *GetRemediationsDefault {
	if code <= 0 {
		code = 500
	}

	return &GetRemediationsDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the get remediations default response
func (o *GetRemediationsDefault) WithStatusCode(code int) *GetRemediationsDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the get remediations default response
func (o *GetRemediationsDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the get remediations default response
func (o *GetRemediationsDefault) WithPayload(payload *models.Error) *GetRemediationsDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get remediations default response
func (o *GetRemediationsDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetRemediationsDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
