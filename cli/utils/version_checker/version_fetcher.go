package version_checker

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

const versionURL = "https://get.keptn.sh/version.json"

type Client struct {
	httpClient *http.Client
	versionUrl string
}

func newClient() *Client {
	client := Client{
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
		versionUrl: versionURL,
	}
	return &client
}

type VersionInfo struct {
	CLIVersionInfo CLIVersionInfo `json:"cli"`
}

type CLIVersionInfo struct {
	StableVersions []string `json:"stable_versions"`
	BetaVersions   []string `json:"beta_versions"`
}

func (client *Client) GetCLIVersionInfo(cliVersion string) (*CLIVersionInfo, error) {

	versionInfo := &VersionInfo{}
	req, err := http.NewRequest("GET", client.versionUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("user-agent", "KeptnCLI/"+cliVersion)
	resp, err := client.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, versionInfo)
	return &versionInfo.CLIVersionInfo, err
}
