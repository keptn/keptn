package zero_downtime

import (
	"context"
	"github.com/stretchr/testify/suite"
	"os"
	"sync"
	"testing"
	"time"
)

const apiProbeInterval = 5 * time.Second
const sequencesInterval = 15 * time.Second

const EnvInstallVersion = "INSTALL_HELM_CHART"
const EnvUpgradeVersion = "UPGRADE_HELM_CHART"

const pathToChart = "https://charts-dev.keptn.sh/packages/"
const rawChart = "?raw=true"

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

type TestSuiteDowntime struct {
	suite.Suite
}

func (suite *TestSuiteDowntime) SetupSuite() {

}

func Test_ZeroDowntime(t *testing.T) {

}

func (suite *TestSuiteDowntime) TestSequences() {

}

func (suite *TestSuiteDowntime) TestWebhook() {

}

func (suite *TestSuiteDowntime) TearDownSuite() {
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
