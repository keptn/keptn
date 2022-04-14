package ZD

import (
	"context"
	"github.com/benbjohnson/clock"
	"testing"
	"time"
)

const apiProbeInterval = 1 * time.Second
const deliveryInterval = 10 * time.Second

func TestZeroDowntime(t *testing.T) {

	ticker := clock.New().Ticker(apiProbeInterval)
	ticker2 := clock.New().Ticker(deliveryInterval)
	ctx, cancel := context.WithCancel(context.Background())
	go rollingUpgrade(cancel, t)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			go TestAPIs(t)
		case <-ticker2.C:
			go TestDelivery(t)
		}

	}
}

func rollingUpgrade(cancel context.CancelFunc, t *testing.T) {
	defer cancel()
	t.Log("rolling")
	time.Sleep(30 * time.Second)

	//TODO setup helm upgrade here
	//_, err := ExecuteCommand(fmt.Sprintf(`helm upgrade -n %s keptn %s --wait --set=control-plane.apiGatewayNginx.type=LoadBalancer --set=control-plane.common.strategy.rollingUpdate.maxUnavailable=0 --set control-plane.resourceService.enabled=true --set control-plane.resourceService.env.DIRECTORY_STAGE_STRUCTURE=true --set control-plane.distributor.image.repository=docker.io/keptndev/distributor --set control-plane.distributor.image.tag=0.14.1-dev-PR-7308.202204010740 --set control-plane.lighthouseService.image.repository=docker.io/warber/lighthouse-service --set control-plane.lighthouseService.image.tag=%s`, GetKeptnNameSpaceFromEnv(), chartURL, lighthouseVersion))
	//assert.Nil(t, err)
}
