package zero_downtime

import (
	"errors"
	"fmt"
	"github.com/keptn/go-utils/pkg/common/retry"
	testutils "github.com/keptn/keptn/test/go-tests"
	"github.com/stretchr/testify/suite"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

const EnvInstallVersion = "INSTALL_HELM_CHART"
const EnvUpgradeVersion = "UPGRADE_HELM_CHART"
const valuesFile = "./assets/test-values.yml"

const shipyard = `--- 
apiVersion: spec.keptn.sh/0.2.0
kind: Shipyard
metadata: 
  name: shipyard-quality-gates
spec: 
  stages: 
    - 
      name: hardening`

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
	token, keptnAPIURL, err := testutils.GetApiCredentials()
	suite.Require().Nil(err)
	suite.T().Log("KEPTN ENDPOINT", keptnAPIURL)
	suite.T().Log("Authenticating keptn CLI")
	err = retry.Retry(func() error {
		out, err := testutils.ExecuteCommand(fmt.Sprintf("keptn auth --endpoint=%s --api-token=%s", keptnAPIURL, token))
		if err != nil {
			return err
		}
		if !strings.Contains(out, "Successfully authenticated") {
			return errors.New("authentication unsuccessful")
		}
		return nil
	}, retry.NumberOfRetries(10))
	suite.Require().Nil(err)

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

// ZDTestTemplate runs a test function in parallel to rolling upgrade and api probing,
// this can be used to create more zero downtime scenarios in the suite
func ZDTestTemplate(t *testing.T, F func(t2 *testing.T, e *ZeroDowntimeEnv), name string) {

	env := SetupZD()

	t.Run("Rolling Upgrade", func(t1 *testing.T) {
		t1.Parallel()
		RollingUpgrade(t1, env)
	})
	env.Wg.Add(1)
	t.Run("API", func(t1 *testing.T) {
		t1.Parallel()

		APIs(t1, env)
		env.Wg.Done()
	})
	env.Wg.Add(1)
	t.Run(name, func(t1 *testing.T) {
		t1.Parallel()
		F(t1, env)
		// The test summary should be printed after the tests have finished and before the test suite returns
		// to avoid failure due to test context expired
		t1.Run("Summary", func(t *testing.T) {
			t1.Log("Test results for ", name)
			PrintSequencesResults(env)
			PrintAPIresults(env)
			env.Wg.Done()
		})
	})

}

func RollingUpgrade(t *testing.T, env *ZeroDowntimeEnv) {
	defer func() {
		close(env.quit)
		t.Log("Rolling upgrade terminated")
		env.Wg.Wait()
	}()
	time.Sleep(30 * time.Second)
	chartPreviousVersion, chartLatestVersion := GetCharts(t)

	t.Log("Upgrade in progress")

	for i := 0; i < env.NrOfUpgrades; i++ {
		chartPath := ""
		var err error
		if i%2 == 0 {
			chartPath = chartLatestVersion
		} else {
			chartPath = chartPreviousVersion
		}
		t.Logf("Upgrading Keptn to %s", chartPath)
		_, err = testutils.ExecuteCommand(
			fmt.Sprintf(
				"helm upgrade keptn -n %s %s --wait --values=%s", testutils.GetKeptnNameSpaceFromEnv(), chartPath, valuesFile))
		if err != nil {
			t.Logf("Encountered error when upgrading keptn: %v", err)

		}
	}
}

func PrintSequencesResults(env *ZeroDowntimeEnv) {
	// print so that the log is shown even in case the test passes with gotestsum
	fmt.Println("-----------------------------------------------")
	fmt.Println("TOTAL SEQUENCES: ", env.FiredSequences)
	fmt.Println("TOTAL SUCCESS ", env.PassedSequences)
	fmt.Println("TOTAL FAILURES ", env.FailedSequences)
	fmt.Println("-----------------------------------------------")

}

// GetCharts returns the versions of helm charts for the rolling upgrade
// these can be set by two environment variables:
// "INSTALL_HELM_CHART" and "UPGRADE_HELM_CHART"
func GetCharts(t *testing.T) (string, string) {
	var install, upgrade string

	if install = os.Getenv(EnvInstallVersion); install == "" {
		t.Errorf("Helm chart unavailable, please set env variable %s", EnvInstallVersion)
	}
	if upgrade = os.Getenv(EnvUpgradeVersion); upgrade == "" {
		t.Errorf("Helm chart unavailable, please set env variable %s", EnvUpgradeVersion)
	}

	return install, upgrade
}

func GetShipyard() (string, error) {
	return testutils.CreateTmpShipyardFile(shipyard)
}
