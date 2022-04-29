package ZD

import (
	"context"
	"fmt"
	testutils "github.com/keptn/keptn/test/go-tests"
	"sync"
	"testing"
	"time"
)

const apiProbeInterval = 5 * time.Second
const sequencesInterval = 15 * time.Second

var chartLatestVersion = "https://github.com/keptn/helm-charts-dev/blob/e242b0ab5d4eba9ec7b170a60e6d6e2421bcea86/packages/keptn-0.15.0-dev.tgz?raw=true"
var chartPreviousVersion = "https://github.com/keptn/helm-charts-dev/blob/366d236e97e147596e332b48d94f44b094fb349a/packages/keptn-0.15.0-dev-PR-7504.tgz?raw=true"

type ZeroDowntimeEnv struct {
	Ctx          context.Context //TODO substitute context & cancel with a quit channel not to store/share context
	Cancel       context.CancelFunc
	NrOfUpgrades int
	Wg           *sync.WaitGroup

	//api test fields
	TotalAPICalls  uint64
	FailedAPICalls uint64
	PassedAPICalls uint64

	//sequence related test fields
	ShipyardFile    string
	ExistingProject string
	FiredSequences  uint64
	FailedSequences uint64
	PassedSequences uint64
	Id              uint64
}

func SetupZD() *ZeroDowntimeEnv {

	zd := ZeroDowntimeEnv{}
	zd.Ctx, zd.Cancel = context.WithCancel(context.Background())
	zd.NrOfUpgrades = 2
	zd.Wg = &sync.WaitGroup{}
	zd.ShipyardFile, _ = GetShipyard()
	zd.TotalAPICalls = 0
	zd.FailedAPICalls = 0
	zd.PassedAPICalls = 0
	zd.FiredSequences = 0
	zd.FailedSequences = 0
	zd.PassedSequences = 0
	zd.Id = 0
	return &zd
}

func Test_ZeroDowntime(t *testing.T) {
	t.Run("Sequences Zero Downtime", TestSequencesZD)
	t.Run("Webhook Zero Downtime", TestWebhookZD)
}

func RollingUpgrade(t *testing.T, env *ZeroDowntimeEnv) {
	defer func() {
		env.Cancel()
		t.Log("Rolling upgrade terminated")
	}()
	t.Log("Upgrade in progress")
	//time.Sleep(1 * time.Minute)
	for i := 0; i < env.NrOfUpgrades; i++ {
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
	}
}
