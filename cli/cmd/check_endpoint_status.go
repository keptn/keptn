package cmd

import (
	"crypto/tls"
	"net/http"

	apiUtils "github.com/keptn/go-utils/pkg/api/utils"
)

var endPointErrorReasons = `
Possible reasons:
* The Keptn API server is currently not available. Check if your Kubernetes cluster is available.
* Your Keptn CLI points to the wrong API server (verify using 'keptn status')`

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
