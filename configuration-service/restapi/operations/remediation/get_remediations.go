// Code generated by go-swagger; DO NOT EDIT.

package remediation

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// GetRemediationsHandlerFunc turns a function with the right signature into a get remediations handler
type GetRemediationsHandlerFunc func(GetRemediationsParams) middleware.Responder

// Handle executing the request and returning a response
func (fn GetRemediationsHandlerFunc) Handle(params GetRemediationsParams) middleware.Responder {
	return fn(params)
}

// GetRemediationsHandler interface for that can handle valid get remediations params
type GetRemediationsHandler interface {
	Handle(GetRemediationsParams) middleware.Responder
}

// NewGetRemediations creates a new http.Handler for the get remediations operation
func NewGetRemediations(ctx *middleware.Context, handler GetRemediationsHandler) *GetRemediations {
	return &GetRemediations{Context: ctx, Handler: handler}
}

/*GetRemediations swagger:route GET /project/{projectName}/stage/{stageName}/service/{serviceName}/remediation remediation getRemediations

Get all open remediations

*/
type GetRemediations struct {
	Context *middleware.Context
	Handler GetRemediationsHandler
}

func (o *GetRemediations) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewGetRemediationsParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
