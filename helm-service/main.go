package main

import (
	"context"
	"fmt"
	"github.com/keptn/keptn/helm-service/controller"
	"github.com/keptn/keptn/helm-service/pkg/configurationchanger"
	"github.com/keptn/keptn/helm-service/pkg/helm"
	"net/url"

	"github.com/keptn/keptn/helm-service/pkg/namespacemanager"
	"log"
	"os"

	keptnutils "github.com/keptn/kubernetes-utils/pkg"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/kelseyhightower/envconfig"
	configutils "github.com/keptn/go-utils/pkg/api/utils"
	keptnapi "github.com/keptn/go-utils/pkg/api/utils"
	keptncommon "github.com/keptn/go-utils/pkg/lib/keptn"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"

	utils "github.com/keptn/go-utils/pkg/api/utils"
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

	if event.Context.GetSource() == serviceName {
		return nil
	}
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

	configServiceURL, err := serviceutils.GetConfigServiceURL()
	if err != nil {
		keptnHandler.Logger.Error(fmt.Sprintf("Error when getting configServiceURL: %s", err.Error()))
		return err
	}

	shipyardControllerURL, err := serviceutils.GetShipyardControllerURL()
	if err != nil {
		keptnHandler.Logger.Error(fmt.Sprintf("Error when getting shipyardControllerURL: %s", err.Error()))
		return err
	}

	//create dependencies

	mesh := mesh.NewIstioMesh()
	keptnHandler.Logger.Debug("Got event of type " + event.Type())

	if event.Type() == keptnv2.GetTriggeredEventType(keptnv2.DeploymentTaskName) {
		deploymentHandler := createDeploymentHandler(configServiceURL, keptnHandler, mesh)
		deploymentHandler.HandleEvent(event)
	} else if event.Type() == keptnv2.GetTriggeredEventType(keptnv2.ReleaseTaskName) {
		releaseHandler := createReleaseHandler(configServiceURL, mesh, keptnHandler)
		releaseHandler.HandleEvent(event)
	} else if event.Type() == keptnv2.GetFinishedEventType(keptnv2.ServiceCreateTaskName) {
		onboardHandler := createOnboardHandler(configServiceURL, shipyardControllerURL, keptnHandler, mesh)
		onboardHandler.HandleEvent(event)
	} else if event.Type() == keptnv2.GetTriggeredEventType(keptnv2.ActionTaskName) {
		actionHandler := createActionTriggeredHandler(configServiceURL, keptnHandler)
		actionHandler.HandleEvent(event)
	} else if event.Type() == keptnv2.GetFinishedEventType(keptnv2.ServiceDeleteTaskName) {
		deleteHandler := createDeleteHandler(configServiceURL, shipyardControllerURL, keptnHandler)
		deleteHandler.HandleEvent(event)
	} else {
		keptnHandler.Logger.Error("Received unexpected keptn event")
	}

	return nil
}

func createKeptnBaseHandler(url *url.URL, keptn *keptnv2.Keptn) controller.Handler {
	namespaceManager := namespacemanager.NewNamespaceManager(keptn.Logger)
	helmExecutor := helm.NewHelmV3Executor(keptn.Logger, namespaceManager)
	keptnHandlerBase := controller.NewHandlerBase(keptn, helmExecutor, url.String())
	return keptnHandlerBase

}

func createDeleteHandler(configServiceURL *url.URL, shipyardControllerURL *url.URL, keptn *keptnv2.Keptn) *controller.DeleteHandler {
	stagesHandler := configutils.NewStageHandler(shipyardControllerURL.String())
	keptnBaseHandler := createKeptnBaseHandler(configServiceURL, keptn)
	deleteHandler := controller.NewDeleteHandler(keptnBaseHandler, stagesHandler, configServiceURL.String())
	return deleteHandler
}

func createActionTriggeredHandler(configServiceURL *url.URL, keptn *keptnv2.Keptn) *controller.ActionTriggeredHandler {
	configChanger := configurationchanger.NewConfigurationChanger(configServiceURL.String())
	keptnBaseHandler := createKeptnBaseHandler(configServiceURL, keptn)
	actionHandler := controller.NewActionTriggeredHandler(keptnBaseHandler, configChanger, configServiceURL.String())
	return actionHandler
}

func createReleaseHandler(url *url.URL, mesh *mesh.IstioMesh, keptn *keptnv2.Keptn) *controller.ReleaseHandler {
	configChanger := configurationchanger.NewConfigurationChanger(url.String())
	chartGenerator := helm.NewGeneratedChartGenerator(mesh, keptn.Logger)
	chartStorer := keptnutils.NewChartStorer(utils.NewResourceHandler(url.String()))
	chartPackager := keptnutils.NewChartPackager()
	keptnBaseHandler := createKeptnBaseHandler(url, keptn)
	releaseHandler := controller.NewReleaseHandler(keptnBaseHandler, mesh, configChanger, chartGenerator, chartStorer, chartPackager, url.String())
	return releaseHandler
}

func createOnboarder(configServiceURL *url.URL, keptn *keptnv2.Keptn, mesh *mesh.IstioMesh) controller.Onboarder {
	namespaceManager := namespacemanager.NewNamespaceManager(keptn.Logger)
	chartStorer := keptnutils.NewChartStorer(utils.NewResourceHandler(configServiceURL.String()))
	chartGenerator := helm.NewGeneratedChartGenerator(mesh, keptn.Logger)
	chartPackager := keptnutils.NewChartPackager()
	keptnBaseHandler := createKeptnBaseHandler(configServiceURL, keptn)
	onBoarder := controller.NewOnboarder(keptnBaseHandler, namespaceManager, chartStorer, chartGenerator, chartPackager)
	return onBoarder
}

func createOnboardHandler(configServiceURL *url.URL, shipyardControllerURL *url.URL, keptn *keptnv2.Keptn, mesh *mesh.IstioMesh) *controller.OnboardHandler {
	projectHandler := keptnapi.NewProjectHandler(shipyardControllerURL.String())
	stagesHandler := configutils.NewStageHandler(shipyardControllerURL.String())
	onBoarder := createOnboarder(configServiceURL, keptn, mesh)
	keptnBaseHandler := createKeptnBaseHandler(configServiceURL, keptn)
	return controller.NewOnboardHandler(keptnBaseHandler, projectHandler, stagesHandler, onBoarder)
}

func createDeploymentHandler(url *url.URL, keptn *keptnv2.Keptn, mesh *mesh.IstioMesh) *controller.DeploymentHandler {
	chartGenerator := helm.NewGeneratedChartGenerator(mesh, keptn.Logger)
	onBoarder := createOnboarder(url, keptn, mesh)
	keptnBaseHandler := createKeptnBaseHandler(url, keptn)
	deploymentHandler := controller.NewDeploymentHandler(keptnBaseHandler, mesh, onBoarder, chartGenerator)
	return deploymentHandler
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
