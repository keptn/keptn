package zero_downtime

import (
	"errors"
	"fmt"
	"github.com/kelseyhightower/envconfig"
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

type ZeroDowntimeEnv struct {
	quit              chan struct{}
	NrOfUpgrades      int           `envconfig:"NUMBER_OF_UPGRADES" default:"3"`
	EnvInstallVersion string        `envconfig:"INSTALL_HELM_CHART"` //for local run you can add a default ref to this string here e.g. default:"https://github.com/keptn/helm-charts-dev/raw/69eea439a26a99ecc163e296860dbb5d43e41600/packages/keptn-0.15.1-dev.tgz"`
	EnvUpgradeVersion string        `envconfig:"UPGRADE_HELM_CHART"` //for local run you can add a default ref to this string here e.g. default:"https://github.com/keptn/helm-charts-dev/raw/gh-pages/packages/keptn-0.15.0-dev.tgz"
	ApiProbeInterval  time.Duration `envconfig:"API_PROBES_INTERVAL" default:"5s"`
	SequencesInterval time.Duration `envconfig:"SEQUENCES_INTERVAL" default:"30s"`
	Wg                *sync.WaitGroup

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
	if err := envconfig.Process("", &zd); err != nil {
		os.Exit(1)
	}

	zd.quit = make(chan struct{})
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
	env.Wg.Add(2)
	t.Run("Rolling Upgrade", func(t1 *testing.T) {
		t1.Parallel()
		RollingUpgrade(t1, env)
	})
	t.Run("API", func(t1 *testing.T) {
		t1.Parallel()
		APIs(t1, env)
		env.Wg.Done()
	})
	t.Run(name, func(t1 *testing.T) {
		t1.Parallel()
		F(t1, env)
		// The test summary should be printed after the tests have finished and before the test suite returns
		// to avoid failure due to test context expired
		t.Run("Summary", func(t *testing.T) {
			t.Log("Test results for ", name)
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
	if env.EnvUpgradeVersion == "" || env.EnvInstallVersion == "" {
		t.Fatal("Test cannot run without setting INSTALL_HELM_CHART and UPGRADE_HELM_CHART")
	}
	t.Log("Upgrade in progress")

	for i := 0; i < env.NrOfUpgrades; i++ {
		time.Sleep(60 * time.Second)
		chartPath := ""
		var err error
		if i%2 == 0 {
			chartPath = env.EnvUpgradeVersion
		} else {
			chartPath = env.EnvInstallVersion
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

func GetShipyard() (string, error) {
	return testutils.CreateTmpShipyardFile(shipyard)
}
