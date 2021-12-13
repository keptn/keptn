// Code generated by go-swagger; DO NOT EDIT.

package stage_resource

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"

	"github.com/keptn/keptn/configuration-service/models"
)

// PutProjectProjectNameStageStageNameResourceHandlerFunc turns a function with the right signature into a put project project name stage stage name resource handler
type PutProjectProjectNameStageStageNameResourceHandlerFunc func(PutProjectProjectNameStageStageNameResourceParams) middleware.Responder

// Handle executing the request and returning a response
func (fn PutProjectProjectNameStageStageNameResourceHandlerFunc) Handle(params PutProjectProjectNameStageStageNameResourceParams) middleware.Responder {
	return fn(params)
}

// PutProjectProjectNameStageStageNameResourceHandler interface for that can handle valid put project project name stage stage name resource params
type PutProjectProjectNameStageStageNameResourceHandler interface {
	Handle(PutProjectProjectNameStageStageNameResourceParams) middleware.Responder
}

// NewPutProjectProjectNameStageStageNameResource creates a new http.Handler for the put project project name stage stage name resource operation
func NewPutProjectProjectNameStageStageNameResource(ctx *middleware.Context, handler PutProjectProjectNameStageStageNameResourceHandler) *PutProjectProjectNameStageStageNameResource {
	return &PutProjectProjectNameStageStageNameResource{Context: ctx, Handler: handler}
}

/* PutProjectProjectNameStageStageNameResource swagger:route PUT /project/{projectName}/stage/{stageName}/resource Stage Resource putProjectProjectNameStageStageNameResource

Update list of stage resources

*/
type PutProjectProjectNameStageStageNameResource struct {
	Context *middleware.Context
	Handler PutProjectProjectNameStageStageNameResourceHandler
}

func (o *PutProjectProjectNameStageStageNameResource) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewPutProjectProjectNameStageStageNameResourceParams()
	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}

// PutProjectProjectNameStageStageNameResourceBody put project project name stage stage name resource body
//
// swagger:model PutProjectProjectNameStageStageNameResourceBody
type PutProjectProjectNameStageStageNameResourceBody struct {

	// resources
	Resources []*models.Resource `json:"resources"`
}

// Validate validates this put project project name stage stage name resource body
func (o *PutProjectProjectNameStageStageNameResourceBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateResources(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *PutProjectProjectNameStageStageNameResourceBody) validateResources(formats strfmt.Registry) error {
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
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("resources" + "." + "resources" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// ContextValidate validate this put project project name stage stage name resource body based on the context it is used
func (o *PutProjectProjectNameStageStageNameResourceBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := o.contextValidateResources(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *PutProjectProjectNameStageStageNameResourceBody) contextValidateResources(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(o.Resources); i++ {

		if o.Resources[i] != nil {
			if err := o.Resources[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("resources" + "." + "resources" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("resources" + "." + "resources" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (o *PutProjectProjectNameStageStageNameResourceBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *PutProjectProjectNameStageStageNameResourceBody) UnmarshalBinary(b []byte) error {
	var res PutProjectProjectNameStageStageNameResourceBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
