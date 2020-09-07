package main

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	keptncommon "github.com/keptn/go-utils/pkg/lib/keptn"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"log"
	"os"
	"strings"

	configmodels "github.com/keptn/go-utils/pkg/api/models"
	configutils "github.com/keptn/go-utils/pkg/api/utils"

	keptnapi "github.com/keptn/go-utils/pkg/api/utils"
	keptn "github.com/keptn/go-utils/pkg/lib"

	cloudevents "github.com/cloudevents/sdk-go/v2"
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

var configServiceURL string

func _main(args []string, env envConfig) int {

	url, err := keptncommon.GetServiceEndpoint(configservice)
	if err != nil {
		log.Fatalf("failed to get service endpoint for %s: %s", configservice, err.Error())
	}
	configServiceURL = url.String()

	ctx := context.Background()
	ctx = cloudevents.WithEncodingStructured(ctx)

	p, err := cloudevents.NewHTTP(cloudevents.WithPath(env.Path), cloudevents.WithPort(env.Port))
	if err != nil {
		log.Fatalf("failed to create client, %v", err)
	}
	c, err := cloudevents.NewClient(p)
	if err != nil {
		log.Fatalf("failed to create client, %v", err)
	}
	log.Fatal(c.StartReceiver(ctx, gotEvent))

	return 0
}

func gotEvent(ctx context.Context, event cloudevents.Event) error {

	serviceName := "shipyard-service"
	keptnHandler, err := keptnv2.NewKeptn(&event, keptncommon.KeptnOpts{
		LoggingOptions: &keptncommon.LoggingOpts{
			EnableWebsocket: true,
			ServiceName:     &serviceName,
		},
	})
	if err != nil {
		fmt.Printf("failed to create Keptn handler: %v", err)
		return err
	}

	if event.Type() == keptn.InternalProjectCreateEventType {
		if err := createProjectAndProcessShipyard(event, keptnHandler.Logger); err != nil {
			keptnHandler.Logger.Error(err.Error())
		}
		keptnHandler.Logger.Terminate("")
	} else if event.Type() == keptn.InternalProjectDeleteEventType {
		if err := deleteProject(event, keptnHandler); err != nil {
			keptnHandler.Logger.Error(err.Error())
		}
		keptnHandler.Logger.Terminate("")
	} else if event.Type() == keptn.InternalServiceCreateEventType {
		terminate, err := createService(event, keptnHandler.Logger)
		if err != nil {
			keptnHandler.Logger.Error(err.Error())
		}
		if terminate {
			keptnHandler.Logger.Terminate("")
		}
		return err
	} else if event.Type() == keptn.InternalServiceDeleteEventType {
		err := deleteService(event, keptnHandler.Logger)
		if err != nil {
			keptnHandler.Logger.Error(err.Error())
		}
		keptnHandler.Logger.Terminate("")
		return err
	} else {
		const errorMsg = "Received unexpected keptn event that cannot be processed"
		keptnHandler.Logger.Terminate(errorMsg)
		return errors.New(errorMsg)
	}
	return nil
}

func deleteService(event cloudevents.Event, logger keptncommon.LoggerInterface) error {
	eventData := keptn.ServiceDeleteEventData{}
	if err := event.DataAs(&eventData); err != nil {
		return err
	}

	stageHandler := configutils.NewStageHandler(configServiceURL)
	stages, err := stageHandler.GetAllStages(eventData.Project)
	if err != nil {
		return fmt.Errorf("Failed to get stages for project %s: %v", eventData.Project, err)
	}

	serviceHandler := configutils.NewServiceHandler(configServiceURL)
	for _, stage := range stages {
		logger.Info("Deleting Keptn service " + eventData.Service + " from stage " + stage.StageName)
		_, err := serviceHandler.DeleteServiceFromStage(eventData.Project, stage.StageName, eventData.Service)
		if err != nil {
			return fmt.Errorf("Failed to create service %s in project %s: %s", eventData.Service, eventData.Project, *err.Message)
		}
	}

	return nil

}

func createService(event cloudevents.Event, logger keptncommon.LoggerInterface) (bool, error) {

	eventData := keptn.ServiceCreateEventData{}
	if err := event.DataAs(&eventData); err != nil {
		return len(eventData.HelmChart) == 0, err
	}

	if !keptncommon.ValididateUnixDirectoryName(eventData.Service) {
		return len(eventData.HelmChart) == 0, errors.New("Service name contains special character(s). " +
			"The service name has to be a valid Unix directory name. For details see " +
			"https://www.cyberciti.biz/faq/linuxunix-rules-for-naming-file-and-directory-names/")
	}

	stageHandler := configutils.NewStageHandler(configServiceURL)
	stages, err := stageHandler.GetAllStages(eventData.Project)
	if err != nil {
		return len(eventData.HelmChart) == 0, fmt.Errorf("Failed to get stages for project %s: %v", eventData.Project, err)
	}

	serviceHandler := configutils.NewServiceHandler(configServiceURL)
	for _, stage := range stages {
		logger.Info("Creating new Keptn service " + eventData.Service + " in stage " + stage.StageName)
		_, err := serviceHandler.CreateServiceInStage(eventData.Project, stage.StageName, eventData.Service)
		if err != nil {
			return len(eventData.HelmChart) == 0, fmt.Errorf("Failed to create service %s in project %s: %s", eventData.Service, eventData.Project, *err.Message)
		}
	}

	return len(eventData.HelmChart) == 0, nil
}

// createProjectAndProcessShipyard creates a project and stages defined in the shipyard
func createProjectAndProcessShipyard(event cloudevents.Event, logger keptncommon.LoggerInterface) error {
	eventData := keptn.ProjectCreateEventData{}
	if err := event.DataAs(&eventData); err != nil {
		return err
	}

	shipyard := keptn.Shipyard{}
	data, err := base64.StdEncoding.DecodeString(eventData.Shipyard)
	if err != nil {
		return fmt.Errorf("Failed to decode shipyard. %s", err.Error())
	}
	err = yaml.Unmarshal(data, &shipyard)
	if err != nil {
		return fmt.Errorf("Failed to unmarshal shipyard. %s", err.Error())
	}
	// create project
	project := configmodels.Project{
		ProjectName:  eventData.Project,
		GitUser:      eventData.GitUser,
		GitToken:     eventData.GitToken,
		GitRemoteURI: eventData.GitRemoteURL,
	}

	if err := createProject(project); err != nil {
		return fmt.Errorf("Failed to create project %s: %s", project.ProjectName, err.Error())
	}
	logger.Info(fmt.Sprintf("Project %s created", project.ProjectName))

	// process shipyard file and create stages
	for _, shipyardStage := range shipyard.Stages {
		if err := createStage(project, shipyardStage.Name); err != nil {
			return fmt.Errorf("Failed to create stage %s: %s", shipyardStage.Name, err.Error())
		}
		logger.Info(fmt.Sprintf("Stage %s created", shipyardStage.Name))
	}
	// store shipyard.yaml
	if err := storeResourceForProject(project.ProjectName, string(data)); err != nil {
		return fmt.Errorf("Failed to store shipyard.yaml for project %s: %s", project.ProjectName, err.Error())
	}

	logger.Info("Project successfully created")
	return nil
}

// deleteProject processes event and deletes project
func deleteProject(event cloudevents.Event, keptnHandler *keptnv2.Keptn) error {
	eventData := keptn.ProjectDeleteEventData{}
	if err := event.DataAs(&eventData); err != nil {
		return err
	}

	project := configmodels.Project{
		ProjectName: eventData.Project,
	}

	keptnHandler.Logger.Info(getDeleteInfoMessage(keptnHandler, eventData.Project))

	// get remote url of project
	projectResp, err := getProject(project)
	if err != nil {
		keptnHandler.Logger.Info(
			fmt.Sprintf("Project %s cannot be retrieved anymore. Any Git upstream of the project will not be deleted.", project.ProjectName))
	} else if projectResp != nil && projectResp.GitRemoteURI != "" {
		keptnHandler.Logger.Info(fmt.Sprintf("The Git upstream of the project will not be deleted: %s", projectResp.GitRemoteURI))
	}

	// delete project
	prjHandler := configutils.NewProjectHandler(configServiceURL)
	_, mErr := prjHandler.DeleteProject(project)
	if mErr != nil {
		return fmt.Errorf("Faild to delete project %s: %s", project.ProjectName, *mErr.Message)
	}

	keptnHandler.Logger.Info("Project successfully deleted")
	return nil
}

// storeResourceForProject stores the resource for a project using the keptn.ResourceHandler
func storeResourceForProject(projectName, shipyard string) error {
	handler := configutils.NewResourceHandler(configServiceURL)
	uri := "shipyard.yaml"
	resource := configmodels.Resource{ResourceURI: &uri, ResourceContent: shipyard}
	if _, err := handler.CreateProjectResources(projectName, []*configmodels.Resource{&resource}); err != nil {
		return fmt.Errorf("Storing %s file failed. %s", *resource.ResourceURI, err.Error())
	}
	return nil
}

// createProject creates a project by using the configuration-service
func createProject(project configmodels.Project) error {
	prjHandler := configutils.NewProjectHandler(configServiceURL)
	if _, errObj := prjHandler.CreateProject(project); errObj != nil {
		return errors.New(*errObj.Message)
	}
	return nil
}

func getDeleteInfoMessage(keptnHandler *keptnv2.Keptn, project string) string {
	res, err := keptnHandler.ResourceHandler.GetProjectResource(project, "shipyard.yaml")
	if err != nil {
		return getShipyardNotAvailableError(project)
	}

	shipyard := &keptn.Shipyard{}
	err = yaml.Unmarshal([]byte(res.ResourceContent), shipyard)
	if err != nil {
		return getShipyardNotAvailableError(project)
	}

	msg := "\n"
	for _, stage := range shipyard.Stages {
		namespace := keptnHandler.Event.GetProject() + "-" + stage.Name
		msg += fmt.Sprintf("- A potentially created namespace %s is not managed by Keptn anymore but is not deleted. "+
			"If you would like to delete this namespace, please execute "+
			"'kubectl delete ns %s'\n", namespace, namespace)
	}
	return strings.TrimSpace(msg)
}

func getShipyardNotAvailableError(project string) string {
	return fmt.Sprintf("Shipyard of project %s cannot be retrieved anymore. "+
		"After deleting the project, the namespaces containing the services are still available. "+
		"This may cause problems if a project with the same name is created later.", project)
}

// getProject returns a project by using the configuration-service
func getProject(project configmodels.Project) (*configmodels.Project, error) {
	prjHandler := configutils.NewProjectHandler(configServiceURL)
	respProject, respError := prjHandler.GetProject(project)
	if respError != nil {
		return nil, errors.New(*respError.Message)
	}
	return respProject, nil
}

// createStage creates a stage by using the configuration-service
func createStage(project configmodels.Project, stage string) error {
	stageHandler := configutils.NewStageHandler(configServiceURL)
	if _, errorObj := stageHandler.CreateStage(project.ProjectName, stage); errorObj != nil {
		return errors.New(*errorObj.Message)
	}
	return nil
}
