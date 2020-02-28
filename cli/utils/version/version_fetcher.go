package version

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

const versionURL = "https://get.keptn.sh/version.json"

type versionInfo struct {
	CLIVersionInfo cliVersionInfo `json:"cli"`
}

type cliVersionInfo struct {
	Stable     []string `json:"stable"`
	Prerelease []string `json:"prerelease"`
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

func (client *versionFetcherClient) getCLIVersionInfo(cliVersion string) (*cliVersionInfo, error) {

	versionInfo := &versionInfo{}
	req, err := http.NewRequest("GET", client.versionUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("user-agent", "keptn/cli:"+cliVersion)
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
