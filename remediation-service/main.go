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

	"github.com/keptn/keptn/remediation-service/actions"
	"github.com/keptn/keptn/remediation-service/pkg/utils"
	"gopkg.in/yaml.v2"

	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/client"
	cloudeventshttp "github.com/cloudevents/sdk-go/pkg/cloudevents/transport/http"
	"github.com/kelseyhightower/envconfig"

	configmodels "github.com/keptn/go-utils/pkg/configuration-service/models"
	configutils "github.com/keptn/go-utils/pkg/configuration-service/utils"
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
const remediationFileName = "remediation.yaml"
const configurationserviceconnection = "CONFIGURATION_SERVICE" //"localhost:6060" // "configuration-service:8080"

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

func deriveProblemData(problem *keptnevents.ProblemEventData) {

	if !isProjectAndStageAvailable(problem) {
		deriveFromTags(problem)
	}
	if !isProjectAndStageAvailable(problem) {
		deriveFromImpactedEntity(problem)
	}
}

func isProjectAndStageAvailable(problem *keptnevents.ProblemEventData) bool {
	return problem.Project != "" && problem.Stage != ""
}

// deriveFromTags allows to derive project, stage, and service information from the tags
// Input example: "Tags:":"service:carts, environment:sockshop-production, [Environment]application:sockshop"
func deriveFromTags(problem *keptnevents.ProblemEventData) {

	tags := strings.Split(problem.Tags, ", ")

	for _, tag := range tags {
		if strings.HasPrefix(tag, "service:") {
			problem.Service = tag[len("service:"):]
		} else if strings.HasPrefix(tag, "environment:") {
			environment := tag[len("environment:"):]
			envSplits := strings.Split(environment, "-")
			if len(envSplits) == 2 {
				problem.Project = envSplits[0]
				problem.Stage = envSplits[1]
			}
		}
	}
}

func deriveFromImpactedEntity(problem *keptnevents.ProblemEventData) {
	releasename, err := getReleasename(problem)
	if err != nil {
		// Ignore error as this format is specific for Prometheus
	}
	problem.Project, problem.Stage, problem.Service = splitReleaseName(releasename)
}

func gotEvent(ctx context.Context, event cloudevents.Event) error {
	var shkeptncontext string
	event.Context.ExtensionAs("shkeptncontext", &shkeptncontext)

	logger := keptnutils.NewLogger(shkeptncontext, event.Context.GetID(), "remediation-service")
	logger.Debug("Received event for shkeptncontext:" + shkeptncontext)

	var problemEvent *keptnevents.ProblemEventData
	if event.Type() == keptnevents.ProblemOpenEventType {
		logger.Debug("Received problem notification")
		problemEvent = &keptnevents.ProblemEventData{}
		if err := event.DataAs(problemEvent); err != nil {
			return err
		}
	}

	deriveProblemData(problemEvent)
	if !isProjectAndStageAvailable(problemEvent) {
		return errors.New("Cannot derive project and stage from tags nor impactedentity")
	}

	if problemEvent.State != "OPEN" {
		logger.Info("Received closed problem")
		if problemEvent.Service != "" {
			utils.SendTestsFinishedEvent(shkeptncontext, problemEvent.Project, problemEvent.Stage, problemEvent.Service)
		}
		return nil
	}

	logger.Debug("Received open problem")

	// valide if remediation should be performed
	resourceHandler := configutils.NewResourceHandler(os.Getenv(configurationserviceconnection))
	autoremediate, err := isRemediationEnabled(resourceHandler, problemEvent.Project, problemEvent.Stage)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to check if remediation is enabled: %s", err.Error()))
		return err
	}

	if autoremediate {
		logger.Info(fmt.Sprintf("Remediation enabled for project %s in stage %s", problemEvent.Project, problemEvent.Stage))
	} else {
		logger.Info(fmt.Sprintf("Remediation disabled for project %s in stage %s", problemEvent.Project, problemEvent.Stage))
		return nil
	}

	// get remediation.yaml
	var resource *configmodels.Resource
	if problemEvent.Service != "" {
		resource, err = resourceHandler.GetServiceResource(problemEvent.Project, problemEvent.Stage,
			problemEvent.Service, remediationFileName)
	} else {
		resource, err = resourceHandler.GetStageResource(problemEvent.Project, problemEvent.Stage, remediationFileName)
	}

	if err != nil {
		logger.Error("Failed to get remediation.yaml file")
		return err
	}
	logger.Debug("remediation.yaml for service found")

	// get remediation action from remediation.yaml
	var remediationData keptnmodels.Remediations
	err = yaml.Unmarshal([]byte(resource.ResourceContent), &remediationData)
	if err != nil {
		logger.Error("Could not unmarshal remediation.yaml")
		return err
	}

	actionExecutors := []actions.ActionExecutor{actions.NewScaler(), actions.NewSlower(),
		actions.NewBlackLister(), actions.NewAborter()}

	for _, remediation := range remediationData.Remediations {
		if strings.HasPrefix(problemEvent.ProblemTitle, remediation.Name) {
			logger.Debug("Remediation for problem found")
			// currently only one remediation action is supported
			for _, a := range actionExecutors {
				if a.GetAction() == remediation.Actions[0].Action {
					if err := a.ExecuteAction(problemEvent, shkeptncontext, remediation.Actions[0]); err != nil {
						logger.Error(err.Error())
						return err
					}
					logger.Info(fmt.Sprintf("Remediation action %s successfully applied",
						remediation.Actions[0].Action))
					return nil
				}
			}
		}
	}

	logger.Info("No remediation action found")
	return nil
}

// initialize helm
func init() {
	_, err := keptnutils.ExecuteCommand("helm", []string{"init", "--client-only"})
	if err != nil {
		log.Fatal(err)
	}
}

// get helm release name by impacted entity
func getReleasename(eventData *keptnevents.ProblemEventData) (string, error) {
	var details string
	err := yaml.Unmarshal(eventData.ProblemDetails, &details)
	if err != nil {
		return "", err
	}
	switch impact := details; {
	case strings.HasPrefix(impact, "Pod"):
		return getReleaseByPodName(eventData.ImpactedEntity)
	case strings.HasPrefix(impact, "Service"):
		return "", errors.New("Service remediation not yet supported")
	case strings.HasPrefix(impact, "Node"):
		return "", errors.New("Node remediation not yet supported")
	default:
		return "", errors.New("could not interpret problem details")
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

func isRemediationEnabled(rh *configutils.ResourceHandler, project string, stage string) (bool, error) {
	keptnHandler := keptnutils.NewKeptnHandler(rh)
	shipyard, err := keptnHandler.GetShipyard(project)
	if err != nil {
		return false, err
	}
	for _, s := range shipyard.Stages {
		if s.Name == stage && s.RemediationStrategy == "automated" {
			return true, nil
		}
	}

	return false, nil
}

// gets Helm release name by pod name
func getReleaseByPodName(podname string) (string, error) {
	client, err := getHelmClient()
	if err != nil {
		return "", fmt.Errorf("failed to get Helm client: %v", err)
	}

	releaseList, err := client.ListReleases()
	if err != nil {
		return "", fmt.Errorf("failed to fetch list of releases: %v", err)
	}

	for _, r := range releaseList.GetReleases() {
		rs, err := client.ReleaseStatus(r.Name)
		if err != nil {
			return "", fmt.Errorf("failed to get release status: %v", err)
		}
		if strings.Contains(rs.Info.Status.Resources, podname) {
			return r.Name, nil
		}
	}

	return "", fmt.Errorf("could not find release for pod %s", podname)
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
