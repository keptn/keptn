package version

import (
	"errors"
	"fmt"
	"time"

	"github.com/keptn/keptn/cli/pkg/logging"

	"github.com/keptn/keptn/cli/pkg/config"

	"github.com/hashicorp/go-version"
)

type availableNewestVersions struct {
	stable     newestVersions
	prerelease newestVersions
}

type newestVersions struct {
	// A version is compatible, if the major and minor version number match the current version
	newestCompatible   *version.Version
	newestIncompatible *version.Version
}

func (v availableNewestVersions) equal(o availableNewestVersions) bool {
	return v.stable.equal(o.stable) && v.prerelease.equal(o.prerelease)
}

func (v newestVersions) equal(o newestVersions) bool {
	return equal(v.newestCompatible, o.newestCompatible) && equal(v.newestIncompatible, o.newestIncompatible)
}

func equal(v *version.Version, o *version.Version) bool {
	if v == nil {
		return o == nil
	}
	if o == nil {
		return v == nil
	}
	return v.Equal(o)
}

func getAvailableVersions(usedVersionString string, versionInfo cliVersionInfo) (availableNewestVersions, error) {

	usedVersion, err := version.NewSemver(usedVersionString)
	if err != nil {
		return availableNewestVersions{}, err
	}
	stable, err := browseVersions(usedVersion, versionInfo.Stable)
	if err != nil {
		return availableNewestVersions{}, err
	}
	prerelease, err := browseVersions(usedVersion, versionInfo.Prerelease)
	if err != nil {
		return availableNewestVersions{}, err
	}
	return availableNewestVersions{stable: stable, prerelease: prerelease}, nil
}

func browseVersions(usedVersion *version.Version, versions []string) (newestVersions, error) {

	newestCompatible := usedVersion
	newestIncompatible := usedVersion

	usedVersionSeg := usedVersion.Segments()
	if len(usedVersionSeg) != 3 {
		return newestVersions{}, errors.New("Unexpected number of segements")
	}

	for _, vString := range versions {
		v, err := version.NewSemver(vString)
		if err != nil {
			return newestVersions{}, err
		}

		vSeg := v.Segments()
		if len(vSeg) != 3 {
			return newestVersions{}, errors.New("Unexpected number of segements")
		}

		if usedVersionSeg[0] == vSeg[0] && usedVersionSeg[1] == vSeg[1] {
			// Major and minor version match (ensures compatible APIs)
			if newestCompatible.Compare(v) < 0 {
				newestCompatible = v
			}
		} else if newestIncompatible.Compare(v) < 0 {
			newestIncompatible = v
		}
	}

	newerVersions := newestVersions{}
	if newestCompatible.Compare(usedVersion) > 0 {
		newerVersions.newestCompatible = newestCompatible
	}
	if newestIncompatible.Compare(usedVersion) > 0 {
		newerVersions.newestIncompatible = newestIncompatible
	}

	return newerVersions, nil
}

const checkInterval = time.Hour * 24

// VersionChecker manages data for checking the version
type VersionChecker struct {
	versionFetcherClient *versionFetcherClient
}

// NewVersionChecker creates a new VersionChecker
func NewVersionChecker() *VersionChecker {
	versionChecker := VersionChecker{}
	versionChecker.versionFetcherClient = newVersionFetcherClient()
	return &versionChecker
}

// getNewerCLIVersion checks for newer CLI versions
func (v *VersionChecker) getNewerCLIVersion(usedVersionString string) (availableNewestVersions, error) {

	cliVersionInfo, err := v.versionFetcherClient.getCLIVersionInfo(usedVersionString)
	if err != nil {
		return availableNewestVersions{}, fmt.Errorf("error when fetching CLI version infos: %v", err)
	}

	res, err := getAvailableVersions(usedVersionString, cliVersionInfo)
	if err != nil {
		return availableNewestVersions{}, fmt.Errorf("error when analyzing the available versions: %v", err)
	}
	return res, nil
}

const newCompatibleVersionMsg = `* Keptn CLI version %s is available! 
 - To update to the latest CLI run this command: curl -sL https://get.keptn.sh | bash
 - For more information, please visit https://keptn.sh/docs/%s/operate/upgrade/`
const newIncompatibleVersionMsg = `* Keptn CLI version %s is available! Please note that this version might be incompatible with your Keptn cluster ` +
	`version and requires to update the cluster too. Please visit https://keptn.sh/docs/%s/operate/upgrade/ for more information.`

// CheckCLIVersion checks whether there is a new CLI version available and prints corresponding
// messages to the stdout
func (v *VersionChecker) CheckCLIVersion(cliVersion string, considerPrevCheck bool) (bool, bool) {

	configMng := config.NewCLIConfigManager()
	cliConfig, err := configMng.LoadCLIConfig()
	if err != nil {
		logging.PrintLog(err.Error(), logging.InfoLevel)
		return false, false
	}

	msgPrinted := false
	if cliConfig.AutomaticVersionCheck && IsOfficialKeptnVersion(cliVersion) {
		checkTime := time.Now()
		if !considerPrevCheck || cliConfig.LastVersionCheck == nil ||
			checkTime.Sub(*cliConfig.LastVersionCheck) >= checkInterval {
			newVersions, err := v.getNewerCLIVersion(cliVersion)
			if err != nil {
				logging.PrintLog(err.Error(), logging.InfoLevel)
				return false, false
			}
			if newVersions.stable.newestCompatible != nil {
				segments := newVersions.stable.newestCompatible.Segments()
				majorMinorXVersion := fmt.Sprintf("%v.%v.x", segments[0], segments[1])
				fmt.Printf(newCompatibleVersionMsg+"\n", newVersions.stable.newestCompatible.String(),
					majorMinorXVersion)
				msgPrinted = true
			}
			if newVersions.prerelease.newestCompatible != nil {
				segments := newVersions.prerelease.newestCompatible.Segments()
				majorMinorXVersion := fmt.Sprintf("%v.%v.x", segments[0], segments[1])
				fmt.Printf(newCompatibleVersionMsg+"\n", newVersions.prerelease.newestCompatible.String(),
					majorMinorXVersion)
				msgPrinted = true
			}
			if newVersions.stable.newestIncompatible != nil {
				segments := newVersions.stable.newestIncompatible.Segments()
				majorMinorXVersion := fmt.Sprintf("%v.%v.x", segments[0], segments[1])
				fmt.Printf(newIncompatibleVersionMsg+"\n", newVersions.stable.newestIncompatible.String(),
					majorMinorXVersion)
				msgPrinted = true
			}
			return true, msgPrinted
		}
	}
	return false, msgPrinted
}

// IsOfficialKeptnVersion checks whether the provided version string follows a Keptn version pattern
func IsOfficialKeptnVersion(versionStr string) bool {
	_, err := version.NewSemver(versionStr)
	return err == nil
}

// GetOfficialKeptnVersion extracts the Keptn version from the provided string
// More precisely, this method returns the segments and prerelease info w/o metadata
func GetOfficialKeptnVersion(versionStr string) (string, error) {
	s, err := version.NewSemver(versionStr)
	if err != nil {
		return "", err
	}
	v := s.String()
	metadata := s.Metadata()
	if metadata != "" {
		metadata = "+" + metadata
	}
	return v[:len(v)-len(metadata)], nil
}
