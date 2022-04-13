package main

import (
	"context"
	"github.com/keptn/keptn/statistics-service/config"
	"github.com/keptn/keptn/statistics-service/db"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/keptn/keptn/statistics-service/api"
	"github.com/keptn/keptn/statistics-service/controller"
	_ "github.com/keptn/keptn/statistics-service/docs" // docs is generated by Swag CLI, you have to import it.
)

// @title Statistics Service API
// @version develop
// @description This is the API documentation of the Statistics Service.

// @securityDefinitions.apiKey ApiKeyAuth
// @in header
// @name x-token

// @contact.name Keptn Team
// @contact.url http://www.keptn.sh

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath /v1

func main() {
	envConfig := config.GetConfig()
	logLevel, err := log.ParseLevel(envConfig.LogLevel)
	if err != nil {
		log.WithError(err).Error("could not parse log level provided by 'LOG_LEVEL' env var")
		log.SetLevel(log.InfoLevel)
	} else {
		log.SetLevel(logLevel)
	}

	if !envConfig.DataMigrationDisabled {
		// data migration
		go func() {
			log.Infof("Migrating data (%d entries every %d seconds)", envConfig.DataMigrationBatchSize, envConfig.DataMigrationIntervalSec)
			migrator := db.NewMigrator(envConfig.DataMigrationBatchSize, time.Second*time.Duration(envConfig.DataMigrationIntervalSec))
			_, err := migrator.Run(context.Background())
			if err != nil {
				log.Errorf("Error during migration: %v", err)
			}
			log.Info("Migration finished")
		}()
	}

	controller.GetStatisticsBucketInstance()

	router := gin.New()
	router.Use(gin.Recovery())
	/// setting up middleware to handle graceful shutdown
	wg := &sync.WaitGroup{}
	router.Use(controller.GracefulShutdownMiddleware(wg))

	if os.Getenv("GIN_MODE") == "release" {
		// disable GIN request logging in release mode
		gin.SetMode("release")
		gin.DefaultWriter = ioutil.Discard
	}

	apiV1 := router.Group("/v1")
	apiV1.GET("/statistics", api.GetStatistics)

	apiV1.POST("/event", api.HandleEvent)
	router.Static("/swagger-ui", "./swagger-ui")

	apiHealth := router.Group("")
	apiHealth.GET("/health", func(c *gin.Context) { c.Status(http.StatusOK) })

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
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
