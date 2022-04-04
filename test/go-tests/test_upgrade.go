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

type TriggeredSequences struct {
	sequences []TriggeredSequence
	mutex     *sync.Mutex
}

func (ts *TriggeredSequences) Add(s TriggeredSequence) {
	ts.mutex.Lock()
	defer ts.mutex.Unlock()
	ts.sequences = append(ts.sequences, s)
}

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
	NrRequests                 int64
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
			go func() {
				var resp *req.Resp
				var err error
				switch t.Method {
				case http.MethodGet:
					resp, err = apiCaller.Get(t.URL, 0)
				case http.MethodPost:
					resp, err = apiCaller.Post(t.URL, t.Payload, 0)
				}

				if err != nil || resp.Response().StatusCode != t.ExpectedStatus {
					tst.Logf("HTTP %s request to %s failed: %v", t.Method, t.URL, err)
					atomic.AddInt64(&t.Result.FailedRequests, 1)
				}
				atomic.AddInt64(&t.NrRequests, 1)
			}()

			<-time.After(t.WaitSecondsBetweenRequests)
		}
	}
}

func (t *HTTPEndpointTest) String() string {
	failureRate := float64(t.Result.FailedRequests) / float64(t.NrRequests)
	return fmt.Sprintf("\n======\nURL: %s\nExecutedRequests: %d\n FailedRequests: %d\n FailureRate: %f\n======\n", t.URL, t.NrRequests, t.Result.FailedRequests, failureRate)
}

type HTTPEndpointTestResult struct {
	FailedRequests int64
}

func Test_UpgradeZeroDowntime(t *testing.T) {
	projectName := "upgrade-zero-downtime6"
	serviceName := "my-service"
	//sequenceName := "evaluation"

	nrOfUpgrades := 2

	nrStages := int64(10)

	shipyard := &v0_2_0.Shipyard{
		ApiVersion: "0.2.3",
		Kind:       "shipyard",
		Metadata:   v0_2_0.Metadata{},
		Spec: v0_2_0.ShipyardSpec{
			Stages: []v0_2_0.Stage{},
		},
	}

	for i := int64(0); i < nrStages; i++ {
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

	chartLatestVersion := "https://github.com/keptn/helm-charts-dev/blob/b5d01b0a4f42404abee23a031fb5a9e693a57486/packages/keptn-0.14.1-dev-PR-7266.tgz?raw=true"
	chartPreviousVersion := "https://github.com/keptn/helm-charts-dev/blob/087273a72ee19dfb71d766ccdc6ebfb3a5ef5dec/packages/keptn-0.14.1-dev-PR-7266.tgz?raw=true"

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

	waitBetweenRequests := 300 * time.Millisecond
	// test api endpoints
	httpEndpointTests := []*HTTPEndpointTest{
		{
			URL:                        "/controlPlane/v1/uniform/registration",
			Method:                     http.MethodGet,
			Payload:                    nil,
			ExpectedStatus:             200,
			WaitSecondsBetweenRequests: waitBetweenRequests,
		},
		{
			URL:                        "/v1/metadata",
			Method:                     http.MethodGet,
			Payload:                    nil,
			ExpectedStatus:             200,
			WaitSecondsBetweenRequests: waitBetweenRequests,
		},
		{
			URL:                        "/mongodb-datastore/event?project=" + projectName,
			Method:                     http.MethodGet,
			Payload:                    nil,
			ExpectedStatus:             200,
			WaitSecondsBetweenRequests: waitBetweenRequests,
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
			WaitSecondsBetweenRequests: waitBetweenRequests,
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
			//_, err = ExecuteCommand(fmt.Sprintf("helm upgrade -n %s keptn %s --wait --set=control-plane.apiGatewayNginx.type=LoadBalancer --set=control-plane.common.strategy.rollingUpdate.maxUnavailable=0 --set control-plane.resourceService.enabled=true --set control-plane.resourceService.env.DIRECTORY_STAGE_STRUCTURE=true", GetKeptnNameSpaceFromEnv(), chartURL))
			_, err = ExecuteCommand(fmt.Sprintf("helm upgrade -n %s keptn %s --wait --set=control-plane.apiGatewayNginx.type=LoadBalancer --set=control-plane.common.strategy.rollingUpdate.maxUnavailable=0 --set control-plane.resourceService.enabled=true --set control-plane.resourceService.env.DIRECTORY_STAGE_STRUCTURE=true --set control-plane.distributor.image.repository=docker.io/keptndev/distributor --set control-plane.distributor.image.tag=0.14.1-dev-PR-7308.202204010740", GetKeptnNameSpaceFromEnv(), chartURL))
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
	var nrTriggeredSequences int64
	triggeredSequences := &TriggeredSequences{
		mutex: &sync.Mutex{},
	}

	triggerSequence := func(sequenceName, sequenceStageName string) {
		var keptnContext string
		var err error
		// trigger an evaluation sequence
		keptnContext, err = TriggerSequence(projectName, serviceName, sequenceStageName, sequenceName, nil)
		atomic.AddInt64(&nrTriggeredSequences, 1)

		if err == nil && keptnContext != "" {
			triggeredSequences.Add(TriggeredSequence{
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
	}

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
				//go func(stage string) {
				//triggerSequence("evaluation", sequenceStageName)
				//}(sequenceStageName)
				// trigger a webhook sequence
				//go func(stage string) {
				triggerSequence("hooks", sequenceStageName)
				//}(sequenceStageName)
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
	assert.Equal(t, int64(len(triggeredSequences.sequences)), nrTriggeredSequences)
	// TODO remove return again
	return
	t.Logf("Triggered %d sequences. Let's check if they have been finished", nrTriggeredSequences)
	var nrFinishedSequences uint64

	checkSequencesWg := &sync.WaitGroup{}

	checkSequencesWg.Add(len(triggeredSequences.sequences))
	for _, triggeredSequence := range triggeredSequences.sequences {
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
