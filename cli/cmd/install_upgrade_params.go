package cmd

import "github.com/keptn/keptn/cli/pkg/version"

type installUpgradeParams struct {
	ConfigFilePath     *string
	KeptnVersion       *string
	PlatformIdentifier *string
	ChartRepoURL       *string
	Namespace          *string
	PatchNamespace     *bool
	SkipUpgradeCheck   *bool
}

func getChartRepoURL(input *string) string {

	const keptnInstallerHelmRepoURL = "https://storage.googleapis.com/keptn-installer/"
	// Determine installer version
	if input != nil && *input != "" {
		return *input
	} else if version.IsOfficialKeptnVersion(Version) {
		version, _ := version.GetOfficialKeptnVersion(Version)
		return keptnInstallerHelmRepoURL + "keptn-" + version + ".tgz"
	}
	return keptnInstallerHelmRepoURL + "latest/keptn-0.1.0.tgz"
}
