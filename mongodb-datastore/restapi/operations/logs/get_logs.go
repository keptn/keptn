// Code generated by go-swagger; DO NOT EDIT.

package logs

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"
	"strconv"

	errors "github.com/go-openapi/errors"
	middleware "github.com/go-openapi/runtime/middleware"
	strfmt "github.com/go-openapi/strfmt"
	swag "github.com/go-openapi/swag"

	models "github.com/keptn/keptn/mongodb-datastore/models"
)

// GetLogsHandlerFunc turns a function with the right signature into a get logs handler
type GetLogsHandlerFunc func(GetLogsParams) middleware.Responder

// Handle executing the request and returning a response
func (fn GetLogsHandlerFunc) Handle(params GetLogsParams) middleware.Responder {
	return fn(params)
}

// GetLogsHandler interface for that can handle valid get logs params
type GetLogsHandler interface {
	Handle(GetLogsParams) middleware.Responder
}

// NewGetLogs creates a new http.Handler for the get logs operation
func NewGetLogs(ctx *middleware.Context, handler GetLogsHandler) *GetLogs {
	return &GetLogs{Context: ctx, Handler: handler}
}

/*GetLogs swagger:route GET /log logs getLogs

gets the logs from the datastore

*/
type GetLogs struct {
	Context *middleware.Context
	Handler GetLogsHandler
}

func (o *GetLogs) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewGetLogsParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}

// GetLogsOKBody get logs o k body
// swagger:model GetLogsOKBody
type GetLogsOKBody struct {

	// logs
	Logs []*models.LogEntry `json:"logs"`

	// Pointer to the next page
	NextPageKey string `json:"nextPageKey,omitempty"`

	// Size of the returned page
	PageSize int64 `json:"pageSize,omitempty"`

	// Total number of logs
	TotalCount int64 `json:"totalCount,omitempty"`
}

// Validate validates this get logs o k body
func (o *GetLogsOKBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateLogs(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *GetLogsOKBody) validateLogs(formats strfmt.Registry) error {

	if swag.IsZero(o.Logs) { // not required
		return nil
	}

	for i := 0; i < len(o.Logs); i++ {
		if swag.IsZero(o.Logs[i]) { // not required
			continue
		}

		if o.Logs[i] != nil {
			if err := o.Logs[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("getLogsOK" + "." + "logs" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (o *GetLogsOKBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *GetLogsOKBody) UnmarshalBinary(b []byte) error {
	var res GetLogsOKBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
