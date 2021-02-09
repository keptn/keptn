// This file is safe to edit. Once it exists it will not be overwritten
package restapi

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	errors "github.com/go-openapi/errors"
	runtime "github.com/go-openapi/runtime"
	middleware "github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"

	keptnapi "github.com/keptn/go-utils/pkg/api/utils"
	"github.com/keptn/keptn/mongodb-datastore/handlers"
	"github.com/keptn/keptn/mongodb-datastore/models"
	"github.com/keptn/keptn/mongodb-datastore/restapi/operations"
	"github.com/keptn/keptn/mongodb-datastore/restapi/operations/event"
)

//go:generate swagger generate server --target ../../mongodb-datastore --name mongodb-datastore --spec ../swagger.yaml

func configureFlags(api *operations.MongodbDatastoreAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.MongodbDatastoreAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	api.EventSaveEventHandler = event.SaveEventHandlerFunc(func(params event.SaveEventParams) middleware.Responder {
		if err := handlers.ProcessEvent(params.Body); err != nil {
			return event.NewSaveEventDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String(err.Error())})
		}
		return event.NewSaveEventCreated()
	})

	api.EventGetEventsHandler = event.GetEventsHandlerFunc(func(params event.GetEventsParams) middleware.Responder {
		events, err := handlers.GetEvents(params)
		if err != nil {
			return event.NewGetEventsDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String(err.Error())})
		}
		return event.NewGetEventsOK().WithPayload(events)
	})

	api.EventGetEventsByTypeHandler = event.GetEventsByTypeHandlerFunc(func(params event.GetEventsByTypeParams) middleware.Responder {
		events, err := handlers.GetEventsByType(params)
		if err != nil {
			return event.NewGetEventsDefault(500).WithPayload(&models.Error{Code: 500, Message: swag.String(err.Error())})
		}
		return event.NewGetEventsByTypeOK().WithPayload(events)
	})

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix"
func configureServer(s *http.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics
func setupGlobalMiddleware(handler http.Handler) http.Handler {

	prefixPath := os.Getenv("PREFIX_PATH")
	if len(prefixPath) > 0 {
		// Set the prefix-path in the swagger.yaml
		input, err := ioutil.ReadFile("swagger-ui/swagger.yaml")
		if err == nil {
			editedSwagger := strings.Replace(string(input), "basePath: /api/mongodb-datastore",
				"basePath: "+prefixPath+"/api/mongodb-datastore", -1)
			err = ioutil.WriteFile("swagger-ui/swagger.yaml", []byte(editedSwagger), 0644)
			if err != nil {
				fmt.Println("Failed to write edited swagger.yaml")
			}
		} else {
			fmt.Println("Failed to set basePath in swagger.yaml")
		}
	}

	go keptnapi.RunHealthEndpoint("10999")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Serving ./swagger-ui/
		if strings.Index(r.URL.Path, "/swagger-ui/") == 0 {
			pathToSwaggerUI := "swagger-ui"
			// in case of local execution, the dir is stored in a parent folder
			if _, err := os.Stat(pathToSwaggerUI); os.IsNotExist(err) {
				pathToSwaggerUI = "../../swagger-ui"
			}
			http.StripPrefix("/swagger-ui/", http.FileServer(http.Dir(pathToSwaggerUI))).ServeHTTP(w, r)
			return
		}
		handler.ServeHTTP(w, r)
	})
}
