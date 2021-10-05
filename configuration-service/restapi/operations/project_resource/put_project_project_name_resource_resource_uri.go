// Code generated by go-swagger; DO NOT EDIT.

package project_resource

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// PutProjectProjectNameResourceResourceURIHandlerFunc turns a function with the right signature into a put project project name resource resource URI handler
type PutProjectProjectNameResourceResourceURIHandlerFunc func(PutProjectProjectNameResourceResourceURIParams) middleware.Responder

// Handle executing the request and returning a response
func (fn PutProjectProjectNameResourceResourceURIHandlerFunc) Handle(params PutProjectProjectNameResourceResourceURIParams) middleware.Responder {
	return fn(params)
}

// PutProjectProjectNameResourceResourceURIHandler interface for that can handle valid put project project name resource resource URI params
type PutProjectProjectNameResourceResourceURIHandler interface {
	Handle(PutProjectProjectNameResourceResourceURIParams) middleware.Responder
}

// NewPutProjectProjectNameResourceResourceURI creates a new http.Handler for the put project project name resource resource URI operation
func NewPutProjectProjectNameResourceResourceURI(ctx *middleware.Context, handler PutProjectProjectNameResourceResourceURIHandler) *PutProjectProjectNameResourceResourceURI {
	return &PutProjectProjectNameResourceResourceURI{Context: ctx, Handler: handler}
}

/* PutProjectProjectNameResourceResourceURI swagger:route PUT /project/{projectName}/resource/{resourceURI} Project Resource putProjectProjectNameResourceResourceUri

Update the specified resource

*/
type PutProjectProjectNameResourceResourceURI struct {
	Context *middleware.Context
	Handler PutProjectProjectNameResourceResourceURIHandler
}

func (o *PutProjectProjectNameResourceResourceURI) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewPutProjectProjectNameResourceResourceURIParams()
	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}
