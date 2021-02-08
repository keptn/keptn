package cmd

import (
	"context"
	"crypto/tls"
	"net/http"
	"time"
)

var maxHTTPTimeout = 5 * time.Second

var endPointErrorReasons = `
Possible reasons:
* The Keptn API server is currently not available. Check if your Kubernetes cluster is available.
* Your Keptn CLI points to the wrong API server (verify using 'keptn status')`

var checkEndPointStatusMock = false

var client = http.Client{
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	},
}

func checkEndPointStatus(endPoint string) error {
	if checkEndPointStatusMock {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), maxHTTPTimeout)
	defer cancel()

	req, err := http.NewRequest("HEAD", endPoint, nil)
	if err != nil {
		return err
	}

	resp, err := client.Do(req.WithContext(ctx))
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}
