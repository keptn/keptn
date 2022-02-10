package main

import (
	"github.com/keptn/keptn/go-sdk/pkg/sdk"
	"github.com/keptn/keptn/remediation-service/handler"
	"github.com/sirupsen/logrus"
	"log"
)

const getActionTriggeredEventType = "sh.keptn.event.get-action.triggered"
const approvalTriggeredEventType = "sh.keptn.event.approval.triggered"

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

	log.Fatal(sdk.NewKeptn(
		serviceName,
		options...,
	).Start())
}
