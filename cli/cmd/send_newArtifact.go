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
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/types"
	"github.com/google/uuid"
	"github.com/keptn/keptn/cli/utils"
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
	Short: "Sends a new-artifact-event to the keptn installation in order to deploy a new artifact" +
		"for the specified service in the provided project.",
	Long: `Sends a new-artifact-event to the keptn installation in order to deploy a new artifact
for the specified service in the provided project.
Therefore, this command takes the project, the service as well as the image and tag of the new artifact.
	
Example:
	keptn new-artifact --project=sockshop --service=carts --image=docker.io/keptnexamples/carts --tag=0.7.0`,
	SilenceUsage: true,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return checkImageAvailability()
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		endPoint, apiToken, err := credentialmanager.GetCreds()
		if err != nil {
			return errors.New(authErrorMsg)
		}

		utils.PrintLog("Starting to send a new-artifact-event to deploy the service "+
			*newArtifact.Service+" in project "+*newArtifact.Project+" in version "+*newArtifact.Image+":"+*newArtifact.Tag, utils.InfoLevel)

		source, _ := url.Parse("https://github.com/keptn/keptn/cli#new-artifact")
		contentType := "application/json"
		event := cloudevents.Event{
			Context: cloudevents.EventContextV02{
				ID:          uuid.New().String(),
				Type:        "sh.keptn.events.new-artefact",
				Source:      types.URLRef{URL: *source},
				ContentType: &contentType,
			}.AsV02(),
			Data: newArtifact,
		}

		eventURL := endPoint
		eventURL.Path = "event"

		utils.PrintLog(fmt.Sprintf("Connecting to server %s", eventURL.String()), utils.VerboseLevel)
		if !mocking {
			responseCE, err := utils.Send(eventURL, event, apiToken)
			if err != nil {
				utils.PrintLog("Send new-artifact was unsuccessful", utils.QuietLevel)
				return err
			}

			// check for responseCE to include token
			if responseCE == nil {
				utils.PrintLog("Response CE is nil", utils.QuietLevel)

				return nil
			}
			if responseCE.Data != nil {
				return websockethelper.PrintWSContentCEResponse(responseCE)
			}
		} else {
			fmt.Println("Skipping send-new artifact due to mocking flag set to true")
		}
		return nil
	},
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
	utils.PrintLog("Availability of provided image cannot be checked.", utils.InfoLevel)
	return nil
}

func init() {
	sendCmd.AddCommand(newArtifactCmd)

	newArtifact.Project = newArtifactCmd.Flags().StringP("project", "", "", "The project containing the service which will be new deployed")
	newArtifactCmd.MarkFlagRequired("project")

	newArtifact.Service = newArtifactCmd.Flags().StringP("service", "", "", "The service which will be new deployed")
	newArtifactCmd.MarkFlagRequired("service")

	newArtifact.Image = newArtifactCmd.Flags().StringP("image", "", "", "The image name, e.g."+
		"docker.io/YOUR_ORG/YOUR_IMAGE or quay.io/YOUR_ORG/YOUR_IMAGE")
	newArtifactCmd.MarkFlagRequired("image")

	newArtifact.Tag = newArtifactCmd.Flags().StringP("tag", "", "", "The tag of the image name")
	newArtifactCmd.MarkFlagRequired("tag")
}
