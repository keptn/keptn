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

	url, err := serviceutils.GetConfigServiceURL()
	if err != nil {
		keptnHandler.Logger.Error(fmt.Sprintf("Error when getting config service url: %s", err.Error()))
		return err
	}

	//create dependencies

	mesh := mesh.NewIstioMesh()
	keptnHandler.Logger.Debug("Got event of type " + event.Type())

	if event.Type() == keptnv2.GetTriggeredEventType(keptnv2.DeploymentTaskName) {
		deploymentHandler := createDeploymentHandler(url, keptnHandler, mesh)
		go deploymentHandler.HandleEvent(event)
	} else if event.Type() == keptnv2.GetTriggeredEventType(keptnv2.ReleaseTaskName) {
		releaseHandler := createReleaseHandler(url, mesh, keptnHandler)
		go releaseHandler.HandleEvent(event)
	} else if event.Type() == keptnv2.GetFinishedEventType(keptnv2.ServiceCreateTaskName) {
		onBoarder := createOnboarder(keptnHandler, url, mesh)
		go onBoarder.HandleEvent(event)
	} else if event.Type() == keptnv2.GetTriggeredEventType(keptnv2.ActionTaskName) {
		actionHandler := createActionTriggeredHandler(url, keptnHandler)
		go actionHandler.HandleEvent(event)
	} else if event.Type() == keptnv2.GetFinishedEventType(keptnv2.ServiceDeleteTaskName) {
		deleteHandler := createDeleteHandler(url, keptnHandler)
		go deleteHandler.HandleEvent(event)
	} else {
		keptnHandler.Logger.Error("Received unexpected keptn event")
	}

	return nil
}

func createDeleteHandler(url *url.URL, keptnHandler *keptnv2.Keptn) *controller.DeleteHandler {
	stagesHandler := configutils.NewStageHandler(url.String())
	deleteHandler := controller.NewDeleteHandler(keptnHandler, stagesHandler, url.String())
	return deleteHandler
}

func createActionTriggeredHandler(url *url.URL, keptnHandler *keptnv2.Keptn) *controller.ActionTriggeredHandler {
	configChanger := configurationchanger.NewConfigurationChanger(url.String())
	actionHandler := controller.NewActionTriggeredHandler(keptnHandler, configChanger, url.String())
	return actionHandler
}

func createReleaseHandler(url *url.URL, mesh *mesh.IstioMesh, keptnHandler *keptnv2.Keptn) *controller.ReleaseHandler {
	configChanger := configurationchanger.NewConfigurationChanger(url.String())
	chartGenerator := helm.NewGeneratedChartGenerator(mesh, keptnHandler.Logger)
	chartStorer := keptnutils.NewChartStorer(utils.NewResourceHandler(url.String()))
	chartPackager := keptnutils.NewChartPackager()
	releaseHandler := controller.NewReleaseHandler(keptnHandler, mesh, configChanger, chartGenerator, chartStorer, chartPackager, url.String())
	return releaseHandler
}

func createOnboarder(keptnHandler *keptnv2.Keptn, url *url.URL, mesh *mesh.IstioMesh) controller.Onboarder {
	namespaceManager := namespacemanager.NewNamespaceManager(keptnHandler.Logger)
	projectHandler := keptnapi.NewProjectHandler(url.String())
	stagesHandler := configutils.NewStageHandler(url.String())
	serviceHandler := configutils.NewServiceHandler(url.String())
	chartStorer := keptnutils.NewChartStorer(utils.NewResourceHandler(url.String()))
	chartGenerator := helm.NewGeneratedChartGenerator(mesh, keptnHandler.Logger)
	chartPackager := keptnutils.NewChartPackager()
	onBoarder := controller.NewOnboarder(keptnHandler, mesh, projectHandler, namespaceManager, stagesHandler, serviceHandler, chartStorer, chartGenerator, chartPackager, url.String())
	return onBoarder
}

func createDeploymentHandler(url *url.URL, keptnHandler *keptnv2.Keptn, mesh *mesh.IstioMesh) *controller.DeploymentHandler {
	projectHandler := keptnapi.NewProjectHandler(url.String())
	namespaceManager := namespacemanager.NewNamespaceManager(keptnHandler.Logger)
	stagesHandler := configutils.NewStageHandler(url.String())
	serviceHandler := configutils.NewServiceHandler(url.String())
	chartStorer := keptnutils.NewChartStorer(utils.NewResourceHandler(url.String()))
	chartGenerator := helm.NewGeneratedChartGenerator(mesh, keptnHandler.Logger)
	chartPackager := keptnutils.NewChartPackager()
	onBoarder := controller.NewOnboarder(keptnHandler, mesh, projectHandler, namespaceManager, stagesHandler, serviceHandler, chartStorer, chartGenerator, chartPackager, url.String())
	deploymentHandler := controller.NewDeploymentHandler(keptnHandler, mesh, onBoarder, chartGenerator, url.String())
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
