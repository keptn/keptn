package zerodowntime

import (
	"context"
	"github.com/stretchr/testify/suite"
	"sync"
	"testing"
	"time"
)

const apiProbeInterval = 5 * time.Second
const sequencesInterval = 15 * time.Second

var chartLatestVersion = "https://github.com/keptn/helm-charts-dev/blob/1c234d5370f76532e0338adb8d135fe6e1d4caf8/packages/keptn-0.15.0-dev.tgz?raw=true"
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

type TestSuiteDowntime struct {
	suite.Suite
}

func (suite *TestSuiteDowntime) SetupSuite() {

}

//Test_ZeroDowntime runs all test suites
func Test_ZeroDowntime(t *testing.T) {
	suite.Run(t, new(TestSuiteDowntime))
}

//func (suite *TestSuiteDowntime) TestSequences() {
//	ZDTestTemplate(suite.T(), Sequences, "Sequences")
//}

func (suite *TestSuiteDowntime) TestWebhook() {
	ZDTestTemplate(suite.T(), Webhook, "Webhook")
}

func (suite *TestSuiteDowntime) TearDownSuite() {
}

func ZDTestTemplate(t *testing.T, F func(t1 *testing.T, e *ZeroDowntimeEnv), name string) {

	env := SetupZD()
	t.Run("Rolling Upgrade", func(t2 *testing.T) {
		t2.Parallel()
		RollingUpgrade(t2, env)
	})

	t.Run("API", func(t2 *testing.T) {
		t2.Parallel()
		APIs(t2, env)
	})

	t.Run(name, func(t2 *testing.T) {
		t2.Parallel()
		F(t2, env)

	})

}

func RollingUpgrade(t *testing.T, env *ZeroDowntimeEnv) {
	defer func() {
		env.Cancel()
		t.Log("Rolling upgrade terminated")
		env.Wg.Wait()

		PrintSequencesResults(t, env)
		PrintAPIresults(t, env)

	}()
	t.Log("Upgrade in progress")
	time.Sleep(1 * time.Minute)
	//for i := 0; i < env.NrOfUpgrades; i++ {
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
