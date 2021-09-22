package main

import (
	"github.com/gin-gonic/gin"
	"github.com/keptn/go-utils/pkg/common/osutils"
	_ "github.com/keptn/keptn/secret-service/docs"
	"github.com/keptn/keptn/secret-service/pkg/backend"
	"github.com/keptn/keptn/secret-service/pkg/controller"
	"github.com/keptn/keptn/secret-service/pkg/handler"
	"github.com/keptn/keptn/secret-service/pkg/repository"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
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
func main() {

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
	apiHealth := engine.Group("")

	// only kubernetes supported, so we hard code it for now
	secretsBackend := backend.CreateBackend("kubernetes")
	secretController := controller.NewSecretController(handler.NewSecretHandler(secretsBackend))
	secretController.Inject(apiV1)

	apiHealth.GET("/health", func(c *gin.Context) { c.Status(http.StatusOK) })

	//go keptnapi.RunHealthEndpoint("10998")

	engine.Static("/swagger-ui", "./swagger-ui")
	err := engine.Run()
	if err != nil {
		log.Fatalf("Unable to start service: %s", err.Error())
	}
}
