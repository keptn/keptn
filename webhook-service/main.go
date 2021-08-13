package main

import (
	"github.com/keptn/keptn/go-sdk/pkg/sdk"
	"github.com/keptn/keptn/webhook-service/handler"
	"log"
)

const eventTypeWildcard = "*"
const serviceName = "webhook-service"

func main() {
	log.Fatal(sdk.NewKeptn(
		serviceName,
		sdk.WithHandler(
			eventTypeWildcard,
			handler.NewTaskHandler(),
			map[string]interface{}{},
		),
	).Start())
}
