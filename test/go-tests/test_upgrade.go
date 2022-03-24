package go_tests

import (
	"fmt"
	"github.com/imroc/req"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"os"
	"sync"
	"testing"
	"time"
)

type HTTPEndpointTest struct {
	URL                        string
	Method                     string
	Payload                    interface{}
	ExpectedStatus             int
	NrRequests                 int
	WaitSecondsBetweenRequests time.Duration
	Result                     HTTPEndpointTestResult
}

func (t *HTTPEndpointTest) Run(wg *sync.WaitGroup) error {
	apiCaller, err := NewAPICaller()
	if err != nil {
		wg.Done()
		return err
	}
	t.Result = HTTPEndpointTestResult{}
	for i := 0; i < t.NrRequests; i++ {
		var resp *req.Resp
		var err error
		switch t.Method {
		case http.MethodGet:
			resp, err = apiCaller.Get(t.URL, 0)
		case http.MethodPost:
			resp, err = apiCaller.Post(t.URL, t.Payload, 0)
		}
		if err != nil || resp.Response().StatusCode != t.ExpectedStatus {
			t.Result.FailedRequests++
		}
		<-time.After(t.WaitSecondsBetweenRequests)
	}
	wg.Done()
	return nil
}

func (t *HTTPEndpointTest) String() string {
	failureRate := t.Result.FailedRequests / t.NrRequests
	return fmt.Sprintf("\n======\nURL: %s\nExecutedRequests: %d\n FailedRequests: %d\n FailureRate: %d\n======\n", t.URL, t.NrRequests, t.Result.FailedRequests, failureRate)
}

type HTTPEndpointTestResult struct {
	FailedRequests int
}

func Test_UpgradeZeroDowntime(t *testing.T) {
	projectName := "upgrade-zero-downtime"
	serviceName := "my-service"
	//stageName := "dev"
	//sequenceName := "evaluation"
	shipyardFile, err := CreateTmpShipyardFile(zeroDownTimeShipyard)
	require.Nil(t, err)
	defer func() {
		err := os.Remove(shipyardFile)
		if err != nil {
			t.Logf("Could not delete file: %s: %v", shipyardFile, err)
		}
	}()

	// check if the project 'state' is already available - if not, delete it before creating it again
	// check if the project is already available - if not, delete it before creating it again
	projectName, err = CreateProject(projectName, shipyardFile)
	require.Nil(t, err)

	output, err := ExecuteCommand(fmt.Sprintf("keptn create service %s --project=%s", serviceName, projectName))

	require.Nil(t, err)
	require.Contains(t, output, "created successfully")

	// test api endpoints
	requestsPerEndpoint := 1
	httpEndpointTests := []*HTTPEndpointTest{
		{
			URL:                        "/controlPlane/v1/uniform/registration",
			Method:                     http.MethodGet,
			Payload:                    nil,
			ExpectedStatus:             200,
			NrRequests:                 requestsPerEndpoint,
			WaitSecondsBetweenRequests: 3,
		},
		{
			URL:                        "/v1/metadata",
			Method:                     http.MethodGet,
			Payload:                    nil,
			ExpectedStatus:             200,
			NrRequests:                 requestsPerEndpoint,
			WaitSecondsBetweenRequests: 3,
		},
		{
			URL:                        "/mongodb-datastore/event?project=" + projectName,
			Method:                     http.MethodGet,
			Payload:                    nil,
			ExpectedStatus:             200,
			NrRequests:                 requestsPerEndpoint,
			WaitSecondsBetweenRequests: 3,
		},
		{
			URL:                        "/configuration-service/v1/project/" + projectName + "/resource",
			Method:                     http.MethodGet,
			Payload:                    nil,
			ExpectedStatus:             200,
			NrRequests:                 requestsPerEndpoint,
			WaitSecondsBetweenRequests: 3,
		},
		{
			URL:                        "/secrets/v1/secret",
			Method:                     http.MethodGet,
			Payload:                    nil,
			ExpectedStatus:             200,
			NrRequests:                 requestsPerEndpoint,
			WaitSecondsBetweenRequests: 3,
		},
	}

	wg := &sync.WaitGroup{}

	wg.Add(len(httpEndpointTests))

	// TODO periodically execute upgrades between different versions

	for index := range httpEndpointTests {
		endpointTest := httpEndpointTests[index]
		go func() {
			err := endpointTest.Run(wg)
			if err != nil {
				t.Logf("encountered error: %v", err)
			}
		}()
	}

	// TODO trigger evaluation sequences

	wg.Wait()

	for _, endpointTest := range httpEndpointTests {
		t.Log(endpointTest.String())
		assert.Zero(t, endpointTest.Result.FailedRequests)
	}

}
