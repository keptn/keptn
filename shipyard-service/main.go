package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/keptn/keptn/shipyard-service/models"

	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/client"
	cloudeventshttp "github.com/cloudevents/sdk-go/pkg/cloudevents/transport/http"
	"github.com/kelseyhightower/envconfig"
	keptnutils "github.com/keptn/go-utils/pkg/utils"
)

const timeout = 60
const configservice = "CONFIGURATION_SERVICE"

type envConfig struct {
	Port int    `envconfig:"RCV_PORT" default:"8080"`
	Path string `envconfig:"RCV_PATH" default:"/"`
}

type createProjectEvent struct {
	Project  string          `json:"project"`
	Upstream string          `json:"upstream"`
	User     string          `json:"user"`
	Token    string          `json:"token"`
	Shipyard []shipyardStage `json:"stages"`
}

type shipyardStage struct {
	Name               string `json:"name"`
	DeplyomentStrategy string `json:"deployment_strategy"`
	TestStrategy       string `json:"test_strategy"`
}

type Client struct {
	httpClient *http.Client
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

func NewClient() *Client {
	client := Client{
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}

	return &client
}

func gotEvent(ctx context.Context, event cloudevents.Event) error {
	var shkeptncontext string
	event.Context.ExtensionAs("shkeptncontext", &shkeptncontext)

	logger := keptnutils.NewLogger(shkeptncontext, event.Context.GetID(), "project-mgt-service")

	if event.Type() == "sh.keptn.events.project.create" {
		eventData := &createProjectEvent{}
		if err := event.DataAs(eventData); err != nil {
			logger.Error(fmt.Sprintf("data of the event is incompatible: %s", err.Error()))
			return err
		}
		return createProjectAndProcessShipyard(*eventData, *logger)
	}

	const errorMsg = "received unexpected keptn event that cannot be processed"
	logger.Error(errorMsg)
	return errors.New(errorMsg)
}

func createProjectAndProcessShipyard(eventData createProjectEvent, logger keptnutils.Logger) error {

	client := NewClient()

	project := models.Project{}
	project.ProjectName = eventData.Project
	if err := client.createProject(project, logger); err != nil {
		logger.Error(fmt.Sprintf("project not created"))
		return err
	}

	for _, shipyardStage := range eventData.Shipyard {
		stage := models.Stage{}
		stage.StageName = shipyardStage.Name

		if err := client.createStage(project, stage, logger); err != nil {
			logger.Error(fmt.Sprintf("stage not created"))
			return err
		}
	}

	shipyard := models.Resource{}

	var resourceURI = "shipyard.yaml"
	shipyard.ResourceURI = &resourceURI
	shipyard.ResourceContent, _ = json.Marshal(eventData.Shipyard)

	resources := []*models.Resource{&shipyard}

	if err := client.storeResource(project, resources, logger); err != nil {
		logger.Error(fmt.Sprintf("resource not stored"))
		return err
	}

	return nil
}

func getEndpoint(logger keptnutils.Logger) (url.URL, error) {
	url, err := url.Parse(os.Getenv(configservice))
	if err != nil {
		logger.Error(fmt.Sprintf("failed to retrieve value from ENVIRONMENT_VARIBLE: %s", configservice))
	}
	return *url, err
}

func postRequest(client *Client, path string, body []byte, logger keptnutils.Logger) (*http.Response, error) {
	endPoint, err := getEndpoint(logger)
	if err != nil {
		const errorMsg = "post request execution aborted"
		logger.Error(errorMsg)
		return nil, errors.New(errorMsg)
	}

	if endPoint.Host == "" {
		const errorMsg = "host of configuration-service not set"
		logger.Error(errorMsg)
		return nil, errors.New(errorMsg)
	}

	eventURL := endPoint
	eventURL.Path = path

	req, err := http.NewRequest("POST", eventURL.String(), bytes.NewReader(body))
	if err != nil {
		const errorMsg = "building new request failed"
		logger.Error(errorMsg)
		return nil, errors.New(errorMsg)
	}

	return client.httpClient.Do(req)
}

func (client *Client) createProject(project models.Project, logger keptnutils.Logger) error {
	data, err := project.MarshalBinary()
	if err != nil {
		logger.Error("project can not be encoded")
		return err
	}

	resp, err := postRequest(client, "/project", data, logger)
	if err != nil {
		logger.Error("request failed - creating project was unsuccessful")
		return err
	}

	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		logger.Info("project created successfully")
	}

	return nil
}

func (client *Client) createStage(project models.Project, stage models.Stage, logger keptnutils.Logger) error {
	data, err := project.MarshalBinary()
	if err != nil {
		logger.Error("stage can not be encoded")
		return err
	}

	resp, err := postRequest(client, fmt.Sprintf("project/%s/stage", project.ProjectName), data, logger)
	if err != nil {
		logger.Error("request failed - creating stage was unsuccessful")
		return err
	}

	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		logger.Info("stage created successfully")
	}

	return nil
}

func (client *Client) storeResource(project models.Project, resources []*models.Resource, logger keptnutils.Logger) error {
	data, err := project.MarshalBinary()
	if err != nil {
		logger.Error("resources can not be encoded")
		return err
	}

	resp, err := postRequest(client, fmt.Sprintf("project/%s/resource", project.ProjectName), data, logger)
	if err != nil {
		logger.Error("request failed - storing resource was unsuccessful")
		return err
	}

	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		logger.Info("resource stored successfully")
	}

	return nil
}
