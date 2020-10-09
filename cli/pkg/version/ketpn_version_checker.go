package version

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-version"
	"github.com/keptn/keptn/cli/pkg/config"
	"github.com/keptn/keptn/cli/pkg/logging"
)

// KeptnVersionChecker implements functions for checking the Keptn-cluster version
type KeptnVersionChecker struct {
	versionFetcherClient *versionFetcherClient
}

// NewKeptnVersionChecker creates a new VersionChecker
func NewKeptnVersionChecker() *KeptnVersionChecker {
	versionChecker := KeptnVersionChecker{}
	versionChecker.versionFetcherClient = newVersionFetcherClient()
	return &versionChecker
}

const newKeptnVersionMsg = `* Keptn version %s is available! Please visit https://keptn.sh/docs/%s/operate/upgrade/ for more information.`

// CheckKeptnVersion checks whether there is a new Keptn version available and prints corresponding
// messages to the stdout
func (c KeptnVersionChecker) CheckKeptnVersion(cliVersion, clusterVersion string, considerPrevCheck bool) (bool, bool) {
	configMng := config.NewCLIConfigManager()
	cliConfig, err := configMng.LoadCLIConfig()
	if err != nil {
		logging.PrintLog(err.Error(), logging.InfoLevel)
		return false, false
	}

	msgPrinted := false
	if cliConfig.AutomaticVersionCheck && IsOfficialKeptnVersion(clusterVersion) {
		checkTime := time.Now()
		if !considerPrevCheck || cliConfig.LastVersionCheck == nil ||
			checkTime.Sub(*cliConfig.LastVersionCheck) >= checkInterval {
			newVersion, err := c.getNewestStableVersion(cliVersion, clusterVersion)
			if err != nil {
				logging.PrintLog(err.Error(), logging.InfoLevel)
				return false, false
			}
			if newVersion != nil {
				segments := newVersion.Segments()
				majorMinorXVersion := fmt.Sprintf("%v.%v.x", segments[0], segments[1])
				fmt.Printf(newKeptnVersionMsg+"\n", newVersion.String(), majorMinorXVersion)
				msgPrinted = true
			}
			return true, msgPrinted
		}
	}
	return false, msgPrinted
}

// getNewestStableVersion returns the newest stable version to which the current version can be upgraded
func (c KeptnVersionChecker) getNewestStableVersion(cliVersion, keptnVersion string) (*version.Version, error) {
	keptnVersionInfo, err := c.versionFetcherClient.getKeptnVersionInfo(cliVersion)
	if err != nil {
		return nil, fmt.Errorf("error when fetching Keptn version infos: %v", err)
	}

	currentVersion, err := version.NewSemver(keptnVersion)
	if err != nil {
		return nil, fmt.Errorf("error when parsing current Keptn version: %v", err)
	}

	var latestVersion *version.Version
	for _, kv := range keptnVersionInfo.Stable {
		availableVersion, err := version.NewSemver(kv.Version)
		if err != nil {
			logging.PrintLog(fmt.Sprintf("error when parsing version %s", kv.Version), logging.InfoLevel)
			continue
		}
		if availableVersion.Compare(currentVersion) > 0 && contains(kv.UpgradableVersions, keptnVersion) {
			if latestVersion == nil || availableVersion.Compare(latestVersion) > 0 {
				latestVersion = availableVersion
			}
		}
	}
	return latestVersion, nil
}

// GetStableVersions returns a list of all stable version to which the current version can be upgraded
func (c KeptnVersionChecker) GetStableVersions(cliVersion, keptnVersion string) ([]string, error) {

	keptnVersionInfo, err := c.versionFetcherClient.getKeptnVersionInfo(cliVersion)
	if err != nil {
		return nil, fmt.Errorf("error when fetching Keptn version infos: %v", err)
	}

	currentVersion, err := version.NewSemver(keptnVersion)
	if err != nil {
		return nil, fmt.Errorf("error when parsing current Keptn version: %v", err)
	}

	upgradeableVersion := make([]string, 0)

	for _, kv := range keptnVersionInfo.Stable {
		availableVersion, err := version.NewSemver(kv.Version)
		if err != nil {
			logging.PrintLog(fmt.Sprintf("error when parsing version %s", kv.Version), logging.InfoLevel)
			continue
		}
		if availableVersion.Compare(currentVersion) > 0 && contains(kv.UpgradableVersions, keptnVersion) {
			upgradeableVersion = append(upgradeableVersion, kv.Version)
		}
	}
	return upgradeableVersion, nil
}

// IsUpgradable checks whether a Keptn version can be upgraded to a new version
func (c KeptnVersionChecker) IsUpgradable(cliVersion, currentKeptnVersion, newDesiredKeptnVersion string) (bool, error) {

	versions, err := c.GetStableVersions(cliVersion, currentKeptnVersion)
	if err != nil {
		return false, err
	}
	return contains(versions, newDesiredKeptnVersion), nil
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if strings.TrimSpace(a) == e {
			return true
		}
	}
	return false
}
