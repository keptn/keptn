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
	Stage   *string `json:"stage"`
	Image   *string `json:"image"`
	Tag     *string `json:"tag"`
}

var newArtifact newArtifactStruct

// newArtifactCmd represents the newArtifact command
var newArtifactCmd = &cobra.Command{
	Use: "new-artifact",
	Short: "Sends a new-artifact-event to the keptn installation in order to deploy a new artifact" +
		"for the specified service in the provided project and stage.",
	Long: `Sends a new-artifact-event to the keptn installation in order to deploy a new artifact
for the specified service in the provided project and stage.
Therefore, this command takes the project containing the service, the name of the service, the stage into which the new artifact is deployed
as well as the image and tag of the new artifact.
	
Example:
	keptn new-artifact --project=sockshop --service=carts --stage=dev --image=docker.io/keptnexamples/carts --tag=0.7.0`,
	RunE: func(cmd *cobra.Command, args []string) error {
		controlEndPoint, apiToken, err := credentialmanager.GetCreds()
		if err != nil {
			return errors.New(authErrorMsg)
		}

		utils.PrintLog("Starting to send a new-artifact-event to deploy the service "+
			*newArtifact.Service+" in project "+*newArtifact.Project+" and stage "+
			*newArtifact.Stage+" in version "+*newArtifact.Image+":"+*newArtifact.Tag, utils.InfoLevel)

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

		eventBrokerHostName := controlEndPoint
		eventBrokerHostName.Host = strings.Replace(eventBrokerHostName.Host, "control", "event-broker-ext", -1)
		eventBrokerHostName.Path = "event"

		utils.PrintLog(fmt.Sprintf("Connecting to server %s", eventBrokerHostName.String()), utils.VerboseLevel)
		if !mocking {
			responseCE, err := utils.Send(eventBrokerHostName, event, apiToken)
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
				return websockethelper.PrintWSContent(responseCE)
			}
		} else {
			fmt.Println("Skipping send-new artifact due to mocking flag set to true")
		}
		return nil
	},
}

func init() {
	sendCmd.AddCommand(newArtifactCmd)

	newArtifact.Project = newArtifactCmd.Flags().StringP("project", "", "", "The project containing the service which will be new deployed")
	newArtifactCmd.MarkFlagRequired("project")

	newArtifact.Service = newArtifactCmd.Flags().StringP("service", "", "", "The service which will be new deployed")
	newArtifactCmd.MarkFlagRequired("service")

	newArtifact.Stage = newArtifactCmd.Flags().StringP("stage", "", "", "The stage into which the new artifact will be deployed")
	newArtifactCmd.MarkFlagRequired("stage")

	newArtifact.Image = newArtifactCmd.Flags().StringP("image", "", "", "The image name, e.g."+
		"docker.io/YOUR_ORG/YOUR_IMAGE or quay.io/YOUR_ORG/YOUR_IMAGE")
	newArtifactCmd.MarkFlagRequired("image")

	newArtifact.Tag = newArtifactCmd.Flags().StringP("tag", "", "", "The tag of the image name")
	newArtifactCmd.MarkFlagRequired("tag")
}
