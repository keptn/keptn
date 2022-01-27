// Code generated by go-swagger; DO NOT EDIT.

package service_default_resource

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// PutProjectProjectNameServiceServiceNameResourceResourceURIHandlerFunc turns a function with the right signature into a put project project name service service name resource resource URI handler
type PutProjectProjectNameServiceServiceNameResourceResourceURIHandlerFunc func(PutProjectProjectNameServiceServiceNameResourceResourceURIParams) middleware.Responder

// Handle executing the request and returning a response
func (fn PutProjectProjectNameServiceServiceNameResourceResourceURIHandlerFunc) Handle(params PutProjectProjectNameServiceServiceNameResourceResourceURIParams) middleware.Responder {
	return fn(params)
}

// PutProjectProjectNameServiceServiceNameResourceResourceURIHandler interface for that can handle valid put project project name service service name resource resource URI params
type PutProjectProjectNameServiceServiceNameResourceResourceURIHandler interface {
	Handle(PutProjectProjectNameServiceServiceNameResourceResourceURIParams) middleware.Responder
}

// NewPutProjectProjectNameServiceServiceNameResourceResourceURI creates a new http.Handler for the put project project name service service name resource resource URI operation
func NewPutProjectProjectNameServiceServiceNameResourceResourceURI(ctx *middleware.Context, handler PutProjectProjectNameServiceServiceNameResourceResourceURIHandler) *PutProjectProjectNameServiceServiceNameResourceResourceURI {
	return &PutProjectProjectNameServiceServiceNameResourceResourceURI{Context: ctx, Handler: handler}
}

/*PutProjectProjectNameServiceServiceNameResourceResourceURI swagger:route PUT /project/{projectName}/service/{serviceName}/resource/{resourceURI} Service Default Resource putProjectProjectNameServiceServiceNameResourceResourceUri

Update the specified default resource for the service

*/
type PutProjectProjectNameServiceServiceNameResourceResourceURI struct {
	Context *middleware.Context
	Handler PutProjectProjectNameServiceServiceNameResourceResourceURIHandler
}

func (o *PutProjectProjectNameServiceServiceNameResourceResourceURI) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewPutProjectProjectNameServiceServiceNameResourceResourceURIParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
