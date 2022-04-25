package main

import (
	"os"

	"github.com/keptn/keptn/go-sdk/pkg/sdk"
	"github.com/keptn/keptn/webhook-service/handler"
	"github.com/keptn/keptn/webhook-service/lib"
	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

const eventTypeWildcard = "*"
const serviceName = "webhook-service"
const envVarLogLevel = "LOG_LEVEL"

func main() {
	if os.Getenv(envVarLogLevel) != "" {
		logLevel, err := log.ParseLevel(os.Getenv(envVarLogLevel))
		if err != nil {
			log.WithError(err).Error("could not parse log level provided by 'LOG_LEVEL' env var")
		} else {
			log.SetLevel(logLevel)
		}
	}
	kubeAPI, err := createKubeAPI()
	if err != nil {
		log.Fatalf("could not create kubernetes client: %s", err.Error())
	}
	secretReader := lib.NewK8sSecretReader(kubeAPI)

	curlExecutor := lib.NewCmdCurlExecutor(
		&lib.OSCmdExecutor{},
		lib.WithDeniedURLs(lib.GetDeniedURLs(lib.GetEnv())),
	)
	taskHandler := handler.NewTaskHandler(&lib.TemplateEngine{}, curlExecutor, secretReader)

	log.Fatal(sdk.NewKeptn(
		serviceName,
		sdk.WithTaskHandler(
			eventTypeWildcard,
			taskHandler,
		),
		sdk.WithAutomaticResponse(false),
		sdk.WithLogger(log.New()),
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
