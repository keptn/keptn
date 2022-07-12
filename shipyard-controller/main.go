package main

import (
	"context"
	"github.com/keptn/keptn/shipyard-controller/leaderelection"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/gin-gonic/gin"
	"github.com/kelseyhightower/envconfig"
	apimodels "github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/go-utils/pkg/common/osutils"
	"github.com/keptn/keptn/shipyard-controller/common"
	"github.com/keptn/keptn/shipyard-controller/config"
	"github.com/keptn/keptn/shipyard-controller/controller"
	"github.com/keptn/keptn/shipyard-controller/db"
	"github.com/keptn/keptn/shipyard-controller/db/migration"
	_ "github.com/keptn/keptn/shipyard-controller/docs"
	"github.com/keptn/keptn/shipyard-controller/handler"
	"github.com/keptn/keptn/shipyard-controller/handler/sequencehooks"
	_ "github.com/keptn/keptn/shipyard-controller/models"
	"github.com/keptn/keptn/shipyard-controller/nats"
	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// @title        Control Plane API
// @version      develop
// @description  This is the API documentation of the Shipyard Controller.

// @securityDefinitions.apiKey  ApiKeyAuth
// @in                          header
// @name                        x-token

// @contact.name  Keptn Team
// @contact.url   http://www.keptn.sh

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath  /v1

const envVarSequenceDispatchIntervalSecDefault = "10s"
const envVarLogsTTLDefault = "120h" // 5 days
const envVarUniformTTLDefault = "1m"
const envVarSequenceWatcherIntervalDefault = "1m"
const envVarTaskStartedWaitDurationDefault = "10m"

func main() {
	kubeAPI, err := createKubeAPI()
	if err != nil {
		log.Fatalf("could not create kubernetes client: %s", err.Error())
	}
	var env config.EnvConfig
	if err := envconfig.Process("", &env); err != nil {
		log.Fatalf("Failed to process env var: %v", err)
	}

	_main(env, kubeAPI)
}

func _main(env config.EnvConfig, kubeAPI kubernetes.Interface) {
	log.SetLevel(log.InfoLevel)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if env.LogLevel != "" {
		logLevel, err := log.ParseLevel(env.LogLevel)
		if err != nil {
			log.WithError(err).Error("could not parse log level provided by 'LOG_LEVEL' env var")
		} else {
			log.SetLevel(logLevel)
		}
	}

	if osutils.GetAndCompareOSEnv("GIN_MODE", "release") {
		// disable GIN request logging in release mode
		gin.SetMode("release")
		gin.DefaultWriter = ioutil.Discard
	}

	csEndpoint, err := url.Parse(env.ConfigurationSvcEndpoint)
	if err != nil {
		log.Fatal(err)
	}

	connectionHandler := nats.NewNatsConnectionHandler(
		ctx,
		env.NatsURL,
	)

	eventSender, err := connectionHandler.GetPublisher()
	if err != nil {
		log.Fatal(err)
	}

	sequenceExecutionRepo := createSequenceExecutionRepo()

	secretStore := createSecretStore(kubeAPI)

	projectMVRepo := createProjectMVRepo()
	projectManager := handler.NewProjectManager(
		common.NewGitConfigurationStore(csEndpoint.String()),
		secretStore,
		projectMVRepo,
		sequenceExecutionRepo,
		createEventsRepo(),
		createSequenceQueueRepo(),
		createEventQueueRepo())

	repositoryProvisioner := handler.NewRepositoryProvisioner(env.AutomaticProvisioningURL, &http.Client{})

	uniformRepo := createUniformRepo()
	err = uniformRepo.SetupTTLIndex(getDurationFromEnvVar(env.UniformIntegrationTTL, envVarUniformTTLDefault))
	if err != nil {
		log.WithError(err).Error("could not setup TTL index for uniform repo entries")
	}

	serviceManager := handler.NewServiceManager(
		projectMVRepo,
		common.NewGitConfigurationStore(csEndpoint.String()),
		uniformRepo,
	)

	stageManager := handler.NewStageManager(projectMVRepo)

	debugManager := handler.NewDebugManager(createEventsRepo(), createStateRepo(), createProjectRepo())

	eventDispatcher := handler.NewEventDispatcher(createEventsRepo(), createEventQueueRepo(), sequenceExecutionRepo, eventSender, time.Duration(env.EventDispatchIntervalSec)*time.Second)
	sequenceDispatcher := handler.NewSequenceDispatcher(
		createEventsRepo(),
		createSequenceQueueRepo(),
		sequenceExecutionRepo,
		getDurationFromEnvVar(env.SequenceDispatchIntervalSec, envVarSequenceDispatchIntervalSecDefault),
		clock.New(),
		common.SDModeRW,
	)

	sequenceTimeoutChannel := make(chan apimodels.SequenceTimeout)

	shipyardRetriever := handler.NewShipyardRetriever(
		common.NewGitConfigurationStore(csEndpoint.String()),
		projectMVRepo,
	)
	shipyardController := handler.GetShipyardControllerInstance(
		ctx,
		eventDispatcher,
		sequenceDispatcher,
		sequenceTimeoutChannel,
		shipyardRetriever,
	)

	engine := gin.Default()

	/// setting up middleware to handle graceful shutdown
	wg := &sync.WaitGroup{}
	engine.Use(handler.GracefulShutdownMiddleware(wg))

	apiV1 := engine.Group("/v1")
	apiHealth := engine.Group("")

	projectService := handler.NewProjectHandler(projectManager, eventSender, env, repositoryProvisioner)

	projectController := controller.NewProjectController(projectService)
	projectController.Inject(apiV1)

	serviceHandler := handler.NewServiceHandler(serviceManager, eventSender, env)
	serviceController := controller.NewServiceController(serviceHandler)
	serviceController.Inject(apiV1)

	eventHandler := handler.NewEventHandler(shipyardController)
	eventController := controller.NewEventController(eventHandler)
	eventController.Inject(apiV1)

	stageHandler := handler.NewStageHandler(stageManager)
	stageController := controller.NewStageController(stageHandler)
	stageController.Inject(apiV1)

	debugEngine := gin.Default()
	apiDebug := debugEngine.Group("/")
	debugHandler := handler.NewDebugHandler(debugManager)
	debugController := controller.NewDebugController(debugHandler)
	debugController.Inject(apiDebug)

	evaluationManager, err := handler.NewEvaluationManager(eventSender, projectMVRepo)
	if err != nil {
		log.Fatal(err)
	}
	evaluationHandler := handler.NewEvaluationHandler(evaluationManager)
	evaluationController := controller.NewEvaluationController(evaluationHandler)
	evaluationController.Inject(apiV1)

	stateHandler := handler.NewStateHandler(db.NewMongoDBStateRepo(db.GetMongoDBConnectionInstance()), shipyardController)
	stateController := controller.NewStateController(stateHandler)
	stateController.Inject(apiV1)

	sequenceStateMaterializedView := sequencehooks.NewSequenceStateMaterializedView(createStateRepo())
	shipyardController.AddSequenceTriggeredHook(sequenceStateMaterializedView)
	shipyardController.AddSequenceStartedHook(sequenceStateMaterializedView)
	shipyardController.AddSequenceWaitingHook(sequenceStateMaterializedView)
	shipyardController.AddSequenceTaskTriggeredHook(sequenceStateMaterializedView)
	shipyardController.AddSequenceTaskTriggeredHook(projectMVRepo)
	shipyardController.AddSequenceTaskStartedHook(sequenceStateMaterializedView)
	shipyardController.AddSequenceTaskStartedHook(projectMVRepo)
	shipyardController.AddSequenceTaskFinishedHook(sequenceStateMaterializedView)
	shipyardController.AddSequenceTaskFinishedHook(projectMVRepo)
	shipyardController.AddSubSequenceFinishedHook(sequenceStateMaterializedView)
	shipyardController.AddSequenceFinishedHook(sequenceStateMaterializedView)
	shipyardController.AddSequenceTimeoutHook(sequenceStateMaterializedView)
	shipyardController.AddSequenceAbortedHook(sequenceStateMaterializedView)
	shipyardController.AddSequenceTimeoutHook(eventDispatcher)
	shipyardController.AddSequencePausedHook(sequenceStateMaterializedView)
	shipyardController.AddSequenceResumedHook(sequenceStateMaterializedView)

	taskStartedWaitDuration := getDurationFromEnvVar(env.TaskStartedWaitDuration, envVarTaskStartedWaitDurationDefault)

	watcher := handler.NewSequenceWatcher(
		sequenceTimeoutChannel,
		createEventsRepo(),
		createEventQueueRepo(),
		createProjectRepo(),
		taskStartedWaitDuration,
		getDurationFromEnvVar(env.SequenceWatcherInterval, envVarSequenceWatcherIntervalDefault),
		clock.New(),
	)

	watcher.Run(ctx)

	uniformHandler := handler.NewUniformIntegrationHandler(uniformRepo)
	uniformController := controller.NewUniformIntegrationController(uniformHandler)
	uniformController.Inject(apiV1)

	sequenceExecutionHandler := handler.NewSequenceExecutionHandler(sequenceExecutionRepo, createProjectRepo())
	sequenceExecutionController := controller.NewSequenceExecutionController(sequenceExecutionHandler)
	sequenceExecutionController.Inject(apiV1)

	logRepo := createLogRepo()
	err = logRepo.SetupTTLIndex(getDurationFromEnvVar(env.LogTTL, envVarLogsTTLDefault))
	if err != nil {
		log.WithError(err).Error("could not setup TTL index for log repo entries")
	}
	logHandler := handler.NewLogHandler(handler.NewLogManager(logRepo))
	logController := controller.NewLogController(logHandler)
	logController.Inject(apiV1)

	log.Info("Migrating project git credentials")
	projectCredentialsMigrator := migration.NewProjectCredentialsMigrator(db.GetMongoDBConnectionInstance(), secretStore)
	err = projectCredentialsMigrator.Transform()
	if err != nil {
		log.Errorf("Unable to transform project git credentials: %v", err)
	}
	log.Info("Finished migrating project git credentials")

	log.Info("Migrating project key format")
	projectsMigrator := migration.NewProjectMVMigrator(db.GetMongoDBConnectionInstance())
	err = projectsMigrator.MigrateKeys()
	if err != nil {
		log.Errorf("Unable to run projects migrator: %v", err)
	}
	log.Info("Finished migrating project key format")

	log.Info("Migrating sequence execution format")
	sequenceExecutionMigrator := migration.NewSequenceExecutionMigrator(db.GetMongoDBConnectionInstance())
	err = sequenceExecutionMigrator.Run()
	if err != nil {
		log.Errorf("Unable to run sequence execution migrator: %v", err)
	}
	log.Info("Finished migrating sequence execution format")

	healthHandler := handler.NewHealthHandler()
	healthController := controller.NewHealthController(healthHandler)
	healthController.Inject(apiHealth)

	engine.Static("/swagger-ui", "./swagger-ui")
	srv := &http.Server{
		Addr:    ":8080",
		Handler: engine,
	}

	if err := connectionHandler.SubscribeToTopics([]string{"sh.keptn.>"}, nats.NewKeptnNatsMessageHandler(shipyardController.HandleIncomingEvent)); err != nil {
		log.Fatalf("Could not subscribe to nats: %v", err)
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.WithError(err).Error("could not start API server")
		}
	}()

	debugEngine.Run("localhost:9090")

	if env.DisableLeaderElection {
		// single shipyard
		shipyardController.StartDispatchers(ctx, common.SDModeRW)
	} else {
		// multiple shipyards
		go leaderelection.LeaderElection(kubeAPI.CoordinationV1(), ctx, shipyardController.StartDispatchers, shipyardController.StopDispatchers)
	}

	operationsEngine := gin.New()

	operationsV1 := operationsEngine.Group("/operations/v1")

	operationsV1.GET("/pre-stop", func(c *gin.Context) {
		log.Debug("PreStop hook has been called.")
		// invoke the cancel() function to shut down the periodically executed
		// tasks such as nats subscription, sequence watcher, sequence dispatcher, event dispatcher
		// this should ensure that no iteration of either of these tasks is attempted to be started right before the termination of the pod
		cancel()
		log.Debugf("PreStop: Sleeping for %d seconds", env.PreStopHookTime)
		<-time.After(time.Duration(env.PreStopHookTime) * time.Second)
		log.Debug("PreStop hook has been finished")
		c.Status(http.StatusOK)
	})

	operationsSrv := &http.Server{
		Addr:    ":8081",
		Handler: operationsEngine,
	}

	go func() {
		if err := operationsSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.WithError(err).Error("could not start API server")
		}
	}()

	GracefulShutdown(wg, srv)
}

func GracefulShutdown(wg *sync.WaitGroup, srv *http.Server) {
	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	wg.Wait()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting")
}

func createProjectMVRepo() *db.MongoDBProjectMVRepo {
	return db.NewProjectMVRepo(db.NewMongoDBKeyEncodingProjectsRepo(db.GetMongoDBConnectionInstance()), db.NewMongoDBEventsRepo(db.GetMongoDBConnectionInstance()))
}

func createUniformRepo() *db.MongoDBUniformRepo {
	return db.NewMongoDBUniformRepo(db.GetMongoDBConnectionInstance())
}

func createStateRepo() *db.MongoDBStateRepo {
	return db.NewMongoDBStateRepo(db.GetMongoDBConnectionInstance())
}

func createProjectRepo() *db.MongoDBKeyEncodingProjectsRepo {
	return db.NewMongoDBKeyEncodingProjectsRepo(db.GetMongoDBConnectionInstance())
}

func createEventsRepo() *db.MongoDBEventsRepo {
	return db.NewMongoDBEventsRepo(db.GetMongoDBConnectionInstance())
}

func createSequenceExecutionRepo() *db.MongoDBSequenceExecutionRepo {
	return db.NewMongoDBSequenceExecutionRepo(db.GetMongoDBConnectionInstance())
}

func createSequenceQueueRepo() *db.MongoDBSequenceQueueRepo {
	return db.NewMongoDBSequenceQueueRepo(db.GetMongoDBConnectionInstance())
}

func createEventQueueRepo() *db.MongoDBEventQueueRepo {
	return db.NewMongoDBEventQueueRepo(db.GetMongoDBConnectionInstance())
}

func createSecretStore(kubeAPI kubernetes.Interface) *common.K8sSecretStore {
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

func getDurationFromEnvVar(durationString, fallbackValue string) time.Duration {
	var duration time.Duration
	var err error
	if durationString != "" {
		duration, err = time.ParseDuration(durationString)
		if err != nil {
			log.Errorf("could not parse log %s env var %s: %s. Will use default value %s", durationString, duration, err.Error(), fallbackValue)
		}
	}

	if duration.Seconds() == 0 {
		duration, err = time.ParseDuration(fallbackValue)
		if err != nil {
			log.Errorf("could not parse default duration string %s. %s will be set to 0", err.Error(), durationString)
			return time.Duration(0)
		}
	}
	return duration
}
