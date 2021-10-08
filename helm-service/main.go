package main

import (
	"context"
	"fmt"
	"github.com/keptn/keptn/helm-service/controller"
	"github.com/keptn/keptn/helm-service/pkg/configurationchanger"
	"github.com/keptn/keptn/helm-service/pkg/helm"
	logger "github.com/sirupsen/logrus"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/url"
	"os/signal"
	"sync"
	"syscall"

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
	Port     int    `envconfig:"RCV_PORT" default:"8080"`
	Path     string `envconfig:"RCV_PATH" default:"/"`
	LogLevel string `envconfig:"LOG_LEVEL" default:"info"`
}

const serviceName = "helm-service"

func main() {
	var env envConfig
	if err := envconfig.Process("", &env); err != nil {
		logger.Fatalf("Failed to process env var: %s", err)
	}

	logger.SetLevel(logger.InfoLevel)

	if os.Getenv(env.LogLevel) != "" {
		logLevel, err := logger.ParseLevel(os.Getenv(env.LogLevel))
		if err != nil {
			logger.WithError(err).Error("could not parse log level provided by 'LOG_LEVEL' env var")
		} else {
			logger.SetLevel(logLevel)
		}
	}

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
		ConfigurationServiceURL: os.Getenv("CONFIGURATION_SERVICE"),
	})
	if err != nil {
		fmt.Println("Could not initialize keptn handler: " + err.Error())
		return err
	}

	configServiceURL, err := serviceutils.GetConfigServiceURL()
	if err != nil {
		logger.WithError(err).Error("Error when getting configServiceURL")
		return err
	}

	shipyardControllerURL, err := serviceutils.GetShipyardControllerURL()
	if err != nil {
		logger.WithError(err).Error("Error when getting shipyardControllerURL")
		return err
	}

	//create dependencies

	mesh := mesh.NewIstioMesh()
	logger.Debug("Got event of type " + event.Type())

	// ToDo: Multithreaded is important here, such that the endpoint responds immediately
	// else we will have deployment handler take 30 seconds, and after that the response will be sent
	ctx.Value("Wg").(*sync.WaitGroup).Add(1)
	if event.Type() == keptnv2.GetTriggeredEventType(keptnv2.DeploymentTaskName) {
		deploymentHandler := createDeploymentHandler(configServiceURL, keptnHandler, mesh)
		go deploymentHandler.HandleEvent(ctx, event)
	} else if event.Type() == keptnv2.GetTriggeredEventType(keptnv2.ReleaseTaskName) {
		releaseHandler := createReleaseHandler(configServiceURL, mesh, keptnHandler)
		go releaseHandler.HandleEvent(ctx, event)
	} else if event.Type() == keptnv2.GetTriggeredEventType(keptnv2.RollbackTaskName) {
		rollbackHandler := createRollbackHandler(configServiceURL, mesh, keptnHandler)
		go rollbackHandler.HandleEvent(ctx, event)
	} else if event.Type() == keptnv2.GetTriggeredEventType(keptnv2.ActionTaskName) {
		actionHandler := createActionTriggeredHandler(configServiceURL, keptnHandler)
		go actionHandler.HandleEvent(ctx, event)
	} else if event.Type() == keptnv2.GetFinishedEventType(keptnv2.ServiceDeleteTaskName) {
		deleteHandler := createDeleteHandler(configServiceURL, shipyardControllerURL, keptnHandler)
		go deleteHandler.HandleEvent(ctx, event)
	} else {
		logger.Error("Received unexpected keptn event")
		ctx.Value("Wg").(*sync.WaitGroup).Done()
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
	actionHandler := controller.NewActionTriggeredHandler(keptnBaseHandler, configChanger)
	return actionHandler
}

func createReleaseHandler(url *url.URL, mesh *mesh.IstioMesh, keptn *keptnv2.Keptn) *controller.ReleaseHandler {
	configChanger := configurationchanger.NewConfigurationChanger(url.String())
	chartGenerator := helm.NewGeneratedChartGenerator(mesh)
	chartStorer := keptnutils.NewChartStorer(utils.NewResourceHandler(url.String()))
	chartPackager := keptnutils.NewChartPackager()
	keptnBaseHandler := createKeptnBaseHandler(url, keptn)
	releaseHandler := controller.NewReleaseHandler(keptnBaseHandler, mesh, configChanger, chartGenerator, chartStorer, chartPackager, url.String())
	return releaseHandler
}

func createRollbackHandler(url *url.URL, mesh *mesh.IstioMesh, keptn *keptnv2.Keptn) *controller.RollbackHandler {
	configChanger := configurationchanger.NewConfigurationChanger(url.String())
	keptnBaseHandler := createKeptnBaseHandler(url, keptn)
	rollbackHandler := controller.NewRollbackHandler(keptnBaseHandler, mesh, configChanger)
	return rollbackHandler
}

func createOnboarder(configServiceURL *url.URL, keptn *keptnv2.Keptn, mesh *mesh.IstioMesh) controller.Onboarder {
	namespaceManager := namespacemanager.NewNamespaceManager(keptn.Logger)
	chartStorer := keptnutils.NewChartStorer(utils.NewResourceHandler(configServiceURL.String()))
	chartGenerator := helm.NewGeneratedChartGenerator(mesh)
	chartPackager := keptnutils.NewChartPackager()
	keptnBaseHandler := createKeptnBaseHandler(configServiceURL, keptn)
	onBoarder := controller.NewOnboarder(keptnBaseHandler, namespaceManager, chartStorer, chartGenerator, chartPackager)
	return onBoarder
}

func createDeploymentHandler(url *url.URL, keptn *keptnv2.Keptn, mesh *mesh.IstioMesh) *controller.DeploymentHandler {
	chartGenerator := helm.NewGeneratedChartGenerator(mesh)
	onBoarder := createOnboarder(url, keptn, mesh)
	keptnBaseHandler := createKeptnBaseHandler(url, keptn)
	deploymentHandler := controller.NewDeploymentHandler(keptnBaseHandler, mesh, onBoarder, chartGenerator)
	return deploymentHandler
}

func _main(args []string, env envConfig) int {

	// Check admin rights
	adminRights, err := hasAdminRights()
	if err != nil {
		log.Fatal(fmt.Sprintf("failed to check whether helm-service has admin right: %v", err))
	}
	if !adminRights {
		log.Println("Warning: helm-service is running without admin RBAC rights. See #3511 for details.")
	}

	ctx := getGracefulContext()

	p, err := cloudevents.NewHTTP(cloudevents.WithPath(env.Path), cloudevents.WithPort(env.Port), cloudevents.WithGetHandlerFunc(keptnapi.HealthEndpointHandler))
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

// hasAdminRights checks if the current pod is assigned the Admin Role
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
	resp, err := clientset.AuthorizationV1().SelfSubjectAccessReviews().Create(context.TODO(), sar, v1.CreateOptions{})
	if err != nil {
		return false, err
	}
	return resp.Status.Allowed, nil
}

func getGracefulContext() context.Context {

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	wg := &sync.WaitGroup{}
	ctx, cancel := context.WithCancel(cloudevents.WithEncodingStructured(context.WithValue(context.Background(), "Wg", wg)))

	go func() {
		<-ch
		log.Fatal("Container termination triggered, starting graceful shutdown")
		cancel()
	}()

	return ctx
}
