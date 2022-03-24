package main

import (
	"context"
	"github.com/kelseyhightower/envconfig"
	"github.com/keptn/keptn/resource-service/common"
	nats2 "github.com/keptn/keptn/resource-service/handler/nats"
	"github.com/keptn/keptn/resource-service/pkg/nats/subscriber"
	"io/ioutil"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/keptn/go-utils/pkg/common/osutils"
	"github.com/keptn/keptn/resource-service/config"
	"github.com/keptn/keptn/resource-service/controller"
	"github.com/keptn/keptn/resource-service/handler"
	log "github.com/sirupsen/logrus"
)

// @title Resource Service API
// @version develop
// @description This is the API documentation of the Resource Service.

// @securityDefinitions.apiKey ApiKeyAuth
// @in header
// @name x-token

// @contact.name Keptn Team
// @contact.url http://www.keptn.sh

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath /v1

const envVarLogLevel = "LOG_LEVEL"
const eventProjectDeleteFinished = "sh.keptn.event.project.delete.finished"

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := envconfig.Process("", &config.Global); err != nil {
		log.Errorf("Failed to process env var: %v", err)
		os.Exit(1)
	}

	log.SetLevel(log.InfoLevel)

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

	engine := gin.Default()
	engine.UnescapePathValues = false // To be compatible with current configuration-service
	engine.UseRawPath = true
	/// setting up middleware to handle graceful shutdown
	wg := &sync.WaitGroup{}
	engine.Use(handler.GracefulShutdownMiddleware(wg))

	apiV1 := engine.Group("/v1")
	apiHealth := engine.Group("")

	kubeAPI, err := createKubeAPI()
	if err != nil {
		log.Fatalf("could not create kubernetes client: %s", err.Error())
	}

	credentialReader := common.NewK8sCredentialReader(kubeAPI)
	fileSystem := common.NewFileSystem(common.GetConfigDir())

	git := common.NewGit(&common.GogitReal{})
	configurationContext := createConfigurationContext(git, fileSystem)

	projectManager := handler.NewProjectManager(git, credentialReader, fileSystem)
	projectHandler := handler.NewProjectHandler(projectManager)
	projectController := controller.NewProjectController(projectHandler)
	projectController.Inject(apiV1)

	stageManager := createStageManager(configurationContext, git, fileSystem, credentialReader)
	stageHandler := handler.NewStageHandler(stageManager)
	stageController := controller.NewStageController(stageHandler)
	stageController.Inject(apiV1)

	serviceManager := handler.NewServiceManager(git, credentialReader, fileSystem, configurationContext)
	serviceHandler := handler.NewServiceHandler(serviceManager)
	serviceController := controller.NewServiceController(serviceHandler)
	serviceController.Inject(apiV1)

	projectResourceManager := handler.NewResourceManager(git, credentialReader, fileSystem, configurationContext)
	projectResourceHandler := handler.NewProjectResourceHandler(projectResourceManager)
	projectResourceController := controller.NewProjectResourceController(projectResourceHandler)
	projectResourceController.Inject(apiV1)

	stageResourceManager := handler.NewResourceManager(git, credentialReader, fileSystem, configurationContext)
	stageResourceHandler := handler.NewStageResourceHandler(stageResourceManager)
	stageResourceController := controller.NewStageResourceController(stageResourceHandler)
	stageResourceController.Inject(apiV1)

	serviceResourceManager := handler.NewResourceManager(git, credentialReader, fileSystem, configurationContext)
	serviceResourceHandler := handler.NewServiceResourceHandler(serviceResourceManager)
	serviceResourceController := controller.NewServiceResourceController(serviceResourceHandler)
	serviceResourceController.Inject(apiV1)

	healthHandler := handler.NewHealthHandler()
	healthController := controller.NewHealthController(healthHandler)
	healthController.Inject(apiHealth)

	engine.Static("/swagger-ui", "./swagger-ui")
	srv := &http.Server{
		Addr:    ":8080",
		Handler: engine,
	}

	sub, err := subscriber.ConnectFromEnv()
	if err != nil {
		log.Fatal(err)
	}

	if err := sub.Subscribe(eventProjectDeleteFinished, nats2.EventHandler(projectManager).Process); err != nil {
		log.Fatal(err)
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.WithError(err).Error("could not start API server")
		}
	}()

	gracefulShutdown(ctx, wg, srv)

}

func createConfigurationContext(git *common.Git, fileSystem *common.FileSystem) handler.IConfigurationContext {
	var configContext handler.IConfigurationContext
	if config.Global.DirectoryStageStructure {
		configContext = handler.NewDirectoryConfigurationContext(git, fileSystem)
	} else {
		configContext = handler.NewBranchConfigurationContext(git, fileSystem)
	}
	return configContext
}

func createStageManager(configurationContext handler.IConfigurationContext, git common.IGit, fileSystem common.IFileSystem, credentialReader common.CredentialReader) handler.IStageManager {
	var stageManager handler.IStageManager
	if config.Global.DirectoryStageStructure {
		stageManager = handler.NewDirectoryStageManager(configurationContext, fileSystem, credentialReader, git)
	} else {
		stageManager = handler.NewStageManager(git, credentialReader)
	}
	return stageManager
}

func gracefulShutdown(ctx context.Context, wg *sync.WaitGroup, srv *http.Server) {
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	wg.Wait()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting")
}

func createKubeAPI() (*kubernetes.Clientset, error) {
	var cfg *rest.Config
	cfg, err := rest.InClusterConfig()

	if err != nil {
		return nil, err
	}

	kubeAPI, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}
	return kubeAPI, nil
}
