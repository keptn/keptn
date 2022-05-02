package zerodowntime

import (
	"fmt"
	testutils "github.com/keptn/keptn/test/go-tests"
	"github.com/stretchr/testify/suite"
	"os"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

const EnvInstallVersion = "INSTALL_HELM_CHART"
const EnvUpgradeVersion = "UPGRADE_HELM_CHART"
const pathToChart = "https://charts-dev.keptn.sh/packages/"
const rawChart = "?raw=true"

const apiProbeInterval = 5 * time.Second
const sequencesInterval = 15 * time.Second

type ZeroDowntimeEnv struct {
	quit         chan struct{}
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
	zd.quit = make(chan struct{})
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

func (env *ZeroDowntimeEnv) gedId() string {
	atomic.AddUint64(&env.Id, 1)
	return fmt.Sprintf("%d", env.Id)
}

func (env *ZeroDowntimeEnv) failSequence() {
	atomic.AddUint64(&env.FailedSequences, 1)
}

func (env *ZeroDowntimeEnv) passSequence() {
	atomic.AddUint64(&env.PassedSequences, 1)
}

func (env *ZeroDowntimeEnv) passFailedSequence() {
	atomic.AddUint64(&env.FailedSequences, ^uint64(1-1))
	env.passSequence()
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

func (suite *TestSuiteDowntime) TestSequences() {
	ZDTestTemplate(suite.T(), Sequences, "Sequences")
}

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
		close(env.quit)
		t.Log("Rolling upgrade terminated")
		env.Wg.Wait()

		PrintSequencesResults(t, env)
		PrintAPIresults(t, env)

	}()

	chartPreviousVersion, chartLatestVersion := GetCharts()

	t.Log("Upgrade in progress")
	time.Sleep(1 * time.Minute)
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

func PrintSequencesResults(t *testing.T, env *ZeroDowntimeEnv) {

	t.Log("-----------------------------------------------")
	t.Log("TOTAL SEQUENCES: ", env.FiredSequences)
	t.Log("TOTAL SUCCESS ", env.PassedSequences)
	t.Log("TOTAL FAILURES ", env.FailedSequences)
	t.Log("-----------------------------------------------")

}

//Returns current test helm charts for the rolling upgrade
func GetCharts() (string, string) {
	var install, upgrade string

	chartInstallVersion := "https://charts-dev.keptn.sh/packages/keptn-0.15.0-dev.tgz?raw=true"
	chartUpgradeVersion := "https://charts-dev.keptn.sh/packages/keptn-0.15.0-dev-PR-7504.tgz?raw=true"

	if install = os.Getenv(EnvInstallVersion); install != "" {
		chartInstallVersion = pathToChart + install + rawChart
	}
	if upgrade = os.Getenv(EnvUpgradeVersion); upgrade != "" {
		chartUpgradeVersion = pathToChart + upgrade + rawChart
	}

	return chartInstallVersion, chartUpgradeVersion
}
