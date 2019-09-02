package main

import (
	"context"
	"errors"
	"log"
	"net/url"
	"os"

	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/client"
	cloudeventshttp "github.com/cloudevents/sdk-go/pkg/cloudevents/transport/http"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/types"
	"github.com/google/uuid"
	"github.com/kelseyhightower/envconfig"
	keptnevents "github.com/keptn/go-utils/pkg/events"
	keptnutils "github.com/keptn/go-utils/pkg/utils"
	"github.com/keptn/keptn/helm-service/controller"
	"github.com/keptn/keptn/helm-service/controller/helm"
	"github.com/keptn/keptn/helm-service/controller/mesh"
)

type envConfig struct {
	// Port on which to listen for cloudevents
	Port int    `envconfig:"RCV_PORT" default:"8080"`
	Path string `envconfig:"RCV_PATH" default:"/"`
}

func main() {
	var env envConfig
	if err := envconfig.Process("", &env); err != nil {
		log.Fatalf("Failed to process env var: %s", err)
	}
	os.Exit(_main(os.Args[1:], env))
}

func getConfigurationServiceURL() string {
	if os.Getenv("env") == "production" {
		return "configuration-service.keptn.svc.cluster.local:8080"
	}
	return "localhost:8080"
}

func getKeptnDomain() (string, error) {
	useInClusterConfig := false
	if os.Getenv("env") == "production" {
		useInClusterConfig = true
	}
	return keptnutils.GetKeptnDomain(useInClusterConfig)
}

func gotEvent(ctx context.Context, event cloudevents.Event) error {
	var shkeptncontext string
	event.Context.ExtensionAs("shkeptncontext", &shkeptncontext)

	logger := keptnutils.NewLogger(shkeptncontext, event.Context.GetID(), "helm-service")
	mesh := mesh.NewIstioMesh()
	canaryLevelGen := helm.NewCanaryOnNamespaceGenerator()

	if event.Type() == keptnevents.ConfigurationChangeEventType {
		configChanger := controller.NewConfigurationChanger(mesh, canaryLevelGen, logger, getConfigurationServiceURL())
		go configChanger.ChangeAndApplyConfiguration(event)
	} else if event.Type() == keptnevents.InternalServiceCreateEventType {
		keptnDomain, err := getKeptnDomain()
		if err == nil {
			onboarder := controller.NewOnboarder(mesh, canaryLevelGen, logger, getConfigurationServiceURL(), keptnDomain)
			go onboarder.DoOnboard(event)
		} else {
			logger.Error("Error when reading the keptn domain")
		}
	} else {
		logger.Error("Received unexpected keptn event")
	}

	return nil
}

func sendDeploymentFinishedEvent(shkeptncontext string, incomingEvent cloudevents.Event) error {

	source, _ := url.Parse("helm-service")
	contentType := "application/json"

	event := cloudevents.Event{
		Context: cloudevents.EventContextV02{
			ID:          uuid.New().String(),
			Type:        "sh.keptn.events.deployment-finished",
			Source:      types.URLRef{URL: *source},
			ContentType: &contentType,
			Extensions:  map[string]interface{}{"shkeptncontext": shkeptncontext},
		}.AsV02(),
		Data: incomingEvent.Data,
	}

	t, err := cloudeventshttp.New(
		cloudeventshttp.WithTarget("http://event-broker.keptn.svc.cluster.local/keptn"),
		cloudeventshttp.WithEncoding(cloudeventshttp.StructuredV02),
	)
	if err != nil {
		return errors.New("Failed to create transport:" + err.Error())
	}

	c, err := client.New(t)
	if err != nil {
		return errors.New("Failed to create HTTP client:" + err.Error())
	}

	if _, err := c.Send(context.Background(), event); err != nil {
		return errors.New("Failed to send cloudevent:, " + err.Error())
	}
	return nil
}

func _main(args []string, env envConfig) int {
	ctx := context.Background()

	t, err := cloudeventshttp.New(
		cloudeventshttp.WithPort(env.Port),
		cloudeventshttp.WithPath(env.Path),
	)

	if err != nil {
		log.Fatalf("failed to create transport, %v", err)
	}
	c, err := client.New(t)
	if err != nil {
		log.Fatalf("failed to create client, %v", err)
	}

	log.Fatalf("failed to start receiver: %s", c.StartReceiver(ctx, gotEvent))

	return 0
}
