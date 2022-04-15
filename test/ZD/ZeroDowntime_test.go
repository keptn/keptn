package ZD

import (
	"context"
	"fmt"
	"github.com/benbjohnson/clock"
	testutils "github.com/keptn/keptn/test/go-tests"
	"github.com/stretchr/testify/assert"
	"os"
	"sync"
	"testing"
	"time"
)

const apiProbeInterval = 1 * time.Second
const deliveryInterval = 10 * time.Second

var chartLatestVersion = "https://github.com/keptn/helm-charts-dev/blob/1efe3dab77da9ea3cf2b7dd5eff4b2fac6f76633/packages/keptn-0.15.0-dev-PR-7266.tgz?raw=true"
var chartPreviousVersion = "https://github.com/keptn/helm-charts-dev/blob/5b4fbc630895a2a71721763110376b452f4c2c67/packages/keptn-0.15.0-dev-PR-7266.tgz?raw=true"

var FailedSequences uint64 = 0
var PassedSequences uint64 = 0
var Id uint64 = 0

var wg = new(sync.WaitGroup)

func TestZeroDowntime(t *testing.T) {
	t.Parallel()
	// run a tests before starting update
	go t.Run("Triggering tests", Test_Evaluation)

	ticker := clock.New().Ticker(apiProbeInterval)
	ticker2 := clock.New().Ticker(deliveryInterval)
	ctx, cancel := context.WithCancel(context.Background())

	//start update
	//go t.Run("Rolling Updates", func(t *testing.T) {
	go rollingUpgrade(cancel, 2, t)
	//})

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				go t.Run("API test", TestAPIs)

			}
		}
	}()

	// run tests meanwhile
	for {
		select {
		case <-ctx.Done():
			err := os.Remove(shipyardFile)
			assert.Nil(t, err)

			wg.Wait()

			t.Log("Total run seq: ", PassedSequences+FailedSequences)
			t.Log("TOTAL SUCCESS ", PassedSequences)
			t.Log("TOTAL FAILURES ", FailedSequences)
			return
		case <-ticker2.C:
			go t.Run("Triggering tests", Test_Evaluation)

		}
	}

}

func rollingUpgrade(cancel context.CancelFunc, nrOfUpgrades int, t *testing.T) {
	defer cancel()
	t.Log("rolling")
	time.Sleep(40 * time.Second)

	//TODO setup helm upgrade here
	for i := 0; i < nrOfUpgrades; i++ {
		chartURL := ""
		//lighthouseVersion := ""
		var err error
		if i%2 == 0 {
			//lighthouseVersion = "v1"
			chartURL = chartLatestVersion
			//_, err = ExecuteCommand(fmt.Sprintf("kubectl -n %s set image deployment.v1.apps/lighthouse-service lighthouse-service=keptndev/lighthouse-service:0.14.0-dev", GetKeptnNameSpaceFromEnv()))
		} else {
			//lighthouseVersion = "v2"
			chartURL = chartPreviousVersion
			//_, err = ExecuteCommand(fmt.Sprintf("kubectl -n %s set image deployment.v1.apps/lighthouse-service lighthouse-service=keptndev/lighthouse-service:0.14.0-dev-PR-7266.202203280650", GetKeptnNameSpaceFromEnv()))
		}
		t.Logf("Upgrading Keptn to %s", chartURL)
		_, err = testutils.ExecuteCommand(fmt.Sprintf("helm upgrade -n %s keptn %s --wait --set=control-plane.apiGatewayNginx.type=LoadBalancer --set=control-plane.common.strategy.rollingUpdate.maxUnavailable=0 --set control-plane.resourceService.enabled=true --set control-plane.resourceService.env.DIRECTORY_STAGE_STRUCTURE=true", testutils.GetKeptnNameSpaceFromEnv(), chartURL))
		//_, err = testutils.ExecuteCommand(fmt.Sprintf(`helm upgrade -n %s keptn %s --wait --set=control-plane.apiGatewayNginx.type=LoadBalancer --set=control-plane.common.strategy.rollingUpdate.maxUnavailable=0 --set control-plane.resourceService.enabled=true --set control-plane.resourceService.env.DIRECTORY_STAGE_STRUCTURE=true --set control-plane.distributor.image.repository=docker.io/keptndev/distributor --set control-plane.distributor.image.tag=0.14.1-dev-PR-7308.202204010740 --set control-plane.lighthouseService.image.repository=docker.io/warber/lighthouse-service --set control-plane.lighthouseService.image.tag=%s`, testutils.GetKeptnNameSpaceFromEnv(), chartURL, lighthouseVersion))
		if err != nil {
			t.Logf("Encountered error when upgrading keptn: %v", err)

		}
		<-time.After(5 * time.Second)
	}

}
