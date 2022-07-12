package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"
	apimodels "github.com/keptn/go-utils/pkg/api/models"
	apiutils "github.com/keptn/go-utils/pkg/api/utils"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/cli/internal"
	"github.com/keptn/keptn/cli/pkg/credentialmanager"
	"github.com/keptn/keptn/cli/pkg/docker"
	"github.com/keptn/keptn/cli/pkg/logging"
	"github.com/spf13/cobra"
)

type deliveryStruct struct {
	Project   *string            `json:"project"`
	Service   *string            `json:"service"`
	Stage     *string            `json:"stage"`
	Image     *string            `json:"image"`
	Sequence  *string            `json:"sequence"`
	Values    *[]string          `json:"values"`
	Labels    *map[string]string `json:"labels"`
	Watch     *bool
	WatchTime *int
	Output    *string
}

var delivery deliveryStruct

var triggerDeliveryCmd = &cobra.Command{
	Use:     "delivery",
	Aliases: []string{"delivery"},
	Short:   "Triggers the delivery of a new artifact for a service in a project",
	Long: `Triggers the delivery of a new artifact for a service in a project.
An "artifact" is the name of a Docker image which can be located at any Docker registry (e.g., DockerHub or Quay).
The new artifact is pushed in the first stage specified in the Shipyard of the project. Afterwards, Keptn takes care
of deploying the artifact to the other stages as well.

Note: The value provided in the --image flag has to contain the full qualified image name (incl. docker registry).
The only exception is "docker.io", as this is the default in Kubernetes.
For pulling an image from a private registry, we would like to refer to the Kubernetes documentation (https://kubernetes.io/docs/tasks/configure-pod-container/pull-image-private-registry/).
`,
	Example:      `keptn trigger delivery --project=<project> --service=<service> --image=<image[:tag]> [--sequence=<sequence>]`,
	SilenceUsage: true,
	Args:         cobra.NoArgs,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return checkImageAvailability(*delivery.Image)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return doTriggerDelivery(delivery)
	},
}

func doTriggerDelivery(deliveryInputData deliveryStruct) error {
	var endPoint url.URL
	var apiToken string
	var err error
	if !mocking {
		endPoint, apiToken, err = credentialmanager.NewCredentialManager(assumeYes).GetCreds(namespace)
	} else {
		endPointPtr, _ := url.Parse(os.Getenv("MOCK_SERVER"))
		endPoint = *endPointPtr
		apiToken = ""
	}

	if err != nil {
		return errors.New(authErrorMsg)
	}

	logging.PrintLog("Starting to deliver the service "+
		*deliveryInputData.Service+" in project "+*deliveryInputData.Project+" with image: "+*deliveryInputData.Image, logging.InfoLevel)

	api, err := internal.APIProvider(endPoint.String(), apiToken)
	if err != nil {
		return err
	}

	project, errObj := api.ProjectsV1().GetProject(apimodels.Project{ProjectName: *deliveryInputData.Project})
	if errObj != nil {
		return fmt.Errorf("error while retrieving information for project %v: %s", *deliveryInputData.Project, *errObj.Message)
	}

	// if no stage has been provided to the delivery command, use the first stage in the shipyard.yaml
	if deliveryInputData.Stage == nil || *deliveryInputData.Stage == "" {
		// retrieve the project information to determine the first stage
		if len(project.Stages) > 0 {
			deliveryInputData.Stage = &project.Stages[0].StageName
		} else {
			return fmt.Errorf("could not start sequence because no stage has been found in project %s", *deliveryInputData.Project)
		}
	}

	projectServices, err := api.ServicesV1().GetAllServices(*deliveryInputData.Project, *deliveryInputData.Stage)
	if err != nil {
		return fmt.Errorf("error while retrieving information for service %s: %s", *deliveryInputData.Service, err.Error())
	}
	if !ServiceInSlice(*deliveryInputData.Service, projectServices) {
		return fmt.Errorf("could not start sequence because service %s has not been found in project %s", *deliveryInputData.Service, *deliveryInputData.Project)
	}

	jsonStr, err := internal.JSONPathToJSONObj(*deliveryInputData.Values)
	if err != nil {
		return fmt.Errorf("error while parsing --values flag %v", err)
	}

	valuesJson := map[string]interface{}{}
	valuesJson["image"] = *deliveryInputData.Image
	err = json.Unmarshal([]byte(jsonStr), &valuesJson)
	if err != nil {
		return fmt.Errorf("error unmarshalling json in project %v: %v", *deliveryInputData.Project, err)
	}

	deploymentEvent := keptnv2.DeploymentTriggeredEventData{
		EventData: keptnv2.EventData{
			Project: *deliveryInputData.Project,
			Stage:   *deliveryInputData.Stage,
			Service: *deliveryInputData.Service,
			Labels:  *deliveryInputData.Labels,
		},
		ConfigurationChange: keptnv2.ConfigurationChange{Values: valuesJson},
	}

	sdkEvent := cloudevents.NewEvent()
	sdkEvent.SetID(uuid.New().String())
	sdkEvent.SetType(keptnv2.GetTriggeredEventType(*deliveryInputData.Stage + "." + *deliveryInputData.Sequence))
	sdkEvent.SetSource("https://github.com/keptn/keptn/cli#configuration-change")
	sdkEvent.SetDataContentType(cloudevents.ApplicationJSON)
	if err := sdkEvent.SetData(cloudevents.ApplicationJSON, deploymentEvent); err != nil {
		return fmt.Errorf("failed to create cloud event %s", err.Error())
	}

	eventByte, err := sdkEvent.MarshalJSON()
	if err != nil {
		return fmt.Errorf("failed to marshal cloud event. %s", err.Error())
	}

	apiEvent := apimodels.KeptnContextExtendedCE{}
	err = json.Unmarshal(eventByte, &apiEvent)
	if err != nil {
		return fmt.Errorf("failed to map cloud event to API event model. %v", err)
	}

	logging.PrintLog(fmt.Sprintf("Connecting to server %s", endPoint.String()), logging.VerboseLevel)

	eventContext, err2 := api.APIV1().SendEvent(apiEvent)
	if err2 != nil {
		logging.PrintLog("trigger delivery was unsuccessful", logging.QuietLevel)
		return fmt.Errorf("trigger delivery was unsuccessful. %s", *err2.Message)
	}

	logging.PrintLog("ID of Keptn context: "+*eventContext.KeptnContext, logging.InfoLevel)

	if *deliveryInputData.Watch {
		filter := apiutils.EventFilter{
			KeptnContext: *eventContext.KeptnContext,
			Project:      *deliveryInputData.Project,
		}
		watcher := NewDefaultWatcher(api.EventsV1(), filter, time.Duration(*deliveryInputData.WatchTime)*time.Second)
		PrintEventWatcher(rootCmd.Context(), watcher, *deliveryInputData.Output, os.Stdout)
	}

	return nil
}

func checkImageAvailability(img string) error {
	trimmedImage := strings.TrimSuffix(img, "/")

	image, tag := docker.SplitImageName(trimmedImage)
	return docker.CheckImageAvailability(image, tag, nil)
}

func init() {
	triggerCmd.AddCommand(triggerDeliveryCmd)

	delivery.Project = triggerDeliveryCmd.Flags().StringP("project", "", "",
		"The project containing the service for which the new artifact will be delivered")
	triggerDeliveryCmd.MarkFlagRequired("project")

	delivery.Service = triggerDeliveryCmd.Flags().StringP("service", "", "",
		"The service which for which the new artifact will be delivered")
	triggerDeliveryCmd.MarkFlagRequired("service")

	delivery.Stage = triggerDeliveryCmd.Flags().StringP("stage", "", "",
		"The stage containing the service for which a new artifact will be delivered")
	delivery.Image = triggerDeliveryCmd.Flags().StringP("image", "", "", `The image name, e.g.
"docker.io/<YOUR_ORG>/<YOUR_IMAGE> or quay.io/<YOUR_ORG>/<YOUR_IMAGE>
"Optionally, you can append a tag using ":<YOUR_TAG>. If no tag is provided, "latest" will be used per default`)

	triggerDeliveryCmd.MarkFlagRequired("image")

	delivery.Labels = triggerDeliveryCmd.Flags().StringToStringP("labels", "l", nil, "Additional labels to be included in the event")
	delivery.Sequence = triggerDeliveryCmd.Flags().StringP("sequence", "", "delivery", "The name of the sequence to be triggered")
	delivery.Values = triggerDeliveryCmd.Flags().StringSlice("values", []string{}, "Values to use for the new artifact to be delivered")
	delivery.Output = AddOutputFormatFlag(triggerDeliveryCmd)
	delivery.Watch = AddWatchFlag(triggerDeliveryCmd)
	delivery.WatchTime = AddWatchTimeFlag(triggerDeliveryCmd)
}
