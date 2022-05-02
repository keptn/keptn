package zero_downtime

import (
	"fmt"
	testutils "github.com/keptn/keptn/test/go-tests"
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
)

const EnvInstallVersion = "INSTALL_HELM_CHART"
const EnvUpgradeVersion = "UPGRADE_HELM_CHART"
const valuesFile = "test-values.yaml"

type ZeroDowntimeEnv struct {
	quit         chan struct{}
	NrOfUpgrades int
}

func SetupZD() *ZeroDowntimeEnv {

	zd := ZeroDowntimeEnv{}
	zd.quit = make(chan struct{})
	zd.NrOfUpgrades = 2
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

func (suite *TestSuiteDowntime) TestSequences() {
	env := SetupZD()
	suite.T().Run("Rolling Upgrade", func(t1 *testing.T) {
		RollingUpgrade(t1, env)
	})
}

func (suite *TestSuiteDowntime) TearDownSuite() {
}

func RollingUpgrade(t *testing.T, env *ZeroDowntimeEnv) {
	defer func() {
		close(env.quit)
		t.Log("Rolling upgrade terminated")
	}()

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
				"helm upgrade -n %s keptn %s --wait --values=%s ", testutils.GetKeptnNameSpaceFromEnv(), chartPath, valuesFile))
		if err != nil {
			t.Logf("Encountered error when upgrading keptn: %v", err)

		}
	}
}

//Returns current test helm charts for the rolling upgrade
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
