package ZD

import (
	"context"
	"fmt"
	"github.com/benbjohnson/clock"
	testutils "github.com/keptn/keptn/test/go-tests"
	"github.com/stretchr/testify/assert"
	"os"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

const apiProbeInterval = 1 * time.Second
const sequencesInterval = 10 * time.Second

var chartLatestVersion = "https://github.com/keptn/helm-charts-dev/blob/1efe3dab77da9ea3cf2b7dd5eff4b2fac6f76633/packages/keptn-0.15.0-dev-PR-7266.tgz?raw=true"
var chartPreviousVersion = "https://github.com/keptn/helm-charts-dev/blob/5b4fbc630895a2a71721763110376b452f4c2c67/packages/keptn-0.15.0-dev-PR-7266.tgz?raw=true"

var TotalAPICalls uint64 = 0
var FiredSequences uint64 = 0
var FailedSequences uint64 = 0
var PassedSequences uint64 = 0
var Id uint64 = 0

var wg = sync.WaitGroup{}

func TestZeroDowntime(t *testing.T) {

	t.Parallel()
	// run a tests before starting update
	go t.Run("Before upgrade sequences", Test_Run)

	apiTicker := clock.New().Ticker(apiProbeInterval)
	seqTicker := clock.New().Ticker(sequencesInterval)
	ctx, cancel := context.WithCancel(context.Background())

	//start update
	go t.Run("rolling update", func(t *testing.T) {
		rollingUpgrade(cancel, 2, t)
	})

	// run tests meanwhile
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-apiTicker.C:
				atomic.AddUint64(&TotalAPICalls, 1)
				go t.Run("API test", TestAPIs)
			}
		}
	}()

	for {
		select {
		case <-ctx.Done():
			wg.Wait()
			err := os.Remove(shipyardFile)
			assert.Nil(t, err)
			t.Run("Summary: ", printResults)
			return
		case <-seqTicker.C:
			wg.Add(2)
			go t.Run("Sequences test", Test_Run)

		}
	}

}

func printResults(t *testing.T) {
	t.Log("-----------------------------------------------")
	t.Log("TOTAL SEQUENCES: ", FiredSequences)
	t.Log("TOTAL SUCCESS ", PassedSequences)
	t.Log("TOTAL FAILURES ", FailedSequences)
	t.Log("-----------------------------------------------")
	t.Log("TOTAL API PROBES", TotalAPICalls)
}

func rollingUpgrade(cancel context.CancelFunc, nrOfUpgrades int, t *testing.T) {
	defer cancel()
	t.Log("Rolling")
	//time.Sleep(20 * time.Second)

	for i := 0; i < nrOfUpgrades; i++ {
		chartURL := ""
		var err error
		if i%2 == 0 {
			chartURL = chartLatestVersion
		} else {
			chartURL = chartPreviousVersion
		}
		t.Logf("Upgrading Keptn to %s", chartURL)
		_, err = testutils.ExecuteCommand(
			fmt.Sprintf(
				"helm upgrade -n %s keptn %s --wait --set=control-plane.apiGatewayNginx.type=LoadBalancer "+
					"--set=control-plane.common.strategy.rollingUpdate.maxUnavailable=0 --set control-plane.resourceService.enabled=true"+
					" --set control-plane.resourceService.env.DIRECTORY_STAGE_STRUCTURE=true", testutils.GetKeptnNameSpaceFromEnv(), chartURL))
		if err != nil {
			t.Logf("Encountered error when upgrading keptn: %v", err)

		}
		<-time.After(60 * time.Second)
	}

}
