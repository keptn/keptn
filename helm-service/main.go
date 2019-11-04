package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/client"
	cloudeventshttp "github.com/cloudevents/sdk-go/pkg/cloudevents/transport/http"
	"github.com/gorilla/websocket"
	"github.com/kelseyhightower/envconfig"
	keptnevents "github.com/keptn/go-utils/pkg/events"
	keptnutils "github.com/keptn/go-utils/pkg/utils"
	"github.com/keptn/keptn/helm-service/controller"
	"github.com/keptn/keptn/helm-service/controller/helm"
	"github.com/keptn/keptn/helm-service/controller/mesh"
	"github.com/keptn/keptn/helm-service/pkg/serviceutils"
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

func getKeptnDomain() (string, error) {
	useInClusterConfig := false
	if os.Getenv("ENVIRONMENT") == "production" {
		useInClusterConfig = true
	}
	return keptnutils.GetKeptnDomain(useInClusterConfig)
}

func gotEvent(ctx context.Context, event cloudevents.Event) error {
	var shkeptncontext string
	event.Context.ExtensionAs("shkeptncontext", &shkeptncontext)

	stdLogger := keptnutils.NewLogger(shkeptncontext, event.Context.GetID(), "helm-service")

	var logger keptnutils.LoggerInterface

	loggingDone := make(chan bool)

	connData := keptnutils.ConnectionData{}
	if err := event.DataAs(&connData); err != nil ||
		connData.ChannelInfo.ChannelID == "" || connData.ChannelInfo.Token == "" {
		logger = stdLogger
		logger.Debug("No Websocket connection data available")
	} else {
		apiServiceURL, err := serviceutils.GetAPIURL()
		if err != nil {
			logger.Error(err.Error())
			return nil
		}
		ws, _, err := keptnutils.OpenWS(connData, *apiServiceURL)
		if err != nil {
			stdLogger.Error(fmt.Sprintf("Opening websocket connection failed. %s", err.Error()))
			return nil
		}
		combinedLogger := keptnutils.NewCombinedLogger(stdLogger, ws, shkeptncontext)
		logger = combinedLogger
		go closeLogger(loggingDone, combinedLogger, ws)
	}

	mesh := mesh.NewIstioMesh()
	var canaryLevelGen helm.CanaryLevelGenerator
	if os.Getenv("CANARY") == "deployment" {
		canaryLevelGen = helm.NewCanaryOnDeploymentGenerator()
	} else {
		canaryLevelGen = helm.NewCanaryOnNamespaceGenerator()
	}

	keptnDomain, err := getKeptnDomain()
	if err != nil {
		logger.Error("Error when reading the keptn domain")
		return nil
	}

	if event.Type() == keptnevents.ConfigurationChangeEventType {
		configChanger := controller.NewConfigurationChanger(mesh, canaryLevelGen, logger, keptnDomain)
		go configChanger.ChangeAndApplyConfiguration(event, loggingDone)
	} else if event.Type() == keptnevents.InternalServiceCreateEventType {
		onboarder := controller.NewOnboarder(mesh, canaryLevelGen, logger, keptnDomain)
		go onboarder.DoOnboard(event, loggingDone)
	} else {
		logger.Error("Received unexpected keptn event")
		loggingDone <- true
	}

	return nil
}

func closeLogger(loggingDone chan bool, combinedLogger *keptnutils.CombinedLogger, ws *websocket.Conn) {
	<-loggingDone
	combinedLogger.Terminate()
	ws.Close()
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
