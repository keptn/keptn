package main

import (
	"context"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/keptn/go-utils/pkg/common/osutils"
	"github.com/keptn/keptn/shipyard-controller/controller"
	"github.com/keptn/keptn/shipyard-controller/db"
	_ "github.com/keptn/keptn/shipyard-controller/docs"
	"github.com/keptn/keptn/shipyard-controller/handler"
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

const envVarLogTTL = "LOG_TTL"
const envVarLogLevel = "LOG_LEVEL"
const envVarLogsTTLDefault = "120h" // 5 days

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
	engine.Use(handler.GracefulShutdownMiddleware(wg))

	apiV1 := engine.Group("/v1")
	apiHealth := engine.Group("")

	logRepo := createLogRepo()
	err := logRepo.SetupTTLIndex(getDurationFromEnvVar(envVarLogTTL, envVarLogsTTLDefault))
	if err != nil {
		log.WithError(err).Error("could not setup TTL index for log repo entries")
	}
	logHandler := handler.NewLogHandler(handler.NewLogManager(logRepo))
	logController := controller.NewLogController(logHandler)
	logController.Inject(apiV1)

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

func createLogRepo() *db.MongoDBLogRepo {
	return db.NewMongoDBLogRepo(db.GetMongoDBConnectionInstance())
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
