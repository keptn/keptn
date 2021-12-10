package main

import (
	"github.com/keptn/keptn/go-sdk/pkg/sdk"
	"github.com/keptn/keptn/remediation-service/handler"
	"github.com/sirupsen/logrus"
	"log"
)

const getActionTriggeredEventType = "sh.keptn.event.get-action.triggered"
const serviceName = "remediation-service"

func main() {
	log.Fatal(sdk.NewKeptn(
		serviceName,
		sdk.WithTaskHandler(
			getActionTriggeredEventType,
			handler.NewGetActionEventHandler()),
		sdk.WithLogger(logrus.New()),
	).Start())
}
