// Code generated by go-swagger; DO NOT EDIT.

package project

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// DeleteProjectProjectNameHandlerFunc turns a function with the right signature into a delete project project name handler
type DeleteProjectProjectNameHandlerFunc func(DeleteProjectProjectNameParams) middleware.Responder

// Handle executing the request and returning a response
func (fn DeleteProjectProjectNameHandlerFunc) Handle(params DeleteProjectProjectNameParams) middleware.Responder {
	return fn(params)
}

// DeleteProjectProjectNameHandler interface for that can handle valid delete project project name params
type DeleteProjectProjectNameHandler interface {
	Handle(DeleteProjectProjectNameParams) middleware.Responder
}

// NewDeleteProjectProjectName creates a new http.Handler for the delete project project name operation
func NewDeleteProjectProjectName(ctx *middleware.Context, handler DeleteProjectProjectNameHandler) *DeleteProjectProjectName {
	return &DeleteProjectProjectName{Context: ctx, Handler: handler}
}

/*DeleteProjectProjectName swagger:route DELETE /project/{projectName} Project deleteProjectProjectName

INTERNAL Endpoint: Delete the specified project

*/
type DeleteProjectProjectName struct {
	Context *middleware.Context
	Handler DeleteProjectProjectNameHandler
}

func (o *DeleteProjectProjectName) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewDeleteProjectProjectNameParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
