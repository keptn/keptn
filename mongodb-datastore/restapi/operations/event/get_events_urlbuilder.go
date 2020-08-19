// Code generated by go-swagger; DO NOT EDIT.

package event

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"errors"
	"net/url"
	golangswaggerpaths "path"

	"github.com/go-openapi/swag"
)

// GetEventsURL generates an URL for the get events operation
type GetEventsURL struct {
	EventID      *string
	FromTime     *string
	KeptnContext *string
	NextPageKey  *string
	PageSize     *int64
	Project      *string
	Result       *string
	Root         *string
	Service      *string
	Source       *string
	Stage        *string
	Type         *string

	_basePath string
	// avoid unkeyed usage
	_ struct{}
}

// WithBasePath sets the base path for this url builder, only required when it's different from the
// base path specified in the swagger spec.
// When the value of the base path is an empty string
func (o *GetEventsURL) WithBasePath(bp string) *GetEventsURL {
	o.SetBasePath(bp)
	return o
}

// SetBasePath sets the base path for this url builder, only required when it's different from the
// base path specified in the swagger spec.
// When the value of the base path is an empty string
func (o *GetEventsURL) SetBasePath(bp string) {
	o._basePath = bp
}

// Build a url path and query string
func (o *GetEventsURL) Build() (*url.URL, error) {
	var _result url.URL

	var _path = "/event"

	_basePath := o._basePath
	if _basePath == "" {
		_basePath = "/"
	}
	_result.Path = golangswaggerpaths.Join(_basePath, _path)

	qs := make(url.Values)

	var eventIDQ string
	if o.EventID != nil {
		eventIDQ = *o.EventID
	}
	if eventIDQ != "" {
		qs.Set("eventID", eventIDQ)
	}

	var fromTimeQ string
	if o.FromTime != nil {
		fromTimeQ = *o.FromTime
	}
	if fromTimeQ != "" {
		qs.Set("fromTime", fromTimeQ)
	}

	var keptnContextQ string
	if o.KeptnContext != nil {
		keptnContextQ = *o.KeptnContext
	}
	if keptnContextQ != "" {
		qs.Set("keptnContext", keptnContextQ)
	}

	var nextPageKeyQ string
	if o.NextPageKey != nil {
		nextPageKeyQ = *o.NextPageKey
	}
	if nextPageKeyQ != "" {
		qs.Set("nextPageKey", nextPageKeyQ)
	}

	var pageSizeQ string
	if o.PageSize != nil {
		pageSizeQ = swag.FormatInt64(*o.PageSize)
	}
	if pageSizeQ != "" {
		qs.Set("pageSize", pageSizeQ)
	}

	var projectQ string
	if o.Project != nil {
		projectQ = *o.Project
	}
	if projectQ != "" {
		qs.Set("project", projectQ)
	}

	var resultQ string
	if o.Result != nil {
		resultQ = *o.Result
	}
	if resultQ != "" {
		qs.Set("result", resultQ)
	}

	var rootQ string
	if o.Root != nil {
		rootQ = *o.Root
	}
	if rootQ != "" {
		qs.Set("root", rootQ)
	}

	var serviceQ string
	if o.Service != nil {
		serviceQ = *o.Service
	}
	if serviceQ != "" {
		qs.Set("service", serviceQ)
	}

	var sourceQ string
	if o.Source != nil {
		sourceQ = *o.Source
	}
	if sourceQ != "" {
		qs.Set("source", sourceQ)
	}

	var stageQ string
	if o.Stage != nil {
		stageQ = *o.Stage
	}
	if stageQ != "" {
		qs.Set("stage", stageQ)
	}

	var typeVarQ string
	if o.Type != nil {
		typeVarQ = *o.Type
	}
	if typeVarQ != "" {
		qs.Set("type", typeVarQ)
	}

	_result.RawQuery = qs.Encode()

	return &_result, nil
}

// Must is a helper function to panic when the url builder returns an error
func (o *GetEventsURL) Must(u *url.URL, err error) *url.URL {
	if err != nil {
		panic(err)
	}
	if u == nil {
		panic("url can't be nil")
	}
	return u
}

// String returns the string representation of the path with query string
func (o *GetEventsURL) String() string {
	return o.Must(o.Build()).String()
}

// BuildFull builds a full url with scheme, host, path and query string
func (o *GetEventsURL) BuildFull(scheme, host string) (*url.URL, error) {
	if scheme == "" {
		return nil, errors.New("scheme is required for a full url on GetEventsURL")
	}
	if host == "" {
		return nil, errors.New("host is required for a full url on GetEventsURL")
	}

	base, err := o.Build()
	if err != nil {
		return nil, err
	}

	base.Scheme = scheme
	base.Host = host
	return base, nil
}

// StringFull returns the string representation of a complete url
func (o *GetEventsURL) StringFull(scheme, host string) string {
	return o.Must(o.BuildFull(scheme, host)).String()
}
