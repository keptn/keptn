// Code generated by go-swagger; DO NOT EDIT.

package service_default_resource

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"
	"strconv"

	errors "github.com/go-openapi/errors"
	middleware "github.com/go-openapi/runtime/middleware"
	strfmt "github.com/go-openapi/strfmt"
	swag "github.com/go-openapi/swag"

	"github.com/keptn/keptn/configuration-service/models"
)

// PostProjectProjectNameServiceServiceNameResourceHandlerFunc turns a function with the right signature into a post project project name service service name resource handler
type PostProjectProjectNameServiceServiceNameResourceHandlerFunc func(PostProjectProjectNameServiceServiceNameResourceParams) middleware.Responder

// Handle executing the request and returning a response
func (fn PostProjectProjectNameServiceServiceNameResourceHandlerFunc) Handle(params PostProjectProjectNameServiceServiceNameResourceParams) middleware.Responder {
	return fn(params)
}

// PostProjectProjectNameServiceServiceNameResourceHandler interface for that can handle valid post project project name service service name resource params
type PostProjectProjectNameServiceServiceNameResourceHandler interface {
	Handle(PostProjectProjectNameServiceServiceNameResourceParams) middleware.Responder
}

// NewPostProjectProjectNameServiceServiceNameResource creates a new http.Handler for the post project project name service service name resource operation
func NewPostProjectProjectNameServiceServiceNameResource(ctx *middleware.Context, handler PostProjectProjectNameServiceServiceNameResourceHandler) *PostProjectProjectNameServiceServiceNameResource {
	return &PostProjectProjectNameServiceServiceNameResource{Context: ctx, Handler: handler}
}

/*PostProjectProjectNameServiceServiceNameResource swagger:route POST /project/{projectName}/service/{serviceName}/resource Service Default Resource postProjectProjectNameServiceServiceNameResource

Create list of default resources for the service used in all stages

*/
type PostProjectProjectNameServiceServiceNameResource struct {
	Context *middleware.Context
	Handler PostProjectProjectNameServiceServiceNameResourceHandler
}

func (o *PostProjectProjectNameServiceServiceNameResource) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewPostProjectProjectNameServiceServiceNameResourceParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}

// PostProjectProjectNameServiceServiceNameResourceBody post project project name service service name resource body
// swagger:model PostProjectProjectNameServiceServiceNameResourceBody
type PostProjectProjectNameServiceServiceNameResourceBody struct {

	// resources
	Resources []*models.Resource `json:"resources"`
}

// Validate validates this post project project name service service name resource body
func (o *PostProjectProjectNameServiceServiceNameResourceBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateResources(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *PostProjectProjectNameServiceServiceNameResourceBody) validateResources(formats strfmt.Registry) error {

	if swag.IsZero(o.Resources) { // not required
		return nil
	}

	for i := 0; i < len(o.Resources); i++ {
		if swag.IsZero(o.Resources[i]) { // not required
			continue
		}

		if o.Resources[i] != nil {
			if err := o.Resources[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("resources" + "." + "resources" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (o *PostProjectProjectNameServiceServiceNameResourceBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *PostProjectProjectNameServiceServiceNameResourceBody) UnmarshalBinary(b []byte) error {
	var res PostProjectProjectNameServiceServiceNameResourceBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
