package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/cloudevents/sdk-go/pkg/cloudevents/types"
	"github.com/google/uuid"

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
const eventbroker = "EVENTBROKER"

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

	var eventData *keptnevents.ProblemEventData
	if event.Type() == keptnevents.ProblemOpenEventType {
		logger.Debug("Received open problem")
		eventData = &keptnevents.ProblemEventData{}
		if err := event.DataAs(eventData); err != nil {
			return err
		}
	}

	var releasename *string
	var err error
	// get helm release name by impacted entity
	switch impact := eventData.ProblemDetails; {
	case strings.HasPrefix(impact, "Pod"):
		releasename, err = getReleaseByPodName(eventData.ImpactedEntity)
		if err != nil {
			fmt.Println("could not get release for pod: ", err)
			return err
		}
	default:
		logger.Error("could not interpret problem details")
		return errors.New("could not interpret problem details")
	}

	projectname, stagename, servicename := splitReleaseName(*releasename)
	resourceURI := remediationfilename

	// get remediation.yaml
	resourceHandler := keptnutils.NewResourceHandler(configurationserviceconnection)
	resource, err := resourceHandler.GetServiceResource(projectname, stagename, servicename, resourceURI)
	if err != nil {
		logger.Error("could not get remediation.yaml file")
		return err
	}
	logger.Debug("remediation.yaml for service found")

	// get remediation action from remediation.yaml
	var remediationData keptnmodels.Remediations
	yaml.Unmarshal([]byte(resource.ResourceContent), &remediationData)

	for _, remediation := range remediationData.Remediations {
		if remediation.Name == eventData.ProblemTitle {
			logger.Debug("Remediation for problem found in remediation.yaml")
			// currently only one remediation action is supported
			if remediation.Actions[0].Action == "scaling" {
				currentReplicaCount, err := getReplicaCount(projectname, stagename, servicename, configurationserviceconnection)
				if err != nil {
					logger.Error("could not get replicaCount")
					return err
				}
				adjustment, err := strconv.Atoi(remediation.Actions[0].Value)
				if err != nil {
					logger.Error("could not get value for scaling remediation action")
				}
				newReplicacount := currentReplicaCount + adjustment
				err = createAndSendEvent(shkeptncontext, projectname, servicename, stagename, newReplicacount)
				if err != nil {
					logger.Error(err.Error())
				}
			}

		}
	}

	return nil
}

// initialize helm
func init() {
	_, err := keptnutils.ExecuteCommand("helm", []string{"init", "--client-only"})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("initialization of Helm done")
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

func getReleaseByPodName(podname string) (*string, error) {
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

func splitReleaseName(releasename string) (project string, stage string, service string) {
	// currently no "-" in project and service name are allowed, thus "-" is used to split
	s := strings.SplitN(releasename, "-", 3)
	project = s[0]
	stage = s[1]
	// remove the "-generated" suffix
	service = strings.Replace(s[2], "-generated", "", 1)
	return project, stage, service
}

func getReplicaCount(project, stage, service, configServiceURL string) (int, error) {
	helmChartName := service + "-generated"
	// Read chart
	chart, err := keptnutils.GetChart(project, service, stage, helmChartName, configServiceURL)
	if err != nil {
		return -1, err
	}

	values := make(map[string]interface{})
	yaml.Unmarshal([]byte(chart.Values.Raw), &values)

	val, ok := values["replicas"].(int)
	if !ok {
		return -1, errors.New("could not get replicas")
	}
	return val, nil
}

func createAndSendEvent(shkeptncontext string, project string,
	service string, stage string, replicas int) error {

	source, _ := url.Parse("remediation-service")
	contentType := "application/json"

	valuesPrimary := make(map[string]interface{})
	valuesPrimary["replicas"] = replicas
	configChangedEvent := keptnevents.ConfigurationChangeEventData{
		Project:       project,
		Service:       service,
		Stage:         stage,
		ValuesPrimary: valuesPrimary,
	}

	event := cloudevents.Event{
		Context: cloudevents.EventContextV02{
			ID:          uuid.New().String(),
			Type:        keptnevents.ConfigurationChangeEventType,
			Source:      types.URLRef{URL: *source},
			ContentType: &contentType,
			Extensions:  map[string]interface{}{"shkeptncontext": shkeptncontext},
		}.AsV02(),
		Data: configChangedEvent,
	}

	return sendEvent(event)
}

func sendEvent(event cloudevents.Event) error {
	endPoint, err := getServiceEndpoint(eventbroker)
	if err != nil {
		return errors.New("Failed to retrieve endpoint of eventbroker. %s" + err.Error())
	}

	if endPoint.Host == "" {
		return errors.New("Host of eventbroker not set")
	}

	transport, err := cloudeventshttp.New(
		cloudeventshttp.WithTarget(endPoint.String()),
		cloudeventshttp.WithEncoding(cloudeventshttp.StructuredV02),
	)
	if err != nil {
		return errors.New("Failed to create transport:" + err.Error())
	}

	c, err := client.New(transport)
	if err != nil {
		return errors.New("Failed to create HTTP client:" + err.Error())
	}

	if _, err := c.Send(context.Background(), event); err != nil {
		return errors.New("Failed to send cloudevent:, " + err.Error())
	}
	return nil
}

// getServiceEndpoint gets an endpoint stored in an environment variable and sets http as default scheme
func getServiceEndpoint(service string) (url.URL, error) {
	url, err := url.Parse(os.Getenv(service))
	if err != nil {
		return *url, fmt.Errorf("Failed to retrieve value from ENVIRONMENT_VARIABLE: %s", service)
	}

	if url.Scheme == "" {
		url.Scheme = "http"
	}

	return *url, nil
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
