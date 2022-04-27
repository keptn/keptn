package ZD

import (
	"context"
	"fmt"
	"github.com/benbjohnson/clock"
	testutils "github.com/keptn/keptn/test/go-tests"
	"github.com/stretchr/testify/suite"
	"sync"
	"testing"
	"time"
)

const apiProbeInterval = 5 * time.Second
const sequencesInterval = 15 * time.Second

var chartLatestVersion = "https://github.com/keptn/helm-charts-dev/blob/1efe3dab77da9ea3cf2b7dd5eff4b2fac6f76633/packages/keptn-0.15.0-dev-PR-7266.tgz?raw=true"
var chartPreviousVersion = "https://github.com/keptn/helm-charts-dev/blob/5b4fbc630895a2a71721763110376b452f4c2c67/packages/keptn-0.15.0-dev-PR-7266.tgz?raw=true"

type ZeroDowntimeEnv struct {
	Ctx             context.Context
	Cancel          context.CancelFunc
	ApiTicker       *clock.Ticker
	SeqTicker       *clock.Ticker
	NrOfUpgrades    int
	Wg              *sync.WaitGroup
	ShipyardFile    string
	ExistingProject string
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

type TestSuiteDowntime struct {
	suite.Suite
	env *ZeroDowntimeEnv
}

func (suite *TestSuiteDowntime) SetupSuite() {
	suite.env = SetupZD()
	var err error
	suite.env.ExistingProject, err = testutils.CreateProject("projectzd", suite.env.ShipyardFile)
	suite.Nil(err)
	_, err = testutils.ExecuteCommand(fmt.Sprintf("keptn create service %s --project=%s", "myservice", suite.env.ExistingProject))
	suite.Nil(err)

}

func Test_ZeroDowntime(t *testing.T) {
	suite.Run(t, new(TestSuiteDowntime))
}

func (suite *TestSuiteDowntime) TestParallelZeroDowntime() {

	suite.T().Run("Rolling Upgrade", func(t2 *testing.T) {
		t2.Parallel()
		rollingUpgrade(t2, suite.env)
	})

	suite.T().Run("API", func(t2 *testing.T) {
		t2.Parallel()
		APIs(t2, suite.env)
	})

	suite.T().Run("Sequences", func(t2 *testing.T) {
		t2.Parallel()
		Sequences(t2, suite.env)
	})

}

func (suite *TestSuiteDowntime) TestResults() {
	suite.T().Run("Summary:", func(t2 *testing.T) {
		PrintSequencesResults(t2, suite.env)
		PrintAPIresults(t2, suite.env)
	})
}

func rollingUpgrade(t *testing.T, env *ZeroDowntimeEnv) {
	defer func() {
		env.Cancel()
		t.Log("Ended")
	}()
	t.Log("Rolling")

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
