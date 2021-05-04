package main

import (
	"github.com/keptn/keptn/remediation-service/handler"
	"github.com/keptn/keptn/remediation-service/pkg/sdk"
	"log"
)

const getActionTriggeredEventType = "sh.keptn.event.get-action.triggered"

func main() {
	keptn := sdk.NewKeptn(
		sdk.GetHTTPClientFromEnv(),
		"remediation-service",
		sdk.WithHandler(handler.NewGetActionEventHandler(), getActionTriggeredEventType),
	)
	log.Fatal(keptn.Start())
}
