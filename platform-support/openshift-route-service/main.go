package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	b64 "encoding/base64"
	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/client"
	cloudeventshttp "github.com/cloudevents/sdk-go/pkg/cloudevents/transport/http"
	"github.com/kelseyhightower/envconfig"
	keptnevents "github.com/keptn/go-utils/pkg/events"
	keptnmodels "github.com/keptn/go-utils/pkg/models"
	keptnutils "github.com/keptn/go-utils/pkg/utils"

	"gopkg.in/yaml.v2"
)

type envConfig struct {
	// Port on which to listen for cloudevents
	Port int    `envconfig:"RCV_PORT" default:"8080"`
	Path string `envconfig:"RCV_PATH" default:"/"`
}

type stage struct {
	Name string `json:"name"`
}
type projectData struct {
	Project string  `json:"project"`
	Stages  []stage `json:"stages"`
}

func main() {
	var env envConfig
	if err := envconfig.Process("", &env); err != nil {
		log.Fatalf("Failed to process env var: %s", err)
	}
	os.Exit(_main(os.Args[1:], env))
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

	log.Printf("will listen on :%d%s\n", env.Port, env.Path)
	log.Fatalf("failed to start receiver: %s", c.StartReceiver(ctx, gotEvent))

	return 0
}

func gotEvent(ctx context.Context, event cloudevents.Event) error {
	var shkeptncontext string
	event.Context.ExtensionAs("shkeptncontext", &shkeptncontext)

	logger := keptnutils.NewLogger(shkeptncontext, event.Context.GetID(), "openshift-route-service")

	logger.Debug(fmt.Sprintf("Got Event Context: %+v", event.Context))

	data := &keptnevents.ProjectCreateEventData{}
	if err := event.DataAs(data); err != nil {
		logger.Error(fmt.Sprintf("Got Data Error: %s", err.Error()))
		return err
	}
	if event.Type() != keptnevents.InternalProjectCreateEventType {
		const errorMsg = "Received unexpected keptn event"
		logger.Error(errorMsg)
		return errors.New(errorMsg)
	}

	go createRoutes(data, logger)
	return nil
}

func createRoutes(data *keptnevents.ProjectCreateEventData, logger *keptnutils.Logger) {
	shipyard := keptnmodels.Shipyard{}
	decodedStr, err := b64.StdEncoding.DecodeString(data.Shipyard)
	if err != nil {
		logger.Error("Could not parse shipyard content: " + err.Error())
		return
	}
	err = yaml.Unmarshal([]byte(decodedStr), &shipyard)
	if err != nil {
		logger.Error("Could not parse shipyard content: " + err.Error())
		return
	}
	for _, stage := range shipyard.Stages {
		exposeRoute(data.Project, stage.Name, logger)
	}
}

func exposeRoute(project string, stage string, logger *keptnutils.Logger) {
	appDomain := os.Getenv("APP_DOMAIN")
	if appDomain == "" {
		logger.Error("No app domain defined. Cannot create route.")
		return
	}
	// oc create route edge istio-wildcard-ingress-secure-keptn --service=istio-ingressgateway --hostname="www.keptn.ingress-gateway.$BASE_URL" --port=http2 --wildcard-policy=Subdomain --insecure-policy='Allow'
	logger.Info("Trying to create route www." + stage + "." + project + "." + appDomain)
	out, err := keptnutils.ExecuteCommand("oc",
		[]string{
			"create",
			"route",
			"edge",
			project + "-" + stage,
			"--service=istio-ingressgateway",
			"--hostname=www." + project + "-" + stage + "." + appDomain,
			"--port=http2",
			"--wildcard-policy=Subdomain",
			"--insecure-policy=Allow",
			"-n",
			"istio-system",
		})
	if err != nil {
		logger.Error("Could not create route for: " + err.Error())
	}
	logger.Info("oc create route output: " + out)
}
