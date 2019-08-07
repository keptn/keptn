// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"net/http"

	errors "github.com/go-openapi/errors"
	runtime "github.com/go-openapi/runtime"
	middleware "github.com/go-openapi/runtime/middleware"

	"github.com/keptn/keptn/configuration-service/restapi/operations"
	"github.com/keptn/keptn/configuration-service/restapi/operations/project"
	"github.com/keptn/keptn/configuration-service/restapi/operations/project_resource"
	"github.com/keptn/keptn/configuration-service/restapi/operations/service"
	"github.com/keptn/keptn/configuration-service/restapi/operations/service_default_resource"
	"github.com/keptn/keptn/configuration-service/restapi/operations/service_resource"
	"github.com/keptn/keptn/configuration-service/restapi/operations/stage"
	"github.com/keptn/keptn/configuration-service/restapi/operations/stage_resource"
)

//go:generate swagger generate server --target ../../configuration-service --name ConfigurationService --spec ../swagger.yaml

func configureFlags(api *operations.ConfigurationServiceAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.ConfigurationServiceAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	api.ProjectDeleteProjectProjectNameHandler = project.DeleteProjectProjectNameHandlerFunc(func(params project.DeleteProjectProjectNameParams) middleware.Responder {
		return middleware.NotImplemented("operation project.DeleteProjectProjectName has not yet been implemented")
	})

	api.ProjectResourceDeleteProjectProjectNameResourceResourceURIHandler = project_resource.DeleteProjectProjectNameResourceResourceURIHandlerFunc(func(params project_resource.DeleteProjectProjectNameResourceResourceURIParams) middleware.Responder {
		return middleware.NotImplemented("operation project_resource.DeleteProjectProjectNameResourceResourceURI has not yet been implemented")
	})

	api.ServiceDefaultResourceDeleteProjectProjectNameServiceServiceNameResourceResourceURIHandler = service_default_resource.DeleteProjectProjectNameServiceServiceNameResourceResourceURIHandlerFunc(func(params service_default_resource.DeleteProjectProjectNameServiceServiceNameResourceResourceURIParams) middleware.Responder {
		return middleware.NotImplemented("operation service_default_resource.DeleteProjectProjectNameServiceServiceNameResourceResourceURI has not yet been implemented")
	})

	api.StageDeleteProjectProjectNameStageStageNameHandler = stage.DeleteProjectProjectNameStageStageNameHandlerFunc(func(params stage.DeleteProjectProjectNameStageStageNameParams) middleware.Responder {
		return middleware.NotImplemented("operation stage.DeleteProjectProjectNameStageStageName has not yet been implemented")
	})

	api.StageResourceDeleteProjectProjectNameStageStageNameResourceResourceURIHandler = stage_resource.DeleteProjectProjectNameStageStageNameResourceResourceURIHandlerFunc(func(params stage_resource.DeleteProjectProjectNameStageStageNameResourceResourceURIParams) middleware.Responder {
		return middleware.NotImplemented("operation stage_resource.DeleteProjectProjectNameStageStageNameResourceResourceURI has not yet been implemented")
	})

	api.ServiceDeleteProjectProjectNameStageStageNameServiceServiceNameHandler = service.DeleteProjectProjectNameStageStageNameServiceServiceNameHandlerFunc(func(params service.DeleteProjectProjectNameStageStageNameServiceServiceNameParams) middleware.Responder {
		return middleware.NotImplemented("operation service.DeleteProjectProjectNameStageStageNameServiceServiceName has not yet been implemented")
	})

	api.ServiceResourceDeleteProjectProjectNameStageStageNameServiceServiceNameResourceResourceURIHandler = service_resource.DeleteProjectProjectNameStageStageNameServiceServiceNameResourceResourceURIHandlerFunc(func(params service_resource.DeleteProjectProjectNameStageStageNameServiceServiceNameResourceResourceURIParams) middleware.Responder {
		return middleware.NotImplemented("operation service_resource.DeleteProjectProjectNameStageStageNameServiceServiceNameResourceResourceURI has not yet been implemented")
	})

	api.ProjectGetProjectHandler = project.GetProjectHandlerFunc(func(params project.GetProjectParams) middleware.Responder {
		return middleware.NotImplemented("operation project.GetProject has not yet been implemented")
	})

	api.ProjectGetProjectProjectNameHandler = project.GetProjectProjectNameHandlerFunc(func(params project.GetProjectProjectNameParams) middleware.Responder {
		return middleware.NotImplemented("operation project.GetProjectProjectName has not yet been implemented")
	})

	api.ProjectResourceGetProjectProjectNameResourceHandler = project_resource.GetProjectProjectNameResourceHandlerFunc(func(params project_resource.GetProjectProjectNameResourceParams) middleware.Responder {
		return middleware.NotImplemented("operation project_resource.GetProjectProjectNameResource has not yet been implemented")
	})

	api.ProjectResourceGetProjectProjectNameResourceResourceURIHandler = project_resource.GetProjectProjectNameResourceResourceURIHandlerFunc(func(params project_resource.GetProjectProjectNameResourceResourceURIParams) middleware.Responder {
		return middleware.NotImplemented("operation project_resource.GetProjectProjectNameResourceResourceURI has not yet been implemented")
	})

	api.ServiceDefaultResourceGetProjectProjectNameServiceServiceNameResourceHandler = service_default_resource.GetProjectProjectNameServiceServiceNameResourceHandlerFunc(func(params service_default_resource.GetProjectProjectNameServiceServiceNameResourceParams) middleware.Responder {
		return middleware.NotImplemented("operation service_default_resource.GetProjectProjectNameServiceServiceNameResource has not yet been implemented")
	})

	api.ServiceDefaultResourceGetProjectProjectNameServiceServiceNameResourceResourceURIHandler = service_default_resource.GetProjectProjectNameServiceServiceNameResourceResourceURIHandlerFunc(func(params service_default_resource.GetProjectProjectNameServiceServiceNameResourceResourceURIParams) middleware.Responder {
		return middleware.NotImplemented("operation service_default_resource.GetProjectProjectNameServiceServiceNameResourceResourceURI has not yet been implemented")
	})

	api.StageGetProjectProjectNameStageHandler = stage.GetProjectProjectNameStageHandlerFunc(func(params stage.GetProjectProjectNameStageParams) middleware.Responder {
		return middleware.NotImplemented("operation stage.GetProjectProjectNameStage has not yet been implemented")
	})

	api.StageGetProjectProjectNameStageStageNameHandler = stage.GetProjectProjectNameStageStageNameHandlerFunc(func(params stage.GetProjectProjectNameStageStageNameParams) middleware.Responder {
		return middleware.NotImplemented("operation stage.GetProjectProjectNameStageStageName has not yet been implemented")
	})

	api.StageResourceGetProjectProjectNameStageStageNameResourceHandler = stage_resource.GetProjectProjectNameStageStageNameResourceHandlerFunc(func(params stage_resource.GetProjectProjectNameStageStageNameResourceParams) middleware.Responder {
		return middleware.NotImplemented("operation stage_resource.GetProjectProjectNameStageStageNameResource has not yet been implemented")
	})

	api.StageResourceGetProjectProjectNameStageStageNameResourceResourceURIHandler = stage_resource.GetProjectProjectNameStageStageNameResourceResourceURIHandlerFunc(func(params stage_resource.GetProjectProjectNameStageStageNameResourceResourceURIParams) middleware.Responder {
		return middleware.NotImplemented("operation stage_resource.GetProjectProjectNameStageStageNameResourceResourceURI has not yet been implemented")
	})

	api.ServiceGetProjectProjectNameStageStageNameServiceHandler = service.GetProjectProjectNameStageStageNameServiceHandlerFunc(func(params service.GetProjectProjectNameStageStageNameServiceParams) middleware.Responder {
		return middleware.NotImplemented("operation service.GetProjectProjectNameStageStageNameService has not yet been implemented")
	})

	api.ServiceGetProjectProjectNameStageStageNameServiceServiceNameHandler = service.GetProjectProjectNameStageStageNameServiceServiceNameHandlerFunc(func(params service.GetProjectProjectNameStageStageNameServiceServiceNameParams) middleware.Responder {
		return middleware.NotImplemented("operation service.GetProjectProjectNameStageStageNameServiceServiceName has not yet been implemented")
	})

	api.ServiceResourceGetProjectProjectNameStageStageNameServiceServiceNameResourceHandler = service_resource.GetProjectProjectNameStageStageNameServiceServiceNameResourceHandlerFunc(func(params service_resource.GetProjectProjectNameStageStageNameServiceServiceNameResourceParams) middleware.Responder {
		return middleware.NotImplemented("operation service_resource.GetProjectProjectNameStageStageNameServiceServiceNameResource has not yet been implemented")
	})

	api.ServiceResourceGetProjectProjectNameStageStageNameServiceServiceNameResourceResourceURIHandler = service_resource.GetProjectProjectNameStageStageNameServiceServiceNameResourceResourceURIHandlerFunc(func(params service_resource.GetProjectProjectNameStageStageNameServiceServiceNameResourceResourceURIParams) middleware.Responder {
		return middleware.NotImplemented("operation service_resource.GetProjectProjectNameStageStageNameServiceServiceNameResourceResourceURI has not yet been implemented")
	})

	api.ProjectPostProjectHandler = project.PostProjectHandlerFunc(func(params project.PostProjectParams) middleware.Responder {
		return middleware.NotImplemented("operation project.PostProject has not yet been implemented")
	})

	api.ProjectResourcePostProjectProjectNameResourceHandler = project_resource.PostProjectProjectNameResourceHandlerFunc(func(params project_resource.PostProjectProjectNameResourceParams) middleware.Responder {
		return middleware.NotImplemented("operation project_resource.PostProjectProjectNameResource has not yet been implemented")
	})

	api.ServiceDefaultResourcePostProjectProjectNameServiceServiceNameResourceHandler = service_default_resource.PostProjectProjectNameServiceServiceNameResourceHandlerFunc(func(params service_default_resource.PostProjectProjectNameServiceServiceNameResourceParams) middleware.Responder {
		return middleware.NotImplemented("operation service_default_resource.PostProjectProjectNameServiceServiceNameResource has not yet been implemented")
	})

	api.StagePostProjectProjectNameStageHandler = stage.PostProjectProjectNameStageHandlerFunc(func(params stage.PostProjectProjectNameStageParams) middleware.Responder {
		return middleware.NotImplemented("operation stage.PostProjectProjectNameStage has not yet been implemented")
	})

	api.StageResourcePostProjectProjectNameStageStageNameResourceHandler = stage_resource.PostProjectProjectNameStageStageNameResourceHandlerFunc(func(params stage_resource.PostProjectProjectNameStageStageNameResourceParams) middleware.Responder {
		return middleware.NotImplemented("operation stage_resource.PostProjectProjectNameStageStageNameResource has not yet been implemented")
	})

	api.ServicePostProjectProjectNameStageStageNameServiceHandler = service.PostProjectProjectNameStageStageNameServiceHandlerFunc(func(params service.PostProjectProjectNameStageStageNameServiceParams) middleware.Responder {
		return middleware.NotImplemented("operation service.PostProjectProjectNameStageStageNameService has not yet been implemented")
	})

	api.ServiceResourcePostProjectProjectNameStageStageNameServiceServiceNameResourceHandler = service_resource.PostProjectProjectNameStageStageNameServiceServiceNameResourceHandlerFunc(func(params service_resource.PostProjectProjectNameStageStageNameServiceServiceNameResourceParams) middleware.Responder {
		return middleware.NotImplemented("operation service_resource.PostProjectProjectNameStageStageNameServiceServiceNameResource has not yet been implemented")
	})

	api.ProjectPutProjectProjectNameHandler = project.PutProjectProjectNameHandlerFunc(func(params project.PutProjectProjectNameParams) middleware.Responder {
		return middleware.NotImplemented("operation project.PutProjectProjectName has not yet been implemented")
	})

	api.ProjectResourcePutProjectProjectNameResourceHandler = project_resource.PutProjectProjectNameResourceHandlerFunc(func(params project_resource.PutProjectProjectNameResourceParams) middleware.Responder {
		return middleware.NotImplemented("operation project_resource.PutProjectProjectNameResource has not yet been implemented")
	})

	api.ProjectResourcePutProjectProjectNameResourceResourceURIHandler = project_resource.PutProjectProjectNameResourceResourceURIHandlerFunc(func(params project_resource.PutProjectProjectNameResourceResourceURIParams) middleware.Responder {
		return middleware.NotImplemented("operation project_resource.PutProjectProjectNameResourceResourceURI has not yet been implemented")
	})

	api.ServiceDefaultResourcePutProjectProjectNameServiceServiceNameResourceHandler = service_default_resource.PutProjectProjectNameServiceServiceNameResourceHandlerFunc(func(params service_default_resource.PutProjectProjectNameServiceServiceNameResourceParams) middleware.Responder {
		return middleware.NotImplemented("operation service_default_resource.PutProjectProjectNameServiceServiceNameResource has not yet been implemented")
	})

	api.ServiceDefaultResourcePutProjectProjectNameServiceServiceNameResourceResourceURIHandler = service_default_resource.PutProjectProjectNameServiceServiceNameResourceResourceURIHandlerFunc(func(params service_default_resource.PutProjectProjectNameServiceServiceNameResourceResourceURIParams) middleware.Responder {
		return middleware.NotImplemented("operation service_default_resource.PutProjectProjectNameServiceServiceNameResourceResourceURI has not yet been implemented")
	})

	api.StagePutProjectProjectNameStageStageNameHandler = stage.PutProjectProjectNameStageStageNameHandlerFunc(func(params stage.PutProjectProjectNameStageStageNameParams) middleware.Responder {
		return middleware.NotImplemented("operation stage.PutProjectProjectNameStageStageName has not yet been implemented")
	})

	api.StageResourcePutProjectProjectNameStageStageNameResourceHandler = stage_resource.PutProjectProjectNameStageStageNameResourceHandlerFunc(func(params stage_resource.PutProjectProjectNameStageStageNameResourceParams) middleware.Responder {
		return middleware.NotImplemented("operation stage_resource.PutProjectProjectNameStageStageNameResource has not yet been implemented")
	})

	api.StageResourcePutProjectProjectNameStageStageNameResourceResourceURIHandler = stage_resource.PutProjectProjectNameStageStageNameResourceResourceURIHandlerFunc(func(params stage_resource.PutProjectProjectNameStageStageNameResourceResourceURIParams) middleware.Responder {
		return middleware.NotImplemented("operation stage_resource.PutProjectProjectNameStageStageNameResourceResourceURI has not yet been implemented")
	})

	api.ServicePutProjectProjectNameStageStageNameServiceServiceNameHandler = service.PutProjectProjectNameStageStageNameServiceServiceNameHandlerFunc(func(params service.PutProjectProjectNameStageStageNameServiceServiceNameParams) middleware.Responder {
		return middleware.NotImplemented("operation service.PutProjectProjectNameStageStageNameServiceServiceName has not yet been implemented")
	})

	api.ServiceResourcePutProjectProjectNameStageStageNameServiceServiceNameResourceHandler = service_resource.PutProjectProjectNameStageStageNameServiceServiceNameResourceHandlerFunc(func(params service_resource.PutProjectProjectNameStageStageNameServiceServiceNameResourceParams) middleware.Responder {
		return middleware.NotImplemented("operation service_resource.PutProjectProjectNameStageStageNameServiceServiceNameResource has not yet been implemented")
	})

	api.ServiceResourcePutProjectProjectNameStageStageNameServiceServiceNameResourceResourceURIHandler = service_resource.PutProjectProjectNameStageStageNameServiceServiceNameResourceResourceURIHandlerFunc(func(params service_resource.PutProjectProjectNameStageStageNameServiceServiceNameResourceResourceURIParams) middleware.Responder {
		return middleware.NotImplemented("operation service_resource.PutProjectProjectNameStageStageNameServiceServiceNameResourceResourceURI has not yet been implemented")
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
	return handler
}
