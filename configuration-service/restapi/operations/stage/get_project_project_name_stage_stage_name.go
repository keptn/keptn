// Code generated by go-swagger; DO NOT EDIT.

package stage

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// GetProjectProjectNameStageStageNameHandlerFunc turns a function with the right signature into a get project project name stage stage name handler
type GetProjectProjectNameStageStageNameHandlerFunc func(GetProjectProjectNameStageStageNameParams) middleware.Responder

// Handle executing the request and returning a response
func (fn GetProjectProjectNameStageStageNameHandlerFunc) Handle(params GetProjectProjectNameStageStageNameParams) middleware.Responder {
	return fn(params)
}

// GetProjectProjectNameStageStageNameHandler interface for that can handle valid get project project name stage stage name params
type GetProjectProjectNameStageStageNameHandler interface {
	Handle(GetProjectProjectNameStageStageNameParams) middleware.Responder
}

// NewGetProjectProjectNameStageStageName creates a new http.Handler for the get project project name stage stage name operation
func NewGetProjectProjectNameStageStageName(ctx *middleware.Context, handler GetProjectProjectNameStageStageNameHandler) *GetProjectProjectNameStageStageName {
	return &GetProjectProjectNameStageStageName{Context: ctx, Handler: handler}
}

/*GetProjectProjectNameStageStageName swagger:route GET /project/{projectName}/stage/{stageName} Stage getProjectProjectNameStageStageName

Get the specified stage

*/
type GetProjectProjectNameStageStageName struct {
	Context *middleware.Context
	Handler GetProjectProjectNameStageStageNameHandler
}

func (o *GetProjectProjectNameStageStageName) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewGetProjectProjectNameStageStageNameParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
