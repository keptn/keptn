// Code generated by go-swagger; DO NOT EDIT.

package event

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"
)

// GetNewArtifactEventsOKCode is the HTTP code returned for type GetNewArtifactEventsOK
const GetNewArtifactEventsOKCode int = 200

/*GetNewArtifactEventsOK ok

swagger:response getNewArtifactEventsOK
*/
type GetNewArtifactEventsOK struct {

	/*
	  In: Body
	*/
	Payload []*GetNewArtifactEventsOKBodyItems0 `json:"body,omitempty"`
}

// NewGetNewArtifactEventsOK creates GetNewArtifactEventsOK with default headers values
func NewGetNewArtifactEventsOK() *GetNewArtifactEventsOK {

	return &GetNewArtifactEventsOK{}
}

// WithPayload adds the payload to the get new artifact events o k response
func (o *GetNewArtifactEventsOK) WithPayload(payload []*GetNewArtifactEventsOKBodyItems0) *GetNewArtifactEventsOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get new artifact events o k response
func (o *GetNewArtifactEventsOK) SetPayload(payload []*GetNewArtifactEventsOKBodyItems0) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetNewArtifactEventsOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	payload := o.Payload
	if payload == nil {
		// return empty array
		payload = make([]*GetNewArtifactEventsOKBodyItems0, 0, 50)
	}

	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}

/*GetNewArtifactEventsDefault error

swagger:response getNewArtifactEventsDefault
*/
type GetNewArtifactEventsDefault struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload *GetNewArtifactEventsDefaultBody `json:"body,omitempty"`
}

// NewGetNewArtifactEventsDefault creates GetNewArtifactEventsDefault with default headers values
func NewGetNewArtifactEventsDefault(code int) *GetNewArtifactEventsDefault {
	if code <= 0 {
		code = 500
	}

	return &GetNewArtifactEventsDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the get new artifact events default response
func (o *GetNewArtifactEventsDefault) WithStatusCode(code int) *GetNewArtifactEventsDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the get new artifact events default response
func (o *GetNewArtifactEventsDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the get new artifact events default response
func (o *GetNewArtifactEventsDefault) WithPayload(payload *GetNewArtifactEventsDefaultBody) *GetNewArtifactEventsDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get new artifact events default response
func (o *GetNewArtifactEventsDefault) SetPayload(payload *GetNewArtifactEventsDefaultBody) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetNewArtifactEventsDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
