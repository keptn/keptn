package version

import (
	"errors"
	"time"

	"github.com/keptn/keptn/cli/utils/config"

	"github.com/hashicorp/go-version"
)

type AvailableVersions struct {
	Stable     VersionInfo
	Prerelease VersionInfo
}

type VersionInfo struct {
	// A version is compatible, if the major and minor version number match the current version
	newestCompatible   *version.Version
	newestIncompatible *version.Version
}

func (v AvailableVersions) Equal(o AvailableVersions) bool {
	return v.Stable.Equal(o.Stable) && v.Prerelease.Equal(o.Prerelease)
}

func (v VersionInfo) Equal(o VersionInfo) bool {
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

func getAvailableVersions(usedVersionString string, versionInfo cliVersionInfo) (AvailableVersions, error) {

	usedVersion, err := version.NewSemver(usedVersionString)
	if err != nil {
		return AvailableVersions{}, err
	}
	stable, err := browseVersions(usedVersion, versionInfo.Stable)
	if err != nil {
		return AvailableVersions{}, err
	}
	prerelease, err := browseVersions(usedVersion, versionInfo.Prerelease)
	if err != nil {
		return AvailableVersions{}, err
	}
	return AvailableVersions{Stable: stable, Prerelease: prerelease}, nil
}

func browseVersions(usedVersion *version.Version, versions []string) (VersionInfo, error) {

	newestCompatible := usedVersion
	newestIncompatible := usedVersion

	usedVersionSeg := usedVersion.Segments()
	if len(usedVersionSeg) != 3 {
		return VersionInfo{}, errors.New("Unexpected number of segements")
	}

	for _, vString := range versions {
		v, err := version.NewSemver(vString)
		if err != nil {
			return VersionInfo{}, err
		}

		vSeg := v.Segments()
		if len(vSeg) != 3 {
			return VersionInfo{}, errors.New("Unexpected number of segements")
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

	newerVersions := VersionInfo{}
	if newestCompatible.Compare(usedVersion) > 0 {
		newerVersions.newestCompatible = newestCompatible
	}
	if newestIncompatible.Compare(usedVersion) > 0 {
		newerVersions.newestIncompatible = newestIncompatible
	}

	return newerVersions, nil
}

const checkInterval = time.Hour * 24

type VersionChecker struct {
	versionFetcherClient *versionFetcherClient
}

func NewVersionChecker() *VersionChecker {
	versionChecker := VersionChecker{}
	versionChecker.versionFetcherClient = newVersionFetcherClient()
	return &versionChecker
}

// GetNewerCLIVersion checks for newer CLI versions if the automatic version check is enabled in the config
func (v *VersionChecker) GetNewerCLIVersion(cliConfig *config.CLIConfig, usedVersionString string) (AvailableVersions, error) {

	if cliConfig.AutomaticVersionCheck {
		checkTime := time.Now()
		if cliConfig.LastVersionCheck == nil || checkTime.Sub(*cliConfig.LastVersionCheck) >= checkInterval {
			cliVersionInfo, err := v.versionFetcherClient.getCLIVersionInfo(usedVersionString)
			if err != nil {
				return AvailableVersions{}, err
			}

			res, err := getAvailableVersions(usedVersionString, *cliVersionInfo)
			if err != nil {
				return AvailableVersions{}, err
			}
			cliConfig.LastVersionCheck = &checkTime
			return res, err
		}
	}

	return AvailableVersions{}, nil
}
