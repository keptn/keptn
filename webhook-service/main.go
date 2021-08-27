package main

import (
	api "github.com/keptn/go-utils/pkg/api/utils"
	"github.com/keptn/keptn/go-sdk/pkg/sdk"
	"github.com/keptn/keptn/webhook-service/handler"
	"github.com/keptn/keptn/webhook-service/lib"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"log"
)

const eventTypeWildcard = "*"
const serviceName = "webhook-service"

func main() {
	kubeAPI, err := createKubeAPI()
	if err != nil {
		log.Fatalf("could not create kubernetes client: %s", err.Error())
	}
	secretReader := lib.NewK8sSecretReader(kubeAPI)
	taskHandler := handler.NewTaskHandler(&lib.TemplateEngine{}, &lib.CmdCurlExecutor{}, secretReader)

	go api.RunHealthEndpoint("10998")
	log.Fatal(sdk.NewKeptn(
		serviceName,
		sdk.WithHandler(
			eventTypeWildcard,
			taskHandler,
		),
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
