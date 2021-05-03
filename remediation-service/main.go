package main

import (
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/keptn/keptn/remediation-service/handler"
	"github.com/keptn/keptn/remediation-service/pkg/sdk"
)

const getActionTriggeredEventType = "sh.keptn.event.get-action.triggered"

func main() {
	httpClient := sdk.GetHTTPClient(cloudevents.WithPath("/"), cloudevents.WithPort(8080))
	keptn := sdk.NewKeptn(httpClient, "remediation-service", sdk.WithHandler(handler.NewGetActionEventHandler(), getActionTriggeredEventType))
	keptn.Start()
}
