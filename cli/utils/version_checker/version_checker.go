package version_checker

import (
	"errors"

	"github.com/hashicorp/go-version"
)

// GetNewerVersion compares the usedVersion with the available versionInfo.
// If the versionInfo contains a newer version, it returns this version.
// If no newer version is available, an empty string is returned.
// Note: The major and minor version number have to match in order to ensure API compatibility
// Currently, the following logic is implemented:
// 1. If versionInfo.StableVersions contain a stable version that is newer than the usedVersion, return this version
// 2. If the usedVersion is a prerelease and the versionInfo.PrereleaseVersions contain a prerelease version that is newer than
// the used usedVersion, return the prerelease version
func GetNewerVersion(usedVersionString string, versionInfo CLIVersionInfo) (string, error) {

	usedVersion, err := version.NewSemver(usedVersionString)
	if err != nil {
		return "", err
	}
	if newerVersion, err := checkForNewerVersion(usedVersion, versionInfo.StableVersions); newerVersion != "" || err != nil {
		return newerVersion, err
	}
	// Check if the usedVersion is a beta version
	if usedVersion.Prerelease() != "" {
		return checkForNewerVersion(usedVersion, versionInfo.PrereleaseVersions)
	}
	return "", nil
}

func checkForNewerVersion(usedVersion *version.Version, versions []string) (string, error) {

	usedVersionSeg := usedVersion.Segments()
	if len(usedVersionSeg) != 3 {
		return "", errors.New("Unexpected number of segements")
	}

	for _, vString := range versions {
		v, err := version.NewSemver(vString)
		if err != nil {
			return "", err
		}

		vSeg := v.Segments()
		if len(vSeg) != 3 {
			return "", errors.New("Unexpected number of segements")
		}

		// Compare the major and minor version to ensure compatible APIs
		if usedVersionSeg[0] == vSeg[0] && usedVersionSeg[1] == vSeg[1] {

			if usedVersion.Compare(v) < 0 {
				return vString, nil
			}
		}
	}
	return "", nil
}
