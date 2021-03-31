package cmd

import (
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
