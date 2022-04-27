package ZD

import (
	"context"
	"github.com/benbjohnson/clock"
	"sync"
	"testing"
	"time"
)

const apiProbeInterval = 5 * time.Second
const sequencesInterval = 15 * time.Second

var chartLatestVersion = "https://github.com/keptn/helm-charts-dev/blob/1efe3dab77da9ea3cf2b7dd5eff4b2fac6f76633/packages/keptn-0.15.0-dev-PR-7266.tgz?raw=true"
var chartPreviousVersion = "https://github.com/keptn/helm-charts-dev/blob/5b4fbc630895a2a71721763110376b452f4c2c67/packages/keptn-0.15.0-dev-PR-7266.tgz?raw=true"

var env = SetupZD()

type ZeroDowntimeEnv struct {
	Ctx             context.Context
	Cancel          context.CancelFunc
	ApiTicker       *clock.Ticker
	SeqTicker       *clock.Ticker
	NrOfUpgrades    int
	Wg              *sync.WaitGroup
	ShipyardFile    string
	TotalAPICalls   uint64
	FailedAPICalls  uint64
	PassedAPICalls  uint64
	FiredSequences  uint64
	FailedSequences uint64
	PassedSequences uint64
	Id              uint64
}

func SetupZD() *ZeroDowntimeEnv {

	zd := ZeroDowntimeEnv{}
	zd.Ctx, zd.Cancel = context.WithCancel(context.Background())
	zd.ApiTicker = clock.New().Ticker(apiProbeInterval)
	zd.SeqTicker = clock.New().Ticker(sequencesInterval)
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

func TestParallelZeroDowntime(t *testing.T) {

	t.Run("Rolling Upgrade", func(t2 *testing.T) {
		t2.Parallel()
		rollingUpgrade(t2)

	})

	t.Run("API", func(t2 *testing.T) {
		t2.Parallel()
		TestAPIs(t2)
	})

	t.Run("Sequences", func(t2 *testing.T) {
		t2.Parallel()
		Test_Sequences(t2)

	})

}

//
//func TestZeroDowntime(t *testing.T) {
//
//	ExistingProject, err := testutils.CreateProject("projectzd", env.ShipyardFile)
//	assert.Nil(t, err)
//	_, err = testutils.ExecuteCommand(fmt.Sprintf("keptn create service %s --project=%s", "myservice", ExistingProject))
//	assert.Nil(t, err)
//
//	t.Parallel()
//	// run a tests before starting update
//	go t.Run("Before upgrade sequences", Test_Sequences)
//
//	//start update
//	go t.Run("rolling update", rollingUpgrade)
//
//	// run tests meanwhile
//	go func() {
//		for {
//			select {
//			case <-env.Ctx.Done():
//				break
//			case <-env.ApiTicker.C:
//				atomic.AddUint64(&env.TotalAPICalls, 1)
//				go t.Run("API test", TestAPIs)
//			}
//		}
//	}()
//
//	for {
//		select {
//		case <-env.Ctx.Done():
//			break
//		case <-env.SeqTicker.C:
//			env.Wg.Add(1)
//			go t.Run("Before upgrade sequences", Test_Sequences)
//		}
//	}
//
//	env.Wg.Wait()
//	t.Run("Summary: ", printResults)
//}

func rollingUpgrade(t *testing.T) {
	defer func() {
		env.Cancel()
		t.Log("Ended")
	}()
	t.Log("Rolling")
	time.Sleep(1 * time.Minute)
	//
	//for i := 0; i < nrOfUpgrades; i++ {
	//	chartURL := ""
	//	var err error
	//	if i%2 == 0 {
	//		chartURL = chartLatestVersion
	//	} else {
	//		chartURL = chartPreviousVersion
	//	}
	//	t.Logf("Upgrading Keptn to %s", chartURL)
	//	_, err = testutils.ExecuteCommand(
	//		fmt.Sprintf(
	//			"helm upgrade -n %s keptn %s --wait --set=control-plane.apiGatewayNginx.type=LoadBalancer "+
	//				"--set=control-plane.common.strategy.rollingUpdate.maxUnavailable=0 --set control-plane.resourceService.enabled=true"+
	//				" --set control-plane.resourceService.env.DIRECTORY_STAGE_STRUCTURE=true", testutils.GetKeptnNameSpaceFromEnv(), chartURL))
	//	if err != nil {
	//		t.Logf("Encountered error when upgrading keptn: %v", err)
	//
	//	}
	//}
}
