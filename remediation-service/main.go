package main

import (
	"log"
	"os"

	"github.com/keptn/keptn/go-sdk/pkg/sdk"
	"github.com/keptn/keptn/remediation-service/handler"
	"github.com/sirupsen/logrus"
)

const getActionTriggeredEventType = "sh.keptn.event.get-action.triggered"
const serviceName = "remediation-service"
const envVarLogLevel = "LOG_LEVEL"

func main() {
	if os.Getenv(envVarLogLevel) != "" {
		logLevel, err := logrus.ParseLevel(os.Getenv(envVarLogLevel))
		if err != nil {
			logrus.WithError(err).Error("could not parse log level provided by 'LOG_LEVEL' env var")
			logrus.SetLevel(logrus.InfoLevel)
		} else {
			logrus.SetLevel(logLevel)
		}
	}

	log.Fatal(sdk.NewKeptn(
		serviceName,
		sdk.WithTaskHandler(
			getActionTriggeredEventType,
			handler.NewGetActionEventHandler()),
		sdk.WithLogger(logrus.New()),
	).Start())
}
