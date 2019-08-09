// Code generated by go-swagger; DO NOT EDIT.

package project

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// GetProjectProjectNameHandlerFunc turns a function with the right signature into a get project project name handler
type GetProjectProjectNameHandlerFunc func(GetProjectProjectNameParams) middleware.Responder

// Handle executing the request and returning a response
func (fn GetProjectProjectNameHandlerFunc) Handle(params GetProjectProjectNameParams) middleware.Responder {
	return fn(params)
}

// GetProjectProjectNameHandler interface for that can handle valid get project project name params
type GetProjectProjectNameHandler interface {
	Handle(GetProjectProjectNameParams) middleware.Responder
}

// NewGetProjectProjectName creates a new http.Handler for the get project project name operation
func NewGetProjectProjectName(ctx *middleware.Context, handler GetProjectProjectNameHandler) *GetProjectProjectName {
	return &GetProjectProjectName{Context: ctx, Handler: handler}
}

/*GetProjectProjectName swagger:route GET /project/{projectName} Project getProjectProjectName

Get the specified project

*/
type GetProjectProjectName struct {
	Context *middleware.Context
	Handler GetProjectProjectNameHandler
}

func (o *GetProjectProjectName) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewGetProjectProjectNameParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
