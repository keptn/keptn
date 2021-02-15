package version

import (
	"io"
	"net/http"
	"testing"

	"github.com/hashicorp/go-version"

	"github.com/stretchr/testify/assert"
)

var versionTests = []struct {
	usedVersionString string
	versionInfo       cliVersionInfo
	res               availableNewestVersions
}{
	{"0.6.0", cliVersionInfo{Stable: []string{"0.7.0", "0.6.1"}}, availableVersionInitHelper("0.6.1", "0.7.0", "", "")},
	{"0.6.0-beta", cliVersionInfo{Stable: []string{"0.7.0", "0.6.0"}, Prerelease: []string{"0.6.0-beta2"}}, availableVersionInitHelper("0.6.0", "0.7.0", "0.6.0-beta2", "")},
	{"0.6.0-beta", cliVersionInfo{Stable: []string{"0.7.0"}, Prerelease: []string{"0.6.0-beta2"}}, availableVersionInitHelper("", "0.7.0", "0.6.0-beta2", "")},
	{"0.6.0-beta", cliVersionInfo{Stable: []string{"0.7.0"}, Prerelease: []string{"0.6.1-beta2"}}, availableVersionInitHelper("", "0.7.0", "0.6.1-beta2", "")},

	{"0.6.0", cliVersionInfo{Stable: []string{"0.7.0"}}, availableVersionInitHelper("", "0.7.0", "", "")},
	{"0.6.0-beta", cliVersionInfo{Stable: []string{"0.7.0"}, Prerelease: []string{"0.6.0-beta2", "0.7.0-beta"}}, availableVersionInitHelper("", "0.7.0", "0.6.0-beta2", "0.7.0-beta")},
	{"0.6.0-alpha", cliVersionInfo{Stable: []string{"0.7.0"}, Prerelease: []string{"0.6.0-beta", "0.7.0-beta"}}, availableVersionInitHelper("", "0.7.0", "0.6.0-beta", "0.7.0-beta")},
	{"0.6.0", cliVersionInfo{Stable: []string{"0.7.0"}, Prerelease: []string{"0.6.0-beta"}}, availableVersionInitHelper("", "0.7.0", "", "")},
	{"0.6.0", cliVersionInfo{Stable: []string{"0.7.0"}, Prerelease: []string{"0.6.1-beta"}}, availableVersionInitHelper("", "0.7.0", "0.6.1-beta", "")},
}

func availableVersionInitHelper(newestCompatibleStable string, newestIncompatibleStable string,
	newestCompatiblePrerelease string, newestIncompatiblePrerelease string) availableNewestVersions {
	var nCS *version.Version
	var nIS *version.Version
	var nCP *version.Version
	var nIP *version.Version
	if newestCompatibleStable != "" {
		nCS, _ = version.NewSemver(newestCompatibleStable)
	}
	if newestIncompatibleStable != "" {
		nIS, _ = version.NewSemver(newestIncompatibleStable)
	}
	if newestCompatiblePrerelease != "" {
		nCP, _ = version.NewSemver(newestCompatiblePrerelease)
	}
	if newestIncompatiblePrerelease != "" {
		nIP, _ = version.NewSemver(newestIncompatiblePrerelease)
	}
	return availableNewestVersions{
		stable:     newestVersions{newestCompatible: nCS, newestIncompatible: nIS},
		prerelease: newestVersions{newestCompatible: nCP, newestIncompatible: nIP},
	}
}

func TestGetNewerVersion(t *testing.T) {
	for _, tt := range versionTests {
		t.Run(tt.usedVersionString, func(t *testing.T) {
			res, err := getAvailableVersions(tt.usedVersionString, tt.versionInfo)
			if err != nil {
				t.Errorf("Unexpected error %v", err)
			}
			if !res.equal(tt.res) {
				t.Errorf("got %v, want %v for %s", res, tt.res, tt.usedVersionString)
			}
		})
	}
}

func TestCheckCLIVersion(t *testing.T) {

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, `{ "cli": { "stable": ["0.5.2", "0.6.1"], "prerelease": ["0.6.0-beta2"] } }`)
	})

	httpClient, url, teardown := testingHTTPClient(handler)
	defer teardown()

	versionChecker := NewVersionChecker()
	versionChecker.versionFetcherClient.httpClient = httpClient
	versionChecker.versionFetcherClient.versionUrl = url

	res, err := versionChecker.getNewerCLIVersion("0.6.0")

	expectedRes := availableVersionInitHelper("0.6.1", "", "", "")

	assert.Equal(t, err, nil, "Unexpected error")
	assert.Equal(t, res.equal(expectedRes), true, "Wrong versions")
}

var iskeptnVersions = []struct {
	in  string
	res bool
}{
	{"master+20191212.1033", false},
	{"0.6.0-beta2", true},
	{"0.6.0", true},
	{"feature-443+20191213.1105", false},
	{"0.6.0-beta2+20191204.1329", true},
	{"0.6.0-beta2+201912044.1329", true},
}

func TestIsOfficialKeptnVersion(t *testing.T) {
	for _, tt := range iskeptnVersions {
		t.Run(tt.in, func(t *testing.T) {
			res := IsOfficialKeptnVersion(tt.in)
			if res != tt.res {
				t.Errorf("got %t, want %t", res, tt.res)
			}
		})
	}
}

var getkeptnVersions = []struct {
	in  string
	res string
}{
	{"master+20191212.1033", ""},
	{"0.6.0-beta2", "0.6.0-beta2"},
	{"0.6.0", "0.6.0"},
	{"feature-443+20191213.1105", ""},
	{"0.6.0-beta2+20191204.1329", "0.6.0-beta2"},
	{"0.6.0-beta2+201912044.1329", "0.6.0-beta2"},
}

func TestGetOfficialKeptnVersion(t *testing.T) {
	for _, tt := range getkeptnVersions {
		t.Run(tt.in, func(t *testing.T) {
			res, _ := GetOfficialKeptnVersion(tt.in)
			if res != tt.res {
				t.Errorf("got %s, want %s", res, tt.res)
			}
		})
	}
}
