package cmd

const keptnInstallerHelmRepoURL = "https://storage.googleapis.com/keptn-installer/"

type installUpgradeParams struct {
	ConfigFilePath     *string
	KeptnVersion       *string
	PlatformIdentifier *string
	ChartRepoURL       *string
	Namespace          *string
}
