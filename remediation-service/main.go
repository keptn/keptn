package main

import (
	"github.com/keptn/keptn/go-sdk/pkg/sdk"
	"github.com/keptn/keptn/remediation-service/handler"
	"github.com/sirupsen/logrus"
)

const getActionTriggeredEventType = "sh.keptn.event.get-action.triggered"
const approvalTriggeredEventType = "sh.keptn.event.approval.triggered"
const evaluationTriggeredEventType = "sh.keptn.event.evaluation.triggered"
const getSLIFinishedEventType = "sh.keptn.event.get-sli.finished"
const monitoringConfigureEventType = "sh.keptn.event.monitoring.configure"

const serviceName = "remediation-service" //TODO change me and deployment names

func main() {

	var options []sdk.KeptnOption
	options = append(options, sdk.WithLogger(logrus.New()),
		sdk.WithAutomaticResponse(false))

	if true { //TODO shall we make the UC configurable? OS variable could go here
		options = append(options, sdk.WithTaskHandler(
			getActionTriggeredEventType,
			handler.NewGetActionEventHandler()))
	}

	if true { //TODO
		options = append(options, sdk.WithTaskHandler(
			approvalTriggeredEventType,
			handler.NewApprovalTriggeredEventHandler()))
	}

	if true { //TODO
		configurationHandler, err := handler.NewConfigureMonitoringHandler()
		if err != nil {
			logrus.Fatalf("could not start configuration handler: %s", err.Error())
		}

		options = append(options,
			sdk.WithTaskHandler(
				evaluationTriggeredEventType,
				configurationHandler), //TODO this is a wrong endpoint change me when lighthouse stops sending triggered events
			sdk.WithTaskHandler(
				getSLIFinishedEventType,
				configurationHandler), // TODO his is a wrong endpoint change me when lighthouse stops sending triggered events
			sdk.WithTaskHandler(
				monitoringConfigureEventType,
				configurationHandler))

	}

	logrus.Fatal(sdk.NewKeptn(
		serviceName,
		options...,
	).Start())
}
