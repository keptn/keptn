package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/client"
	cloudeventshttp "github.com/cloudevents/sdk-go/pkg/cloudevents/transport/http"
	"github.com/kelseyhightower/envconfig"

	keptnevents "github.com/keptn/go-utils/pkg/events"
	keptnmodels "github.com/keptn/go-utils/pkg/models"
	keptnutils "github.com/keptn/go-utils/pkg/utils"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/helm/pkg/helm"
	helm_env "k8s.io/helm/pkg/helm/environment"
	"k8s.io/helm/pkg/helm/portforwarder"
	"k8s.io/helm/pkg/kube"
)

// const timeout = 60
// const configservice = "CONFIGURATION_SERVICE"
// const eventbroker = "EVENTBROKER"
// const api = "API"
const tillernamespace = "kube-system"
const remediationfilename = "remediation.yaml"
const configurationserviceconnection = "localhost:6060" // "configuration-service:8080"

type envConfig struct {
	Port int    `envconfig:"RCV_PORT" default:"8080"`
	Path string `envconfig:"RCV_PATH" default:"/"`
}

func main() {
	var env envConfig
	if err := envconfig.Process("", &env); err != nil {
		log.Fatalf("failed to process env var: %s", err)
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
		log.Fatalf("failed to create transport: %v", err)
	}

	c, err := client.New(t)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	log.Fatalf("failed to start receiver: %s", c.StartReceiver(ctx, gotEvent))

	return 0
}

func gotEvent(ctx context.Context, event cloudevents.Event) error {
	var shkeptncontext string
	event.Context.ExtensionAs("shkeptncontext", &shkeptncontext)

	logger := keptnutils.NewLogger(shkeptncontext, event.Context.GetID(), "remediation-service")
	logger.Debug("Received event for shkeptncontext:" + shkeptncontext)

	//fmt.Println(event.Data)

	var eventData *keptnevents.ProblemEventData
	if event.Type() == keptnevents.ProblemOpenEventType {
		logger.Debug("Received open problem")
		eventData = &keptnevents.ProblemEventData{}
		if err := event.DataAs(eventData); err != nil {
			return err
		}
		fmt.Println(eventData)
	}

	//impactedEntity := eventData.ImpactedEntity //"carts-blue-856559f565-2knz7"

	// assume impactedEntity is pod
	var releasename *string
	var err error
	switch impact := eventData.ProblemDetails; {
	case strings.HasPrefix(impact, "Pod"):
		fmt.Println("pod")
		releasename, err = getReleaseOfPodName(eventData.ImpactedEntity)
		if err != nil {
			fmt.Println("could not get releaes for pod: ", err)
			return err
		}
	default:
		return errors.New("no remediation supported")
	}

	fmt.Println("release name: ", *releasename)

	// rockshop-production-carts
	projectname, stagename, servicename := splitReleaseNames("sockshop-production-carts")
	resourceURI := remediationfilename

	// get remediation.yaml
	resourceHandler := keptnutils.NewResourceHandler(configurationserviceconnection)

	fmt.Println(projectname, "/", stagename, "/", servicename, "/", resourceURI)

	resource, err := resourceHandler.GetServiceResource(projectname, stagename, servicename, resourceURI)
	if err != nil {
		fmt.Println("could not get resource: ", err)
		return err
	}
	fmt.Println("content\n", resource.ResourceContent)

	// get remediation action
	var remediationData keptnmodels.Remediations
	yaml.Unmarshal([]byte(resource.ResourceContent), &remediationData)

	fmt.Println(remediationData.Remediations[0].Actions[0].Action)

	return nil
}

func init() {
	fmt.Println("init")
	_, err := keptnutils.ExecuteCommand("helm", []string{"init", "--client-only"})
	if err != nil {
		log.Fatal(err)
	}
}

// getKubeClient creates a Kubernetes config and client for a given kubeconfig context.
func getKubeClient(context string, kubeconfig string) (*rest.Config, kubernetes.Interface, error) {
	var config *rest.Config
	var err error

	config, err = getK8sConfig()
	if err != nil {
		return nil, nil, fmt.Errorf("could not get config for Kubernetes client: %s", err)
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, nil, fmt.Errorf("could not get Kubernetes client: %s", err)
	}
	return config, client, nil
}

// configForContext creates a Kubernetes REST client configuration for a given kubeconfig context.
func configForContext(context string, kubeconfig string) (*rest.Config, error) {
	config, err := kube.GetConfig(context, kubeconfig).ClientConfig()
	if err != nil {
		return nil, fmt.Errorf("could not get Kubernetes config for context %q: %s", context, err)
	}
	return config, nil
}

// UserHomeDir retries home directory of user
func UserHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")
}

func getK8sConfig() (*rest.Config, error) {
	var useInClusterConfig bool
	if os.Getenv("env") == "production" {
		useInClusterConfig = true
	} else {
		useInClusterConfig = false
	}
	var config *rest.Config
	var err error
	if useInClusterConfig {
		config, err = rest.InClusterConfig()
	} else {
		kubeconfig := filepath.Join(
			UserHomeDir(), ".kube", "config",
		)
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return nil, err
		}
	}

	return config, nil
}

func getHelmClient() (*helm.Client, error) {
	var settings helm_env.EnvSettings
	var tillerTunnel *kube.Tunnel

	if settings.TillerHost == "" {
		config, client, err := getKubeClient(settings.KubeContext, settings.KubeConfig)
		if err != nil {
			return nil, err
		}
		settings.TillerNamespace = tillernamespace
		tillerTunnel, err = portforwarder.New(settings.TillerNamespace, client, config)
		if err != nil {
			return nil, err
		}

		settings.TillerHost = fmt.Sprintf("127.0.0.1:%d", tillerTunnel.Local)
		//debug("Created tunnel using local port: '%d'\n", tillerTunnel.Local)
	}
	options := []helm.Option{helm.Host(settings.TillerHost), helm.ConnectTimeout(5)}
	client := helm.NewClient(options...)
	return client, nil

}

func getReleaseOfPodName(podname string) (*string, error) {
	client, err := getHelmClient()
	if err != nil {
		fmt.Println("could not get Helm client")
		return nil, err
	}

	releaseList, err := client.ListReleases()
	if err != nil {
		fmt.Println("could not fetch list of releases")
		return nil, err
	}

	for _, r := range releaseList.GetReleases() {
		rs, err := client.ReleaseStatus(r.Name, helm.StatusReleaseVersion(0))
		if err != nil {
			fmt.Println("release status error: ", err)
			return nil, err
		}
		if strings.Contains(rs.Info.Status.Resources, podname) {
			return &r.Name, nil
		}
	}

	return nil, errors.New("could not find release for pod")
}

func splitReleaseNames(releasename string) (project string, stage string, service string) {
	// currently no "-" in project and service name are allowed, thus "-" is used to split
	s := strings.SplitN(releasename, "-", 3)
	project = s[0]
	stage = s[1]
	service = s[2]
	return project, stage, service
}

// func retrieveResourceForProject(projectName string, resourceURI string, logger keptnutils.Logger) (*keptnutils.Resource, error) {
// 	eventURL, err := getServiceEndpoint(configservice)
// 	resourceHandler := keptnutils.NewResourceHandler(eventURL.Host)

// 	resource, err := resourceHandler.GetProjectResource(projectName, resourceURI)
// 	if err != nil {
// 		return nil, fmt.Errorf("Failed to retrieve resource %s", resourceURI, err.Error())
// 	}

// 	logger.Info(resource.ResourceContent)

// 	return resource, nil
// }
