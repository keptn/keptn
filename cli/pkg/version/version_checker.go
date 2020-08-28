package version

import (
	"errors"
	"fmt"
	"os"
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

// getNewerCLIVersion checks for newer CLI versions if the automatic version check is enabled in the config
func (v *VersionChecker) getNewerCLIVersion(cliConfig config.CLIConfig, usedVersionString string) (availableNewestVersions, error) {

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

const newCompatibleVersionMsg = `keptn version %s is available! Please visit https://keptn.sh for more information.`
const newIncompatibleVersionMsg = `keptn version %s is available! Please note that this version might be incompatible with your Keptn cluster ` +
	`version and requires to update the cluster too. Please visit https://keptn.sh for more information.`
const disableMsg = `To disable this notice, run: '%s set config AutomaticVersionCheck false'`

// CheckCLIVersion checks whether there is a new CLI version available and prints corresponding
// messages to the stdout
func (v *VersionChecker) CheckCLIVersion(cliVersion string, considerPrevCheck bool) {

	configMng := config.NewCLIConfigManager()
	cliConfig, err := configMng.LoadCLIConfig()
	if err != nil {
		logging.PrintLog(err.Error(), logging.InfoLevel)
		return
	}

	if cliConfig.AutomaticVersionCheck && IsOfficialKeptnVersion(cliVersion) {
		checkTime := time.Now()
		if !considerPrevCheck || cliConfig.LastVersionCheck == nil ||
			checkTime.Sub(*cliConfig.LastVersionCheck) >= checkInterval {
			newVersions, err := v.getNewerCLIVersion(cliConfig, cliVersion)
			if err != nil {
				logging.PrintLog(err.Error(), logging.InfoLevel)
				return
			}
			msgPrinted := false
			if newVersions.stable.newestCompatible != nil {
				fmt.Printf(newCompatibleVersionMsg+"\n", newVersions.stable.newestCompatible.String())
				msgPrinted = true
			}
			if newVersions.prerelease.newestCompatible != nil {
				fmt.Printf(newCompatibleVersionMsg+"\n", newVersions.prerelease.newestCompatible)
				msgPrinted = true
			}
			if newVersions.stable.newestIncompatible != nil {
				fmt.Printf(newIncompatibleVersionMsg+"\n", newVersions.stable.newestIncompatible.String())
				msgPrinted = true
			}
			if msgPrinted && considerPrevCheck {
				if len(os.Args) > 0 {
					fmt.Printf(disableMsg+"\n", os.Args[0])
				} else {
					fmt.Printf(disableMsg+"\n", "keptn")
				}
			}

			cliConfig.LastVersionCheck = &checkTime
			if err := configMng.StoreCLIConfig(cliConfig); err != nil {
				logging.PrintLog(err.Error(), logging.InfoLevel)
				return
			}
		}
	}
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
