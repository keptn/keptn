package git

import (
	"encoding/json"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"io/ioutil"
	"net/http"
	"strings"
)

type gitHubApi []struct {
	Ref    string `json:"ref"`
	NodeID string `json:"node_id"`
	URL    string `json:"url"`
	Object struct {
		Sha  string `json:"sha"`
		Type string `json:"type"`
		URL  string `json:"url"`
	} `json:"object"`
}

func CloneGitHubUrl(filePath string, url string, tag string) error {
	_, err := git.PlainClone(filePath, false, &git.CloneOptions{
		URL:           url,
		ReferenceName: plumbing.ReferenceName(tag),
	})
	if err != nil {
		return err
	}
	return nil
}

func GetGitHubRefs(url string, version string) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	var gitHubApiResponse gitHubApi
	_ = json.Unmarshal(bodyBytes, &gitHubApiResponse)

	for _, refsTags := range gitHubApiResponse {
		if strings.Contains(refsTags.Ref, version) {
			return refsTags.Ref, nil
		}
	}
	return "", nil
}
