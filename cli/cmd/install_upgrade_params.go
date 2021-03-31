package cmd

import "github.com/keptn/keptn/cli/pkg/version"

const keptnInstallerHelmRepoURL = "https://storage.googleapis.com/keptn-installer/"

type installUpgradeParams struct {
	ConfigFilePath     *string
	KeptnVersion       *string
	PlatformIdentifier *string
	ChartRepoURL       *string
	Namespace          *string
	PatchNamespace     *bool
	SkipUpgradeCheck   *bool
}

func getKeptnHelmChartRepoURL(input *string) string {
	// Determine installer version
	if input != nil && *input != "" {
		return *input
	} else if version.IsOfficialKeptnVersion(Version) {
		version, _ := version.GetOfficialKeptnVersion(Version)
		return keptnInstallerHelmRepoURL + "keptn-" + version + ".tgz"
	}
	return keptnInstallerHelmRepoURL + "latest/keptn-0.1.0.tgz"
}

func getExecutionPlaneServiceChartRepoURL(input *string, serviceName string) string {
	// Determine installer version
	if input != nil && *input != "" {
		return *input
	} else if version.IsOfficialKeptnVersion(Version) {
		version, _ := version.GetOfficialKeptnVersion(Version)
		return keptnInstallerHelmRepoURL + serviceName + "-" + version + ".tgz"
	}
	return keptnInstallerHelmRepoURL + "latest/" + serviceName + "-0.1.0.tgz"
}
