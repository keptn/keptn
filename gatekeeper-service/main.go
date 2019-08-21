package main

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/client"
	cloudeventshttp "github.com/cloudevents/sdk-go/pkg/cloudevents/transport/http"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/types"
	"github.com/google/uuid"
	"github.com/kelseyhightower/envconfig"
	keptnutils "github.com/keptn/go-utils/pkg/utils"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	githttp "gopkg.in/src-d/go-git.v4/plumbing/transport/http"
	"gopkg.in/yaml.v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const configservice = "CONFIGURATION_SERVICE"
const eventbroker = "EVENTBROKER"

type envConfig struct {
	// Port on which to listen for cloudevents
	Port int    `envconfig:"RCV_PORT" default:"8080"`
	Path string `envconfig:"RCV_PATH" default:"/"`
}

type partialShipyard struct {
	Stages []partialStage `json:"stages"`
}

type partialStage struct {
	Name string `json:"name"`
}

type evaluationDoneEvent struct {
	GitHubOrg          string `json:"githuborg"`
	Project            string `json:"project"`
	TestStrategy       string `json:"teststrategy"`
	DeploymentStrategy string `json:"deploymentstrategy"`
	Stage              string `json:"stage"`
	Service            string `json:"service"`
	Image              string `json:"image"`
	Tag                string `json:"tag"`
	EvaluationPassed   bool   `json:"evaluationpassed"`
}

func main() {
	var env envConfig
	if err := envconfig.Process("", &env); err != nil {
		log.Fatalf("Failed to process env var: %s", err)
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
		log.Fatalf("failed to create transport, %v", err)
	}
	c, err := client.New(t)
	if err != nil {
		log.Fatalf("failed to create client, %v", err)
	}

	log.Fatalf("failed to start receiver: %s", c.StartReceiver(ctx, gotEvent))

	return 0
}

func gotEvent(ctx context.Context, event cloudevents.Event) error {
	var shkeptncontext string
	event.Context.ExtensionAs("shkeptncontext", &shkeptncontext)

	logger := keptnutils.NewLogger(shkeptncontext, event.Context.GetID(), "gatekeeper-service")

	data := &evaluationDoneEvent{}
	if err := event.DataAs(data); err != nil {
		logger.Error(fmt.Sprintf("Got Data Error: %s", err.Error()))
		return err
	}

	if event.Type() != "sh.keptn.events.evaluation-done" {
		const errorMsg = "Received unexpected keptn event"
		logger.Error(errorMsg)
		return errors.New(errorMsg)
	}

	go doGateKeeping(event, shkeptncontext, *data, logger)

	return nil
}

func doGateKeeping(event cloudevents.Event, shkeptncontext string, data evaluationDoneEvent, logger *keptnutils.Logger) {

	if data.EvaluationPassed {

		logger.Info("Evaluation is passed")

		nextStage, err := getNextStage(data, logger)
		if err != nil {
			logger.Error(fmt.Sprintf("Error obtaining the next stage: %s", err.Error()))
			return
		}

		if nextStage != "" {
			logger.Info("Promote artifact to stage " + nextStage)
			if err := sendNewArtifactEvent(shkeptncontext, event, nextStage, logger); err != nil {
				logger.Error(fmt.Sprintf("Error sending new artifact event: %s", err.Error()))
				return
			}
		} else {
			logger.Info("No further stage available: End of promotion")
		}
	} else {
		logger.Info("Evaluation not passed, hence do not promote artifact to next stage")

		if strings.ToLower(data.DeploymentStrategy) == "blue_green_service" {
			logger.Info("Reroute traffic to old version")

			repo, err := keptnutils.Checkout(data.GitHubOrg, data.Project, data.Stage)
			if err != nil {
				logger.Error(fmt.Sprintf("Error occured checking out the config repo: %s", err.Error()))
				return
			}

			err = switchWeights(repo, data)
			if err != nil {
				logger.Error(fmt.Sprintf("Error rerouting traffic: %s", err.Error()))
				return
			}

			if err := sendConfigurationChangedEvent(shkeptncontext, event, logger); err != nil {
				logger.Error(fmt.Sprintf("Error sending configuration changed event: %s", err.Error()))
				return
			}
		}
	}
}

func getNextStage(data evaluationDoneEvent, logger *keptnutils.Logger) (string, error) {

	resource, err := retrieveResourceForProject(data.Project, "shipyard.yaml", logger)
	if err != nil {
		return "", err
	}

	var shipyard partialShipyard
	err = yaml.Unmarshal([]byte(resource.ResourceContent), &shipyard)
	if err != nil {
		return "", err
	}

	currentFound := false
	for _, stage := range shipyard.Stages {
		if currentFound {
			// Here, we return the next stage
			return stage.Name, nil
		}
		if stage.Name == data.Stage {
			currentFound = true
		}
	}
	return "", nil
}

func switchWeights(repo *git.Repository, data evaluationDoneEvent) error {

	filePathWithinPrj := "helm-chart/templates/" + data.Service + "-istio-virtual-service.yaml"
	filePath := data.Project + "/" + filePathWithinPrj
	yamlFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	var f interface{}
	err = yaml.Unmarshal(yamlFile, &f)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	spec := f.(map[interface{}]interface{})["spec"]
	http := spec.(map[interface{}]interface{})["http"].([]interface{})[0]
	routes := http.(map[interface{}]interface{})["route"]
	route1 := routes.([]interface{})[0]
	route2 := routes.([]interface{})[1]

	tmp := route1.(map[interface{}]interface{})["weight"]
	route1.(map[interface{}]interface{})["weight"] = route2.(map[interface{}]interface{})["weight"].(int)
	route2.(map[interface{}]interface{})["weight"] = tmp

	out, err := yaml.Marshal(f)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filePath, out, 0644)
	if err != nil {
		return err
	}

	w, err := repo.Worktree()
	if err != nil {
		return err
	}
	_, err = w.Add(filePathWithinPrj)
	if err != nil {
		return err
	}
	_, err = w.Commit("[keptn]: Switched blue green due to failed evaluation.", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "keptn",
			Email: "keptn@dynatrace.com",
			When:  time.Now(),
		},
	})
	if err != nil {
		return err
	}

	user, token, err := getGithubCredentials()
	if err != nil {
		return err
	}

	err = repo.Push(&git.PushOptions{
		Auth: &githttp.BasicAuth{
			Username: user,
			Password: token,
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func sendNewArtifactEvent(shkeptncontext string, incomingEvent cloudevents.Event, nextStage string, logger *keptnutils.Logger) error {

	source, _ := url.Parse("jmeter-service")
	contentType := "application/json"

	// Set next stage
	var newArtifactEvent interface{}
	if err := incomingEvent.DataAs(&newArtifactEvent); err != nil {
		logger.Error(fmt.Sprintf("Got Data Error: %s", err.Error()))
		return err
	}
	newArtifactEvent.(map[string]interface{})["stage"] = nextStage

	event := cloudevents.Event{
		Context: cloudevents.EventContextV02{
			ID:          uuid.New().String(),
			Type:        "sh.keptn.events.new-artifact",
			Source:      types.URLRef{URL: *source},
			ContentType: &contentType,
			Extensions:  map[string]interface{}{"shkeptncontext": shkeptncontext},
		}.AsV02(),
		Data: newArtifactEvent,
	}

	t, err := cloudeventshttp.New(
		cloudeventshttp.WithTarget("http://event-broker.keptn.svc.cluster.local/keptn"),
		cloudeventshttp.WithEncoding(cloudeventshttp.StructuredV02),
	)
	if err != nil {
		return errors.New("Failed to create transport:" + err.Error())
	}

	c, err := client.New(t)
	if err != nil {
		return errors.New("Failed to create HTTP client:" + err.Error())
	}

	if _, err := c.Send(context.Background(), event); err != nil {
		return errors.New("Failed to send cloudevent:, " + err.Error())
	}
	return nil
}

func sendConfigurationChangedEvent(shkeptncontext string, incomingEvent cloudevents.Event, logger *keptnutils.Logger) error {

	source, _ := url.Parse("jmeter-service")
	contentType := "application/json"

	// Remove test strategy
	var configurationChangedEventData interface{}
	if err := incomingEvent.DataAs(&configurationChangedEventData); err != nil {
		logger.Error(fmt.Sprintf("Got Data Error: %s", err.Error()))
		return err
	}
	configurationChangedEventData.(map[string]interface{})["teststrategy"] = ""

	event := cloudevents.Event{
		Context: cloudevents.EventContextV02{
			ID:          uuid.New().String(),
			Type:        "sh.keptn.events.configuration-changed",
			Source:      types.URLRef{URL: *source},
			ContentType: &contentType,
			Extensions:  map[string]interface{}{"shkeptncontext": shkeptncontext},
		}.AsV02(),
		Data: configurationChangedEventData,
	}

	t, err := cloudeventshttp.New(
		cloudeventshttp.WithTarget("http://event-broker.keptn.svc.cluster.local/keptn"),
		cloudeventshttp.WithEncoding(cloudeventshttp.StructuredV02),
	)
	if err != nil {
		return errors.New("Failed to create transport:" + err.Error())
	}

	c, err := client.New(t)
	if err != nil {
		return errors.New("Failed to create HTTP client:" + err.Error())
	}

	if _, err := c.Send(context.Background(), event); err != nil {
		return errors.New("Failed to send cloudevent:, " + err.Error())
	}
	return nil
}

// getGithubCredentials reads the secret 'github-credentials' and returns the user, token, and potentially any occured error
func getGithubCredentials() (string, string, error) {

	api, err := keptnutils.GetKubeAPI(true)
	if err != nil {
		return "", "", err
	}

	getOptions := metav1.GetOptions{}
	secret, err := api.Secrets("keptn").Get("github-credentials", getOptions)
	if err != nil {
		return "", "", err
	}

	return string(secret.Data["user"]), string(secret.Data["token"]), nil
}

// retrieveResourceForProject retrieves a resource stored at a project entity
func retrieveResourceForProject(projectName string, resourceURI string, logger *keptnutils.Logger) (*keptnutils.Resource, error) {
	eventURL, err := getServiceEndpoint(configservice)
	resourceHandler := keptnutils.NewResourceHandler(eventURL.Host)

	resource, err := resourceHandler.GetProjectResource(projectName, resourceURI)
	if err != nil {
		return nil, fmt.Errorf("Failed to retrieve resource %s", resourceURI, err.Error())
	}

	logger.Info(resource.ResourceContent)

	return resource, nil
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
