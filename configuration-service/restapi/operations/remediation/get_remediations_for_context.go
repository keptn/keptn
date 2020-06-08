// Code generated by go-swagger; DO NOT EDIT.

package remediation

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// GetRemediationsForContextHandlerFunc turns a function with the right signature into a get remediations for context handler
type GetRemediationsForContextHandlerFunc func(GetRemediationsForContextParams) middleware.Responder

// Handle executing the request and returning a response
func (fn GetRemediationsForContextHandlerFunc) Handle(params GetRemediationsForContextParams) middleware.Responder {
	return fn(params)
}

// GetRemediationsForContextHandler interface for that can handle valid get remediations for context params
type GetRemediationsForContextHandler interface {
	Handle(GetRemediationsForContextParams) middleware.Responder
}

// NewGetRemediationsForContext creates a new http.Handler for the get remediations for context operation
func NewGetRemediationsForContext(ctx *middleware.Context, handler GetRemediationsForContextHandler) *GetRemediationsForContext {
	return &GetRemediationsForContext{Context: ctx, Handler: handler}
}

/*GetRemediationsForContext swagger:route GET /project/{projectName}/stage/{stageName}/service/{serviceName}/remediation/{keptnContext} remediation getRemediationsForContext

Get open remediations by KeptnContext

*/
type GetRemediationsForContext struct {
	Context *middleware.Context
	Handler GetRemediationsForContextHandler
}

func (o *GetRemediationsForContext) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewGetRemediationsForContextParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
