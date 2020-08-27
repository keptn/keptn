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
	Stable     []VersionWithUpgradePath `json:"stable"`
	Prerelease []VersionWithUpgradePath `json:"prerelease"`
}

type VersionWithUpgradePath struct {
	Version            string   `json:"version"`
	UpgradableVersions []string `json:"upgradableVersions"`
}

type versionFetcherClient struct {
	httpClient *http.Client
	versionUrl string
}

func newVersionFetcherClient() *versionFetcherClient {
	client := versionFetcherClient{
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
		versionUrl: versionURL,
	}
	return &client
}

func (client *versionFetcherClient) getCLIVersionInfo(cliVersion string) (cliVersionInfo, error) {
	v, err := client.getVersionInfo(cliVersion)
	return v.CLIVersionInfo, err
}

func (client *versionFetcherClient) getKeptnVersionInfo(cliVersion string) (keptnVersionInfo, error) {
	v, err := client.getVersionInfo(cliVersion)
	return v.KeptnVersionInfo, err
}

func (client *versionFetcherClient) getVersionInfo(cliVersion string) (versionInfo, error) {
	versionInfo := versionInfo{}
	req, err := http.NewRequest("GET", client.versionUrl, nil)
	if err != nil {
		return versionInfo, err
	}
	req.Header.Set("user-agent", "keptn/cli:"+cliVersion)
	resp, err := client.httpClient.Do(req)
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
