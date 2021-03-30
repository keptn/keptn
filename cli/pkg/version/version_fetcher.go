package version

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

const versionURL = "https://get.keptn.sh/version.json"

type versionInfo struct {
	CLIVersionInfo   cliVersionInfo   `json:"cli"`
	KeptnVersionInfo keptnVersionInfo `json:"keptn"`
}

type cliVersionInfo struct {
	Stable     []string `json:"stable"`
	Prerelease []string `json:"prerelease"`
}

type keptnVersionInfo struct {
	Stable     []versionWithUpgradePath `json:"stable"`
	Prerelease []versionWithUpgradePath `json:"prerelease"`
}

type versionWithUpgradePath struct {
	Version            string   `json:"version"`
	UpgradableVersions []string `json:"upgradableVersions"`
}

type VersionFetcherClient struct {
	HttpClient *http.Client
	VersionUrl string
}

func newVersionFetcherClient() *VersionFetcherClient {
	client := VersionFetcherClient{
		HttpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
		VersionUrl: versionURL,
	}
	return &client
}

func (client *VersionFetcherClient) getCLIVersionInfo(cliVersion string) (cliVersionInfo, error) {
	v, err := client.getVersionInfo(cliVersion)
	return v.CLIVersionInfo, err
}

func (client *VersionFetcherClient) getKeptnVersionInfo(cliVersion string) (keptnVersionInfo, error) {
	v, err := client.getVersionInfo(cliVersion)
	return v.KeptnVersionInfo, err
}

func (client *VersionFetcherClient) getVersionInfo(cliVersion string) (versionInfo, error) {
	versionInfo := versionInfo{}
	req, err := http.NewRequest("GET", client.VersionUrl, nil)
	if err != nil {
		return versionInfo, err
	}
	req.Header.Set("user-agent", "keptn/cli:"+cliVersion)
	resp, err := client.HttpClient.Do(req)
	if err != nil {
		return versionInfo, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return versionInfo, err
	}
	err = json.Unmarshal(body, &versionInfo)
	return versionInfo, err
}
