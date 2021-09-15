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
	"errors"
	"fmt"
	"github.com/keptn/keptn/cli/pkg/credentialmanager"
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
var upgraderCmd = NewUpgraderCommand(version.NewKeptnVersionChecker())

func NewUpgraderCommand(vChecker *version.KeptnVersionChecker) *cobra.Command {
	upgradeCmd := &cobra.Command{
		Use:   "upgrade",
		Short: "Upgrades Keptn on a Kubernetes cluster and supports upgrading the shipyard of a project to a new specification.",
		Long: `The Keptn CLI allows upgrading Keptn on any Kubernetes derivative to which your kube config is pointing to, and on OpenShift. Also, it supports upgrading the shipyard of an existing project to a new specification.

For more information, please follow the installation guide [Upgrade Keptn](https://keptn.sh/docs/` + keptnReleaseDocsURL + `/operate/upgrade/)
`,
		Example: `keptn upgrade # upgrades Keptn

keptn upgrade --platform=openshift # upgrades Keptn on OpenShift

keptn upgrade --platform=kubernetes # upgrades Keptn on the Kubernetes cluster
`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := doUpgradePreRunCheck(vChecker); err != nil {
				return err
			}
			if !mocking {
				if *upgradeParams.PatchNamespace {
					return patchNamespace()
				}
				return doUpgrade()
			}
			fmt.Println("Skipping upgrade due to mocking flag")
			return nil
		},
	}

	return upgradeCmd
}

func doUpgradePreRunCheck(vChecker *version.KeptnVersionChecker) error {
	if *upgradeParams.PatchNamespace {
		return nil
	}

	chartRepoURL := getKeptnHelmChartRepoURL(upgradeParams.ChartRepoURL)

	var err error
	if keptnUpgradeChart, err = helm.NewHelper().DownloadChart(chartRepoURL); err != nil {
		return err
	}

	if !*upgradeParams.SkipUpgradeCheck {
		res, err := isUpgradeCompatible(vChecker)
		if err != nil {
			return err
		}
		if !res {
			installedKeptnVersion, err := getInstalledKeptnVersion()
			if err != nil {
				return err
			}
			if installedKeptnVersion == getAppVersion(keptnUpgradeChart) {
				vChecker := version.NewVersionChecker()
				cliVersionCheck, _ := vChecker.CheckCLIVersion(Version, false)
				if cliVersionCheck {
					return fmt.Errorf("Please upgrade Keptn CLI to upgrade your Keptn Cluster!")
				}
				return fmt.Errorf("Unable to check for upgrades due to aforementioned error")
			}
			return fmt.Errorf("No upgrade path exists from Keptn version %s to %s",
				installedKeptnVersion, getAppVersion(keptnUpgradeChart))
		}
	} else {
		logging.PrintLog("Skipping upgrade compatibility check!", logging.InfoLevel)
	}

	logging.PrintLog(fmt.Sprintf("Helm Chart used for Keptn upgrade: %s", chartRepoURL), logging.InfoLevel)

	cm := credentialmanager.NewCredentialManager(assumeYes)
	currentKeptnCLIContext := cm.GetCurrentKeptnCLIConfig().CurrentContext
	currentKubernetesContext := cm.GetCurrentKubeConfig().CurrentContext

	if currentKeptnCLIContext != currentKubernetesContext {
		return fmt.Errorf("your current Keptn CLI context '%s' does not match current Kubeconfig '%s'. Please ensure your kubectl CLI is connected to '%s' before upgrading your Keptn cluster", currentKeptnCLIContext, currentKubernetesContext, currentKubernetesContext)
	}

	platformManager, err := platform.NewPlatformManager(*upgradeParams.PlatformIdentifier, cm)
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
		err = platformManager.ReadCreds(assumeYes)
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

	// Check if statistics service is already running and NOT deployed by helm (https://github.com/keptn/keptn/issues/3399)
	statisticsDeploymentAvailable, err := kube.CheckDeploymentAvailable("statistics-service", namespace)
	if err != nil {
		return err
	}
	if statisticsDeploymentAvailable {
		statisticsServiceManagedByHelm, err := kube.CheckDeploymentManagedByHelm("statistics-service", namespace)
		if err != nil {
			return err
		}
		if !statisticsServiceManagedByHelm {
			return errors.New("deployment for statistics-service is already running and not managed by Helm. Please uninstall it")
		}
	}

	return nil
}

func getInstalledKeptnVersion() (string, error) {
	if mocking {
		// return a fake version
		return "0.7.0", nil
	}
	lastRelease, err := getLatestKeptnRelease()
	if err != nil {
		return "", err
	}
	return lastRelease.Chart.Metadata.AppVersion, nil
}

func getAppVersion(ch *chart.Chart) string {
	return ch.Metadata.AppVersion
}

func isUpgradeCompatible(versionChecker *version.KeptnVersionChecker) (bool, error) {
	installedVersion, err := getInstalledKeptnVersion()
	if err != nil {
		return false, err
	}
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

	// iterate over releases and find the one with status = deployed
	for _, r := range releases {
		if r.Info.Status == release.StatusDeployed {
			return r, nil
		}
	}

	return nil, fmt.Errorf("Found %d releases, but none of them is currently deployed", len(releases))

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
	upgradeParams.SkipUpgradeCheck = upgraderCmd.Flags().BoolP("skip-upgrade-check", "", false, "Skip upgrade compatibility check, useful for nightly version upgrades or upgrades to preview versions")
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

	// check if the helm-service and jmeter-service are part of the previous installation
	// if yes, they need to be installed separately as they have moved to their own charts
	helmHelper := helm.NewHelper()

	if err := helmHelper.UpgradeChart(keptnUpgradeChart, keptnReleaseName, keptnNamespace, nil); err != nil {
		msg := fmt.Sprintf("Could not complete Keptn upgrade: %s \nFor troubleshooting, please check the status of the keptn deployment by executing the following command: \n\nkubectl get pods -n %s\n", err.Error(), keptnNamespace)
		return errors.New(msg)
	}

	logging.PrintLog("Upgrading of Keptn control plane has been successful.", logging.InfoLevel)

	values, err := helmHelper.GetValues(keptnReleaseName, keptnNamespace)
	if err != nil {
		return fmt.Errorf("Could not determine configuration of current Keptn installation: %s", err.Error())
	}
	if isContinuousDeliveryEnabled(values) {
		logging.PrintLog("Upgrading execution plane services for continuous-delivery use case.", logging.InfoLevel)
		continuousDeliveryServiceCharts, err := fetchContinuousDeliveryCharts(*helmHelper, upgradeParams.ChartRepoURL)
		if err != nil {
			return fmt.Errorf("Could not fetch continuous delivery execution service charts: %s \n", err.Error())
		}

		for _, serviceChart := range continuousDeliveryServiceCharts {
			if err := helmHelper.UpgradeChart(serviceChart, serviceChart.Name(), keptnNamespace, values); err != nil {
				msg := fmt.Sprintf("Could not complete upgrade of Keptn execution plane services: %s \nFor troubleshooting, please check the status of the keptn deployment by executing the following command: \n\nkubectl get pods -n %s\n", err.Error(), keptnNamespace)
				return errors.New(msg)
			}
		}
	}

	logging.PrintLog("Keptn has been successfully upgraded on your cluster.", logging.InfoLevel)
	// when upgrading from 0.7.x to 0.8.x, display information about how to upgrade projects to the new shipyard format
	if strings.HasPrefix(installedKeptnVersion, "0.7") && strings.HasPrefix(getAppVersion(keptnUpgradeChart), "0.8") {
		logging.PrintLog("Please upgrade your projects using keptn upgrade project.", logging.InfoLevel)
		logging.PrintLog("For detailed instructions about upgrading your projects, head to keptn.sh/docs/0.8.x/operate/upgrade", logging.InfoLevel)
	}
	return nil
}

func isContinuousDeliveryEnabled(configValues map[string]interface{}) bool {
	if continuousDeliveryConfig, ok := configValues["continuous-delivery"].(map[string]interface{}); ok {
		if isEnabled, ok := continuousDeliveryConfig["enabled"].(bool); ok {
			return isEnabled
		} else if isEnabled, ok := continuousDeliveryConfig["enabled"].(string); ok {
			return isEnabled == "true"
		}
	}
	return false
}

func patchNamespace() error {
	err := keptnutils.PatchKeptnManagedNamespace(false, namespace)
	if err != nil {
		return err
	}
	logging.PrintLog(namespace+" namespace has been successfully patched on your cluster.", logging.InfoLevel)
	return nil
}
