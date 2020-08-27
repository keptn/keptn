package version

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-version"
	"github.com/keptn/keptn/cli/pkg/logging"
)

type KeptnVersionChecker struct {
	versionFetcherClient *versionFetcherClient
}

// NewVersionChecker creates a new VersionChecker
func NewKeptnVersionChecker() *KeptnVersionChecker {
	versionChecker := KeptnVersionChecker{}
	versionChecker.versionFetcherClient = newVersionFetcherClient()
	return &versionChecker
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
