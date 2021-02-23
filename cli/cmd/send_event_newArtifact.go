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

// NOTE: THIS COMMAND WILL BE REMOVED, THUS THE WHOLE FILE WILL BE REMOVED IN A FUTURE RELEASE

import (
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
	Deprecated:   `Use "keptn trigger delivery" instead`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		delivery := deliveryStruct{
			Project:   newArtifact.Project,
			Service:   newArtifact.Service,
			Stage:     newArtifact.Stage,
			Image:     newArtifact.Image,
			Tag:       newArtifact.Tag,
			Sequence:  newArtifact.Sequence,
			Watch:     newArtifact.Watch,
			WatchTime: newArtifact.WatchTime,
			Output:    newArtifact.Output,
		}
		return doTriggerDeliveryPreRunCheck(delivery)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		delivery := deliveryStruct{
			Project:   newArtifact.Project,
			Service:   newArtifact.Service,
			Stage:     newArtifact.Stage,
			Image:     newArtifact.Image,
			Tag:       newArtifact.Tag,
			Sequence:  newArtifact.Sequence,
			Watch:     newArtifact.Watch,
			WatchTime: newArtifact.WatchTime,
			Output:    newArtifact.Output,
		}
		return doTriggerDelivery(delivery)
	},
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
