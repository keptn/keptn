package main

import (
	"github.com/gin-gonic/gin"
	keptncommon "github.com/keptn/go-utils/pkg/lib/keptn"
	"github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/shipyard-controller/common"
	"github.com/keptn/keptn/shipyard-controller/controller"
	"github.com/keptn/keptn/shipyard-controller/db"
	"github.com/keptn/keptn/shipyard-controller/docs"
	"github.com/keptn/keptn/shipyard-controller/handler"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
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

	kubeAPI, err := createKubeAPI()
	if err != nil {
		log.Fatalf("could not create kubernetes client: %s", err.Error())
	}

	eventSender, err := v0_2_0.NewHTTPEventSender("")
	if err != nil {
		log.Fatal(err)
	}

	logger := keptncommon.NewLogger("", "", "shipyard-controller")

	projectManager := handler.NewProjectManager(
		common.NewGitConfigurationStore(csEndpoint.String()),
		createSecretStore(kubeAPI),
		createMaterializedView(logger),
		createTaskSequenceRepo(logger),
		createEventsRepo(logger))

	serviceManager := handler.NewServiceManager(
		createMaterializedView(logger),
		common.NewGitConfigurationStore(csEndpoint.String()),
		logger)

	stageManager := handler.NewStageManager(createMaterializedView(logger), logger)

	engine := gin.Default()
	apiV1 := engine.Group("/v1")
	projectService := handler.NewProjectHandler(projectManager, eventSender)
	projectController := controller.NewProjectController(projectService)
	projectController.Inject(apiV1)

	serviceHandler := handler.NewServiceHandler(serviceManager, logger)
	serviceController := controller.NewServiceController(serviceHandler)
	serviceController.Inject(apiV1)

	eventHandler := handler.NewEventHandler()
	eventController := controller.NewEventController(eventHandler)
	eventController.Inject(apiV1)

	stageHandler := handler.NewStageHandler(stageManager)
	stageController := controller.NewStageController(stageHandler)
	stageController.Inject(apiV1)

	evaluationManager, err := handler.NewEvaluationManager(eventSender, createMaterializedView(logger), logger)
	if err != nil {
		log.Fatal(err)
	}
	evaluationHandler := handler.NewEvaluationHandler(evaluationManager)
	evaluationController := controller.NewEvaluationController(evaluationHandler)
	evaluationController.Inject(apiV1)

	engine.Static("/swagger-ui", "./swagger-ui")
	engine.Run()
}

func createMaterializedView(logger *keptncommon.Logger) *db.ProjectsMaterializedView {
	projectesMaterializedView := &db.ProjectsMaterializedView{
		ProjectRepo:     createProjectRepo(logger),
		EventsRetriever: createEventsRepo(logger),
		Logger:          logger,
	}
	return projectesMaterializedView
}

func createProjectRepo(logger *keptncommon.Logger) *db.MongoDBProjectsRepo {
	return &db.MongoDBProjectsRepo{Logger: logger}
}

func createEventsRepo(logger *keptncommon.Logger) *db.MongoDBEventsRepo {
	return &db.MongoDBEventsRepo{Logger: logger}
}

func createTaskSequenceRepo(logger *keptncommon.Logger) *db.TaskSequenceMongoDBRepo {
	return &db.TaskSequenceMongoDBRepo{Logger: logger}
}

func createSecretStore(kubeAPI *kubernetes.Clientset) *common.K8sSecretStore {
	return common.NewK8sSecretStore(kubeAPI)
}

// GetKubeAPI godoc
func createKubeAPI() (*kubernetes.Clientset, error) {
	var config *rest.Config
	config, err := rest.InClusterConfig()

	if err != nil {
		return nil, err
	}

	kubeAPI, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return kubeAPI, nil
}
