// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"net/http"
	"strings"

	errors "github.com/go-openapi/errors"
	loads "github.com/go-openapi/loads"
	runtime "github.com/go-openapi/runtime"
	middleware "github.com/go-openapi/runtime/middleware"
	security "github.com/go-openapi/runtime/security"
	spec "github.com/go-openapi/spec"
	strfmt "github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"

	"github.com/keptn/keptn/mongodb-datastore/restapi/operations/event"
	"github.com/keptn/keptn/mongodb-datastore/restapi/operations/logs"
)

// NewMongodbDatastoreAPI creates a new MongodbDatastore instance
func NewMongodbDatastoreAPI(spec *loads.Document) *MongodbDatastoreAPI {
	return &MongodbDatastoreAPI{
		handlers:            make(map[string]map[string]http.Handler),
		formats:             strfmt.Default,
		defaultConsumes:     "application/json",
		defaultProduces:     "application/json",
		customConsumers:     make(map[string]runtime.Consumer),
		customProducers:     make(map[string]runtime.Producer),
		ServerShutdown:      func() {},
		spec:                spec,
		ServeError:          errors.ServeError,
		BasicAuthenticator:  security.BasicAuth,
		APIKeyAuthenticator: security.APIKeyAuth,
		BearerAuthenticator: security.BearerAuth,
		JSONConsumer:        runtime.JSONConsumer(),
		JSONProducer:        runtime.JSONProducer(),
		EventGetEventHandler: event.GetEventHandlerFunc(func(params event.GetEventParams) middleware.Responder {
			return middleware.NotImplemented("operation EventGetEvent has not yet been implemented")
		}),
		EventGetEventsHandler: event.GetEventsHandlerFunc(func(params event.GetEventsParams) middleware.Responder {
			return middleware.NotImplemented("operation EventGetEvents has not yet been implemented")
		}),
		LogsGetLogsHandler: logs.GetLogsHandlerFunc(func(params logs.GetLogsParams) middleware.Responder {
			return middleware.NotImplemented("operation LogsGetLogs has not yet been implemented")
		}),
		EventGetNewArtifactEventsHandler: event.GetNewArtifactEventsHandlerFunc(func(params event.GetNewArtifactEventsParams) middleware.Responder {
			return middleware.NotImplemented("operation EventGetNewArtifactEvents has not yet been implemented")
		}),
		EventSaveEventHandler: event.SaveEventHandlerFunc(func(params event.SaveEventParams) middleware.Responder {
			return middleware.NotImplemented("operation EventSaveEvent has not yet been implemented")
		}),
		LogsSaveLogHandler: logs.SaveLogHandlerFunc(func(params logs.SaveLogParams) middleware.Responder {
			return middleware.NotImplemented("operation LogsSaveLog has not yet been implemented")
		}),
	}
}

/*MongodbDatastoreAPI the mongodb datastore API */
type MongodbDatastoreAPI struct {
	spec            *loads.Document
	context         *middleware.Context
	handlers        map[string]map[string]http.Handler
	formats         strfmt.Registry
	customConsumers map[string]runtime.Consumer
	customProducers map[string]runtime.Producer
	defaultConsumes string
	defaultProduces string
	Middleware      func(middleware.Builder) http.Handler

	// BasicAuthenticator generates a runtime.Authenticator from the supplied basic auth function.
	// It has a default implementation in the security package, however you can replace it for your particular usage.
	BasicAuthenticator func(security.UserPassAuthentication) runtime.Authenticator
	// APIKeyAuthenticator generates a runtime.Authenticator from the supplied token auth function.
	// It has a default implementation in the security package, however you can replace it for your particular usage.
	APIKeyAuthenticator func(string, string, security.TokenAuthentication) runtime.Authenticator
	// BearerAuthenticator generates a runtime.Authenticator from the supplied bearer token auth function.
	// It has a default implementation in the security package, however you can replace it for your particular usage.
	BearerAuthenticator func(string, security.ScopedTokenAuthentication) runtime.Authenticator

	// JSONConsumer registers a consumer for a "application/cloudevents+json" mime type
	JSONConsumer runtime.Consumer

	// JSONProducer registers a producer for a "application/cloudevents+json" mime type
	JSONProducer runtime.Producer

	// EventGetEventHandler sets the operation handler for the get event operation
	EventGetEventHandler event.GetEventHandler
	// EventGetEventsHandler sets the operation handler for the get events operation
	EventGetEventsHandler event.GetEventsHandler
	// LogsGetLogsHandler sets the operation handler for the get logs operation
	LogsGetLogsHandler logs.GetLogsHandler
	// EventGetNewArtifactEventsHandler sets the operation handler for the get new artifact events operation
	EventGetNewArtifactEventsHandler event.GetNewArtifactEventsHandler
	// EventSaveEventHandler sets the operation handler for the save event operation
	EventSaveEventHandler event.SaveEventHandler
	// LogsSaveLogHandler sets the operation handler for the save log operation
	LogsSaveLogHandler logs.SaveLogHandler

	// ServeError is called when an error is received, there is a default handler
	// but you can set your own with this
	ServeError func(http.ResponseWriter, *http.Request, error)

	// ServerShutdown is called when the HTTP(S) server is shut down and done
	// handling all active connections and does not accept connections any more
	ServerShutdown func()

	// Custom command line argument groups with their descriptions
	CommandLineOptionsGroups []swag.CommandLineOptionsGroup

	// User defined logger function.
	Logger func(string, ...interface{})
}

// SetDefaultProduces sets the default produces media type
func (o *MongodbDatastoreAPI) SetDefaultProduces(mediaType string) {
	o.defaultProduces = mediaType
}

// SetDefaultConsumes returns the default consumes media type
func (o *MongodbDatastoreAPI) SetDefaultConsumes(mediaType string) {
	o.defaultConsumes = mediaType
}

// SetSpec sets a spec that will be served for the clients.
func (o *MongodbDatastoreAPI) SetSpec(spec *loads.Document) {
	o.spec = spec
}

// DefaultProduces returns the default produces media type
func (o *MongodbDatastoreAPI) DefaultProduces() string {
	return o.defaultProduces
}

// DefaultConsumes returns the default consumes media type
func (o *MongodbDatastoreAPI) DefaultConsumes() string {
	return o.defaultConsumes
}

// Formats returns the registered string formats
func (o *MongodbDatastoreAPI) Formats() strfmt.Registry {
	return o.formats
}

// RegisterFormat registers a custom format validator
func (o *MongodbDatastoreAPI) RegisterFormat(name string, format strfmt.Format, validator strfmt.Validator) {
	o.formats.Add(name, format, validator)
}

// Validate validates the registrations in the MongodbDatastoreAPI
func (o *MongodbDatastoreAPI) Validate() error {
	var unregistered []string

	if o.JSONConsumer == nil {
		unregistered = append(unregistered, "JSONConsumer")
	}

	if o.JSONProducer == nil {
		unregistered = append(unregistered, "JSONProducer")
	}

	if o.EventGetEventHandler == nil {
		unregistered = append(unregistered, "event.GetEventHandler")
	}

	if o.EventGetEventsHandler == nil {
		unregistered = append(unregistered, "event.GetEventsHandler")
	}

	if o.LogsGetLogsHandler == nil {
		unregistered = append(unregistered, "logs.GetLogsHandler")
	}

	if o.EventGetNewArtifactEventsHandler == nil {
		unregistered = append(unregistered, "event.GetNewArtifactEventsHandler")
	}

	if o.EventSaveEventHandler == nil {
		unregistered = append(unregistered, "event.SaveEventHandler")
	}

	if o.LogsSaveLogHandler == nil {
		unregistered = append(unregistered, "logs.SaveLogHandler")
	}

	if len(unregistered) > 0 {
		return fmt.Errorf("missing registration: %s", strings.Join(unregistered, ", "))
	}

	return nil
}

// ServeErrorFor gets a error handler for a given operation id
func (o *MongodbDatastoreAPI) ServeErrorFor(operationID string) func(http.ResponseWriter, *http.Request, error) {
	return o.ServeError
}

// AuthenticatorsFor gets the authenticators for the specified security schemes
func (o *MongodbDatastoreAPI) AuthenticatorsFor(schemes map[string]spec.SecurityScheme) map[string]runtime.Authenticator {

	return nil

}

// Authorizer returns the registered authorizer
func (o *MongodbDatastoreAPI) Authorizer() runtime.Authorizer {

	return nil

}

// ConsumersFor gets the consumers for the specified media types
func (o *MongodbDatastoreAPI) ConsumersFor(mediaTypes []string) map[string]runtime.Consumer {

	result := make(map[string]runtime.Consumer)
	for _, mt := range mediaTypes {
		switch mt {

		case "application/cloudevents+json":
			result["application/cloudevents+json"] = o.JSONConsumer

		case "application/json":
			result["application/json"] = o.JSONConsumer

		}

		if c, ok := o.customConsumers[mt]; ok {
			result[mt] = c
		}
	}
	return result

}

// ProducersFor gets the producers for the specified media types
func (o *MongodbDatastoreAPI) ProducersFor(mediaTypes []string) map[string]runtime.Producer {

	result := make(map[string]runtime.Producer)
	for _, mt := range mediaTypes {
		switch mt {

		case "application/cloudevents+json":
			result["application/cloudevents+json"] = o.JSONProducer

		case "application/json":
			result["application/json"] = o.JSONProducer

		}

		if p, ok := o.customProducers[mt]; ok {
			result[mt] = p
		}
	}
	return result

}

// HandlerFor gets a http.Handler for the provided operation method and path
func (o *MongodbDatastoreAPI) HandlerFor(method, path string) (http.Handler, bool) {
	if o.handlers == nil {
		return nil, false
	}
	um := strings.ToUpper(method)
	if _, ok := o.handlers[um]; !ok {
		return nil, false
	}
	if path == "/" {
		path = ""
	}
	h, ok := o.handlers[um][path]
	return h, ok
}

// Context returns the middleware context for the mongodb datastore API
func (o *MongodbDatastoreAPI) Context() *middleware.Context {
	if o.context == nil {
		o.context = middleware.NewRoutableContext(o.spec, o, nil)
	}

	return o.context
}

func (o *MongodbDatastoreAPI) initHandlerCache() {
	o.Context() // don't care about the result, just that the initialization happened

	if o.handlers == nil {
		o.handlers = make(map[string]map[string]http.Handler)
	}

	if o.handlers["GET"] == nil {
		o.handlers["GET"] = make(map[string]http.Handler)
	}
	o.handlers["GET"]["/events/{id}"] = event.NewGetEvent(o.context, o.EventGetEventHandler)

	if o.handlers["GET"] == nil {
		o.handlers["GET"] = make(map[string]http.Handler)
	}
	o.handlers["GET"]["/events"] = event.NewGetEvents(o.context, o.EventGetEventsHandler)

	if o.handlers["GET"] == nil {
		o.handlers["GET"] = make(map[string]http.Handler)
	}
	o.handlers["GET"]["/logs"] = logs.NewGetLogs(o.context, o.LogsGetLogsHandler)

	if o.handlers["GET"] == nil {
		o.handlers["GET"] = make(map[string]http.Handler)
	}
	o.handlers["GET"]["/events/newartifact"] = event.NewGetNewArtifactEvents(o.context, o.EventGetNewArtifactEventsHandler)

	if o.handlers["POST"] == nil {
		o.handlers["POST"] = make(map[string]http.Handler)
	}
	o.handlers["POST"]["/events"] = event.NewSaveEvent(o.context, o.EventSaveEventHandler)

	if o.handlers["POST"] == nil {
		o.handlers["POST"] = make(map[string]http.Handler)
	}
	o.handlers["POST"]["/logs"] = logs.NewSaveLog(o.context, o.LogsSaveLogHandler)

}

// Serve creates a http handler to serve the API over HTTP
// can be used directly in http.ListenAndServe(":8000", api.Serve(nil))
func (o *MongodbDatastoreAPI) Serve(builder middleware.Builder) http.Handler {
	o.Init()

	if o.Middleware != nil {
		return o.Middleware(builder)
	}
	return o.context.APIHandler(builder)
}

// Init allows you to just initialize the handler cache, you can then recompose the middleware as you see fit
func (o *MongodbDatastoreAPI) Init() {
	if len(o.handlers) == 0 {
		o.initHandlerCache()
	}
}

// RegisterConsumer allows you to add (or override) a consumer for a media type.
func (o *MongodbDatastoreAPI) RegisterConsumer(mediaType string, consumer runtime.Consumer) {
	o.customConsumers[mediaType] = consumer
}

// RegisterProducer allows you to add (or override) a producer for a media type.
func (o *MongodbDatastoreAPI) RegisterProducer(mediaType string, producer runtime.Producer) {
	o.customProducers[mediaType] = producer
}
