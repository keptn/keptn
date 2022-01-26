package main

import (
	"context"
	"github.com/keptn/keptn/resource-service/common"
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

func main() {
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

	projectManager := handler.NewProjectManager(git, credentialReader, fileSystem)
	projectHandler := handler.NewProjectHandler(projectManager)
	projectController := controller.NewProjectController(projectHandler)
	projectController.Inject(apiV1)

	stageManager := handler.NewStageManager(git, credentialReader)
	stageHandler := handler.NewStageHandler(stageManager)
	stageController := controller.NewStageController(stageHandler)
	stageController.Inject(apiV1)

	configurationContext := getConfigurationContext(git, fileSystem)

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

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.WithError(err).Error("could not start API server")
		}
	}()

	GracefulShutdown(wg, srv)

}

func getConfigurationContext(git *common.Git, fileSystem *common.FileSystem) *handler.BranchConfigurationContext {
	stageContext := handler.NewBranchConfigurationContext(git, fileSystem)
	return stageContext
}

func GracefulShutdown(wg *sync.WaitGroup, srv *http.Server) {
	// Wait for interrupt signal to gracefully shutdown the server
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
