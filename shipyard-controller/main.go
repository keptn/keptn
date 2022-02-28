package main

import (
	"context"
	"github.com/google/uuid"
	"github.com/keptn/keptn/shipyard-controller/nats"
	"io/ioutil"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/client-go/kubernetes/typed/coordination/v1"
	"k8s.io/client-go/tools/leaderelection"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/gin-gonic/gin"
	"github.com/keptn/go-utils/pkg/common/osutils"
	keptncommon "github.com/keptn/go-utils/pkg/lib/keptn"
	"github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/shipyard-controller/common"
	"github.com/keptn/keptn/shipyard-controller/controller"
	"github.com/keptn/keptn/shipyard-controller/db"
	"github.com/keptn/keptn/shipyard-controller/db/migration"
	_ "github.com/keptn/keptn/shipyard-controller/docs"
	"github.com/keptn/keptn/shipyard-controller/handler"
	"github.com/keptn/keptn/shipyard-controller/handler/sequencehooks"
	"github.com/keptn/keptn/shipyard-controller/models"
	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
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
const envVarUniformIntegrationTTL = "UNIFORM_INTEGRATION_TTL"
const envVarNatsURL = "NATS_URL"
const envVarLogTTL = "LOG_TTL"
const envVarLogLevel = "LOG_LEVEL"
const envVarEventDispatchIntervalSecDefault = "10"
const envVarSequenceDispatchIntervalSecDefault = "10s"
const envVarLogsTTLDefault = "120h" // 5 days
const envVarUniformTTLDefault = "1m"
const envVarTaskStartedWaitDurationDefault = "10m"
const envVarNatsURLDefault = "nats://keptn-nats"

func main() {

	log.SetLevel(log.InfoLevel)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if os.Getenv(envVarLogLevel) != "" {
		logLevel, err := log.ParseLevel(os.Getenv(envVarLogLevel))
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

	projectMVRepo := createProjectMVRepo()
	projectManager := handler.NewProjectManager(
		common.NewGitConfigurationStore(csEndpoint.String()),
		createSecretStore(kubeAPI),
		projectMVRepo,
		createTaskSequenceRepo(),
		createEventsRepo(),
		createSequenceQueueRepo(),
		createEventQueueRepo())

	uniformRepo := createUniformRepo()
	err = uniformRepo.SetupTTLIndex(getDurationFromEnvVar(envVarUniformIntegrationTTL, envVarUniformTTLDefault))
	if err != nil {
		log.WithError(err).Error("could not setup TTL index for uniform repo entries")
	}

	serviceManager := handler.NewServiceManager(
		projectMVRepo,
		common.NewGitConfigurationStore(csEndpoint.String()),
		uniformRepo,
	)

	stageManager := handler.NewStageManager(projectMVRepo)

	eventDispatcher := handler.NewEventDispatcher(createEventsRepo(), createEventQueueRepo(), createTaskSequenceRepo(), eventSender, time.Duration(eventDispatcherSyncInterval)*time.Second)
	sequenceDispatcher := handler.NewSequenceDispatcher(
		createEventsRepo(),
		createEventQueueRepo(),
		createSequenceQueueRepo(),
		createTaskSequenceRepo(),
		getDurationFromEnvVar(envVarSequenceDispatchIntervalSec, envVarSequenceDispatchIntervalSecDefault),
		clock.New(),
	)

	sequenceTimeoutChannel := make(chan models.SequenceTimeout)

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
	shipyardController.AddSequencePausedHook(eventDispatcher)
	shipyardController.AddSequenceResumedHook(sequenceStateMaterializedView)
	shipyardController.AddSequenceResumedHook(eventDispatcher)

	taskStartedWaitDuration := getDurationFromEnvVar(envVarTaskStartedWaitDuration, envVarTaskStartedWaitDurationDefault)

	watcher := handler.NewSequenceWatcher(
		sequenceTimeoutChannel,
		createEventsRepo(),
		createEventQueueRepo(),
		createProjectRepo(),
		taskStartedWaitDuration,
		1*time.Minute,
		clock.New(),
	)

	watcher.Run(ctx)

	uniformHandler := handler.NewUniformIntegrationHandler(uniformRepo)
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

	log.Info("Migrating project key format")
	projectsMigrator := migration.NewProjectMVMigrator(db.GetMongoDBConnectionInstance())
	err = projectsMigrator.MigrateKeys()
	if err != nil {
		log.Errorf("Unable to run projects migrator: %v", err)
	}
	log.Info("Finished migrating project key format")

	healthHandler := handler.NewHealthHandler()
	healthController := controller.NewHealthController(healthHandler)
	healthController.Inject(apiHealth)

	engine.Static("/swagger-ui", "./swagger-ui")
	srv := &http.Server{
		Addr:    ":8080",
		Handler: engine,
	}

	connectionHandler := nats.NewNatsConnectionHandler(
		ctx,
		getNatsURLFromEnvVar(),
		nats.NewKeptnNatsMessageHandler(shipyardController.HandleIncomingEvent),
	)

	if err := connectionHandler.SubscribeToTopics([]string{"sh.keptn.>"}); err != nil {
		log.Fatalf("Could not subscribe to nats: %v", err)
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.WithError(err).Error("could not start API server")
		}
	}()

	if os.Getenv(handler.EnvVarDisableLeaderElection) == "true" {
		// single shipyard
		shipyardController.StartDispatchers(ctx)
	} else {
		// multiple shipyards
		LeaderElection(kubeAPI.CoordinationV1(), ctx, shipyardController.StartDispatchers, shipyardController.StopDispatchers)
	}

	GracefulShutdown(ctx, wg, srv)

}

func LeaderElection(client v1.CoordinationV1Interface, ctx context.Context, start func(ctx context.Context), stop func()) {
	myID := uuid.New().String()
	// we use the Lease lock type since edits to Leases are less common
	// and fewer objects in the cluster watch "all Leases".
	lock := &resourcelock.LeaseLock{
		LeaseMeta: metav1.ObjectMeta{
			Name:      "shipyard-controller-dispatcher",
			Namespace: common.GetKeptnNamespace(),
		},
		Client: client,
		LockConfig: resourcelock.ResourceLockConfig{
			Identity: myID,
		},
	}

	// start the leader election code loop
	leaderelection.RunOrDie(ctx, leaderelection.LeaderElectionConfig{
		Lock: lock,
		// IMPORTANT: you MUST ensure that any code you have that
		// is protected by the lease must terminate **before**
		// you call cancel. Otherwise, you could have a background
		// loop still running and another process could
		// get elected before your background loop finished, violating
		// the stated goal of the lease.
		ReleaseOnCancel: true,
		LeaseDuration:   60 * time.Second,
		RenewDeadline:   15 * time.Second,
		RetryPeriod:     5 * time.Second,
		Callbacks: leaderelection.LeaderCallbacks{
			OnStartedLeading: func(ctx context.Context) {
				// we're notified when we start - this is where you would
				// usually put your code
				start(ctx)
			},
			OnStoppedLeading: func() {
				// we can do cleanup here
				log.Infof("leader lost: %s", myID)
				stop()
			},
			OnNewLeader: func(identity string) {
				// we're notified when a new leader is elected
				if identity == myID {
					// I just got the lock
					return
				}
				log.Infof("new leader elected: %s", identity)
				stop()
			},
		},
	})
}

func GracefulShutdown(ctx context.Context, wg *sync.WaitGroup, srv *http.Server) {
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

func getNatsURLFromEnvVar() string {
	if natsURL, ok := os.LookupEnv(envVarNatsURL); ok && natsURL != "" {
		return natsURL
	}
	return envVarNatsURLDefault
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
