package version

import (
	"io"
	"net/http"
	"testing"

	"github.com/hashicorp/go-version"
	"github.com/stretchr/testify/assert"
)

func TestGetNewStableVersions(t *testing.T) {

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, versionJSONTest)
	})

	httpClient, url, teardown := testingHTTPClient(handler)
	defer teardown()

	versionChecker := NewKeptnVersionChecker()
	versionChecker.versionFetcherClient.httpClient = httpClient
	versionChecker.versionFetcherClient.versionUrl = url

	res, err := versionChecker.GetStableVersions("0.7.0", "0.7.0")

	assert.Equal(t, err, nil, "Unexpected error")
	assert.Equal(t, res, []string{"0.7.1", "0.8.0"}, "Expected 2 new versions")
}

var isUpgradableTests = []struct {
	currentKeptnVersion    string
	newDesiredKeptnVersion string
	res                    bool
}{
	{"0.7.0", "0.7.1", true},
	{"0.7.0", "0.7.2", false},
	{"0.8.0", "0.8.1", false},
	{"0.8.0", "0.9.0", true},
	{"0.9.0", "0.8.1", false},
}

func TestIsUpgradable(t *testing.T) {

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, versionJSONTest)
	})

	httpClient, url, teardown := testingHTTPClient(handler)
	defer teardown()

	versionChecker := NewKeptnVersionChecker()
	versionChecker.versionFetcherClient.httpClient = httpClient
	versionChecker.versionFetcherClient.versionUrl = url

	for _, tt := range isUpgradableTests {
		t.Run(tt.currentKeptnVersion, func(t *testing.T) {
			res, err := versionChecker.IsUpgradable("0.7.0", tt.currentKeptnVersion, tt.newDesiredKeptnVersion)
			if err != nil {
				t.Errorf("Unexpected error %v", err)
			}
			if res != tt.res {
				t.Errorf("got %t, want %t for %s", res, tt.res, tt.currentKeptnVersion)
			}
		})
	}
}

var checkKeptnVersionTests = []struct {
	currentKeptnVersion    string
	newKeptnVersion string
}{
	{"0.7.0", "0.8.0"},
	{"0.7.1", "0.8.0"},
	{"0.8.0", "0.9.0"},
}

func TestCheckKeptnVersion(t *testing.T) {

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, versionJSONTest)
	})

	httpClient, url, teardown := testingHTTPClient(handler)
	defer teardown()

	versionChecker := NewKeptnVersionChecker()
	versionChecker.versionFetcherClient.httpClient = httpClient
	versionChecker.versionFetcherClient.versionUrl = url

	for _, tt := range checkKeptnVersionTests {
		t.Run(tt.currentKeptnVersion, func(t *testing.T) {
			res, err := versionChecker.getNewestStableVersion(tt.currentKeptnVersion, tt.currentKeptnVersion)
			expectedRes, err := version.NewVersion(tt.newKeptnVersion)
			assert.Equal(t, err, nil, "Unexpected error")
			assert.Equal(t, res.Equal(expectedRes), true, "Wrong version")
		})
	}
}
