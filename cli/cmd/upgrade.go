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

	"helm.sh/helm/v3/pkg/release"

	"github.com/keptn/keptn/cli/pkg/version"

	"github.com/keptn/keptn/cli/pkg/helm"
	"github.com/keptn/keptn/cli/pkg/platform"

	"github.com/keptn/keptn/cli/pkg/kube"
	keptnutils "github.com/keptn/kubernetes-utils/pkg"
	"helm.sh/helm/v3/pkg/chart"

	"github.com/keptn/keptn/cli/pkg/logging"
	"github.com/spf13/cobra"
)

const keptnReleaseName = "keptn"

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
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := doUpgradePreRunCheck(); err != nil {
			return err
		}
		if !mocking {
			return doUpgrade()
		}
		fmt.Println("Skipping upgrade due to mocking flag")
		return nil
	},
}

func doUpgradePreRunCheck() error {
		if *upgradeParams.PatchNamespace {
			return nil
		}

		chartRepoURL := getChartRepoURL(upgradeParams.ChartRepoURL)

	var err error
	if keptnUpgradeChart, err = helm.NewHelper().DownloadChart(chartRepoURL); err != nil {
		return err
	}

	res, err := isUpgradeCompatible()
	if err != nil {
		return err
	}
	if !res {
		installedKeptnVerison, err := getInstalledKeptnVersion()
		if err != nil {
			return err
		}
		return fmt.Errorf("No upgrade path exists from Keptn version %s to %s",
			installedKeptnVerison, getAppVersion(keptnUpgradeChart))
	}

	logging.PrintLog(fmt.Sprintf("Helm Chart used for Keptn upgrade: %s", chartRepoURL), logging.InfoLevel)

	platformManager, err := platform.NewPlatformManager(*upgradeParams.PlatformIdentifier)
	if err != nil {
		return err
	}

	if !mocking {
		if err := platformManager.CheckRequirements(); err != nil {
			return err
		}
	}

	if upgradeParams.ConfigFilePath != nil && *upgradeParams.ConfigFilePath != "" {
		// Config was provided in form of a file
		if err := platformManager.ParseConfig(*upgradeParams.ConfigFilePath); err != nil {
			return err
		}

		// Check whether the authentication at the cluster is valid
		if err := platformManager.CheckCreds(); err != nil {
			return err
		}
	} else {
		err = platformManager.ReadCreds()
		if err != nil {
			return err
		}
	}

	// check if Kubernetes server version is compatible (except OpenShift)
	if *upgradeParams.PlatformIdentifier != platform.OpenShiftIdentifier {
		if err := kube.CheckKubeServerVersion(KubeServerVersionConstraints); err != nil {
			logging.PrintLog(err.Error(), logging.VerboseLevel)
			logging.PrintLog("See https://keptn.sh/docs/"+keptnReleaseDocsURL+"/operate/k8s_support/ for details.", logging.VerboseLevel)
			return fmt.Errorf("Failed to check kubernetes server version: %w", err)
		}
	}

	return nil
		return nil
	},
}

func getInstalledKeptnVersion() (string, error) {
	lastRelease, err := getLatestKeptnRelease()
	if err != nil {
		return "", err
	}
	return lastRelease.Chart.Metadata.AppVersion, nil
}

func getAppVersion(ch *chart.Chart) string {
	return ch.Metadata.AppVersion
}

func isUpgradeCompatible() (bool, error) {
	installedVersion, err := getInstalledKeptnVersion()
	if err != nil {
		return false, err
	}
	versionChecker := version.NewKeptnVersionChecker()
	return versionChecker.IsUpgradable(Version, installedVersion, getAppVersion(keptnUpgradeChart))
}

func getLatestKeptnRelease() (*release.Release, error) {
	keptnNamespace := namespace
	releases, err := helm.NewHelper().GetHistory(keptnReleaseName, keptnNamespace)
	if err != nil {
		return nil, fmt.Errorf("Failed to check if Keptn release is available in namespace %s: %v", keptnNamespace, err)
	}
	if len(releases) == 0 {
		return nil, fmt.Errorf("No Keptn release found in namespace %s: %v", keptnNamespace, err)
	}

	return releases[len(releases)-1], nil
}

func init() {
	rootCmd.AddCommand(upgraderCmd)
	upgradeParams = installUpgradeParams{}

	upgradeParams.PlatformIdentifier = upgraderCmd.Flags().StringP("platform", "p", "kubernetes",
		"The platform to run Keptn on ["+platform.KubernetesIdentifier+","+platform.OpenShiftIdentifier+"]")

	upgradeParams.ChartRepoURL = upgraderCmd.Flags().StringP("chart-repo", "",
		"", "URL of the Keptn Helm Chart repository")
	upgraderCmd.Flags().MarkHidden("chart-repo")
	upgradeParams.PatchNamespace = upgraderCmd.Flags().BoolP("patch-namespace", "", false, "Patch the namespace with the annotation & label 'keptn.sh/managed-by: keptn'")
}

func doUpgrade() error {
	keptnNamespace := namespace

	installedKeptnVersion, err := getInstalledKeptnVersion()
	if err != nil {
		return err
	}

	fmt.Printf("Do you want to upgrade Keptn version %s to %s? (y/n)\n", installedKeptnVersion, getAppVersion(keptnUpgradeChart))

	reader := bufio.NewReader(os.Stdin)
	in, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	in = strings.ToLower(strings.TrimSpace(in))
	if !(in == "y" || in == "yes") {
		return fmt.Errorf("Stopping upgrade.")
	}

	if err := helm.NewHelper().UpgradeChart(keptnUpgradeChart, keptnReleaseName, keptnNamespace, nil); err != nil {
		logging.PrintLog("Could not complete Keptn upgrade: "+err.Error(), logging.InfoLevel)
		return err
	}

	logging.PrintLog("Keptn has been successfully upgraded on your cluster.", logging.InfoLevel)
	return nil
}

func patchNamespace() error {
	err := keptnutils.PatchKeptnManagedNamespace(false, namespace)
	if err != nil {
		return err
	}
	logging.PrintLog(namespace+" namespace has been successfully patched on your cluster.", logging.InfoLevel)
	return nil
}
