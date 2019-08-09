// Code generated by go-swagger; DO NOT EDIT.

package service_resource

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"
	"strconv"

	errors "github.com/go-openapi/errors"
	middleware "github.com/go-openapi/runtime/middleware"
	strfmt "github.com/go-openapi/strfmt"
	swag "github.com/go-openapi/swag"

	models "github.com/keptn/keptn/configuration-service/models"
)

// PostProjectProjectNameStageStageNameServiceServiceNameResourceHandlerFunc turns a function with the right signature into a post project project name stage stage name service service name resource handler
type PostProjectProjectNameStageStageNameServiceServiceNameResourceHandlerFunc func(PostProjectProjectNameStageStageNameServiceServiceNameResourceParams) middleware.Responder

// Handle executing the request and returning a response
func (fn PostProjectProjectNameStageStageNameServiceServiceNameResourceHandlerFunc) Handle(params PostProjectProjectNameStageStageNameServiceServiceNameResourceParams) middleware.Responder {
	return fn(params)
}

// PostProjectProjectNameStageStageNameServiceServiceNameResourceHandler interface for that can handle valid post project project name stage stage name service service name resource params
type PostProjectProjectNameStageStageNameServiceServiceNameResourceHandler interface {
	Handle(PostProjectProjectNameStageStageNameServiceServiceNameResourceParams) middleware.Responder
}

// NewPostProjectProjectNameStageStageNameServiceServiceNameResource creates a new http.Handler for the post project project name stage stage name service service name resource operation
func NewPostProjectProjectNameStageStageNameServiceServiceNameResource(ctx *middleware.Context, handler PostProjectProjectNameStageStageNameServiceServiceNameResourceHandler) *PostProjectProjectNameStageStageNameServiceServiceNameResource {
	return &PostProjectProjectNameStageStageNameServiceServiceNameResource{Context: ctx, Handler: handler}
}

/*PostProjectProjectNameStageStageNameServiceServiceNameResource swagger:route POST /project/{projectName}/stage/{stageName}/service/{serviceName}/resource Service Resource postProjectProjectNameStageStageNameServiceServiceNameResource

Create list of new resources for the service

*/
type PostProjectProjectNameStageStageNameServiceServiceNameResource struct {
	Context *middleware.Context
	Handler PostProjectProjectNameStageStageNameServiceServiceNameResourceHandler
}

func (o *PostProjectProjectNameStageStageNameServiceServiceNameResource) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewPostProjectProjectNameStageStageNameServiceServiceNameResourceParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}

// PostProjectProjectNameStageStageNameServiceServiceNameResourceBody post project project name stage stage name service service name resource body
// swagger:model PostProjectProjectNameStageStageNameServiceServiceNameResourceBody
type PostProjectProjectNameStageStageNameServiceServiceNameResourceBody struct {

	// resources
	Resources []*models.Resource `json:"resources"`
}

// Validate validates this post project project name stage stage name service service name resource body
func (o *PostProjectProjectNameStageStageNameServiceServiceNameResourceBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateResources(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *PostProjectProjectNameStageStageNameServiceServiceNameResourceBody) validateResources(formats strfmt.Registry) error {

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
func (o *PostProjectProjectNameStageStageNameServiceServiceNameResourceBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *PostProjectProjectNameStageStageNameServiceServiceNameResourceBody) UnmarshalBinary(b []byte) error {
	var res PostProjectProjectNameStageStageNameServiceServiceNameResourceBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
