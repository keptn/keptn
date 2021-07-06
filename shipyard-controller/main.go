package main

import (
	"context"
	"github.com/benbjohnson/clock"
	"github.com/gin-gonic/gin"
	"github.com/keptn/go-utils/pkg/common/osutils"
	keptncommon "github.com/keptn/go-utils/pkg/lib/keptn"
	"github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/shipyard-controller/common"
	"github.com/keptn/keptn/shipyard-controller/controller"
	"github.com/keptn/keptn/shipyard-controller/db"
	_ "github.com/keptn/keptn/shipyard-controller/docs"
	"github.com/keptn/keptn/shipyard-controller/handler"
	"github.com/keptn/keptn/shipyard-controller/handler/sequencehooks"
	"github.com/keptn/keptn/shipyard-controller/models"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"os"
	"strconv"
	"time"
)

// @title Control Plane API
// @version develop
// @description This is the API documentation of the Shipyard Controller.

// @securityDefinitions.apiKey ApiKeyAuth
// @in header
// @name x-token

// @contact.name Keptn Team
// @contact.url http://www.keptn.sh

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath /v1

const envVarConfigurationSvcEndpoint = "CONFIGURATION_SERVICE"
const envVarEventDispatchIntervalSec = "EVENT_DISPATCH_INTERVAL_SEC"
const envVarSequenceDispatchIntervalSec = "SEQUENCE_DISPATCH_INTERVAL_SEC"
const envVarTaskStartedWaitDuration = "TASK_STARTED_WAIT_DURATION"
const envVarLogTTL = "LOG_TTL"
const envVarEventDispatchIntervalSecDefault = "10"
const envVarSequenceDispatchIntervalSecDefault = "10s"
const envVarLogsTTLDefault = "120h" // 5 days
const envVarTaskStartedWaitDurationDefault = "10m"

func main() {
	log.SetLevel(log.InfoLevel)

	if osutils.GetAndCompareOSEnv("GIN_MODE", "release") {
		// disable GIN request logging in release mode
		gin.SetMode("release")
		gin.DefaultWriter = ioutil.Discard
	}

	eventDispatcherSyncInterval, err := strconv.Atoi(osutils.GetOSEnvOrDefault(envVarEventDispatchIntervalSec, envVarEventDispatchIntervalSecDefault))
	if err != nil {
		log.Fatalf("Unexpected value of EVENT_DISPATCH_INTERVAL_SEC environment variable. Need to be a number")
	}

	csEndpoint, err := keptncommon.GetServiceEndpoint(envVarConfigurationSvcEndpoint)
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

	projectManager := handler.NewProjectManager(
		common.NewGitConfigurationStore(csEndpoint.String()),
		createSecretStore(kubeAPI),
		createMaterializedView(),
		createTaskSequenceRepo(),
		createEventsRepo(),
		createSequenceQueueRepo(),
		createEventQueueRepo())

	serviceManager := handler.NewServiceManager(
		createMaterializedView(),
		common.NewGitConfigurationStore(csEndpoint.String()),
	)

	stageManager := handler.NewStageManager(createMaterializedView())

	eventDispatcher := handler.NewEventDispatcher(createEventsRepo(), createEventQueueRepo(), eventSender, time.Duration(eventDispatcherSyncInterval)*time.Second)
	sequenceDispatcherChannel := make(chan models.Event)
	sequenceDispatcher := handler.NewSequenceDispatcher(
		createEventsRepo(),
		createSequenceQueueRepo(),
		createTaskSequenceRepo(),
		getDurationFromEnvVar(envVarSequenceDispatchIntervalSec, envVarSequenceDispatchIntervalSecDefault),
		sequenceDispatcherChannel,
		clock.New(),
	)

	cancelSequenceChannel := make(chan common.SequenceCancellation)
	shipyardController := handler.GetShipyardControllerInstance(eventDispatcher, sequenceDispatcher, sequenceDispatcherChannel, cancelSequenceChannel)

	engine := gin.Default()
	apiV1 := engine.Group("/v1")
	apiHealth := engine.Group("")

	projectService := handler.NewProjectHandler(projectManager, eventSender)
	projectController := controller.NewProjectController(projectService)
	projectController.Inject(apiV1)

	serviceHandler := handler.NewServiceHandler(serviceManager, eventSender)
	serviceController := controller.NewServiceController(serviceHandler)
	serviceController.Inject(apiV1)

	eventHandler := handler.NewEventHandler(shipyardController)
	eventController := controller.NewEventController(eventHandler)
	eventController.Inject(apiV1)

	stageHandler := handler.NewStageHandler(stageManager)
	stageController := controller.NewStageController(stageHandler)
	stageController.Inject(apiV1)

	evaluationManager, err := handler.NewEvaluationManager(eventSender, createMaterializedView())
	if err != nil {
		log.Fatal(err)
	}
	evaluationHandler := handler.NewEvaluationHandler(evaluationManager)
	evaluationController := controller.NewEvaluationController(evaluationHandler)
	evaluationController.Inject(apiV1)

	stateHandler := handler.NewStateHandler(db.NewMongoDBStateRepo(db.GetMongoDBConnectionInstance()))
	stateController := controller.NewStateController(stateHandler)
	stateController.Inject(apiV1)

	sequenceStateMaterializedView := sequencehooks.NewSequenceStateMaterializedView(createStateRepo())
	shipyardController.AddSequenceTriggeredHook(sequenceStateMaterializedView)
	shipyardController.AddSequenceStartedHook(sequenceStateMaterializedView)
	shipyardController.AddSequenceTaskTriggeredHook(sequenceStateMaterializedView)
	shipyardController.AddSequenceTaskStartedHook(sequenceStateMaterializedView)
	shipyardController.AddSequenceTaskFinishedHook(sequenceStateMaterializedView)
	shipyardController.AddSubSequenceFinishedHook(sequenceStateMaterializedView)
	shipyardController.AddSequenceFinishedHook(sequenceStateMaterializedView)
	shipyardController.AddSequenceTimeoutHook(sequenceStateMaterializedView)

	taskStartedWaitDuration := getDurationFromEnvVar(envVarTaskStartedWaitDuration, envVarTaskStartedWaitDurationDefault)

	watcher := handler.NewSequenceWatcher(
		cancelSequenceChannel,
		createEventsRepo(),
		createProjectRepo(),
		taskStartedWaitDuration,
		1*time.Minute,
		clock.New(),
	)

	watcher.Run(context.Background())

	uniformRepo := createUniformRepo()
	uniformManager := handler.NewUniformIntegrationManager(uniformRepo)
	uniformHandler := handler.NewUniformIntegrationHandler(uniformManager)
	uniformController := controller.NewUniformIntegrationController(uniformHandler)
	uniformController.Inject(apiV1)

	logRepo := createLogRepo()
	err = logRepo.SetupTTLIndex(getDurationFromEnvVar(envVarLogTTL, envVarLogsTTLDefault))
	if err != nil {
		log.WithError(err).Error("could not setup TTL index for log repo entries")
	}
	logHandler := handler.NewLogHandler(handler.NewLogManager(logRepo))
	logController := controller.NewLogController(logHandler)
	logController.Inject(apiV1)

	sequenceMigrator := handler.NewSequenceMigrator(createEventsRepo(), createStateRepo(), createProjectRepo())
	sequenceMigrator.Run()

	healthHandler := handler.NewHealthHandler()
	healthController := controller.NewHealthController(healthHandler)
	healthController.Inject(apiHealth)

	engine.Static("/swagger-ui", "./swagger-ui")
	engine.Run()
}

func createMaterializedView() *db.ProjectsMaterializedView {
	projectesMaterializedView := &db.ProjectsMaterializedView{
		ProjectRepo:     createProjectRepo(),
		EventsRetriever: createEventsRepo(),
	}
	return projectesMaterializedView
}

func createUniformRepo() *db.MongoDBUniformRepo {
	return db.NewMongoDBUniformRepo(db.GetMongoDBConnectionInstance())
}

func createStateRepo() *db.MongoDBStateRepo {
	return db.NewMongoDBStateRepo(db.GetMongoDBConnectionInstance())
}

func createProjectRepo() *db.MongoDBProjectsRepo {
	return db.NewMongoDBProjectsRepo(db.GetMongoDBConnectionInstance())
}

func createEventsRepo() *db.MongoDBEventsRepo {
	return db.NewMongoDBEventsRepo(db.GetMongoDBConnectionInstance())
}

func createSequenceQueueRepo() *db.MongoDBSequenceQueueRepo {
	return db.NewMongoDBSequenceQueueRepo(db.GetMongoDBConnectionInstance())
}

func createEventQueueRepo() *db.MongoDBEventQueueRepo {
	return db.NewMongoDBEventQueueRepo(db.GetMongoDBConnectionInstance())
}

func createTaskSequenceRepo() *db.TaskSequenceMongoDBRepo {
	return db.NewTaskSequenceMongoDBRepo(db.GetMongoDBConnectionInstance())
}

func createSecretStore(kubeAPI *kubernetes.Clientset) *common.K8sSecretStore {
	return common.NewK8sSecretStore(kubeAPI)
}

func createLogRepo() *db.MongoDBLogRepo {
	return db.NewMongoDBLogRepo(db.GetMongoDBConnectionInstance())
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

func getDurationFromEnvVar(envVar, fallbackValue string) time.Duration {
	durationString := os.Getenv(envVar)
	var duration time.Duration
	var err error
	if durationString != "" {
		duration, err = time.ParseDuration(durationString)
		if err != nil {
			log.Errorf("could not parse log %s env var %s: %s. Will use default value %s", envVar, duration, err.Error(), fallbackValue)
		}
	}

	if duration.Seconds() == 0 {
		duration, err = time.ParseDuration(fallbackValue)
		if err != nil {
			log.Errorf("could not parse default duration string %s. %s will be set to 0", err.Error(), envVar)
			return time.Duration(0)
		}
	}

	return duration
}
