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
	"github.com/spf13/cobra"
)

type newArtifactStruct struct {
	Service     *string `json:"service"`
	Stage       *string `json:"stage"`
	DockerImage *string `json:"dockerimage"`
}

var newArtifact newArtifactStruct

// newArtifactCmd represents the newArtifact command
var newArtifactCmd = &cobra.Command{
	Use:   "new-artifact",
	Short: "Sends a new artifact event to deploy it into the specified stage.",
	Long: `Sends a new artifact event to deploy it into the specified stage.
Therefore, this command takes the name of the service, the stage into which the new artifact is deployed
and the location of the new artifact (currently this is a Docker image on Docker Hub).
Going forward we generalize the source supplier of the artifact.
Internally this command sends a cloud event to the external event-broker.
	
Example:
	keptn new-artifact --service=yourService --stage=dev --docker-image=yourOrg/yourService:latest`,
	RunE: func(cmd *cobra.Command, args []string) error {
		controlEndPoint, apiToken, err := credentialmanager.GetCreds()
		if err != nil {
			return errors.New(authErrorMsg)
		}

		fmt.Println("Starting to send a new articat event to deploy the new artifact ", newArtifact.Service, " into ", newArtifact.Stage)

		source, _ := url.Parse("https://github.com/keptn/keptn/cli#new-artifact")
		contentType := "application/json"
		event := cloudevents.Event{
			Context: cloudevents.EventContextV02{
				ID:          uuid.New().String(),
				Type:        "new-artifact",
				Source:      types.URLRef{URL: *source},
				ContentType: &contentType,
			}.AsV02(),
			Data: newArtifact,
		}

		eventBrokerHostName := controlEndPoint
		eventBrokerHostName.Host = strings.Replace(eventBrokerHostName.Host, "control", "event-broker-ext", -1)
		eventBrokerHostName.Path = "newartifact"

		_, err = utils.Send(eventBrokerHostName, event, apiToken, utils.AddAuthorizationHeader)
		if err != nil {
			fmt.Println("Send new artifact was unsuccessful")
			return err
		}

		fmt.Println("Successfully sent new artifact")
		return nil
	},
}

func init() {
	eventCmd.AddCommand(newArtifactCmd)

	newArtifact.Service = newArtifactCmd.Flags().StringP("service", "", "", "The service which should be new deployed")
	newArtifactCmd.MarkFlagRequired("service")
	newArtifact.Stage = newArtifactCmd.Flags().StringP("stage", "", "", "The stage into which the new artifact is deployed")
	newArtifactCmd.MarkFlagRequired("stage")
	newArtifact.DockerImage = newArtifactCmd.Flags().StringP("docker-image", "d", "", "The fully qualified Docker image (organization, name and tag)")
	newArtifactCmd.MarkFlagRequired("docker-image")
}
