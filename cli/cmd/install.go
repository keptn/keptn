//go:build !nokubectl
// +build !nokubectl

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
	"fmt"
	"github.com/keptn/keptn/cli/pkg/platform"
	"github.com/spf13/cobra"
)

type installCmdParams struct {
	installUpgradeParams
	UseCaseInput             *string
	UseCase                  usecase
	EndPointServiceTypeInput *string
	EndPointServiceType      endpointServiceType
	HideSensitiveData        *bool
}

var installParams installCmdParams

// installCmd represents the version command
var installCmd = NewInstallCmd()

func NewInstallCmd() *cobra.Command {
	return &cobra.Command{
		Deprecated:   fmt.Sprintf(MsgDeprecatedUseHelm, Version),
		Use:          "install",
		Args:         cobra.NoArgs,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("this command is deprecated! "+MsgDeprecatedUseHelm, Version)
		},
	}
}

func init() {
	rootCmd.AddCommand(installCmd)
	installParams = installCmdParams{}

	installParams.PlatformIdentifier = installCmd.Flags().StringP("platform", "p", "kubernetes",
		"The platform to run Keptn on ["+platform.KubernetesIdentifier+","+platform.OpenShiftIdentifier+"]")

	installParams.ConfigFilePath = installCmd.Flags().StringP("creds", "c", "",
		"Specify a JSON file containing cluster information needed for the installation. This allows skipping user prompts to execute a *silent* Keptn installation.")

	installParams.UseCaseInput = installCmd.Flags().StringP("use-case", "u", "",
		"Use --use-case=continuous-delivery to install the execution plane for continuous delivery. Without this flag, your Keptn is capable of the quality gate and automated remediations use-case.")

	installParams.EndPointServiceTypeInput = installCmd.Flags().StringP("endpoint-service-type", "",
		ClusterIP.String(), "Installation options for the endpoint-service type ["+ClusterIP.String()+","+
			LoadBalancer.String()+","+NodePort.String()+"]")

	installParams.ChartRepoURL = installCmd.Flags().StringP("chart-repo", "",
		"", "URL of the Keptn Helm Chart repository")

	installCmd.PersistentFlags().BoolVarP(&insecureSkipTLSVerify, "insecure-skip-tls-verify", "s",
		false, "Skip tls verification for kubectl commands")
	installParams.HideSensitiveData = installCmd.Flags().BoolP("hide-sensitive-data", "", false,
		"Hide the sensitive data like api-tokens and endpoints in post-installation output.")

	installCmd.Flags().MarkHidden("platform")
	installCmd.Flags().MarkHidden("creds")
	installCmd.Flags().MarkHidden("use-case")
	installCmd.Flags().MarkHidden("endpoint-service-type")
	installCmd.Flags().MarkHidden("chart-repo")
	installCmd.Flags().MarkHidden("hide-sensitive-data")
	installCmd.PersistentFlags().MarkHidden("insecure-skip-tls-verify")

}
