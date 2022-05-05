package zerodowntime

import (
	"fmt"
	"github.com/benbjohnson/clock"
	testutils "github.com/keptn/keptn/test/go-tests"
	"github.com/steinfletcher/apitest"
	jsonpath "github.com/steinfletcher/apitest-jsonpath"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

type TestSuiteAPI struct {
	suite.Suite
	env         *ZeroDowntimeEnv
	token       string
	keptnAPIURL string
}

// Tests are run before they start
func (suite *TestSuiteAPI) SetupSuite() {
	var err error
	suite.token, suite.keptnAPIURL, err = testutils.GetApiCredentials()
	suite.Assert().Nil(err)
}

//Test_API can be used to test a single call to all the tests in the API test suite
func Test_API(t *testing.T) {
	Env := SetupZD()
	var err error
	Env.ExistingProject, err = testutils.CreateProject("projectzd", Env.ShipyardFile)
	assert.Nil(t, err)
	_, err = testutils.ExecuteCommand(fmt.Sprintf("keptn create service %s --project=%s", "myservice", Env.ExistingProject))
	assert.Nil(t, err)

	s := &TestSuiteAPI{
		env: Env,
	}
	suite.Run(t, s)

}

// APIs is called in the zero downtime test suite
func APIs(t *testing.T, env *ZeroDowntimeEnv) {
	wgAPI := sync.WaitGroup{}
	apiTicker := clock.New().Ticker(apiProbeInterval)

Loop:
	for {
		select {
		case <-env.Ctx.Done():
			break Loop
		case <-apiTicker.C:
			wgAPI.Add(1)
			apisuite := &TestSuiteAPI{
				env: env,
			}
			go func() {
				suite.Run(t, apisuite)
				wgAPI.Done()
			}()

		}
	}
	wgAPI.Wait()
}

//to generate html report we can add 	.Report(apitest.SequenceDiagram())

func (suite *TestSuiteAPI) Test_API_Service() {

	started := time.Now()
	apiURL := suite.keptnAPIURL + "/v1"

	api := apitest.New("Test api-service auth").EnableNetworking(getClient(1)).
		Observe(func(res *http.Response, req *http.Request, apiTest *apitest.APITest) {
			suite.logResult(res, apiTest, http.StatusOK, started)
		})

	api.Post(apiURL + "/auth").Headers(map[string]string{"x-token": suite.token}).
		Expect(suite.T()).Status(http.StatusOK).End()

	api.Get(apiURL + "/metadata").Headers(map[string]string{"x-token": suite.token}).
		Expect(suite.T()).Status(http.StatusOK).End()

}

func (suite *TestSuiteAPI) Test_Statistic_Service() {

	started := time.Now()
	apiURL := suite.keptnAPIURL + "/statistics/v1"

	api := apitest.New("Test statistics-service").EnableNetworking(getClient(1)).
		Observe(func(res *http.Response, req *http.Request, apiTest *apitest.APITest) {
			suite.logResult(res, apiTest, http.StatusNotFound, started)
		})

	api.Get(apiURL+"/statistics").
		Query("from", "1648190000").Query("to", "1648195292").
		Headers(map[string]string{"x-token": suite.token}).
		Expect(suite.T()).Status(http.StatusNotFound).
		Assert(jsonpath.Equal(`$.message`, "no statistics found for selected time frame")).End()

}

func (suite *TestSuiteAPI) Test_Secret_Service() {

	started := time.Now()
	apiURL := suite.keptnAPIURL + "/secrets/v1"

	api := apitest.New("Test secret-service").EnableNetworking(getClient(1)).
		Observe(func(res *http.Response, req *http.Request, apiTest *apitest.APITest) {
			suite.logResult(res, apiTest, http.StatusOK, started)
		})

	api.Get(apiURL + "/scope").
		Headers(map[string]string{"x-token": suite.token}).
		Expect(suite.T()).Status(http.StatusOK).Assert(jsonpath.Contains(`$.scopes`, "keptn-default")).End()

}

func (suite *TestSuiteAPI) Test_Configuration_Service() {

	started := time.Now()
	apiURL := suite.keptnAPIURL + "/configuration-service/v1"

	api := apitest.New("Test configuration-service: not existing project").
		Observe(func(res *http.Response, req *http.Request, apiTest *apitest.APITest) {
			suite.logResult(res, apiTest, http.StatusNotFound, started)
		}).EnableNetworking(getClient(1))

	req := api.Get(apiURL + "/project/unexisting-project/resource").
		Headers(map[string]string{"x-token": suite.token}).
		Expect(suite.T()).
		Status(http.StatusNotFound)

	if res, err := testutils.CompareServiceNameWithDeploymentName("configuration-service", "configuration-service"); err == nil && res {
		req.Assert(jsonpath.Equal(`$.message`, "Project does not exist")).End()
	} else {
		req.Assert(jsonpath.Equal(`$.message`, "Could not find credentials for upstream repository")).End()
	}

}

func (suite *TestSuiteAPI) Test_ControlPlane() {
	started := time.Now()
	apiURL := suite.keptnAPIURL + "/controlPlane/v1"

	api := apitest.New("Test control-plane: check uniform").EnableNetworking(getClient(1)).
		Observe(func(res *http.Response, req *http.Request, apiTest *apitest.APITest) {
			suite.logResult(res, apiTest, http.StatusOK, started)
		})

	api.Get(apiURL+"/uniform/registration").Query("name", "lighthouse-service").
		Headers(map[string]string{"x-token": suite.token}).
		Expect(suite.T()).Status(http.StatusOK).Assert(jsonpath.Equal(`$[0].name`, "lighthouse-service")).End()

}

func (suite *TestSuiteAPI) Test_MongoDB() {
	started := time.Now()
	apiURL := suite.keptnAPIURL + "/mongodb-datastore"

	api := apitest.New("Test mongo-datastore: not existing project").EnableNetworking(getClient(1)).
		Observe(func(res *http.Response, req *http.Request, apiTest *apitest.APITest) {
			suite.logResult(res, apiTest, http.StatusOK, started)
		})

	api.Get(apiURL+"/event").Query("project", "keptn").Query("pageSize", "20").
		Headers(map[string]string{"x-token": suite.token}).
		Expect(suite.T()).Status(http.StatusOK).Body(`{"events":[], "pageSize":20}`).End()

}

func (suite *TestSuiteAPI) logResult(res *http.Response, apiTest *apitest.APITest, expected int, started time.Time) {
	atomic.AddUint64(&(suite.env.TotalAPICalls), 1)
	finished := time.Now()
	//remove headers cookies and auth from request
	censoredReq := strings.Split(strings.Split(fmt.Sprintf("%+v", *apiTest.Request()), "{interceptor:<nil>")[1], "headers:")[0]

	suite.T().Log("\n", "Test API configuration: ", censoredReq)
	suite.T().Logf("Response Body: %v", res.Body)
	suite.T().Logf("Duration: %s", finished.Sub(started))
	if res.StatusCode != expected {
		atomic.AddUint64(&suite.env.FailedAPICalls, 1)
	} else {
		atomic.AddUint64(&suite.env.PassedAPICalls, 1)
	}
}

func getClient(sec time.Duration) *http.Client {
	return &http.Client{
		Timeout: time.Second * sec,
	}
}

func PrintAPIresults(t *testing.T, env *ZeroDowntimeEnv) {
	t.Log("-----------------------------------------------")
	t.Log("TOTAL API PROBES", env.TotalAPICalls)
	t.Log("TOTAL PROBES SUCCEEDED", env.PassedAPICalls)
	t.Log("TOTAL PROBES FAILED", env.FailedAPICalls)
	t.Log("-----------------------------------------------")
}
