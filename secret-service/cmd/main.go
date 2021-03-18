package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/keptn/keptn/secret-service/internal/backend"
	"github.com/keptn/keptn/secret-service/internal/controller"
	"github.com/keptn/keptn/secret-service/internal/handler"
)

func main() {

	fmt.Print("Registered Backends: ")
	fmt.Println(backend.GetRegisteredBackends())

	engine := gin.Default()
	apiV1 := engine.Group("/v1")

	backend := backend.CreateBackend("kubernetes") //only kubernetes supported, so we hard code it for now
	secretController := controller.NewSecretController(handler.NewSecretHandler(backend))

	secretController.Inject(apiV1)
	engine.Run()

}
