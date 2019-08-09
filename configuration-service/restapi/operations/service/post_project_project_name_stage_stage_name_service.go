// Code generated by go-swagger; DO NOT EDIT.

package service

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// PostProjectProjectNameStageStageNameServiceHandlerFunc turns a function with the right signature into a post project project name stage stage name service handler
type PostProjectProjectNameStageStageNameServiceHandlerFunc func(PostProjectProjectNameStageStageNameServiceParams) middleware.Responder

// Handle executing the request and returning a response
func (fn PostProjectProjectNameStageStageNameServiceHandlerFunc) Handle(params PostProjectProjectNameStageStageNameServiceParams) middleware.Responder {
	return fn(params)
}

// PostProjectProjectNameStageStageNameServiceHandler interface for that can handle valid post project project name stage stage name service params
type PostProjectProjectNameStageStageNameServiceHandler interface {
	Handle(PostProjectProjectNameStageStageNameServiceParams) middleware.Responder
}

// NewPostProjectProjectNameStageStageNameService creates a new http.Handler for the post project project name stage stage name service operation
func NewPostProjectProjectNameStageStageNameService(ctx *middleware.Context, handler PostProjectProjectNameStageStageNameServiceHandler) *PostProjectProjectNameStageStageNameService {
	return &PostProjectProjectNameStageStageNameService{Context: ctx, Handler: handler}
}

/*PostProjectProjectNameStageStageNameService swagger:route POST /project/{projectName}/stage/{stageName}/service Service postProjectProjectNameStageStageNameService

Create a new service by service name

*/
type PostProjectProjectNameStageStageNameService struct {
	Context *middleware.Context
	Handler PostProjectProjectNameStageStageNameServiceHandler
}

func (o *PostProjectProjectNameStageStageNameService) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewPostProjectProjectNameStageStageNameServiceParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
