package cmd

import (
	"crypto/tls"
	"net/http"

	apiUtils "github.com/keptn/go-utils/pkg/api/utils"
)

var checkEndPointStatusMock = false

var client = http.Client{
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		DialContext:     apiUtils.ResolveXipIoWithContext,
	},
}

func checkEndPointStatus(endPoint string) error {
	if checkEndPointStatusMock {
		return nil
	}
	req, err := http.NewRequest("HEAD", endPoint, nil)
	if err != nil {
		return err
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}
