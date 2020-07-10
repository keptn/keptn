package main

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	configmodels "github.com/keptn/go-utils/pkg/api/models"
	configutils "github.com/keptn/go-utils/pkg/api/utils"

	keptn "github.com/keptn/go-utils/pkg/lib"
	keptnapi "github.com/keptn/go-utils/pkg/api/utils"

	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/client"
	cloudeventshttp "github.com/cloudevents/sdk-go/pkg/cloudevents/transport/http"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/types"
	"github.com/google/uuid"
	"github.com/kelseyhightower/envconfig"

	"gopkg.in/yaml.v2"
)

const configservice = "CONFIGURATION_SERVICE"
const api = "API"

type envConfig struct {
	Port int    `envconfig:"RCV_PORT" default:"8080"`
	Path string `envconfig:"RCV_PATH" default:"/"`
}

func main() {
	var env envConfig
	if err := envconfig.Process("", &env); err != nil {
		log.Fatalf("failed to process env var: %s", err)
	}
	go keptnapi.RunHealthEndpoint("10999")
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

	logger := keptn.NewLogger(shkeptncontext, event.Context.GetID(), "shipyard-service")

	// open websocket connection to api component
	endPoint, err := keptn.GetServiceEndpoint(api)
	if err != nil {
		return err
	}

	if endPoint.Host == "" {
		const errorMsg = "Host of api not set"
		logger.Error(errorMsg)
		return errors.New(errorMsg)
	}

	connData := &keptn.ConnectionData{}
	if err := event.DataAs(connData); err != nil {
		logger.Error(fmt.Sprintf("Data of the event is incompatible. %s", err.Error()))
		return err
	}

	ws, _, err := keptn.OpenWS(*connData, endPoint)
	if err != nil {
		logger.Error(fmt.Sprintf("Opening websocket connection failed. %s", err.Error()))
		return err
	}
	defer ws.Close()

	if event.Type() == keptn.InternalProjectCreateEventType {
		_, err := createProjectAndProcessShipyard(event, *logger, ws)
		if err := closeWebsocketWithMessage(event, err, "Shipyard successfully processed", *logger, ws); err != nil {
			return err
		}
		return nil
	} else if event.Type() == keptn.InternalProjectDeleteEventType {
		err := deleteProject(event, *logger, ws)
		if err := closeWebsocketWithMessage(event, err, "Project successfully deleted", *logger, ws); err != nil {
			return err
		}
		return nil
	}

	const errorMsg = "Received unexpected keptn event that cannot be processed"
	if err := keptn.WriteWSLog(ws, createEventCopy(event, "sh.keptn.events.log"), errorMsg, true, "INFO"); err != nil {
		logger.Error(fmt.Sprintf("Could not write log to websocket. %s", err.Error()))
	}
	logger.Error(errorMsg)
	return errors.New(errorMsg)
}

// createProjectAndProcessShipyard creates a project and stages defined in the shipyard
func createProjectAndProcessShipyard(event cloudevents.Event, logger keptn.Logger, ws *websocket.Conn) (*configmodels.Version, error) {
	eventData := keptn.ProjectCreateEventData{}
	if err := event.DataAs(&eventData); err != nil {
		return nil, err
	}

	shipyard := keptn.Shipyard{}
	data, err := base64.StdEncoding.DecodeString(eventData.Shipyard)
	if err != nil {
		logger.Error(fmt.Sprintf("Could not decode shipyard. %s", err.Error()))
		return nil, err
	}
	err = yaml.Unmarshal(data, &shipyard)
	if err != nil {
		logger.Error(fmt.Sprintf("Could not unmarshal shipyard. %s", err.Error()))
		return nil, err
	}
	// create project
	project := configmodels.Project{
		ProjectName:  eventData.Project,
		GitUser:      eventData.GitUser,
		GitToken:     eventData.GitToken,
		GitRemoteURI: eventData.GitRemoteURL,
	}

	if err := createProject(project, logger); err != nil {
		return nil, fmt.Errorf("Creating project %s failed. %s", project.ProjectName, err.Error())
	}
	if err := keptn.WriteWSLog(ws, createEventCopy(event, "sh.keptn.events.log"), fmt.Sprintf("Project %s created", project.ProjectName), false, "INFO"); err != nil {
		logger.Error(fmt.Sprintf("Could not write log to websocket. %s", err.Error()))
	}

	// process shipyard file and create stages
	for _, shipyardStage := range shipyard.Stages {
		if err := createStage(project, shipyardStage.Name, logger); err != nil {
			return nil, fmt.Errorf("Creating stage %s failed. %s", shipyardStage.Name, err.Error())
		}
		if err := keptn.WriteWSLog(ws, createEventCopy(event, "sh.keptn.events.log"), fmt.Sprintf("Stage %s created", shipyardStage.Name), false, "INFO"); err != nil {
			logger.Error(fmt.Sprintf("Could not write log to websocket. %s", err.Error()))
		}
	}
	// store shipyard.yaml
	return storeResourceForProject(project.ProjectName, string(data), logger)
}

// deleteProject processes event and deletes project
func deleteProject(event cloudevents.Event, logger keptn.Logger, ws *websocket.Conn) error {
	eventData := keptn.ProjectDeleteEventData{}
	if err := event.DataAs(&eventData); err != nil {
		return err
	}

	project := configmodels.Project{
		ProjectName: eventData.Project,
	}

	var msg string
	keptnHandler, err := keptn.NewKeptn(&event, keptn.KeptnOpts{})
	if err != nil {
		logger.Error("Could not initialize Keptn handler: " + err.Error())
		return err
	}
	msg, err = getDeleteInfoMessage(keptnHandler)
	if err != nil {
		msg = fmt.Sprintf("Shipyard of project %s cannot be retrieved anymore. ", project.ProjectName)
		msg += "After deleting the project, the namespaces containing the services are still available. " +
			"This may cause problems if a project with the same name is created later."
	}
	if err := keptn.WriteWSLog(ws, createEventCopy(event, "sh.keptn.events.log"),
		msg, false, "INFO"); err != nil {
		logger.Error(fmt.Sprintf("Could not write log to websocket. %s", err.Error()))
	}

	// get remote url of project
	projectResp, err := getProject(project, logger)
	if err != nil {
		if err := keptn.WriteWSLog(ws, createEventCopy(event, "sh.keptn.events.log"),
			fmt.Sprintf("Project %s cannot be retrieved anymore. Any Git upstream of the project will not be deleted.", project.ProjectName),
			false, "INFO"); err != nil {
			logger.Error(fmt.Sprintf("Could not write log to websocket. %s", err.Error()))
		}
	} else if projectResp != nil && projectResp.GitRemoteURI != "" {
		if err := keptn.WriteWSLog(ws, createEventCopy(event, "sh.keptn.events.log"),
			fmt.Sprintf("The Git upstream of the project will not be deleted: %s", projectResp.GitRemoteURI), false, "INFO"); err != nil {
			logger.Error(fmt.Sprintf("Could not write log to websocket. %s", err.Error()))
		}
	}

	configServiceURL, err := keptn.GetServiceEndpoint(configservice)
	if err != nil {
		logger.Error(fmt.Sprintf("Could not get service endpoint for %s: %s", configservice, err.Error()))
		return err
	}
	// delete project
	prjHandler := configutils.NewProjectHandler(configServiceURL.String())
	_, mErr := prjHandler.DeleteProject(project)
	if err != nil {
		return fmt.Errorf("Deleting project %s failed. %s", project.ProjectName, *mErr.Message)
	}
	logger.Info(fmt.Sprintf("Project %s deleted", project.ProjectName))

	return nil
}

// storeResourceForProject stores the resource for a project using the keptn.ResourceHandler
func storeResourceForProject(projectName, shipyard string, logger keptn.Logger) (*configmodels.Version, error) {
	configServiceURL, err := keptn.GetServiceEndpoint(configservice)
	if err != nil {
		logger.Error(fmt.Sprintf("Could not get service endpoint for %s: %s", configservice, err.Error()))
		return nil, err
	}
	handler := configutils.NewResourceHandler(configServiceURL.String())
	uri := "shipyard.yaml"
	resource := configmodels.Resource{ResourceURI: &uri, ResourceContent: shipyard}
	versionStr, err := handler.CreateProjectResources(projectName, []*configmodels.Resource{&resource})
	if err != nil {
		return nil, fmt.Errorf("Storing %s file failed. %s", *resource.ResourceURI, err.Error())
	}

	logger.Info(fmt.Sprintf("Resource %s successfully stored", *resource.ResourceURI))
	return &configmodels.Version{Version: versionStr}, nil
}

// closeWebsocketWithMessage sends a log message to the websocket
func closeWebsocketWithMessage(event cloudevents.Event, err error, message string,
	logger keptn.Logger, ws *websocket.Conn) error {
	var webSocketMessage = message

	if err != nil { // error
		webSocketMessage = fmt.Sprintf("%s.", err.Error())
		logger.Error(webSocketMessage)
	} else { // success
		logger.Info(message)
	}

	if err := keptn.WriteWSLog(ws, createEventCopy(event, "sh.keptn.events.log"), webSocketMessage, true, "INFO"); err != nil {
		logger.Error(fmt.Sprintf("Could not write log to websocket. %s", err.Error()))
	}
	return err
}

// createProject creates a project by using the configuration-service
func createProject(project configmodels.Project, logger keptn.Logger) error {
	configServiceURL, err := keptn.GetServiceEndpoint(configservice)
	if err != nil {
		logger.Error(fmt.Sprintf("Could not get service endpoint for %s: %s", configservice, err.Error()))
		return err
	}

	prjHandler := configutils.NewProjectHandler(configServiceURL.String())
	_, errObj := prjHandler.CreateProject(project)

	if errObj == nil {
		logger.Info("Project successfully created")
		return nil
	}
	return errors.New(*errObj.Message)
}

func getDeleteInfoMessage(keptnHandler *keptn.Keptn) (string, error) {

	shipyard, err := keptnHandler.GetShipyard()
	if err != nil {
		return "", fmt.Errorf("error when getting shipyard: %v", err)
	}
	msg := ""
	for _, stage := range shipyard.Stages {
		namespace := keptnHandler.KeptnBase.Project + "-" + stage.Name
		msg += fmt.Sprintf("Namespace %s is not managed by Keptn anymore and not deleted. This may cause problems if "+
			"a project with the same name is created later. "+
			"If you would like to delete the namespace, please execute "+
			"'kubectl delete ns %s'\n", namespace, namespace)
	}
	return strings.TrimSpace(msg), nil
}

// getProject returns a project by using the configuration-service
func getProject(project configmodels.Project, logger keptn.Logger) (*configmodels.Project, error) {
	configServiceURL, err := keptn.GetServiceEndpoint(configservice)
	if err != nil {
		logger.Error(fmt.Sprintf("Could not get service endpoint for %s: %s", configservice, err.Error()))
		return nil, err
	}

	prjHandler := configutils.NewProjectHandler(configServiceURL.String())
	respProject, respError := prjHandler.GetProject(project)
	if respError != nil {
		return nil, fmt.Errorf("Error in getting project: %s", project.ProjectName)
	}

	return respProject, nil
}

// createStage creates a stage by using the configuration-service
func createStage(project configmodels.Project, stage string, logger keptn.Logger) error {

	configServiceURL, err := keptn.GetServiceEndpoint(configservice)
	if err != nil {
		logger.Error(fmt.Sprintf("Could not get service endpoint for %s: %s", configservice, err.Error()))
		return err
	}

	stageHandler := configutils.NewStageHandler(configServiceURL.String())
	_, errorObj := stageHandler.CreateStage(project.ProjectName, stage)

	if errorObj == nil {
		logger.Info("Stage successfully created")
		return nil
	} else if errorObj != nil {
		return errors.New(*errorObj.Message)
	}

	return fmt.Errorf("Error in creating new stage: %s", err.Error())
}

// createEventCopy creates a deep copy of a CloudEvent
func createEventCopy(eventSource cloudevents.Event, eventType string) cloudevents.Event {
	var shkeptncontext string
	eventSource.Context.ExtensionAs("shkeptncontext", &shkeptncontext)
	var shkeptnphaseid string
	eventSource.Context.ExtensionAs("shkeptnphaseid", &shkeptnphaseid)
	var shkeptnphase string
	eventSource.Context.ExtensionAs("shkeptnphase", &shkeptnphase)
	var shkeptnstepid string
	eventSource.Context.ExtensionAs("shkeptnstepid", &shkeptnstepid)
	var shkeptnstep string
	eventSource.Context.ExtensionAs("shkeptnstep", &shkeptnstep)

	source, _ := url.Parse("shipyard-service")
	contentType := "application/json"

	event := cloudevents.Event{
		Context: cloudevents.EventContextV02{
			ID:          uuid.New().String(),
			Time:        &types.Timestamp{Time: time.Now()},
			Type:        eventType,
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
	}

	return event
}
