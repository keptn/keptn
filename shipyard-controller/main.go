package main

import (
	"github.com/gin-gonic/gin"
	keptncommon "github.com/keptn/go-utils/pkg/lib/keptn"
	"github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/shipyard-controller/controller"
	"github.com/keptn/keptn/shipyard-controller/db"
	"github.com/keptn/keptn/shipyard-controller/docs"
	"github.com/keptn/keptn/shipyard-controller/handler"
	"log"
	"os"
)

// @title Shipyard Controller API
// @version 1.0
// @description This is the API documentation of the Shipyard Controller.

// @securityDefinitions.apiKey ApiKeyAuth
// @in header
// @name x-token

// @contact.name Keptn Team
// @contact.url http://www.keptn.sh

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath /v1
func main() {

	if os.Getenv("GIN_MODE") == "release" {
		docs.SwaggerInfo.Version = os.Getenv("version")
		docs.SwaggerInfo.BasePath = "/api/shipyard-controller/v1"
		docs.SwaggerInfo.Schemes = []string{"https"}
	}

	csEndpoint, err := keptncommon.GetServiceEndpoint("CONFIGURATION_SERVICE")
	if err != nil {
		log.Fatalf("could not get configuration-service URL: %s", err.Error())
	}
	logger := keptncommon.NewLogger("", "", "shipyard-controller")

	secretStore, err := handler.NewK8sSecretStore()
	if err != nil {
		log.Fatal(err)
	}
	configurationStore := handler.NewGitConfigurationStore(csEndpoint.String())
	eventRepository := &db.MongoDBEventsRepo{Logger: logger}
	projectRepository := &db.MongoDBProjectsRepo{Logger: logger}
	taskSequenceRepository := &db.TaskSequenceMongoDBRepo{Logger: logger}

	projectManager := handler.NewProjectManager(
		configurationStore,
		secretStore,
		projectRepository,
		taskSequenceRepository,
		eventRepository)

	eventSender, err := v0_2_0.NewHTTPEventSender("")
	if err != nil {
		log.Fatal(err)
	}

	engine := gin.Default()
	apiV1 := engine.Group("/v1")
	projectService := handler.NewProjectHandler(projectManager, eventSender)
	projectController := controller.NewProjectController(projectService)
	projectController.Inject(apiV1)

	serviceHandler := handler.NewServiceHandler()
	serviceController := controller.NewServiceController(serviceHandler)
	serviceController.Inject(apiV1)

	eventHandler := handler.NewEventHandler()
	eventController := controller.NewEventController(eventHandler)
	eventController.Inject(apiV1)

	engine.Static("/swagger-ui", "./swagger-ui")
	engine.Run()
}
