package go_tests

import (
	"context"
	"fmt"
	"github.com/imroc/req"
	"github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/util/rand"
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

func (t *HTTPEndpointTest) Run(ctx context.Context, wg *sync.WaitGroup) error {
	apiCaller, err := NewAPICaller()
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			wg.Done()
			return nil
		default:
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
			t.NrRequests++
			<-time.After(t.WaitSecondsBetweenRequests)
		}
	}
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
	stageName := "dev"
	//sequenceName := "evaluation"
	shipyardFile, err := CreateTmpShipyardFile(zeroDownTimeShipyard)
	require.Nil(t, err)
	defer func() {
		err := os.Remove(shipyardFile)
		if err != nil {
			t.Logf("Could not delete file: %s: %v", shipyardFile, err)
		}
	}()

	chartLatestVersion := "https://github.com/keptn/helm-charts-dev/blob/gh-pages/packages/keptn-0.14.0-dev.tgz?raw=true"
	chartPreviousVersion := "https://github.com/keptn/helm-charts-dev/blob/027ebfdf98176047fdcf6c80f8aa9599a9c66b4e/packages/keptn-0.14.0-dev.tgz?raw=true"

	// check if the project 'state' is already available - if not, delete it before creating it again
	// check if the project is already available - if not, delete it before creating it again
	projectName, err = CreateProject(projectName, shipyardFile)
	require.Nil(t, err)

	output, err := ExecuteCommand(fmt.Sprintf("keptn create service %s --project=%s", serviceName, projectName))

	require.Nil(t, err)
	require.Contains(t, output, "created successfully")

	// test api endpoints
	httpEndpointTests := []*HTTPEndpointTest{
		{
			URL:                        "/controlPlane/v1/uniform/registration",
			Method:                     http.MethodGet,
			Payload:                    nil,
			ExpectedStatus:             200,
			WaitSecondsBetweenRequests: 3,
		},
		{
			URL:                        "/v1/metadata",
			Method:                     http.MethodGet,
			Payload:                    nil,
			ExpectedStatus:             200,
			WaitSecondsBetweenRequests: 3,
		},
		{
			URL:                        "/mongodb-datastore/event?project=" + projectName,
			Method:                     http.MethodGet,
			Payload:                    nil,
			ExpectedStatus:             200,
			WaitSecondsBetweenRequests: 3,
		},
		{
			URL:                        "/configuration-service/v1/project/" + projectName + "/resource",
			Method:                     http.MethodGet,
			Payload:                    nil,
			ExpectedStatus:             200,
			WaitSecondsBetweenRequests: 3,
		},
		{
			URL:                        "/secrets/v1/secret",
			Method:                     http.MethodGet,
			Payload:                    nil,
			ExpectedStatus:             200,
			WaitSecondsBetweenRequests: 3,
		},
	}

	ctx, cancel := context.WithCancel(context.TODO())

	// periodically execute upgrades between different versions

	upgradesWaitGroup := &sync.WaitGroup{}

	upgradesWaitGroup.Add(1)
	nrOfUpgrades := 3
	go func() {
		for i := 0; i < nrOfUpgrades; i++ {
			chartURL := ""
			if i%2 == 0 {
				chartURL = chartLatestVersion
			} else {
				chartURL = chartPreviousVersion
			}
			_, err := ExecuteCommand(fmt.Sprintf("helm upgrade -n %s keptn %s --wait --set=\"control-plane.apiGatewayNginx.type=LoadBalancer", GetKeptnNameSpaceFromEnv(), chartURL))
			if err != nil {
				t.Logf("Encountered error when upgrading keptn: %v", err)
			}
			// wait for a random number of seconds (0-10s)  before performing the next update
			<-time.After(time.Duration(rand.Intn(10)) * time.Second)
		}
		upgradesWaitGroup.Done()
	}()

	endpointTestsWaitGroup := &sync.WaitGroup{}
	endpointTestsWaitGroup.Add(len(httpEndpointTests))

	for index := range httpEndpointTests {
		endpointTest := httpEndpointTests[index]
		go func() {
			err := endpointTest.Run(ctx, endpointTestsWaitGroup)
			if err != nil {
				t.Logf("encountered error: %v", err)
			}
		}()
	}

	// trigger sequences -> keep track of how many sequences we have triggered

	triggerSequenceWaitGroup := &sync.WaitGroup{}
	nrTriggeredSequences := 0
	keptnContextIDs := []string{}
	go func() {
		for {
			select {
			case <-ctx.Done():
				triggerSequenceWaitGroup.Done()
				return
			default:
				keptnContext, _ := TriggerSequence(projectName, serviceName, stageName, "evaluation", nil)
				nrTriggeredSequences++
				keptnContextIDs = append(keptnContextIDs, keptnContext)
				// wait some time before triggering the next sequence
				<-time.After(time.Duration(rand.Intn(10)) * time.Second)
			}
		}
	}()

	// wait for the upgrades to be finished
	upgradesWaitGroup.Wait()
	cancel()

	// wait for the endpoint tests to wrap up
	endpointTestsWaitGroup.Wait()

	// wait for the sequences to wrap up
	triggerSequenceWaitGroup.Wait()

	// get the results of the endpoint tests
	for _, endpointTest := range httpEndpointTests {
		t.Log(endpointTest.String())
		assert.Zero(t, endpointTest.Result.FailedRequests)
	}

	// get the number of dev.delivery.finished events -> this should eventually match the number of
	assert.Equal(t, len(keptnContextIDs), nrTriggeredSequences)

	for _, keptnContext := range keptnContextIDs {
		assert.Eventually(t, func() bool {
			evaluationFinishedEvent, err := GetLatestEventOfType(keptnContext, projectName, stageName, v0_2_0.GetFinishedEventType("dev.evaluation"))
			if evaluationFinishedEvent == nil || err != nil {
				return false
			}
			return true
		}, 2*time.Minute, 10*time.Second)
	}

}
