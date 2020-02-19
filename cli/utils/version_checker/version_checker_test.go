package version_checker

import (
	"testing"
)

var versionTests = []struct {
	usedVersionString string
	versionInfo       CLIVersionInfo
	res               string
}{
	{"0.6.0", CLIVersionInfo{StableVersions: []string{"0.7.0", "0.6.1"}}, "0.6.1"},
	{"0.6.0-beta", CLIVersionInfo{StableVersions: []string{"0.7.0", "0.6.0"}, PrereleaseVersions: []string{"0.6.0-beta2"}}, "0.6.0"},
	{"0.6.0-beta", CLIVersionInfo{StableVersions: []string{"0.7.0", "0.6.0"}, PrereleaseVersions: []string{"0.6.0-beta2"}}, "0.6.0"},
	{"0.6.0-beta", CLIVersionInfo{StableVersions: []string{"0.7.0"}, PrereleaseVersions: []string{"0.6.0-beta2"}}, "0.6.0-beta2"},
	{"0.6.0-beta", CLIVersionInfo{StableVersions: []string{"0.7.0"}, PrereleaseVersions: []string{"0.6.1-beta2"}}, "0.6.1-beta2"},

	{"0.6.0", CLIVersionInfo{StableVersions: []string{"0.7.0"}}, ""},
	{"0.6.0-beta", CLIVersionInfo{StableVersions: []string{"0.7.0"}, PrereleaseVersions: []string{"0.6.0-beta", "0.7.0-beta"}}, ""},
	{"0.6.0-alpha", CLIVersionInfo{StableVersions: []string{"0.7.0"}, PrereleaseVersions: []string{"0.6.0-beta", "0.7.0-beta"}}, "0.6.0-beta"},
}

func TestGetNewerVersion(t *testing.T) {
	for _, tt := range versionTests {
		t.Run(tt.usedVersionString, func(t *testing.T) {
			res, err := GetNewerVersion(tt.usedVersionString, tt.versionInfo)
			if err != nil {
				t.Errorf("Unexpected error %v", err)
			}
			if res != tt.res {
				t.Errorf("got %s, want %s for %s", res, tt.res, tt.usedVersionString)
			}
		})
	}
}
