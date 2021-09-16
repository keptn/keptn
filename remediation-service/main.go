package main

import (
	api "github.com/keptn/go-utils/pkg/api/utils"
	"github.com/keptn/keptn/go-sdk/pkg/sdk"
	"github.com/keptn/keptn/remediation-service/handler"
	"log"
)

const getActionTriggeredEventType = "sh.keptn.event.get-action.triggered"
const serviceName = "remediation-service"

func main() {

	go api.RunHealthEndpoint("10998")
	log.Fatal(sdk.NewKeptn(
		serviceName,
		sdk.WithHandler(
			getActionTriggeredEventType,
			handler.NewGetActionEventHandler()),
	).Start())
}
