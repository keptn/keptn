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
