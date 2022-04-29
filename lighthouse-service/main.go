package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/kelseyhightower/envconfig"
	"github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/cp-connector/pkg/api"
	"github.com/keptn/keptn/cp-connector/pkg/controlplane"
	"github.com/keptn/keptn/cp-connector/pkg/nats"
	"github.com/keptn/keptn/lighthouse-service/event_handler"
	"github.com/pkg/errors"
)

type envConfig struct {
	// Port on which to listen for cloudevents
	Port                    int    `envconfig:"RCV_PORT" default:"8080"`
	Path                    string `envconfig:"RCV_PATH" default:"/"`
	ConfigurationServiceURL string `envconfig:"CONFIGURATION_SERVICE" default:"http://configuration-service:8080"`
	K8SDeploymentName       string `envconfig:"K8S_DEPLOYMENT_NAME" default:""`
	K8SPodName              string `envconfig:"K8S_POD_NAME" default:""`
	K8SNamespace            string `envconfig:"K8S_NAMESPACE" default:""`
	K8SNodeName             string `envconfig:"K8S_NODE_NAME" default:""`
}

func main() {
	var env envConfig
	if err := envconfig.Process("", &env); err != nil {
		log.Fatalf("Failed to process env var: %s", err)
	}

	api, err := api.NewInternal(nil)
	if err != nil {
		log.Fatal(err)
	}

	subscriptionSource := controlplane.NewUniformSubscriptionSource(api.UniformV1())
	natsConnector, err := nats.ConnectFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	eventSource := controlplane.NewNATSEventSource(natsConnector)

	controlPlane := controlplane.New(subscriptionSource, eventSource)
	err = controlPlane.Register(context.TODO(), LighthouseService{env})
	if err != nil {
		log.Fatal(err)
	}
}

type LighthouseService struct {
	env envConfig
}

func (l LighthouseService) OnEvent(ctx context.Context, event models.KeptnContextExtendedCE) error {
	ce := v0_2_0.ToCloudEvent(event)
	handler, err := event_handler.NewEventHandler(ctx, ce)

	if err != nil {
		log.Println(err.Error())
		return errors.Wrap(controlplane.ErrEventHandleFatal, err.Error())
	}
	if handler != nil {
		return handler.HandleEvent(ctx)
	}

	return nil
}

func (l LighthouseService) RegistrationData() controlplane.RegistrationData {
	return controlplane.RegistrationData{
		Name: "lighthouse-service-2",
		MetaData: models.MetaData{
			Hostname:           l.env.K8SNodeName,
			IntegrationVersion: "dev",
			DistributorVersion: "0.15.0",
			Location:           "control-plane",
			KubernetesMetaData: models.KubernetesMetaData{
				Namespace:      l.env.K8SNamespace,
				PodName:        l.env.K8SPodName,
				DeploymentName: l.env.K8SDeploymentName,
			},
		},
		//TODO: read initial event subscription data from env
		Subscriptions: []models.EventSubscription{
			{
				Event:  "sh.keptn.event.evaluation.triggered",
				Filter: models.EventSubscriptionFilter{},
			},
			{
				Event:  "sh.keptn.event.get-sli.finished",
				Filter: models.EventSubscriptionFilter{},
			},
			{
				Event:  "sh.keptn.event.monitoring.configure",
				Filter: models.EventSubscriptionFilter{},
			},
		},
	}
}

func RunHealthEndpoint(port string) {

	http.HandleFunc("/health", healthHandler)
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		log.Println(err)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	status := StatusBody{Status: "OK"}

	body, err := status.ToJSON()
	if err != nil {
		log.Println(err)
	}

	w.Header().Set("content-type", "application/json")

	_, err = w.Write(body)
	if err != nil {
		log.Println(err)
	}
}

type StatusBody struct {
	Status string `json:"status"`
}

// ToJSON converts object to JSON string
func (s *StatusBody) ToJSON() ([]byte, error) {
	return json.Marshal(s)
}
