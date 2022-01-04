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

// @title Control Plane API
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
	/// setting up middlewere to handle graceful shutdown
	wg := &sync.WaitGroup{}

	apiV1 := engine.Group("/v1")

	kubeAPI, err := createKubeAPI()
	if err != nil {
		log.Fatalf("could not create kubernetes client: %s", err.Error())
	}

	credentialReader := common.NewK8sCredentialReader(kubeAPI)
	fileWriter := common.NewFileSystem(common.GetConfigDir())

	projectResourceManager := handler.NewResourceManager(nil, credentialReader, fileWriter)
	projectResourceHandler := handler.NewProjectResourceHandler(projectResourceManager)
	projectResourceController := controller.NewProjectResourceController(projectResourceHandler)
	projectResourceController.Inject(apiV1)

	stageResourceManager := handler.NewResourceManager(nil, credentialReader, fileWriter)
	stageResourceHandler := handler.NewStageResourceHandler(stageResourceManager)
	stageResourceController := controller.NewStageResourceController(stageResourceHandler)
	stageResourceController.Inject(apiV1)

	serviceResourceManager := handler.NewResourceManager(nil, credentialReader, fileWriter)
	serviceResourceHandler := handler.NewServiceResourceHandler(serviceResourceManager)
	serviceResourceController := controller.NewServiceResourceController(serviceResourceHandler)
	serviceResourceController.Inject(apiV1)

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
