package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"
	apimodels "github.com/keptn/go-utils/pkg/api/models"
	apiutils "github.com/keptn/go-utils/pkg/api/utils"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/keptn/keptn/cli/pkg/credentialmanager"
	"github.com/keptn/keptn/cli/pkg/docker"
	"github.com/keptn/keptn/cli/pkg/logging"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"net/url"
	"os"
	"strings"
	"time"
)

type deliveryStruct struct {
	Project   *string            `json:"project"`
	Service   *string            `json:"service"`
	Stage     *string            `json:"stage"`
	Image     *string            `json:"image"`
	Tag       *string            `json:"tag"`
	Sequence  *string            `json:"sequence"`
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
	Example:      `keptn trigger delivery --project=<project> --service=<service> --image=<image> --tag=<tag> [--sequence=<sequence>]`,
	SilenceUsage: true,
	Args: func(cmd *cobra.Command, args []string) error {
		return nil
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return doTriggerDeliveryPreRunCheck(delivery)
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
	*deliveryInputData.Image = strings.Split(*deliveryInputData.Image, ":")[0]

	logging.PrintLog("Starting to deliver the service "+
		*deliveryInputData.Service+" in project "+*deliveryInputData.Project+" in version "+*deliveryInputData.Image+":"+*deliveryInputData.Tag, logging.InfoLevel)

	if endPointErr := checkEndPointStatus(endPoint.String()); endPointErr != nil {
		return fmt.Errorf("Error connecting to server: %s"+endPointErrorReasons,
			endPointErr)
	}

	resourceHandler := apiutils.NewAuthenticatedResourceHandler(endPoint.String(), apiToken, "x-token", nil, endPoint.Scheme)
	shipyardResource, err := resourceHandler.GetProjectResource(*deliveryInputData.Project, "shipyard.yaml")
	if err != nil {
		return fmt.Errorf("Error while retrieving shipyard.yaml for project %s: %s:", *deliveryInputData.Project, err.Error())
	}

	shipyard := &keptnv2.Shipyard{}

	if err := yaml.Unmarshal([]byte(shipyardResource.ResourceContent), shipyard); err != nil {
		return fmt.Errorf("Error while decoding shipyard.yaml for project %s: %s", *deliveryInputData.Project, err.Error())
	}

	// if no stage has been provided to the delivery command, use the first stage in the shipyard.yaml
	if deliveryInputData.Stage == nil || *deliveryInputData.Stage == "" {
		if len(shipyard.Spec.Stages) > 0 {
			deliveryInputData.Stage = &shipyard.Spec.Stages[0].Name
		} else {
			return fmt.Errorf("Could not start sequence because no stage has been found in the shipyard.yaml of project %s", *deliveryInputData.Project)
		}
	}

	deploymentEvent := keptnv2.DeploymentTriggeredEventData{
		EventData: keptnv2.EventData{
			Project: *deliveryInputData.Project,
			Stage:   *deliveryInputData.Stage,
			Service: *deliveryInputData.Service,
			Labels:  *deliveryInputData.Labels,
		},
		ConfigurationChange: keptnv2.ConfigurationChange{
			Values: map[string]interface{}{
				"image": *deliveryInputData.Image + ":" + *deliveryInputData.Tag,
			},
		},
	}

	source, _ := url.Parse("https://github.com/keptn/keptn/cli#configuration-change")

	sdkEvent := cloudevents.NewEvent()
	sdkEvent.SetID(uuid.New().String())
	sdkEvent.SetType(keptnv2.GetTriggeredEventType(*deliveryInputData.Stage + "." + *deliveryInputData.Sequence))
	sdkEvent.SetSource(source.String())
	sdkEvent.SetDataContentType(cloudevents.ApplicationJSON)
	sdkEvent.SetData(cloudevents.ApplicationJSON, deploymentEvent)

	eventByte, err := sdkEvent.MarshalJSON()
	if err != nil {
		return fmt.Errorf("Failed to marshal cloud event. %s", err.Error())
	}

	apiEvent := apimodels.KeptnContextExtendedCE{}
	err = json.Unmarshal(eventByte, &apiEvent)
	if err != nil {
		return fmt.Errorf("Failed to map cloud event to API event model. %s", err.Error())
	}

	apiHandler := apiutils.NewAuthenticatedAPIHandler(endPoint.String(), apiToken, "x-token", nil, endPoint.Scheme)

	logging.PrintLog(fmt.Sprintf("Connecting to server %s", endPoint.String()), logging.VerboseLevel)

	eventContext, err2 := apiHandler.SendEvent(apiEvent)
	if err2 != nil {
		logging.PrintLog("trigger delivery was unsuccessful", logging.QuietLevel)
		return fmt.Errorf("trigger delivery was unsuccessful. %s", *err2.Message)
	}

	if *deliveryInputData.Watch {
		eventHandler := apiutils.NewAuthenticatedEventHandler(endPoint.String(), apiToken, "x-token", nil, endPoint.Scheme)
		filter := apiutils.EventFilter{
			KeptnContext: *eventContext.KeptnContext,
			Project:      *deliveryInputData.Project,
		}
		watcher := NewDefaultWatcher(eventHandler, filter, time.Duration(*deliveryInputData.WatchTime)*time.Second)
		PrintEventWatcher(rootCmd.Context(), watcher, *deliveryInputData.Output, os.Stdout)
	}
	return nil
}

func doTriggerDeliveryPreRunCheck(deliveryInputData deliveryStruct) error {
	trimmedImage := strings.TrimSuffix(*deliveryInputData.Image, "/")
	deliveryInputData.Image = &trimmedImage

	if deliveryInputData.Tag == nil || *deliveryInputData.Tag == "" {
		*deliveryInputData.Image, *deliveryInputData.Tag = docker.SplitImageName(*deliveryInputData.Image)
	}
	return docker.CheckImageAvailability(*deliveryInputData.Image, *deliveryInputData.Tag, nil)
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
"Optionally, you can append a tag using ":<YOUR_TAG>`)

	triggerDeliveryCmd.MarkFlagRequired("image")

	delivery.Labels = triggerDeliveryCmd.Flags().StringToStringP("labels", "l", nil, "Additional labels to be included in the event")

	delivery.Tag = triggerDeliveryCmd.Flags().StringP("tag", "", "", `The tag of the image. If no tag is specified, the "latest" tag is used`)

	delivery.Sequence = triggerDeliveryCmd.Flags().StringP("sequence", "", "delivery", "The name of the sequence to be triggered")

	delivery.Output = AddOutputFormatFlag(triggerDeliveryCmd)
	delivery.Watch = AddWatchFlag(triggerDeliveryCmd)
	delivery.WatchTime = AddWatchTimeFlag(triggerDeliveryCmd)

}
