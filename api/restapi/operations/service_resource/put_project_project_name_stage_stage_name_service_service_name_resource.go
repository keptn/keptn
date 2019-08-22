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

	models "github.com/keptn/keptn/api/models"
)

// PutProjectProjectNameStageStageNameServiceServiceNameResourceHandlerFunc turns a function with the right signature into a put project project name stage stage name service service name resource handler
type PutProjectProjectNameStageStageNameServiceServiceNameResourceHandlerFunc func(PutProjectProjectNameStageStageNameServiceServiceNameResourceParams, *models.Principal) middleware.Responder

// Handle executing the request and returning a response
func (fn PutProjectProjectNameStageStageNameServiceServiceNameResourceHandlerFunc) Handle(params PutProjectProjectNameStageStageNameServiceServiceNameResourceParams, principal *models.Principal) middleware.Responder {
	return fn(params, principal)
}

// PutProjectProjectNameStageStageNameServiceServiceNameResourceHandler interface for that can handle valid put project project name stage stage name service service name resource params
type PutProjectProjectNameStageStageNameServiceServiceNameResourceHandler interface {
	Handle(PutProjectProjectNameStageStageNameServiceServiceNameResourceParams, *models.Principal) middleware.Responder
}

// NewPutProjectProjectNameStageStageNameServiceServiceNameResource creates a new http.Handler for the put project project name stage stage name service service name resource operation
func NewPutProjectProjectNameStageStageNameServiceServiceNameResource(ctx *middleware.Context, handler PutProjectProjectNameStageStageNameServiceServiceNameResourceHandler) *PutProjectProjectNameStageStageNameServiceServiceNameResource {
	return &PutProjectProjectNameStageStageNameServiceServiceNameResource{Context: ctx, Handler: handler}
}

/*PutProjectProjectNameStageStageNameServiceServiceNameResource swagger:route PUT /project/{projectName}/stage/{stageName}/service/{serviceName}/resource Service Resource putProjectProjectNameStageStageNameServiceServiceNameResource

Update list of service resources

*/
type PutProjectProjectNameStageStageNameServiceServiceNameResource struct {
	Context *middleware.Context
	Handler PutProjectProjectNameStageStageNameServiceServiceNameResourceHandler
}

func (o *PutProjectProjectNameStageStageNameServiceServiceNameResource) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewPutProjectProjectNameStageStageNameServiceServiceNameResourceParams()

	uprinc, aCtx, err := o.Context.Authorize(r, route)
	if err != nil {
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}
	if aCtx != nil {
		r = aCtx
	}
	var principal *models.Principal
	if uprinc != nil {
		principal = uprinc.(*models.Principal) // this is really a models.Principal, I promise
	}

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params, principal) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}

// PutProjectProjectNameStageStageNameServiceServiceNameResourceBody put project project name stage stage name service service name resource body
// swagger:model PutProjectProjectNameStageStageNameServiceServiceNameResourceBody
type PutProjectProjectNameStageStageNameServiceServiceNameResourceBody struct {

	// resources
	Resources []*models.Resource `json:"resources"`
}

// Validate validates this put project project name stage stage name service service name resource body
func (o *PutProjectProjectNameStageStageNameServiceServiceNameResourceBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateResources(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *PutProjectProjectNameStageStageNameServiceServiceNameResourceBody) validateResources(formats strfmt.Registry) error {

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
func (o *PutProjectProjectNameStageStageNameServiceServiceNameResourceBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *PutProjectProjectNameStageStageNameServiceServiceNameResourceBody) UnmarshalBinary(b []byte) error {
	var res PutProjectProjectNameStageStageNameServiceServiceNameResourceBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
