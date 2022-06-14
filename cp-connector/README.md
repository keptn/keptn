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
	"github.com/keptn/keptn/cp-connector/pkg/eventsource"
	"github.com/keptn/keptn/cp-connector/pkg/logforwarder"
	"github.com/keptn/keptn/cp-connector/pkg/nats"
	"github.com/keptn/keptn/cp-connector/pkg/subscriptionsource"
	"github.com/keptn/keptn/cp-connector/pkg/types"
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
	subscriptionSource := subscriptionsource.New(keptnAPI.UniformV1())
	
	// 2.i inject your favourite logger as follows
    //subscriptionsource.WithLogger(mylogger)
	
	// 3. create an event source (either NATS of HTTP,...)
	natsConnector:= nats.New("nats://localhost:4222")
	if err != nil {
		log.Fatal(err)
	}

	eventSource := eventsource.New(natsConnector)
	logForwarder := logforwarder.New(keptnAPI.LogsV1())

	// 4. create control plane object and register yourself as an "integration"
	//NOTE: if log forwarding is not needed in your service, pass `nil` instead of the `logForwarder`
	controlPlane := controlplane.New(subscriptionSource, eventSource, logForwarder)
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
		Subscriptions: []models.EventSubscription{},
	}
}

```