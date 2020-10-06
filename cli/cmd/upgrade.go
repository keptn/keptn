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
	"github.com/keptn/keptn/cli/pkg/logging"
	"github.com/keptn/keptn/cli/pkg/version"

	"github.com/spf13/cobra"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/release"
	"helm.sh/helm/v3/pkg/strvals"
)

const keptnReleaseName = "keptn"

var upgradeParams installUpgradeParams
var keptnUpgradeChart *chart.Chart
var upgradeValues map[string]interface{}


// upgraderCmd represents the upgrade command
var upgraderCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrades Keptn on a Kubernetes cluster",
	Long: `The Keptn CLI allows upgrading Keptn on any Kubernetes derivate to which your kube config is pointing to, and on OpenShift.

For more information, please follow the installation guide [Upgrade Keptn](https://keptn.sh/docs/` + keptnReleaseDocsURL + `/operate/upgrade/)
`,
	Example: `keptn upgrade --values=./my-values.yaml`,
	SilenceUsage: true,
	PreRunE: func(cmd *cobra.Command, args []string) error {

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

		var valuesFile map[string]interface{}
		if upgradeParams.ValuesFile != nil && *upgradeParams.ValuesFile != "" {
			err = getAndParseYaml(*upgradeParams.ValuesFile, &valuesFile)
			if err != nil {
				return fmt.Errorf("Failed to read and parse values file - %s", err.Error())
			}
		}
		latestReleaseValues, err := getLatestKeptnReleaseValues()
		if err != nil {
			return fmt.Errorf("Failed to get values of installed release - %s", err.Error())

		}
		upgradeValues = mergeMaps(latestReleaseValues, valuesFile)

		if upgradeParams.Values != nil {
			for _, value := range *upgradeParams.Values {
				if err := strvals.ParseInto(value, upgradeValues); err != nil {
					return fmt.Errorf("Failed to parse --set data - %s", err.Error())
				}
			}
		}

		return checkInput(upgradeParams, upgradeValues)
	},
	RunE: func(cmd *cobra.Command, args []string) error {

		if !mocking {
			return doUpgrade()
		}
		fmt.Println("Skipping upgrade due to mocking flag")
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

func getLatestKeptnReleaseValues() (map[string]interface{}, error) {
	lastRelease, err := getLatestKeptnRelease()
	if err != nil {
		return nil, err
	}
	return helm.NewHelper().GetValues(lastRelease.Name, lastRelease.Namespace)
}

func getLatestKeptnRelease() (*release.Release, error) {
	keptnNamespace := *upgradeParams.Namespace
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

	upgradeParams.ChartRepoURL = upgraderCmd.Flags().StringP("chart-repo", "",
		"", "URL of the Keptn Helm Chart repository")
	upgraderCmd.Flags().MarkHidden("chart-repo")

	upgradeParams.Namespace = upgraderCmd.Flags().StringP("namespace", "n", "keptn",
		"Specify the namespace where Keptn should be upgraded (default keptn).")

	upgradeParams.ValuesFile = upgraderCmd.Flags().StringP("values", "f", "",
		"Specify values in a YAML file or a URL.")

	upgradeParams.Values = upgraderCmd.Flags().StringArrayP("set", "", []string{},
		"Set values on the command line (can specify multiple or separate values with commas: key1=val1,key2=val2)")
}

func doUpgrade() error {
	keptnNamespace := *upgradeParams.Namespace

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

	if err := helm.NewHelper().UpgradeChart(keptnUpgradeChart, keptnReleaseName, keptnNamespace, upgradeValues); err != nil {
		logging.PrintLog("Could not complete Keptn upgrade: "+err.Error(), logging.InfoLevel)
		return err
	}

	logging.PrintLog("Keptn has been successfully upgraded on your cluster.", logging.InfoLevel)
	return nil
}
