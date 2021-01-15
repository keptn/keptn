// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ghodss/yaml"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/keptn/keptn/cli/pkg/docker"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"

	apimodels "github.com/keptn/go-utils/pkg/api/models"
	apiutils "github.com/keptn/go-utils/pkg/api/utils"
	"github.com/keptn/keptn/cli/pkg/credentialmanager"
	"github.com/keptn/keptn/cli/pkg/logging"
	"github.com/spf13/cobra"
)

type newArtifactStruct struct {
	Project   *string `json:"project"`
	Service   *string `json:"service"`
	Stage     *string `json:"stage"`
	Image     *string `json:"image"`
	Tag       *string `json:"tag"`
	Sequence  *string `json:"sequence"`
	Watch     *bool
	WatchTime *int
	Output    *string
}

var newArtifact newArtifactStruct

// newArtifactCmd represents the newArtifact command
var newArtifactCmd = &cobra.Command{
	Use: "new-artifact",
	Short: "Sends a new-artifact event to Keptn in order to deploy a new artifact " +
		"for the specified service in the provided project",
	Long: `Sends a new-artifact event to Keptn in order to deploy a new artifact for the specified service in the provided project.
Therefore, this command takes the project, service, image, and tag of the new artifact.

* The artifact is the name of a image, which can be located at DockerHub, Quay, or any other registry storing docker images. 
* The new artifact is pushed in the first stage specified in the Shipyard of the project. Afterwards, Keptn takes care of deploying this new artifact to the other stages.

**Notes:**
* The value provided in the *image* flag has to contain the full path to your Docker registry. The only exception is *docker.io* because this is the default in Kubernetes and, hence, can be omitted.
* This command does not send the actual Docker image to Keptn, just the image name and tag. Instead, Keptn uses Kubernetes functionalities for pulling this image.
For pulling an image from a private registry, we would like to refer to the Kubernetes documentation (https://kubernetes.io/docs/tasks/configure-pod-container/pull-image-private-registry/).
`,
	Example:      `keptn send event new-artifact --project=sockshop --service=carts --stage=dev --image=docker.io/keptnexamples/carts --tag=0.7.0 --sequence=artifact-delivery`,
	SilenceUsage: true,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return doSendEventNewArtifactPreRunCheck()
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := doSendEventNewArtifactPreRunCheck(); err != nil {
			return err
		}
		var endPoint url.URL
		var apiToken string
		var err error
		if !mocking {
			endPoint, apiToken, err = credentialmanager.NewCredentialManager(false).GetCreds(namespace)
		} else {
			endPointPtr, _ := url.Parse(os.Getenv("MOCK_SERVER"))
			endPoint = *endPointPtr
			apiToken = ""
		}

		if err != nil {
			return errors.New(authErrorMsg)
		}

		logging.PrintLog("Starting to send a new-artifact-event to deploy the service "+
			*newArtifact.Service+" in project "+*newArtifact.Project+" in version "+*newArtifact.Image+":"+*newArtifact.Tag, logging.InfoLevel)

		if endPointErr := checkEndPointStatus(endPoint.String()); endPointErr != nil {
			return fmt.Errorf("Error connecting to server: %s"+endPointErrorReasons,
				endPointErr)
		}

		resourceHandler := apiutils.NewAuthenticatedResourceHandler(endPoint.String(), apiToken, "x-token", nil, endPoint.Scheme)
		shipyardResource, err := resourceHandler.GetProjectResource(*newArtifact.Project, "shipyard.yaml")
		if err != nil {
			return fmt.Errorf("Error while retrieving shipyard.yaml for project %s: %s:", *newArtifact.Project, err.Error())
		}

		shipyard := &keptnv2.Shipyard{}

		if err := yaml.Unmarshal([]byte(shipyardResource.ResourceContent), shipyard); err != nil {
			return fmt.Errorf("Error while decoding shipyard.yaml for project %s: %s", *newArtifact.Project, err.Error())
		}

		// if no stage has been provided to the new-artifact command, use the first stage in the shipyard.yaml
		if newArtifact.Stage == nil || *newArtifact.Stage == "" {
			if len(shipyard.Spec.Stages) > 0 {
				newArtifact.Stage = &shipyard.Spec.Stages[0].Name
			} else {
				return fmt.Errorf("Could not start sequence because no stage has been found in the shipyard.yaml of project %s", *newArtifact.Project)
			}
		}

		deploymentEvent := keptnv2.DeploymentTriggeredEventData{
			EventData: keptnv2.EventData{
				Project: *newArtifact.Project,
				Stage:   *newArtifact.Stage,
				Service: *newArtifact.Service,
			},
			ConfigurationChange: keptnv2.ConfigurationChange{
				Values: map[string]interface{}{
					"image": *newArtifact.Image + ":" + *newArtifact.Tag,
				},
			},
		}

		source, _ := url.Parse("https://github.com/keptn/keptn/cli#configuration-change")

		sdkEvent := cloudevents.NewEvent()
		sdkEvent.SetID(uuid.New().String())
		sdkEvent.SetType(keptnv2.GetTriggeredEventType(*newArtifact.Stage + "." + *newArtifact.Sequence))
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
			logging.PrintLog("Send new-artifact was unsuccessful", logging.QuietLevel)
			return fmt.Errorf("Send new-artifact was unsuccessful. %s", *err2.Message)
		}

		if *newArtifact.Watch {
			eventHandler := apiutils.NewAuthenticatedEventHandler(endPoint.String(), apiToken, "x-token", nil, endPoint.Scheme)
			filter := apiutils.EventFilter{
				KeptnContext: *eventContext.KeptnContext,
				Project:      *newArtifact.Project,
			}
			watcher := NewDefaultWatcher(eventHandler, filter, time.Duration(*newArtifact.WatchTime)*time.Second)
			PrintEventWatcher(watcher, *newArtifact.Output, os.Stdout)
		}
		return nil
	},
}

func doSendEventNewArtifactPreRunCheck() error {
	trimmedImage := strings.TrimSuffix(*newArtifact.Image, "/")
	newArtifact.Image = &trimmedImage

	if newArtifact.Tag == nil || *newArtifact.Tag == "" {
		*newArtifact.Image, *newArtifact.Tag = docker.SplitImageName(*newArtifact.Image)
	}
	return docker.CheckImageAvailability(*newArtifact.Image, *newArtifact.Tag, nil)
}

func init() {
	sendEventCmd.AddCommand(newArtifactCmd)

	newArtifact.Project = newArtifactCmd.Flags().StringP("project", "", "",
		"The project containing the service which will be new deployed")
	newArtifactCmd.MarkFlagRequired("project")

	newArtifact.Service = newArtifactCmd.Flags().StringP("service", "", "",
		"The service which will be new deployed")
	newArtifactCmd.MarkFlagRequired("service")

	newArtifact.Stage = newArtifactCmd.Flags().StringP("stage", "", "",
		"The stage containing the service to be deployed")

	newArtifact.Image = newArtifactCmd.Flags().StringP("image", "", "", "The image name, e.g."+
		"docker.io/YOUR_ORG/YOUR_IMAGE or quay.io/YOUR_ORG/YOUR_IMAGE. "+
		"Optionally, you can directly append the tag using \":YOUR_TAG\"")
	newArtifactCmd.MarkFlagRequired("image")

	newArtifact.Tag = newArtifactCmd.Flags().StringP("tag", "", "", "The tag of the image. "+
		"If no tag is specified, the \"latest\" tag is used.")

	newArtifact.Sequence = newArtifactCmd.Flags().StringP("sequence", "", "", "The name of the sequence to be triggered")
	newArtifactCmd.MarkFlagRequired("sequence")

	newArtifact.Output = AddOutputFormatFlag(newArtifactCmd)
	newArtifact.Watch = AddWatchFlag(newArtifactCmd)
	newArtifact.WatchTime = AddWatchTimeFlag(newArtifactCmd)

}
