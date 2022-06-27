# cp-connector

`cp-connector` is a **GO** library that can be used to implement a Keptn Service.

## Example

```go
package main

import (
	"context"
	"github.com/benbjohnson/clock"
	"github.com/keptn/go-utils/pkg/api/models"
	api "github.com/keptn/go-utils/pkg/api/utils"
	"github.com/keptn/keptn/cp-connector/pkg/controlplane"
	"github.com/keptn/keptn/cp-connector/pkg/eventsource/http"
	"github.com/keptn/keptn/cp-connector/pkg/logforwarder"
	"github.com/keptn/keptn/cp-connector/pkg/subscriptionsource"
	"github.com/keptn/keptn/cp-connector/pkg/types"
	"github.com/sirupsen/logrus"
	"log"
	"time"
)

const (
	Endpoint = ""
	Token    = ""
)

func main() {
	if Endpoint == "" || Token == "" {
		log.Fatal("Please set Keptn API endpoint and API Token to use this example")
	}

	// Create your favorite logger (e.g. logrus)
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)

	// 1. create keptn client
	keptnAPI, err := api.New(Endpoint, api.WithAuthToken(Token))
	if err != nil {
		log.Fatal(err)
	}

	// 2. create a subscription, event source and log forwarder
	subscriptionSource := subscriptionsource.New(keptnAPI.UniformV1(), subscriptionsource.WithLogger(logger))
	eventSource := http.New(clock.New(), http.NewEventAPI(keptnAPI.ShipyardControlV1(), keptnAPI.APIV1()))
	logForwarder := logforwarder.New(keptnAPI.LogsV1())

	// 3. create control plane and start it
	controlPlane := controlplane.New(subscriptionSource, eventSource, logForwarder, controlplane.WithLogger(logger))
	if err := controlplane.RunWithGracefulShutdown(controlPlane, LocalService{}, time.Second*10); err != nil {
		log.Fatal(err)
	}
}

// 4. Example integration
type LocalService struct{}

func (e LocalService) OnEvent(ctx context.Context, event models.KeptnContextExtendedCE) error {
	// handle the event and grab a sender to send back started / finished events to keptn
	// sendEvent := ctx.Value(controlplane.EventSenderKeyType{}).(types.EventSender)
	return nil
}

// RegistrationData is used for initial registration to the Keptn control plane
func (e LocalService) RegistrationData() types.RegistrationData {
	return types.RegistrationData{
		Name: "local-service",
		MetaData: models.MetaData{
			Hostname:           "localhost",
			IntegrationVersion: "dev",
			DistributorVersion: "0.15.0",
			Location:           "local",
			KubernetesMetaData: models.KubernetesMetaData{
				Namespace:      "keptn",
				PodName:        "wuppi",
				DeploymentName: "wuppi",
			},
		},
		Subscriptions: []models.EventSubscription{
			{
				Event: "sh.keptn.event.echo.triggered", // intitially subscribe to an event
			},
		},
	}
}
```