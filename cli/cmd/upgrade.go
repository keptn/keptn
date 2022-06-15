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
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/keptn/keptn/cli/internal"

	"github.com/keptn/keptn/cli/pkg/common"
	"github.com/keptn/keptn/cli/pkg/credentialmanager"

	"helm.sh/helm/v3/pkg/release"

	"github.com/keptn/keptn/cli/pkg/version"

	"github.com/keptn/keptn/cli/pkg/helm"
	"github.com/keptn/keptn/cli/pkg/kube"
	keptnutils "github.com/keptn/kubernetes-utils/pkg"
	"helm.sh/helm/v3/pkg/chart"

	"github.com/keptn/keptn/cli/pkg/logging"
	"github.com/spf13/cobra"
)

const keptnReleaseName = "keptn"

var keptnUpgradeChart *chart.Chart

// installCmd represents the version command
var upgraderCmd = NewUpgraderCommand(version.NewKeptnVersionChecker())

func NewUpgraderCommand(vChecker *version.KeptnVersionChecker) *cobra.Command {
	upgradeCmd := &cobra.Command{
		Use:   "upgrade",
		Args:  cobra.NoArgs,
		Short: "Upgrades Keptn on a Kubernetes cluster.",
		Long: `The Keptn CLI allows upgrading Keptn on any Kubernetes derivative to which your kube config is pointing to, and on OpenShift.

For more information, please follow the installation guide [Upgrade Keptn](https://keptn.sh/docs/` + getReleaseDocsURL() + `/operate/upgrade/)
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
				return doUpgrade()
			}
			fmt.Println("Skipping upgrade due to mocking flag")
			return nil
		},
	}

	return upgradeCmd
}

func doUpgradePreRunCheck(vChecker *version.KeptnVersionChecker) error {

	// Check if statistics service is already running and NOT deployed by helm (https://github.com/keptn/keptn/issues/3399)
	statisticsDeploymentAvailable, err := kube.CheckDeploymentAvailable("statistics-service", namespace)
	if err != nil {
		return err
	}
	if statisticsDeploymentAvailable {
		statisticsServiceManagedByHelm, err := kube.CheckDeploymentManagedByHelm("statistics-service", namespace)
		if err != nil {
			return internal.OnAPIError(err)
		}
		if !statisticsServiceManagedByHelm {
			return errors.New("deployment for statistics-service is already running and not managed by Helm. Please uninstall it")
		}
	}

	if err = addWarningNonExistingProjectUpstream(); err != nil {
		return err
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

func addWarningNonExistingProjectUpstream() error {
	endPoint, apiToken, err := credentialmanager.NewCredentialManager(assumeYes).GetCreds(namespace)
	if err != nil {
		return errors.New(authErrorMsg)
	}

	api, err := internal.APIProvider(endPoint.String(), apiToken)
	if err != nil {
		return err
	}

	logging.PrintLog(fmt.Sprintf("Connecting to server %s", endPoint.String()), logging.VerboseLevel)

	projects, err := api.ProjectsV1().GetAllProjects()
	if err != nil {
		return fmt.Errorf("failed to get all projects from namespace %s", namespace)
	}

	missingUpstreamProjects := make([]string, 0)
	for _, project := range projects {
		if project.GitRemoteURI == "" {
			missingUpstreamProjects = append(missingUpstreamProjects, project.ProjectName)
		}
	}

	if len(missingUpstreamProjects) > 0 {
		fmt.Print("WARNING: the following projects have no Git upstream configured:\n")
		for _, projectName := range missingUpstreamProjects {
			fmt.Printf("  - %s\n", projectName)
		}
		fmt.Print("Please consider setting a Git upstream repository using:\n")
		fmt.Print("  keptn update project PROJECT_NAME --git-user=GIT_USER --git-token=GIT_TOKEN --git-remote-url=GIT_REMOTE_URL\n\n")
	}

	return nil
}

func init() {
	rootCmd.AddCommand(upgraderCmd)
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

	if strings.HasPrefix(getAppVersion(keptnUpgradeChart), "0.11") {
		fmt.Printf("CAUTION: While upgrading Keptn from version %s to version %s there is a possibility of data loss due to moving to a new database model.\n", installedKeptnVersion, getAppVersion(keptnUpgradeChart))
		fmt.Printf("Please backup your data before proceeding to the next step. Information about backing up and restoring the data is described here: https://keptn.sh/docs/0.11.x/operate/upgrade/\n")

		userConfirmation := common.NewUserInput().AskBool("Did you create a backup of the database or do you want to proceed without it?", &common.UserInputOptions{AssumeYes: assumeYes})

		if !userConfirmation {
			return fmt.Errorf("Stopping upgrade.")
		}
	}

	// check if the helm-service and jmeter-service are part of the previous installation
	// if yes, they need to be installed separately as they have moved to their own charts
	helmHelper := helm.NewHelper()

	// fetch user-defined values of the previous Keptn installation and provide those to the upgrade command
	// this will ensure that options such as 'apiGatewayNginx.type' will stay the same, but newly introduced values will correctly be set to their default
	previousValues, err := helmHelper.GetValues(keptnReleaseName, keptnNamespace)
	if err != nil {
		return fmt.Errorf("Could not complete Keptn upgrade: %s", err.Error())
	}

	if err := helmHelper.UpgradeChart(keptnUpgradeChart, keptnReleaseName, keptnNamespace, previousValues); err != nil {
		msg := fmt.Sprintf("Could not complete Keptn upgrade: %s \nFor troubleshooting, please check the status of the keptn deployment by executing the following command: \n\nkubectl get pods -n %s\n", err.Error(), keptnNamespace)
		return errors.New(msg)
	}

	logging.PrintLog("Keptn has been successfully upgraded on your cluster.", logging.InfoLevel)

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
