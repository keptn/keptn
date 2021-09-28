package main

import (
	api "github.com/keptn/go-utils/pkg/api/utils"
	"github.com/keptn/keptn/go-sdk/pkg/sdk"
	"github.com/keptn/keptn/webhook-service/handler"
	"github.com/keptn/keptn/webhook-service/lib"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"log"
	"os"
)

const eventTypeWildcard = "*"
const serviceName = "webhook-service"

func main() {
	kubeAPI, err := createKubeAPI()
	if err != nil {
		log.Fatalf("could not create kubernetes client: %s", err.Error())
	}
	secretReader := lib.NewK8sSecretReader(kubeAPI)

	kubeAPIHostIP := os.Getenv("KUBERNETES_SERVICE_HOST")
	kubeAPIPort := os.Getenv("KUBERNETES_SERVICE_PORT")

	curlExecutor := lib.NewCmdCurlExecutor(
		&lib.OSCmdExecutor{},
		lib.WithUnAllowedURLs(
			[]string{
				kubeAPIHostIP + ":" + kubeAPIPort,
				"kubernetes" + ":" + kubeAPIPort,
				"kubernetes.default" + ":" + kubeAPIPort,
				"kubernetes.default.svc.cluster.local" + ":" + kubeAPIPort,
			},
		),
	)
	taskHandler := handler.NewTaskHandler(&lib.TemplateEngine{}, curlExecutor, secretReader)

	go api.RunHealthEndpoint("10998")
	log.Fatal(sdk.NewKeptn(
		serviceName,
		sdk.WithHandler(
			eventTypeWildcard,
			taskHandler,
		),
		sdk.WithAutomaticResponse(false),
	).Start())
}

func createKubeAPI() (*kubernetes.Clientset, error) {
	var config *rest.Config
	config, err := rest.InClusterConfig()

	if err != nil {
		return nil, err
	}

	kubeAPI, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return kubeAPI, nil
}
