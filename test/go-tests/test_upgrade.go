package go_tests

import (
	"context"
	"fmt"
	"github.com/imroc/req"
	"github.com/keptn/go-utils/pkg/api/models"
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

func (t *HTTPEndpointTest) Run(ctx context.Context, wg *sync.WaitGroup, tst *testing.T) error {
	apiCaller, err := NewAPICaller()
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			wg.Done()
			tst.Logf("Finished tests for Endpoint %s", t.URL)
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
				tst.Logf("HTTP %s request to %s failed", t.Method, t.URL)
				t.Result.FailedRequests++
			}
			t.NrRequests++
			<-time.After(t.WaitSecondsBetweenRequests)
		}
	}
}

func (t *HTTPEndpointTest) String() string {
	failureRate := float64(t.Result.FailedRequests) / float64(t.NrRequests)
	return fmt.Sprintf("\n======\nURL: %s\nExecutedRequests: %d\n FailedRequests: %d\n FailureRate: %f\n======\n", t.URL, t.NrRequests, t.Result.FailedRequests, failureRate)
}

type HTTPEndpointTestResult struct {
	FailedRequests int
}

func Test_UpgradeZeroDowntime(t *testing.T) {
	projectName := "upgrade-zero-downtime2"
	serviceName := "my-service"
	stageName := "dev"
	//sequenceName := "evaluation"

	nrOfUpgrades := 2

	shipyardFile, err := CreateTmpShipyardFile(zeroDownTimeShipyard)
	require.Nil(t, err)
	defer func() {
		err := os.Remove(shipyardFile)
		if err != nil {
			t.Logf("Could not delete file: %s: %v", shipyardFile, err)
		}
	}()

	chartLatestVersion := "https://github.com/keptn/helm-charts-dev/blob/gh-pages/packages/keptn-0.14.0-dev-PR-7266.tgz?raw=true"
	chartPreviousVersion := "https://github.com/keptn/helm-charts-dev/blob/f098f43b33540d5bfc822a0038c0c21ccffe9335/packages/keptn-0.14.0-dev-PR-7266.tgz?raw=true"

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
			WaitSecondsBetweenRequests: 3 * time.Second,
		},
		{
			URL:                        "/v1/metadata",
			Method:                     http.MethodGet,
			Payload:                    nil,
			ExpectedStatus:             200,
			WaitSecondsBetweenRequests: 3 * time.Second,
		},
		{
			URL:                        "/mongodb-datastore/event?project=" + projectName,
			Method:                     http.MethodGet,
			Payload:                    nil,
			ExpectedStatus:             200,
			WaitSecondsBetweenRequests: 3 * time.Second,
		},
		{
			URL:                        "/configuration-service/v1/project/" + projectName + "/resource",
			Method:                     http.MethodGet,
			Payload:                    nil,
			ExpectedStatus:             200,
			WaitSecondsBetweenRequests: 3 * time.Second,
		},
		{
			URL:                        "/secrets/v1/secret",
			Method:                     http.MethodGet,
			Payload:                    nil,
			ExpectedStatus:             200,
			WaitSecondsBetweenRequests: 3 * time.Second,
		},
	}

	ctx, cancel := context.WithCancel(context.TODO())

	// periodically execute upgrades between different versions

	upgradesWaitGroup := &sync.WaitGroup{}

	upgradesWaitGroup.Add(1)

	go func() {
		for i := 0; i < nrOfUpgrades; i++ {
			chartURL := ""
			if i%2 == 0 {
				chartURL = chartLatestVersion
			} else {
				chartURL = chartPreviousVersion
			}
			t.Logf("Upgrading Keptn to %s", chartURL)
			_, err := ExecuteCommand(fmt.Sprintf("helm upgrade -n %s keptn %s --wait --set=control-plane.apiGatewayNginx.type=LoadBalancer --set control-plane.resourceService.enabled=true", GetKeptnNameSpaceFromEnv(), chartURL))
			if err != nil {
				t.Logf("Encountered error when upgrading keptn: %v", err)
			}
		}
		upgradesWaitGroup.Done()
	}()

	endpointTestsWaitGroup := &sync.WaitGroup{}
	endpointTestsWaitGroup.Add(len(httpEndpointTests))

	for index := range httpEndpointTests {
		endpointTest := httpEndpointTests[index]
		go func() {
			err := endpointTest.Run(ctx, endpointTestsWaitGroup, t)
			if err != nil {
				t.Logf("encountered error: %v", err)
			}
		}()
	}

	// trigger sequences -> keep track of how many sequences we have triggered

	triggerSequenceWaitGroup := &sync.WaitGroup{}
	triggerSequenceWaitGroup.Add(1)
	nrTriggeredSequences := 0
	keptnContextIDs := []string{}
	go func() {
		for {
			select {
			case <-ctx.Done():
				triggerSequenceWaitGroup.Done()
				t.Log("Finished triggering sequences")
				return
			default:
				keptnContext, err := TriggerSequence(projectName, serviceName, stageName, "evaluation", nil)
				nrTriggeredSequences++
				if err == nil {
					keptnContextIDs = append(keptnContextIDs, keptnContext)
					t.Logf("Triggered new evaluation sequence with KeptnContext %s", keptnContext)
				}
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

	t.Logf("Triggered %d sequences. Let's check if they have been finished", nrTriggeredSequences)
	nrFinishedSequences := 0
	for _, keptnContext := range keptnContextIDs {
		t.Logf("Checking if sequence %s has been finished", keptnContext)
		var evaluationFinishedEvent *models.KeptnContextExtendedCE
		assert.Eventually(t, func() bool {
			evaluationFinishedEvent, err = GetLatestEventOfType(keptnContext, projectName, stageName, v0_2_0.GetFinishedEventType("dev.evaluation"))
			if evaluationFinishedEvent == nil || err != nil {
				return false
			}
			return true
		}, 10*time.Minute, 10*time.Second)
		if evaluationFinishedEvent != nil {
			nrFinishedSequences++
		}
	}
	t.Logf("Finished sequences: %d/%d", nrFinishedSequences, nrTriggeredSequences)
}
