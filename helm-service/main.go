package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/client"
	cloudeventshttp "github.com/cloudevents/sdk-go/pkg/cloudevents/transport/http"
	"github.com/kelseyhightower/envconfig"
	keptnevents "github.com/keptn/go-utils/pkg/lib"
	"github.com/keptn/keptn/helm-service/controller"
	"github.com/keptn/keptn/helm-service/controller/mesh"
	"github.com/keptn/keptn/helm-service/pkg/serviceutils"
)

type envConfig struct {
	// Port on which to listen for cloudevents
	Port int    `envconfig:"RCV_PORT" default:"8080"`
	Path string `envconfig:"RCV_PATH" default:"/"`
}

const serviceName = "helm-service"

func main() {
	var env envConfig
	if err := envconfig.Process("", &env); err != nil {
		log.Fatalf("Failed to process env var: %s", err)
	}
	os.Exit(_main(os.Args[1:], env))
}

func getIngressHostnameSuffix() string {
	if os.Getenv("INGRESS_HOSTNAME_SUFFIX") != "" {
		return os.Getenv("INGRESS_HOSTNAME_SUFFIX")
	}
	return "svc.cluster.local"
}

func getIngressProtocol() string {
	if os.Getenv("INGRESS_PROTOCOL") != "" {
		return strings.ToLower(os.Getenv("INGRESS_PROTOCOL"))
	}
	return "http"
}

func getIngressPort() string {
	if os.Getenv("INGRESS_PORT") != "" {
		return os.Getenv("INGRESS_PORT")
	}
	return "80"
}

func gotEvent(ctx context.Context, event cloudevents.Event) error {
	serviceName := serviceName
	var shkeptncontext string
	event.Context.ExtensionAs("shkeptncontext", &shkeptncontext)

	keptnHandler, err := keptnevents.NewKeptn(&event, keptnevents.KeptnOpts{
		LoggingOptions: &keptnevents.LoggingOpts{
			EnableWebsocket: true,
			ServiceName:     &serviceName,
		},
	})
	if err != nil {
		fmt.Println("Could not initialize keptn handler: " + err.Error())
		return err
	}

	var logger keptnevents.LoggerInterface
	loggingDone := make(chan bool)
	go closeLogger(loggingDone, keptnHandler.Logger)

	mesh := mesh.NewIstioMesh()

	ingressHostnameSuffix := getIngressHostnameSuffix()
	ingressProtocol := getIngressProtocol()
	ingressPort := getIngressPort()

	url, err := serviceutils.GetConfigServiceURL()
	if err != nil {
		keptnHandler.Logger.Error(fmt.Sprintf("Error when getting config service url: %s", err.Error()))
		loggingDone <- true
		return err
	}

	keptnHandler.Logger.Debug("Got event of type " + event.Type())

	if event.Type() == keptnevents.ConfigurationChangeEventType {
		configChanger := controller.NewConfigurationChanger(mesh, keptnHandler, ingressHostnameSuffix, url.String(), ingressProtocol, ingressPort)
		go configChanger.ChangeAndApplyConfiguration(event, loggingDone)
	} else if event.Type() == keptnevents.InternalServiceCreateEventType {
		onboarder := controller.NewOnboarder(mesh, keptnHandler, ingressHostnameSuffix, url.String(), ingressProtocol, ingressPort)
		go onboarder.DoOnboard(event, loggingDone)
	} else if event.Type() == keptnevents.ActionTriggeredEventType {
		actionHandler := controller.NewActionTriggeredHandler(keptnHandler, url.String())
		go actionHandler.HandleEvent(event, loggingDone)
	} else {
		logger.Error("Received unexpected keptn event")
		loggingDone <- true
	}

	return nil
}

func closeLogger(loggingDone chan bool, logger keptnevents.LoggerInterface) {
	<-loggingDone
	if combinedLogger, ok := logger.(*keptnevents.CombinedLogger); ok {
		combinedLogger.Terminate()
	}
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
