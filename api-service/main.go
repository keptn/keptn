package main

import (
	"context"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/benbjohnson/clock"
	"github.com/gin-gonic/gin"
	"github.com/keptn/go-utils/pkg/common/osutils"
	"github.com/keptn/keptn/api-service/backend"
	"github.com/keptn/keptn/api-service/controller"
	_ "github.com/keptn/keptn/api-service/docs"
	"github.com/keptn/keptn/api-service/handler"
	"github.com/keptn/keptn/api-service/repository"
	log "github.com/sirupsen/logrus"
)

// @title Secret Service API
// @version develop
// @description This is the API documentation of the Secret Service.

// @securityDefinitions.apiKey ApiKeyAuth
// @in header
// @name x-token

// @contact.name Keptn Team
// @contact.url http://www.keptn.sh

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath /v1

const envVarLogLevel = "LOG_LEVEL"
const authRequestsPerSecond = "MAX_AUTH_REQUESTS_PER_SECOND"
const authRequestMaxBurst = "MAX_AUTH_REQUESTS_BURST"

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

	requestsPerSecond, err := strconv.ParseFloat(os.Getenv(authRequestsPerSecond), 32)
	if err != nil {
		log.WithError(err).Error("could not parse max auth requests per second provided by 'MAX_AUTH_REQUESTS_PER_SECOND' env var")
	}
	requestsMaxBurst, err := strconv.Atoi(os.Getenv(authRequestMaxBurst))
	if err != nil {
		log.WithError(err).Error("could not parse max auth requests burst provided by 'MAX_AUTH_REQUESTS_BURST' env var")
	}

	if _, err := os.Stat(repository.ScopesConfigurationFile); os.IsNotExist(err) {
		log.Fatalf("Scopes configuration file not found: %s", repository.ScopesConfigurationFile)
	}

	if osutils.GetAndCompareOSEnv("GIN_MODE", "release") {
		// disable GIN request logging in release mode
		gin.SetMode("release")
		gin.DefaultWriter = ioutil.Discard
	}

	log.Infof("Registered Backends: %v", backend.GetRegisteredBackends())

	engine := gin.Default()
	apiV1 := engine.Group("/v1")

	// only kubernetes supported, so we hard code it for now
	secretsBackend := backend.CreateBackend("kubernetes")
	secretController := controller.NewSecretController(handler.NewSecretHandler(secretsBackend))
	secretController.Inject(apiV1)

	metadataController := controller.NewMetadataController(handler.NewMetadataHandler())
	metadataController.Inject(apiV1)

	authController := controller.NewAuthController(handler.NewAuthHandler(requestsPerSecond, requestsMaxBurst, clock.New()))
	authController.Inject(apiV1)

	eventController := controller.NewEventController(handler.NewEventHandler())
	eventController.Inject(apiV1)

	engine.GET("/health", func(c *gin.Context) { c.Status(http.StatusOK) })

	engine.Static("/swagger-ui", "./swagger-ui")
	srv := &http.Server{
		Addr:    ":8080",
		Handler: engine,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Unable to start service: %s", err.Error())
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting")
}
