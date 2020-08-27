package version

import (
	"io"
	"net/http"
	"testing"

	"github.com/magiconair/properties/assert"
)

func TestGetNewStableVersions(t *testing.T) {

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, versionJsonTest)
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
	{"0.9.0", "0.8.1", false},
}

func TestIsUpgradable(t *testing.T) {

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, versionJsonTest)
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
