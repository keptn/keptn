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
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/keptn/keptn/cli/pkg/helm"
	"github.com/keptn/keptn/cli/pkg/platform"

	"github.com/keptn/keptn/cli/pkg/version"

	"github.com/keptn/keptn/cli/pkg/kube"
	"helm.sh/helm/v3/pkg/chart"

	"github.com/keptn/keptn/cli/pkg/logging"
	"github.com/spf13/cobra"
)

var upgradeParams installUpgradeParams
var keptnUpgradeChart *chart.Chart

// installCmd represents the version command
var upgraderCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrades Keptn on a Kubernetes cluster",
	Long: `The Keptn CLI allows upgrading Keptn on any Kubernetes derivate to which your kube config is pointing to, and on OpenShift.

For more information, please follow the installation guide [Upgrade Keptn](https://keptn.sh/docs/` + keptnReleaseDocsURL + `/operate/upgrade/)
`,
	Example: `keptn upgrade                                                        # upgrades Keptn

keptn upgrade --platform=openshift # upgrades Keptn on Openshift

keptn upgrade --platform=kubernetes # upgrades Keptn on the Kubernetes cluster
`,
	SilenceUsage: true,
	PreRunE: func(cmd *cobra.Command, args []string) error {

		var chartRepoURL string
		// Determine installer version
		if installParams.ChartRepoURL != nil && *installParams.ChartRepoURL != "" {
			chartRepoURL = *installParams.ChartRepoURL
		} else if version.IsOfficialKeptnVersion(Version) {
			version, _ := version.GetOfficialKeptnVersion(Version)
			chartRepoURL = keptnInstallerHelmRepoURL + "keptn-" + version + ".tgz"
		} else {
			chartRepoURL = keptnInstallerHelmRepoURL + "latest/keptn-0.1.0.tgz"
		}

		var err error
		if keptnChart, err = helm.NewHelmHelper().DownloadChart(chartRepoURL); err != nil {
			return err
		}

		logging.PrintLog(fmt.Sprintf("Helm Chart used for Keptn installation: %s", chartRepoURL), logging.InfoLevel)

		installPlatformManager, err := platform.NewPlatformManager(*installParams.PlatformIdentifier)
		if err != nil {
			return err
		}

		if !mocking {
			if err := installPlatformManager.CheckRequirements(); err != nil {
				return err
			}
		}

		if installParams.ConfigFilePath != nil && *installParams.ConfigFilePath != "" {
			// Config was provided in form of a file
			if err := installPlatformManager.ParseConfig(*installParams.ConfigFilePath); err != nil {
				return err
			}

			// Check whether the authentication at the cluster is valid
			if err := installPlatformManager.CheckCreds(); err != nil {
				return err
			}
		} else {
			err = installPlatformManager.ReadCreds()
			if err != nil {
				return err
			}
		}

		// check if Kubernetes server version is compatible (except OpenShift)
		if *installParams.PlatformIdentifier != platform.OpenShiftIdentifier {
			if err := kube.CheckKubeServerVersion(KubeServerVersionConstraints); err != nil {
				logging.PrintLog(err.Error(), logging.VerboseLevel)
				logging.PrintLog("See https://keptn.sh/docs/"+keptnReleaseDocsURL+"/operate/k8s_support/ for details.", logging.VerboseLevel)
				return fmt.Errorf("Failed to check kubernetes server version: %w", err)
			}
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {

		logging.PrintLog("Upgrading Keptn ...", logging.InfoLevel)

		if !mocking {
			return doUpgrade()
		}
		fmt.Println("Skipping upgrade due to mocking flag")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(upgraderCmd)
	upgradeParams = installUpgradeParams{}

	upgradeParams.PlatformIdentifier = upgraderCmd.Flags().StringP("platform", "p", "kubernetes",
		"The platform to run Keptn on ["+platform.KubernetesIdentifier+","+platform.OpenShiftIdentifier+"]")

	upgradeParams.ChartRepoURL = upgraderCmd.Flags().StringP("chart-repo", "",
		"", "URL of the Keptn Helm Chart repository")
	upgraderCmd.Flags().MarkHidden("chart-repo")

	upgradeParams.Namespace = upgraderCmd.Flags().StringP("namespace", "n", "keptn",
		"Specify the namespace Keptn should be installed in (default keptn).")
}

// Preconditions: 1. Already authenticated against the cluster.
func doUpgrade() error {
	const keptnReleaseName = "keptn"
	keptnNamespace := *installParams.Namespace
	helper := helm.NewHelmHelper()
	releases, err := helper.GetHistory(keptnReleaseName, keptnNamespace)
	if err != nil {
		return fmt.Errorf("Failed to check if Keptn release is available in namespace %s: %v", keptnNamespace, err)
	}
	if len(releases) == 0 {
		return fmt.Errorf("No Keptn release found in namespace %s: %v", keptnNamespace, err)
	}

	lastRelease := releases[len(releases)-1]

	fmt.Printf("Existing Keptn installation found in namespace %s\n with version %v", keptnNamespace, lastRelease.Chart.Metadata.AppVersion)
	fmt.Println()
	fmt.Println("Do you want to upgrade this installation? (y/n)")

	reader := bufio.NewReader(os.Stdin)
	in, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	in = strings.ToLower(strings.TrimSpace(in))
	if !(in == "y" || in == "yes") {
		return fmt.Errorf("Stopping upgrade.")
	}

	if err := helper.UpgradeChart(keptnUpgradeChart, keptnReleaseName, keptnNamespace, nil); err != nil {
		logging.PrintLog("Could not complete Keptn upgrade: "+err.Error(), logging.InfoLevel)
		return err
	}

	logging.PrintLog("Keptn has been successfully upgraded on your cluster.", logging.InfoLevel)
	return nil
}
