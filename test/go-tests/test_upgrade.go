package go_tests

import (
	"context"
	"fmt"
	"github.com/imroc/req"
	"github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
	"k8s.io/apimachinery/pkg/util/rand"
	"net/http"
	"os"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

type TriggeredSequence struct {
	keptnContext string
	stage        string
	sequenceName string
}

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
	projectName := "upgrade-zero-downtime5"
	serviceName := "my-service"
	//sequenceName := "evaluation"

	nrOfUpgrades := 2

	nrStages := 20

	shipyard := &v0_2_0.Shipyard{
		ApiVersion: "0.2.3",
		Kind:       "shipyard",
		Metadata:   v0_2_0.Metadata{},
		Spec: v0_2_0.ShipyardSpec{
			Stages: []v0_2_0.Stage{},
		},
	}

	for i := 0; i < nrStages; i++ {
		newStageName := fmt.Sprintf("dev-%d", i)
		newStage := v0_2_0.Stage{
			Name: newStageName,
			Sequences: []v0_2_0.Sequence{
				{
					Name: "hooks",
					Tasks: []v0_2_0.Task{
						{
							Name: "mytask",
						},
					},
				},
			},
		}

		shipyard.Spec.Stages = append(shipyard.Spec.Stages, newStage)
	}

	shipyardFileContent, _ := yaml.Marshal(shipyard)

	shipyardFile, err := CreateTmpShipyardFile(string(shipyardFileContent))
	require.Nil(t, err)
	defer func() {
		err := os.Remove(shipyardFile)
		if err != nil {
			t.Logf("Could not delete file: %s: %v", shipyardFile, err)
		}
	}()

	chartLatestVersion := "https://github.com/keptn/helm-charts-dev/blob/gh-pages/packages/keptn-0.14.0-dev-PR-7266.tgz?raw=true"
	chartPreviousVersion := "https://github.com/keptn/helm-charts-dev/blob/0a8347beae62a203bf0bf534cde716314a2e51c2/packages/keptn-0.14.0-dev-PR-7266.tgz?raw=true"

	projectName, err = CreateProject(projectName, shipyardFile)
	require.Nil(t, err)

	output, err := ExecuteCommand(fmt.Sprintf("keptn create service %s --project=%s", serviceName, projectName))

	require.Nil(t, err)
	require.Contains(t, output, "created successfully")

	// set up webhook

	// create subscriptions for the webhook-service
	taskTypes := []string{"mytask"}

	webhookYamlWithSubscriptionIDs := webhookSimpleYaml
	webhookYamlWithSubscriptionIDs = getWebhookYamlWithSubscriptionIDs(t, taskTypes, projectName, webhookYamlWithSubscriptionIDs)

	// wait some time to make sure the webhook service has pulled the updated subscription
	<-time.After(20 * time.Second) // sorry :(

	// now, let's add a webhook.yaml file to our service
	webhookFilePath, err := CreateTmpFile("webhook.yaml", webhookYamlWithSubscriptionIDs)
	require.Nil(t, err)
	defer func() {
		err := os.Remove(webhookFilePath)
		if err != nil {
			t.Logf("Could not delete tmp file: %s", err.Error())
		}
	}()

	t.Log("Adding webhook.yaml to our service")
	_, err = ExecuteCommand(fmt.Sprintf("keptn add-resource --project=%s --service=%s --resource=%s --resourceUri=webhook/webhook.yaml --all-stages", projectName, serviceName, webhookFilePath))

	require.Nil(t, err)

	// test api endpoints
	httpEndpointTests := []*HTTPEndpointTest{
		{
			URL:                        "/controlPlane/v1/uniform/registration",
			Method:                     http.MethodGet,
			Payload:                    nil,
			ExpectedStatus:             200,
			WaitSecondsBetweenRequests: 1 * time.Second,
		},
		{
			URL:                        "/v1/metadata",
			Method:                     http.MethodGet,
			Payload:                    nil,
			ExpectedStatus:             200,
			WaitSecondsBetweenRequests: 1 * time.Second,
		},
		{
			URL:                        "/mongodb-datastore/event?project=" + projectName,
			Method:                     http.MethodGet,
			Payload:                    nil,
			ExpectedStatus:             200,
			WaitSecondsBetweenRequests: 1 * time.Second,
		},
		{
			URL:                        "/configuration-service/v1/project/" + projectName + "/resource",
			Method:                     http.MethodGet,
			Payload:                    nil,
			ExpectedStatus:             200,
			WaitSecondsBetweenRequests: 1 * time.Second,
		},
		{
			URL:                        "/secrets/v1/secret",
			Method:                     http.MethodGet,
			Payload:                    nil,
			ExpectedStatus:             200,
			WaitSecondsBetweenRequests: 1 * time.Second,
		},
	}

	ctx, cancel := context.WithCancel(context.TODO())

	// periodically execute upgrades between different versions

	upgradesWaitGroup := &sync.WaitGroup{}

	upgradesWaitGroup.Add(1)

	go func() {
		for i := 0; i < nrOfUpgrades; i++ {
			chartURL := ""
			var err error
			if i%2 == 0 {
				chartURL = chartLatestVersion
				//_, err = ExecuteCommand(fmt.Sprintf("kubectl -n %s set image deployment.v1.apps/lighthouse-service lighthouse-service=keptndev/lighthouse-service:0.14.0-dev", GetKeptnNameSpaceFromEnv()))
			} else {
				chartURL = chartPreviousVersion
				//_, err = ExecuteCommand(fmt.Sprintf("kubectl -n %s set image deployment.v1.apps/lighthouse-service lighthouse-service=keptndev/lighthouse-service:0.14.0-dev-PR-7266.202203280650", GetKeptnNameSpaceFromEnv()))
			}
			t.Logf("Upgrading Keptn to %s", chartURL)
			_, err = ExecuteCommand(fmt.Sprintf("helm upgrade -n %s keptn %s --wait --set=control-plane.apiGatewayNginx.type=LoadBalancer --set=control-plane.common.strategy.rollingUpdate.maxUnavailable=0 --set control-plane.resourceService.enabled=true --set control-plane.resourceService.env.DIRECTORY_STAGE_STRUCTURE=true", GetKeptnNameSpaceFromEnv(), chartURL))
			if err != nil {
				t.Logf("Encountered error when upgrading keptn: %v", err)
			}
			<-time.After(5 * time.Second)
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
	triggeredSequences := []TriggeredSequence{}
	go func() {
		for {
			select {
			case <-ctx.Done():
				triggerSequenceWaitGroup.Done()
				t.Log("Finished triggering sequences")
				return
			default:
				stageNr := nrTriggeredSequences % nrStages
				sequenceStageName := fmt.Sprintf("dev-%d", stageNr)
				var keptnContext string
				var err error
				// trigger an evaluation sequence
				keptnContext, err = TriggerSequence(projectName, serviceName, sequenceStageName, "evaluation", nil)
				nrTriggeredSequences++
				if err == nil && keptnContext != "" {
					triggeredSequences = append(triggeredSequences, TriggeredSequence{
						keptnContext: keptnContext,
						stage:        sequenceStageName,
						sequenceName: "evaluation",
					})
				} else {
					if err != nil {
						t.Logf("Could not trigger evaluation sequence: %v", err)
					} else {
						t.Log("Could not trigger evaluation sequence: did not get keptnContext")
					}
				}
				// trigger a webhook sequence
				//keptnContext, err = TriggerSequence(projectName, serviceName, sequenceStageName, "hooks", nil)
				//nrTriggeredSequences++
				//if err == nil {
				//	triggeredSequences = append(triggeredSequences, TriggeredSequence{
				//		keptnContext: keptnContext,
				//		stage:        sequenceStageName,
				//		sequenceName: "hooks",
				//	})
				//} else {
				//	t.Logf("Could not trigger hooks sequence: %v", err)
				//}
				// wait some time before triggering the next sequence
				<-time.After(time.Duration(100+rand.Intn(900)) * time.Millisecond)
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
	assert.Equal(t, len(triggeredSequences), nrTriggeredSequences)

	t.Logf("Triggered %d sequences. Let's check if they have been finished", nrTriggeredSequences)
	var nrFinishedSequences uint64

	checkSequencesWg := &sync.WaitGroup{}

	checkSequencesWg.Add(len(triggeredSequences))
	for _, triggeredSequence := range triggeredSequences {
		go func(sequence TriggeredSequence) {
			var sequenceFinishedEvent *models.KeptnContextExtendedCE

			stageSequenceName := fmt.Sprintf("%s.%s", sequence.stage, sequence.sequenceName)
			assert.Eventually(t, func() bool {
				sequenceFinishedEvent, err = GetLatestEventOfType(sequence.keptnContext, projectName, sequence.stage, v0_2_0.GetFinishedEventType(stageSequenceName))
				if sequenceFinishedEvent == nil || err != nil {
					return false
				}
				return true
			}, 10*time.Minute, 10*time.Second)
			if sequenceFinishedEvent != nil {
				atomic.AddUint64(&nrFinishedSequences, 1)
			} else {
				t.Logf("Sequence %s in stage %s has not been finished", sequence.keptnContext, sequence.stage)
			}
			checkSequencesWg.Done()
		}(triggeredSequence)
	}

	checkSequencesWg.Wait()
	t.Logf("Finished sequences: %d/%d", nrFinishedSequences, nrTriggeredSequences)
}
