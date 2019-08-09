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

	"github.com/keptn/keptn/configuration-service/models"

	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/client"
	cloudeventshttp "github.com/cloudevents/sdk-go/pkg/cloudevents/transport/http"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/types"
	"github.com/google/uuid"
	"github.com/kelseyhightower/envconfig"
	keptnutils "github.com/keptn/go-utils/pkg/utils"
)

const timeout = 60
const configservice = "CONFIGURATION_SERVICE"
const eventbroker = "EVENTBROKER"

type envConfig struct {
	Port int    `envconfig:"RCV_PORT" default:"8080"`
	Path string `envconfig:"RCV_PATH" default:"/"`
}

type createProjectEventData struct {
	Project      string          `json:"project"`
	GitRemoteURL string          `json:"gitremoteurl"`
	GitUser      string          `json:"gituser"`
	GitToken     string          `json:"gittoken"`
	Shipyard     []shipyardStage `json:"stages"`
}

type shipyardStage struct {
	Name               string `json:"name"`
	DeplyomentStrategy string `json:"deployment_strategy"`
	TestStrategy       string `json:"test_strategy"`
}

type doneEventData struct {
	Result  string `json:"result"`
	Message string `json:"message"`
	Version string `json:"version"`
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

func newClient() *Client {
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

	logger := keptnutils.NewLogger(shkeptncontext, event.Context.GetID(), "shipyard-service")

	//if event.Type() == "sh.keptn.internal.events.project.create" {
	if event.Type() == "create.project" {
			eventData := &createProjectEventData{}
		if err := event.DataAs(eventData); err != nil {
			logger.Error(fmt.Sprintf("data of the event is incompatible: %s", err.Error()))
			return err
		}

		err := createProjectAndProcessShipyard(*eventData, *logger)
		if err := respondWithDoneEvent(event, err, *logger); err != nil {
			logger.Error(fmt.Sprintf("no sh.keptn.event.done sent: %s", err.Error()))
			return err
		}
		return nil
	}

	const errorMsg = "received unexpected keptn event that cannot be processed"
	logger.Error(errorMsg)
	return errors.New(errorMsg)
}

func respondWithDoneEvent(event cloudevents.Event, err error, logger keptnutils.Logger) error {
	if err != nil {
		logger.Error(err.Error())
		if err := sendDoneEvent(event, "error", err.Error()); err != nil {
			logger.Error(err.Error())
		}
		return err
	}

	logger.Info("Create project event processed successfully")
	if err := sendDoneEvent(event, "success", "Create project event processed successfully"); err != nil {
		logger.Error(err.Error())
		return err
	}
	return nil
}

func createProjectAndProcessShipyard(eventData createProjectEventData, logger keptnutils.Logger) error {

	client := newClient()

	project := models.Project{}
	project.ProjectName = eventData.Project
	if err := client.createProject(project, logger); err != nil {
		logger.Error(err.Error())
		return errors.New("processing shipyard failed at creating project")
	}

	for _, shipyardStage := range eventData.Shipyard {
		stage := models.Stage{}
		stage.StageName = shipyardStage.Name

		if err := client.createStage(project, stage, logger); err != nil {
			logger.Error(err.Error())
			return errors.New("processing shipyard failed at creating stage: " + stage.StageName)
		}
	}

	shipyard := models.Resource{}

	var resourceURI = "shipyard.yaml"
	shipyard.ResourceURI = &resourceURI
	shipyard.ResourceContent, _ = json.Marshal(eventData.Shipyard)

	resources := []*models.Resource{&shipyard}

	if err := client.storeResource(project, resources, logger); err != nil {
		logger.Error(err.Error())
		return errors.New("processing shipyard failed at storing shipyard.yaml")
	}

	return nil
}

func getServiceEndpoint(service string) (url.URL, error) {
	url, err := url.Parse(os.Getenv(service))
	if err != nil {
		return *url, fmt.Errorf("failed to retrieve value from ENVIRONMENT_VARIBLE: %s", service)
	}
	return *url, nil
}

func postRequest(client *Client, path string, body []byte) (*http.Response, error) {
	endPoint, err := getServiceEndpoint(configservice)
	if err != nil {
		return nil, err
	}

	if endPoint.Host == "" {
		return nil, errors.New("host of configuration-service not set")
	}

	eventURL := endPoint
	eventURL.Path = path

	req, err := http.NewRequest("POST", eventURL.String(), bytes.NewReader(body))
	if err != nil {
		return nil, errors.New("failed to build new request")
	}

	return client.httpClient.Do(req)
}

func (client *Client) createProject(project models.Project, logger keptnutils.Logger) error {
	data, err := project.MarshalBinary()
	if err != nil {
		logger.Error(err.Error())
		return errors.New("failed to create project")
	}

	resp, err := postRequest(client, "/project", data)
	if err != nil {
		logger.Error(err.Error())
		return errors.New("failed to create project")
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
		logger.Error(err.Error())
		return errors.New("failed to create stage")
	}

	resp, err := postRequest(client, fmt.Sprintf("project/%s/stage", project.ProjectName), data)
	if err != nil {
		logger.Error(err.Error())
		return errors.New("failed to create stage")
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
		logger.Error(err.Error())
		return errors.New("failed to store resource")
	}

	resp, err := postRequest(client, fmt.Sprintf("project/%s/resource", project.ProjectName), data)
	if err != nil {
		logger.Error(err.Error())
		return errors.New("failed to store resource")
	}

	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		logger.Info("resource stored successfully")
	}

	return nil
}

func sendDoneEvent(receivedEvent cloudevents.Event, result string, message string) error {

	var shkeptncontext string
	receivedEvent.Context.ExtensionAs("shkeptncontext", &shkeptncontext)
	var shkeptnphaseid string
	receivedEvent.Context.ExtensionAs("shkeptnphaseid", &shkeptnphaseid)
	var shkeptnphase string
	receivedEvent.Context.ExtensionAs("shkeptnphase", &shkeptnphase)
	var shkeptnstepid string
	receivedEvent.Context.ExtensionAs("shkeptnstepid", &shkeptnstepid)
	var shkeptnstep string
	receivedEvent.Context.ExtensionAs("shkeptnstep", &shkeptnstep)

	source, _ := url.Parse("shipyard-service")
	contentType := "application/json"

	eventData := doneEventData{
		Result:  result,
		Message: message,
	}

	event := cloudevents.Event{
		Context: cloudevents.EventContextV02{
			ID:          uuid.New().String(),
			Type:        "sh.keptn.events.done",
			Source:      types.URLRef{URL: *source},
			ContentType: &contentType,
			Extensions: map[string]interface{}{
				"shkeptncontext": shkeptncontext,
				"shkeptnphaseid": shkeptnphaseid,
				"shkeptnphase":   shkeptnphase,
				"shkeptnstepid":  shkeptnstepid,
				"shkeptnstep":    shkeptnstep,
			},
		}.AsV02(),
		Data: eventData,
	}

	endPoint, err := getServiceEndpoint(eventbroker)
	if err != nil {
		return errors.New("failed to retrieve endpoint of eventbroker: %s" + err.Error())
	}

	if endPoint.Host == "" {
		return errors.New("host of eventbroker not set")
	}

	transport, err := cloudeventshttp.New(
		cloudeventshttp.WithTarget(endPoint.String()),
		cloudeventshttp.WithEncoding(cloudeventshttp.StructuredV02),
	)
	if err != nil {
		return errors.New("failed to create transport: " + err.Error())
	}

	client, err := client.New(transport)
	if err != nil {
		return errors.New("failed to create HTTP client: " + err.Error())
	}

	if _, err := client.Send(context.Background(), event); err != nil {
		return errors.New("failed to send cloudevent sh.keptn.events.done: " + err.Error())
	}

	return nil
}
