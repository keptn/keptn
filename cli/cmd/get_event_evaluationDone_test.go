package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/keptn/keptn/cli/pkg/logging"
)

const event1 = `{
  "contenttype": "application/json",
  "data": {
    "deploymentstrategy": "",
    "evaluationdetails": {
      "comparedEvents": [
        "3552946d-20c6-4f2d-8137-2dde96903dbc",
        "61342db8-5775-4b17-b7fe-0968bd4b7d0c"
      ],
      "indicatorResults": [
        {
          "score": 1,
          "status": "pass",
          "targets": [
            {
              "criteria": "<=+75%",
              "targetValue": 57.87294070918944,
              "violated": false
            },
            {
              "criteria": "<7500",
              "targetValue": 7500,
              "violated": false
            }
          ],
          "value": {
            "metric": "response_time_p95",
            "success": true,
            "value": 32.844805419323706
          }
        }
      ],
      "result": "pass",
      "score": 100,
      "sloFileContent": "LS0tCnNwZWNfdmVyc2lvbjogIjAuMS4xIgpjb21wYXJpc29uOgogIGFnZ3JlZ2F0ZV9mdW5jdGlvbjogImF2ZyIKICBjb21wYXJlX3dpdGg6ICJzZXZlcmFsX3Jlc3VsdHMiCiAgaW5jbHVkZV9yZXN1bHRfd2l0aF9zY29yZTogInBhc3MiCiAgbnVtYmVyX29mX2NvbXBhcmlzb25fcmVzdWx0czogMgpmaWx0ZXI6Cm9iamVjdGl2ZXM6CiAgLSBzbGk6ICJyZXNwb25zZV90aW1lX3A5NSIKICAgIGtleV9zbGk6IGZhbHNlCiAgICBwYXNzOiAgICAgICAgICAgICAjIHBhc3MgaWYgKHJlbGF0aXZlIGNoYW5nZSA8PSA3NSUgQU5EIGFic29sdXRlIHZhbHVlIGlzIDwgNzVtcykKICAgICAgLSBjcml0ZXJpYToKICAgICAgICAgIC0gIjw9Kzc1JSIgICMgcmVsYXRpdmUgdmFsdWVzIHJlcXVpcmUgYSBwcmVmaXhlZCBzaWduIChwbHVzIG9yIG1pbnVzKQogICAgICAgICAgLSAiPDc1MDAiICAgICAjIGFic29sdXRlIHZhbHVlcyBvbmx5IHJlcXVpcmUgYSBsb2dpY2FsIG9wZXJhdG9yCiAgICB3ZWlnaHQ6IDEKdG90YWxfc2NvcmU6CiAgcGFzczogIjkwJSIKICB3YXJuaW5nOiAiNzUlIg==",
      "timeEnd": "2020-11-12T15:25:10.000Z",
      "timeStart": "2020-11-12T15:20:10.000Z"
    },
    "labels": {
      "DtCreds": "dynatrace"
    },
    "project": "musicshop",
    "result": "pass",
    "service": "frontend",
    "stage": "hardening",
    "teststrategy": "manual"
  },
  "id": "03ae2098-55de-4990-953f-e4e3c4f58d34",
  "shkeptncontext": "8929e5e5-3826-488f-9257-708bfa974909",
  "source": "some-service",
  "specversion": "0.2",
  "time": "2020-11-12T15:25:10.600Z",
  "type": "sh.keptn.events.evaluation.finished"
}`

const responseEnvelope = `{
    "events": [
       %s
    ],
	"nextPageKey": "0",
    "pageSize": 1,
    "totalCount": 1
}`

var response = fmt.Sprintf(responseEnvelope, event1)

func init() {
	logging.InitLoggers(os.Stdout, os.Stdout, os.Stderr)
}

type fakeCredentialManager struct {
	fakeServer httptest.Server
}

func (f fakeCredentialManager) SetCreds(endPoint url.URL, apiToken string, namespace string) error {
	panic("implement me")
}

func (f fakeCredentialManager) GetCreds(namespace string) (url.URL, string, error) {

	url, _ := url.Parse(f.fakeServer.URL)
	return *url, "", nil

}

func (f fakeCredentialManager) SetInstallCreds(creds string) error {
	panic("implement me")
}

func (f fakeCredentialManager) GetInstallCreds() (string, error) {
	panic("implement me")
}

func (f fakeCredentialManager) getLinuxApiTokenFile(namespace string) string {
	panic("implement me")
}

// TestEvaluationDoneGetEvent tests the evaluation-done command
func TestEvaluationDoneGetEvent(t *testing.T) {

	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Content-Type", "application/json")

			if strings.Contains(r.RequestURI, "/event") {
				w.WriteHeader(200)
				w.Write([]byte(response))
				return
			}

			w.WriteHeader(500)
			w.Write([]byte(`{}`))
		}),
	)
	defer ts.Close()

	credentialManager = &fakeCredentialManager{fakeServer: *ts}
	out = new(bytes.Buffer)

	err := do(nil, nil)

	var want interface{}
	json.Unmarshal([]byte(event1), &want)
	var got interface{}
	json.Unmarshal([]byte(out.(*bytes.Buffer).String()), &got)

	assert.Nil(t, err)
	assert.Equal(t, want, got)

}
