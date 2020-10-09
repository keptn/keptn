package main

import (
	"context"
	"fmt"
	"log"
	"os"

	keptnutils "github.com/keptn/kubernetes-utils/pkg"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/kelseyhightower/envconfig"
	keptnapi "github.com/keptn/go-utils/pkg/api/utils"
	keptncommon "github.com/keptn/go-utils/pkg/lib/keptn"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/helm-service/controller"
	"github.com/keptn/keptn/helm-service/pkg/mesh"
	"github.com/keptn/keptn/helm-service/pkg/serviceutils"
	authorizationv1 "k8s.io/api/authorization/v1"
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
	go keptnapi.RunHealthEndpoint("10999")
	os.Exit(_main(os.Args[1:], env))
}

func gotEvent(ctx context.Context, event cloudevents.Event) error {
	serviceName := serviceName
	var shkeptncontext string
	event.Context.ExtensionAs("shkeptncontext", &shkeptncontext)

	keptnHandler, err := keptnv2.NewKeptn(&event, keptncommon.KeptnOpts{
		LoggingOptions: &keptncommon.LoggingOpts{
			EnableWebsocket: true,
			ServiceName:     &serviceName,
		},
		EventBrokerURL:          os.Getenv("EVENTBROKER"),
		ConfigurationServiceURL: os.Getenv("CONFIGURATION_SERVICE"),
	})
	if err != nil {
		fmt.Println("Could not initialize keptn handler: " + err.Error())
		return err
	}

	mesh := mesh.NewIstioMesh()

	url, err := serviceutils.GetConfigServiceURL()
	if err != nil {
		keptnHandler.Logger.Error(fmt.Sprintf("Error when getting config service url: %s", err.Error()))
		closeLogger(keptnHandler)
		return err
	}

	keptnHandler.Logger.Debug("Got event of type " + event.Type())

	if event.Type() == keptnv2.GetTriggeredEventType(keptnv2.DeploymentTaskName) {
		deploymentHandler := controller.NewDeploymentHandler(keptnHandler, mesh, url.String())
		go deploymentHandler.HandleEvent(event, closeLogger)
	} else if event.Type() == keptnv2.GetTriggeredEventType(keptnv2.ReleaseTaskName) {
		releaseHandler := controller.NewReleaseHandler(keptnHandler, mesh, url.String())
		go releaseHandler.HandleEvent(event, closeLogger)
	} else if event.Type() == keptnv2.GetFinishedEventType(keptnv2.ServiceCreateTaskName){
		onboarder := controller.NewOnboarder(keptnHandler, mesh, url.String())
		go onboarder.HandleEvent(event, closeLogger)
	} else if event.Type() == keptnv2.GetTriggeredEventType(keptnv2.ActionTaskName) {
		actionHandler := controller.NewActionTriggeredHandler(keptnHandler, url.String())
		go actionHandler.HandleEvent(event, closeLogger)
	} else if event.Type() == keptnv2.GetFinishedEventType(keptnv2.ServiceDeleteTaskName) {
		deleteHandler := controller.NewDeleteHandler(keptnHandler, url.String())
		go deleteHandler.HandleEvent(event, closeLogger)
	} else {
		keptnHandler.Logger.Error("Received unexpected keptn event")
		closeLogger(keptnHandler)
	}

	return nil
}

func closeLogger(keptnHandler *keptnv2.Keptn) {
	if combinedLogger, ok := keptnHandler.Logger.(*keptncommon.CombinedLogger); ok {
		combinedLogger.Terminate("")
	}
}

func _main(args []string, env envConfig) int {

	// Check admin rights
	adminRights, err := hasAdminRights()
	if err != nil {
		log.Fatal(fmt.Sprintf("failed to check wheter helm-service has admin right: %v", err))
	}
	if !adminRights {
		log.Fatal("helm-service has insufficient RBAC rights.")
	}

	ctx := context.Background()
	ctx = cloudevents.WithEncodingStructured(ctx)

	p, err := cloudevents.NewHTTP(cloudevents.WithPath(env.Path), cloudevents.WithPort(env.Port))
	if err != nil {
		log.Fatalf("failed to create client, %v", err)
	}
	c, err := cloudevents.NewClient(p)
	if err != nil {
		log.Fatalf("failed to create client, %v", err)
	}
	log.Fatal(c.StartReceiver(ctx, gotEvent))

	return 0
}

func hasAdminRights() (bool, error) {
	clientset, err := keptnutils.GetClientset(true)
	if err != nil {
		return false, err
	}
	sar := &authorizationv1.SelfSubjectAccessReview{
		Spec: authorizationv1.SelfSubjectAccessReviewSpec{
			ResourceAttributes: &authorizationv1.ResourceAttributes{},
		},
	}
	resp, err := clientset.AuthorizationV1().SelfSubjectAccessReviews().Create(sar)
	if err != nil {
		return false, err
	}
	return resp.Status.Allowed, nil
}
