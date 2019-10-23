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
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/types"
	"github.com/google/uuid"
	apimodels "github.com/keptn/go-utils/pkg/api/models"
	apiutils "github.com/keptn/go-utils/pkg/api/utils"
	keptnevents "github.com/keptn/go-utils/pkg/events"
	"github.com/keptn/keptn/cli/pkg/logging"
	"github.com/keptn/keptn/cli/utils/credentialmanager"
	"github.com/keptn/keptn/cli/utils/websockethelper"
	"github.com/spf13/cobra"
)

type newArtifactStruct struct {
	Project *string `json:"project"`
	Service *string `json:"service"`
	Image   *string `json:"image"`
	Tag     *string `json:"tag"`
}

var newArtifact newArtifactStruct

// newArtifactCmd represents the newArtifact command
var newArtifactCmd = &cobra.Command{
	Use: "new-artifact",
	Short: "Sends a new-artifact event to Keptn in order to deploy a new artifact" +
		"for the specified service in the provided project.",
	Long: `Sends a new-artifact event to Keptn in order to deploy a new artifact
for the specified service in the provided project.
Therefore, this command takes the project, service, image, and tag of the new artifact.
	
Example:
	keptn send event new-artifact --project=sockshop --service=carts --image=docker.io/keptnexamples/carts --tag=0.7.0`,
	SilenceUsage: true,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		trimmedImage := strings.TrimSuffix(*newArtifact.Image, "/")
		newArtifact.Image = &trimmedImage
		setTag()
		return checkImageAvailability()
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		endPoint, apiToken, err := credentialmanager.GetCreds()
		if err != nil {
			return errors.New(authErrorMsg)
		}

		logging.PrintLog("Starting to send a new-artifact-event to deploy the service "+
			*newArtifact.Service+" in project "+*newArtifact.Project+" in version "+*newArtifact.Image+":"+*newArtifact.Tag, logging.InfoLevel)

		valuesCanary := make(map[string]interface{})
		valuesCanary["image"] = *newArtifact.Image + ":" + *newArtifact.Tag
		canary := keptnevents.Canary{Action: keptnevents.Set, Value: 100}
		configChangedEvent := keptnevents.ConfigurationChangeEventData{
			Project:      *newArtifact.Project,
			Service:      *newArtifact.Service,
			Stage:        "", // If the stage is empty, the first stage is inserted by the helm-service
			ValuesCanary: valuesCanary,
			Canary:       &canary,
		}

		source, _ := url.Parse("https://github.com/keptn/keptn/cli#configuration-change")
		contentType := "application/json"
		sdkEvent := cloudevents.Event{
			Context: cloudevents.EventContextV02{
				ID:          uuid.New().String(),
				Type:        keptnevents.ConfigurationChangeEventType,
				Source:      types.URLRef{URL: *source},
				ContentType: &contentType,
			}.AsV02(),
			Data: configChangedEvent,
		}

		eventHandler := apiutils.NewAuthenticatedEventHandler(endPoint.String(), apiToken, "x-token", nil, "https")
		logging.PrintLog(fmt.Sprintf("Connecting to server %s", endPoint.String()), logging.VerboseLevel)

		eventByte, err := sdkEvent.MarshalJSON()
		apiEvent := apimodels.Event{}
		json.Unmarshal(eventByte, &apiEvent)

		if !mocking {
			channelInfo, err := eventHandler.SendEvent(apiEvent)
			if err != nil {
				logging.PrintLog("Send new-artifact was unsuccessful", logging.QuietLevel)
				return fmt.Errorf("Send new-artifact was unsuccessful. %s", *err.Message)
			}

			// if ChannelInfo is available, open WebSocket communication
			if channelInfo != nil {
				return websockethelper.PrintWSContentEventContext(channelInfo, endPoint)
			}

			return nil
		}

		fmt.Println("Skipping send new-artifact due to mocking flag set to true")
		return nil
	},
}

func setTag() {

	if newArtifact.Tag != nil && *newArtifact.Tag != "" {
		// The tag is already set
		return
	}

	// Get image name without Docker-organization
	splitsIntoImage := strings.Split(*newArtifact.Image, "/")
	imageName := splitsIntoImage[len(splitsIntoImage)-1]

	splitsIntoTag := strings.Split(imageName, ":")
	if len(splitsIntoTag) == 2 {
		// Tag is provided in the image name
		tag := splitsIntoTag[len(splitsIntoTag)-1]
		newArtifact.Tag = &tag
		imageWithoutTag := strings.TrimSuffix(*newArtifact.Image, ":"+*newArtifact.Tag)
		newArtifact.Image = &imageWithoutTag
		return
	}
	// Otherwise use latest tag
	latest := "latest"
	newArtifact.Tag = &latest
}

func checkImageAvailability() error {

	if strings.HasPrefix(*newArtifact.Image, "docker.io/") {
		resp, err := http.Get("https://index.docker.io/v1/repositories/" +
			strings.TrimPrefix(*newArtifact.Image, "docker.io/") + "/tags/" + *newArtifact.Tag)
		if err != nil {
			return err
		}
		if resp.StatusCode == http.StatusOK {
			return nil
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return errors.New("Provided image not found: " + string(body))
	} else if strings.HasPrefix(*newArtifact.Image, "quay.io/") {
		resp, err := http.Get("https://quay.io/api/v1/repository/" +
			strings.TrimPrefix(*newArtifact.Image, "quay.io/") + "/tag/" + *newArtifact.Tag + "/images")
		if err != nil {
			return err
		}
		if resp.StatusCode == http.StatusOK {
			return nil
		}
		return errors.New("Provided image not found: " + resp.Status)
	}
	logging.PrintLog("Availability of provided image cannot be checked.", logging.InfoLevel)
	return nil
}

func init() {
	sendEventCmd.AddCommand(newArtifactCmd)

	newArtifact.Project = newArtifactCmd.Flags().StringP("project", "", "",
		"The project containing the service which will be new deployed")
	newArtifactCmd.MarkFlagRequired("project")

	newArtifact.Service = newArtifactCmd.Flags().StringP("service", "", "",
		"The service which will be new deployed")
	newArtifactCmd.MarkFlagRequired("service")

	newArtifact.Image = newArtifactCmd.Flags().StringP("image", "", "", "The image name, e.g."+
		"docker.io/YOUR_ORG/YOUR_IMAGE or quay.io/YOUR_ORG/YOUR_IMAGE. "+
		"Optionally, you can directly append the tag using \":YOUR_TAG\"")
	newArtifactCmd.MarkFlagRequired("image")

	newArtifact.Tag = newArtifactCmd.Flags().StringP("tag", "", "", "The tag of the image. "+
		"If no tag is specified, the \"latest\" tag is used.")
}
