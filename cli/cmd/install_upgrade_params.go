package cmd

import (
	"github.com/keptn/keptn/cli/pkg/version"
	"strings"
)

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
		// replace the name of the keptn chart with the name of the service name.
		// e.g. https://keptn.sh/keptn-0.8.1.tgz => https://keptn.sh/helm-service-0.8.1.tgz
		return getServiceChartURLFromKeptnChartURL(*input, serviceName)
	} else if version.IsOfficialKeptnVersion(Version) {
		version, _ := version.GetOfficialKeptnVersion(Version)
		return keptnInstallerHelmRepoURL + serviceName + "-" + version + ".tgz"
	}
	return keptnInstallerHelmRepoURL + "latest/" + serviceName + "-0.1.0.tgz"
}

func getServiceChartURLFromKeptnChartURL(input string, serviceName string) string {
	serviceChartURL := input
	split := strings.Split(serviceChartURL, "/")
	serviceChartURL = strings.TrimSuffix(serviceChartURL, split[len(split)-1])
	chartName := strings.Replace(split[len(split)-1], "keptn", serviceName, 1)
	serviceChartURL = serviceChartURL + chartName
	return serviceChartURL
}
