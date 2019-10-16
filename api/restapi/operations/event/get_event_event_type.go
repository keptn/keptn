// Code generated by go-swagger; DO NOT EDIT.

package event

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"

	models "github.com/keptn/keptn/api/models"
)

// GetEventEventTypeHandlerFunc turns a function with the right signature into a get event event type handler
type GetEventEventTypeHandlerFunc func(GetEventEventTypeParams, *models.Principal) middleware.Responder

// Handle executing the request and returning a response
func (fn GetEventEventTypeHandlerFunc) Handle(params GetEventEventTypeParams, principal *models.Principal) middleware.Responder {
	return fn(params, principal)
}

// GetEventEventTypeHandler interface for that can handle valid get event event type params
type GetEventEventTypeHandler interface {
	Handle(GetEventEventTypeParams, *models.Principal) middleware.Responder
}

// NewGetEventEventType creates a new http.Handler for the get event event type operation
func NewGetEventEventType(ctx *middleware.Context, handler GetEventEventTypeHandler) *GetEventEventType {
	return &GetEventEventType{Context: ctx, Handler: handler}
}

/*GetEventEventType swagger:route GET /event/{eventType} Event getEventEventType

Get the specified event

*/
type GetEventEventType struct {
	Context *middleware.Context
	Handler GetEventEventTypeHandler
}

func (o *GetEventEventType) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewGetEventEventTypeParams()

	uprinc, aCtx, err := o.Context.Authorize(r, route)
	if err != nil {
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}
	if aCtx != nil {
		r = aCtx
	}
	var principal *models.Principal
	if uprinc != nil {
		principal = uprinc.(*models.Principal) // this is really a models.Principal, I promise
	}

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params, principal) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
