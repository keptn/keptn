// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"github.com/keptn/keptn/configuration-service/restapi/operations/event"
	"github.com/keptn/keptn/configuration-service/restapi/operations/service_approval"
	"github.com/keptn/keptn/configuration-service/restapi/operations/services"
	"net/http"
	"os"
	"os/exec"
	"strings"

	errors "github.com/go-openapi/errors"
	runtime "github.com/go-openapi/runtime"
	utils "github.com/keptn/go-utils/pkg/lib"
	handlers "github.com/keptn/keptn/configuration-service/handlers"
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

	api.ProjectDeleteProjectProjectNameHandler = project.DeleteProjectProjectNameHandlerFunc(handlers.DeleteProjectProjectNameHandlerFunc)

	api.ProjectResourceDeleteProjectProjectNameResourceResourceURIHandler = project_resource.DeleteProjectProjectNameResourceResourceURIHandlerFunc(handlers.DeleteProjectProjectNameResourceResourceURIHandlerFunc)

	api.ServiceDefaultResourceDeleteProjectProjectNameServiceServiceNameResourceResourceURIHandler = service_default_resource.DeleteProjectProjectNameServiceServiceNameResourceResourceURIHandlerFunc(handlers.DeleteProjectProjectNameServiceServiceNameResourceResourceURIHandlerFunc)

	api.StageDeleteProjectProjectNameStageStageNameHandler = stage.DeleteProjectProjectNameStageStageNameHandlerFunc(handlers.DeleteProjectProjectNameStageStageNameHandlerFunc)

	api.StageResourceDeleteProjectProjectNameStageStageNameResourceResourceURIHandler = stage_resource.DeleteProjectProjectNameStageStageNameResourceResourceURIHandlerFunc(handlers.DeleteProjectProjectNameStageStageNameResourceResourceURIHandlerFunc)

	api.ServiceDeleteProjectProjectNameStageStageNameServiceServiceNameHandler = service.DeleteProjectProjectNameStageStageNameServiceServiceNameHandlerFunc(handlers.DeleteProjectProjectNameStageStageNameServiceServiceNameHandlerFunc)

	api.ServiceResourceDeleteProjectProjectNameStageStageNameServiceServiceNameResourceResourceURIHandler = service_resource.DeleteProjectProjectNameStageStageNameServiceServiceNameResourceResourceURIHandlerFunc(handlers.DeleteProjectProjectNameStageStageNameServiceServiceNameResourceResourceURIHandlerFunc)

	api.ProjectGetProjectHandler = project.GetProjectHandlerFunc(handlers.GetProjectHandlerFunc)

	api.ProjectGetProjectProjectNameHandler = project.GetProjectProjectNameHandlerFunc(handlers.GetProjectProjectNameHandlerFunc)

	api.ProjectResourceGetProjectProjectNameResourceHandler = project_resource.GetProjectProjectNameResourceHandlerFunc(handlers.GetProjectProjectNameResourceHandlerFunc)

	api.ProjectResourceGetProjectProjectNameResourceResourceURIHandler = project_resource.GetProjectProjectNameResourceResourceURIHandlerFunc(handlers.GetProjectProjectNameResourceResourceURIHandlerFunc)

	api.ServiceDefaultResourceGetProjectProjectNameServiceServiceNameResourceHandler = service_default_resource.GetProjectProjectNameServiceServiceNameResourceHandlerFunc(handlers.GetProjectProjectNameServiceServiceNameResourceHandlerFunc)

	api.ServiceDefaultResourceGetProjectProjectNameServiceServiceNameResourceResourceURIHandler = service_default_resource.GetProjectProjectNameServiceServiceNameResourceResourceURIHandlerFunc(handlers.GetProjectProjectNameServiceServiceNameResourceResourceURIHandlerFunc)

	api.StageGetProjectProjectNameStageHandler = stage.GetProjectProjectNameStageHandlerFunc(handlers.GetProjectProjectNameStageHandlerFunc)

	api.StageGetProjectProjectNameStageStageNameHandler = stage.GetProjectProjectNameStageStageNameHandlerFunc(handlers.GetProjectProjectNameStageStageNameHandlerFunc)

	api.StageResourceGetProjectProjectNameStageStageNameResourceHandler = stage_resource.GetProjectProjectNameStageStageNameResourceHandlerFunc(handlers.GetProjectProjectNameStageStageNameResourceHandlerFunc)

	api.StageResourceGetProjectProjectNameStageStageNameResourceResourceURIHandler = stage_resource.GetProjectProjectNameStageStageNameResourceResourceURIHandlerFunc(handlers.GetProjectProjectNameStageStageNameResourceResourceURIHandlerFunc)

	api.ServiceGetProjectProjectNameStageStageNameServiceHandler = service.GetProjectProjectNameStageStageNameServiceHandlerFunc(handlers.GetProjectProjectNameStageStageNameServiceHandlerFunc)

	api.ServiceGetProjectProjectNameStageStageNameServiceServiceNameHandler = service.GetProjectProjectNameStageStageNameServiceServiceNameHandlerFunc(handlers.GetProjectProjectNameStageStageNameServiceServiceNameHandlerFunc)

	api.ServiceResourceGetProjectProjectNameStageStageNameServiceServiceNameResourceHandler = service_resource.GetProjectProjectNameStageStageNameServiceServiceNameResourceHandlerFunc(handlers.GetProjectProjectNameStageStageNameServiceServiceNameResourceHandlerFunc)

	api.ServiceResourceGetProjectProjectNameStageStageNameServiceServiceNameResourceResourceURIHandler = service_resource.GetProjectProjectNameStageStageNameServiceServiceNameResourceResourceURIHandlerFunc(handlers.GetProjectProjectNameStageStageNameServiceServiceNameResourceResourceURIHandlerFunc)

	api.ProjectPostProjectHandler = project.PostProjectHandlerFunc(handlers.PostProjectHandlerFunc)

	api.ProjectResourcePostProjectProjectNameResourceHandler = project_resource.PostProjectProjectNameResourceHandlerFunc(handlers.PostProjectProjectNameResourceHandlerFunc)

	api.ServiceDefaultResourcePostProjectProjectNameServiceServiceNameResourceHandler = service_default_resource.PostProjectProjectNameServiceServiceNameResourceHandlerFunc(handlers.PostProjectProjectNameServiceServiceNameResourceHandlerFunc)

	api.StagePostProjectProjectNameStageHandler = stage.PostProjectProjectNameStageHandlerFunc(handlers.PostProjectProjectNameStageHandlerFunc)

	api.StageResourcePostProjectProjectNameStageStageNameResourceHandler = stage_resource.PostProjectProjectNameStageStageNameResourceHandlerFunc(handlers.PostProjectProjectNameStageStageNameResourceHandlerFunc)

	api.ServicePostProjectProjectNameStageStageNameServiceHandler = service.PostProjectProjectNameStageStageNameServiceHandlerFunc(handlers.PostProjectProjectNameStageStageNameServiceHandlerFunc)

	api.ServiceResourcePostProjectProjectNameStageStageNameServiceServiceNameResourceHandler = service_resource.PostProjectProjectNameStageStageNameServiceServiceNameResourceHandlerFunc(handlers.PostProjectProjectNameStageStageNameServiceServiceNameResourceHandlerFunc)

	api.ProjectPutProjectProjectNameHandler = project.PutProjectProjectNameHandlerFunc(handlers.PutProjectProjectNameHandlerFunc)

	api.ProjectResourcePutProjectProjectNameResourceHandler = project_resource.PutProjectProjectNameResourceHandlerFunc(handlers.PutProjectProjectNameResourceHandlerFunc)

	api.ProjectResourcePutProjectProjectNameResourceResourceURIHandler = project_resource.PutProjectProjectNameResourceResourceURIHandlerFunc(handlers.PutProjectProjectNameResourceResourceURIHandlerFunc)

	api.ServiceDefaultResourcePutProjectProjectNameServiceServiceNameResourceHandler = service_default_resource.PutProjectProjectNameServiceServiceNameResourceHandlerFunc(handlers.PutProjectProjectNameServiceServiceNameResourceHandlerFunc)

	api.ServiceDefaultResourcePutProjectProjectNameServiceServiceNameResourceResourceURIHandler = service_default_resource.PutProjectProjectNameServiceServiceNameResourceResourceURIHandlerFunc(handlers.PutProjectProjectNameServiceServiceNameResourceResourceURIHandlerFunc)

	api.StagePutProjectProjectNameStageStageNameHandler = stage.PutProjectProjectNameStageStageNameHandlerFunc(handlers.PutProjectProjectNameStageStageNameHandlerFunc)

	api.StageResourcePutProjectProjectNameStageStageNameResourceHandler = stage_resource.PutProjectProjectNameStageStageNameResourceHandlerFunc(handlers.PutProjectProjectNameStageStageNameResourceHandlerFunc)

	api.StageResourcePutProjectProjectNameStageStageNameResourceResourceURIHandler = stage_resource.PutProjectProjectNameStageStageNameResourceResourceURIHandlerFunc(handlers.PutProjectProjectNameStageStageNameResourceResourceURIHandlerFunc)

	api.ServicePutProjectProjectNameStageStageNameServiceServiceNameHandler = service.PutProjectProjectNameStageStageNameServiceServiceNameHandlerFunc(handlers.PutProjectProjectNameStageStageNameServiceServiceNameHandlerFunc)

	api.ServiceResourcePutProjectProjectNameStageStageNameServiceServiceNameResourceHandler = service_resource.PutProjectProjectNameStageStageNameServiceServiceNameResourceHandlerFunc(handlers.PutProjectProjectNameStageStageNameServiceServiceNameResourceHandlerFunc)

	api.ServiceResourcePutProjectProjectNameStageStageNameServiceServiceNameResourceResourceURIHandler = service_resource.PutProjectProjectNameStageStageNameServiceServiceNameResourceResourceURIHandlerFunc(handlers.PutProjectProjectNameStageStageNameServiceServiceNameResourceResourceURIHandlerFunc)

	api.EventHandleEventHandler = event.HandleEventHandlerFunc(handlers.HandleEventHandlerFunc)

	api.ServicesGetServicesHandler = services.GetServicesHandlerFunc(handlers.GetServices)

	api.ServicesGetServiceHandler = services.GetServiceHandlerFunc(handlers.GetService)

	api.ServiceApprovalGetServiceApprovalsHandler = service_approval.GetServiceApprovalsHandlerFunc(handlers.GetServiceApprovals)

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
	logger := utils.NewLogger("", "", "configuration-service")
	if os.Getenv("env") == "production" {
		///////// initialize git ////////////
		logger.Debug("Configuring git user.email")
		cmd := exec.Command("git", "config", "--global", "user.email", "keptn@keptn.com")
		_, err := cmd.Output()
		if err != nil {
			logger.Error("Could not configure git user.email: " + err.Error())
		}
		logger.Debug("Configuring git user.name")
		cmd = exec.Command("git", "config", "--global", "user.name", "keptn")
		_, err = cmd.Output()
		if err != nil {
			logger.Error("Could not configure git user.name: " + err.Error())
		}
		////////////////////////////////////
	}
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Shortcut helpers for swagger-ui
		if r.URL.Path == "/swagger-ui" {
			http.Redirect(w, r, "/swagger-ui/", http.StatusFound)
			return
		}
		// Serving ./swagger-ui/
		if strings.Index(r.URL.Path, "/swagger-ui/") == 0 {
			http.StripPrefix("/swagger-ui/", http.FileServer(http.Dir("swagger-ui"))).ServeHTTP(w, r)
			return
		}

		handler.ServeHTTP(w, r)
	})
}
