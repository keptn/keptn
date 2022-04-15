package ZD

import (
	"context"
	"github.com/benbjohnson/clock"
	"github.com/stretchr/testify/assert"
	"os"
	"sync"
	"testing"
	"time"
)

const apiProbeInterval = 1 * time.Second
const deliveryInterval = 10 * time.Second

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
	go rollingUpgrade(cancel, t)

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

func rollingUpgrade(cancel context.CancelFunc, t *testing.T) {
	defer cancel()
	t.Log("rolling")
	time.Sleep(40 * time.Second)

	//TODO setup helm upgrade here
	//_, err := ExecuteCommand(fmt.Sprintf(`helm upgrade -n %s keptn %s --wait --set=control-plane.apiGatewayNginx.type=LoadBalancer --set=control-plane.common.strategy.rollingUpdate.maxUnavailable=0 --set control-plane.resourceService.enabled=true --set control-plane.resourceService.env.DIRECTORY_STAGE_STRUCTURE=true --set control-plane.distributor.image.repository=docker.io/keptndev/distributor --set control-plane.distributor.image.tag=0.14.1-dev-PR-7308.202204010740 --set control-plane.lighthouseService.image.repository=docker.io/warber/lighthouse-service --set control-plane.lighthouseService.image.tag=%s`, GetKeptnNameSpaceFromEnv(), chartURL, lighthouseVersion))
	//assert.Nil(t, err)

	//	chartLatestVersion := "https://github.com/keptn/helm-charts-dev/blob/1efe3dab77da9ea3cf2b7dd5eff4b2fac6f76633/packages/keptn-0.15.0-dev-PR-7266.tgz?raw=true"
	//	chartPreviousVersion := "https://github.com/keptn/helm-charts-dev/blob/5b4fbc630895a2a71721763110376b452f4c2c67/packages/keptn-0.15.0-dev-PR-7266.tgz?raw=true"
}
