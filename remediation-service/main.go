package main

import (
	"github.com/keptn/keptn/remediation-service/handler"
	"github.com/keptn/keptn/remediation-service/internal/sdk"
	"log"
)

const getActionTriggeredEventType = "sh.keptn.event.get-action.triggered"
const serviceName = "remediation-service"

func main() {
	log.Fatal(sdk.NewKeptn(
		sdk.GetHTTPClientFromEnv(),
		serviceName,
		sdk.WithHandler(getActionTriggeredEventType, handler.NewGetActionEventHandler()),
	).Start())
}
