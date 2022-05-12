# cp-connector

`cp-connector` is a **GO** library that can be used to implement a Keptn Service.

## Example

```go
package main

import (
	"context"
	"fmt"
	"github.com/keptn/go-utils/pkg/api/models"
	api "github.com/keptn/go-utils/pkg/api/utils"
	"github.com/keptn/keptn/cp-connector/pkg/controlplane"
	"github.com/keptn/keptn/cp-connector/pkg/nats"
	"log"
)

const (
	Endpoint = ""
	Token    = ""
)

func main() {

	if Endpoint == "" || Token == "" {
		log.Fatal("Please set Keptn API endpoint and API Token to use this example")
	}

	// 1. create keptn client
	keptnAPI, err := api.New(Endpoint, api.WithAuthToken(Token))
	if err != nil {
		log.Fatal(err)
	}

	// 2. create a subscription source
	subscriptionSource := controlplane.NewUniformSubscriptionSource(keptnAPI.UniformV1())

	// 3. create an event source (either NATS of HTTP,...)
	natsConnector, err := nats.Connect("nats://localhost:4222")
	if err != nil {
		log.Fatal(err)
	}

	eventSource := controlplane.NewNATSEventSource(natsConnector)

	// 4. create control plane object and register yourself as an "integration"
	//NOTE: if log forwarding is not needed in your service, pass `nil` instead of `keptnAPI.LogsV1()`
	controlPlane := controlplane.New(subscriptionSource, eventSource, keptnAPI.LogsV1())
	err = controlPlane.Register(context.TODO(), LocalService{})
	if err != nil {
		log.Fatal(err)
	}
}

// 4. Example integration
type LocalService struct{}

func (e LocalService) OnEvent(ctx context.Context, event models.KeptnContextExtendedCE) error {
	fmt.Println("Got an event: " + *event.Type + " :)")

	return nil
}

func (e LocalService) RegistrationData() controlplane.RegistrationData {
	return controlplane.RegistrationData{
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
		Subscriptions: []models.EventSubscription{},
	}
}

```