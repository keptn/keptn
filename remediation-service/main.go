package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	kyaml "k8s.io/apimachinery/pkg/util/yaml"

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

const tillernamespace = "kube-system"
const remediationfilename = "remediation.yaml"
const configurationserviceconnection = "CONFIGURATION_SERVICE" //"localhost:6060" // "configuration-service:8080"
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
		logger.Debug("Received problem notification")
		eventData = &keptnevents.ProblemEventData{}
		if err := event.DataAs(eventData); err != nil {
			return err
		}
	}

	releasename, err := getReleasename(eventData)
	if err != nil {
		logger.Error("could not get release name")
		return err
	}

	projectname, stagename, servicename := splitReleaseName(*releasename)

	if eventData.State != "OPEN" {
		logger.Debug("Received closed problem")
		sendTestsFinishedEvent(shkeptncontext, projectname, stagename, servicename)
		return nil
	}

	logger.Debug("Received open problem")

	resourceURI := remediationfilename

	// valide if remediation should be performed
	resourceHandler := keptnutils.NewResourceHandler(os.Getenv(configurationserviceconnection))
	autoremediate := isRemediationEnabled(resourceHandler, projectname, stagename)
	logger.Debug(fmt.Sprintf("remediation enabled for project and stage: %t", autoremediate))

	if !autoremediate {
		return errors.New("remediation not enabled for project and stage")
	}

	// get remediation.yaml
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
				currentReplicaCount, err := getReplicaCount(logger, projectname, stagename, servicename, os.Getenv(configurationserviceconnection))
				if err != nil {
					logger.Error("could not get replica count")
					return err
				}
				adjustment, err := strconv.Atoi(remediation.Actions[0].Value)
				if err != nil {
					logger.Error("could not get value for scaling remediation action")
				}
				propertyPath := "spec.replicas"
				propertyValue := currentReplicaCount + adjustment
				err = createAndSendEvent(shkeptncontext, projectname, servicename, stagename, propertyPath, propertyValue)
				if err != nil {
					logger.Error(err.Error())
				}
			}

		}
	}

	logger.Debug("remediation service finished")

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

// get helm release name by impacted entity
func getReleasename(eventData *keptnevents.ProblemEventData) (*string, error) {
	switch impact := eventData.ProblemDetails; {
	case strings.HasPrefix(impact, "Pod"):
		return getReleaseByPodName(eventData.ImpactedEntity)
	case strings.HasPrefix(impact, "Service"):
		return nil, errors.New("Service remediation not yet supported")
	case strings.HasPrefix(impact, "Node"):
		return nil, errors.New("Node remediation not yet supported")
	default:
		return nil, errors.New("could not interpret problem details")
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
	if os.Getenv("ENVIRONMENT") == "production" {
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
	}
	options := []helm.Option{helm.Host(settings.TillerHost), helm.ConnectTimeout(5)}
	client := helm.NewClient(options...)
	return client, nil

}

func isRemediationEnabled(rh *keptnutils.ResourceHandler, project string, stage string) bool {
	keptnHandler := keptnutils.NewKeptnHandler(rh)
	fmt.Println("get shipyard for ", project)
	shipyard, err := keptnHandler.GetShipyard(project)
	if err != nil {
		return false
	}
	fmt.Println("shipyard: ", shipyard)
	for _, s := range shipyard.Stages {
		if s.Name == stage && s.RemediationStrategy == "automated" {
			return true
		}
	}

	return false
}

// gets Helm release name by pod name
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

// splits helm release name into project, stage and service
func splitReleaseName(releasename string) (project string, stage string, service string) {
	// currently no "-" in project and stage name are allowed, thus "-" is used to split
	s := strings.SplitN(releasename, "-", 3)
	project = s[0]
	stage = s[1]
	// remove the "-generated" suffix
	service = strings.Replace(s[2], "-generated", "", 1)
	return project, stage, service
}

// gets replica count from helm chart
func getReplicaCount(logger *keptnutils.Logger, project, stage, service, configServiceURL string) (int, error) {
	helmChartName := service + "-generated"
	// Read chart
	chart, err := keptnutils.GetChart(project, service, stage, helmChartName, configServiceURL)
	if err != nil {
		return -1, err
	}
	logger.Debug("chart found for affected service")

	for _, template := range chart.Templates {
		dec := kyaml.NewYAMLToJSONDecoder(bytes.NewReader(template.Data))
		for {
			var document interface{}
			err := dec.Decode(&document)
			if err == io.EOF {
				break
			}
			if err != nil {
				return -1, err
			}

			doc, err := json.Marshal(document)
			if err != nil {
				return -1, err
			}

			var depl appsv1.Deployment
			if err := json.Unmarshal(doc, &depl); err == nil && keptnutils.IsDeployment(&depl) {
				// It is a deployment
				fmt.Println(depl.Spec.Replicas)
				return int(*depl.Spec.Replicas), nil
			}
		}
	}

	return -1, errors.New("could not find replica count")
}

// creates and sends cloud event
func createAndSendEvent(shkeptncontext string, project string,
	service string, stage string, propertyPath string, propertyValue interface{}) error {

	source, _ := url.Parse("https://github.com/keptn/keptn/remediation-service")
	contentType := "application/json"

	pc := keptnevents.PropertyChange{
		PropertyPath: propertyPath,
		Value:        propertyValue,
	}

	configChangedEvent := keptnevents.ConfigurationChangeEventData{
		Project:           project,
		Service:           service,
		Stage:             stage,
		DeploymentChanges: []keptnevents.PropertyChange{pc},
	}

	event := cloudevents.Event{
		Context: cloudevents.EventContextV02{
			ID:          uuid.New().String(),
			Time:        &types.Timestamp{Time: time.Now()},
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

// sendTestsFinishedEvent sends a Cloud Event of type sh.keptn.events.tests-finished to the event broker
func sendTestsFinishedEvent(shkeptncontext string, project string, stage string, service string) error {

	source, _ := url.Parse("remediation-service")
	contentType := "application/json"

	testFinishedData := keptnevents.TestsFinishedEventData{}
	testFinishedData.Project = project
	testFinishedData.Stage = stage
	testFinishedData.Service = service
	testFinishedData.TestStrategy = "real-user"

	event := cloudevents.Event{
		Context: cloudevents.EventContextV02{
			ID:          uuid.New().String(),
			Time:        &types.Timestamp{Time: time.Now()},
			Type:        "sh.keptn.events.tests-finished",
			Source:      types.URLRef{URL: *source},
			ContentType: &contentType,
			Extensions:  map[string]interface{}{"shkeptncontext": shkeptncontext},
		}.AsV02(),
		Data: testFinishedData,
	}

	return sendEvent(event)
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
