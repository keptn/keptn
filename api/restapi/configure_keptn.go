// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"fmt"
	"github.com/benbjohnson/clock"
	"github.com/kelseyhightower/envconfig"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"

	"github.com/keptn/keptn/api/handlers"
	custommiddleware "github.com/keptn/keptn/api/middleware"
	"github.com/keptn/keptn/api/models"
	"github.com/keptn/keptn/api/restapi/operations"
	"github.com/keptn/keptn/api/restapi/operations/auth"
	"github.com/keptn/keptn/api/restapi/operations/event"
	"github.com/keptn/keptn/api/restapi/operations/metadata"
)

//go:generate swagger generate server --target ../../api --name Keptn --spec ../swagger.yaml --principal models.Principal

type EnvConfig struct {
	AuthRequestsPerSecond float64 `envconfig:"MAX_AUTH_REQUESTS_PER_SECOND" default:"1"`
	AuthRequestMaxBurst   int     `envconfig:"MAX_AUTH_REQUESTS_BURST" default:"2"`
}

func configureFlags(api *operations.KeptnAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func getEnvConfig() (*EnvConfig, error) {
	env := &EnvConfig{}
	if err := envconfig.Process("", env); err != nil {
		return nil, err
	}
	return env, nil
}

func configureAPI(api *operations.KeptnAPI) http.Handler {
	env, err := getEnvConfig()
	if err != nil {
		log.WithError(err).Error("Failed to process env var")
		os.Exit(1)
	}

	/// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	// Applies when the "x-token" header is set
	tokenValidator := &custommiddleware.BasicTokenValidator{}
	api.KeyAuth = tokenValidator.ValidateToken

	// Set your custom authorizer if needed. Default one is security.Authorized()
	// Expected interface runtime.Authorizer
	//
	// Example:
	// api.APIAuthorizer = security.Authorized()
	api.AuthAuthHandler = auth.AuthHandlerFunc(func(params auth.AuthParams, principal *models.Principal) middleware.Responder {
		return auth.NewAuthOK()
	})

	api.EventPostEventHandler = event.PostEventHandlerFunc(handlers.PostEventHandlerFunc)
	//api.EventGetEventHandler = event.GetEventHandlerFunc(handlers.GetEventHandlerFunc)

	// Metadata endpoint
	api.MetadataMetadataHandler = metadata.MetadataHandlerFunc(handlers.GetMetadataHandlerFunc)

	//api.EvaluationTriggerEvaluationHandler = evaluation.TriggerEvaluationHandlerFunc(handlers.TriggerEvaluationHandlerFunc)

	rateLimiter := custommiddleware.NewRateLimiter(env.AuthRequestsPerSecond, env.AuthRequestMaxBurst, tokenValidator, clock.New())
	api.AddMiddlewareFor(http.MethodPost, "/auth", rateLimiter.Handle)
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
			editedSwagger := strings.Replace(string(input), "basePath: /api/v1",
				"basePath: "+prefixPath+"/api/v1", -1)
			err = ioutil.WriteFile("swagger-ui/swagger.yaml", []byte(editedSwagger), 0644)
			if err != nil {
				fmt.Println("Failed to write edited swagger.yaml")
			}
		} else {
			fmt.Println("Failed to set basePath in swagger.yaml")
		}

		// Set the prefix-path in the index.html
		input, err = ioutil.ReadFile("swagger-ui/index.html")
		if err == nil {
			editedSwagger := strings.Replace(string(input), "const prefixPath = \"\"",
				"const prefixPath = \""+prefixPath+"\"", -1)
			err = ioutil.WriteFile("swagger-ui/index.html", []byte(editedSwagger), 0644)
			if err != nil {
				fmt.Println("Failed to write edited index.html")
			}
		} else {
			fmt.Println("Failed to set basePath in index.html")
		}
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Index(r.URL.Path, "/swagger-ui/") == 0 {
			http.StripPrefix("/swagger-ui/", http.FileServer(http.Dir("swagger-ui"))).ServeHTTP(w, r)
			return
		}
		if strings.Index(r.URL.Path, "/health") == 0 {
			w.WriteHeader(http.StatusOK)
			return
		}

		handler.ServeHTTP(w, r)
	})
}
