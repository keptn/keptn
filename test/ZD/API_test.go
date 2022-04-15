package ZD

import (
	testutils "github.com/keptn/keptn/test/go-tests"
	"github.com/steinfletcher/apitest"
	jsonpath "github.com/steinfletcher/apitest-jsonpath"
	"github.com/stretchr/testify/suite"
	"net/http"
	"testing"
	"time"
)

type TestSuiteEnv struct {
	suite.Suite
	token       string
	keptnAPIURL string
}

// Tests are run before they start
func (suite *TestSuiteEnv) SetupSuite() {
	var err error
	suite.token, suite.keptnAPIURL, err = testutils.GetApiCredentials()
	suite.Assert().Nil(err)
}

// Running after each test
func (suite *TestSuiteEnv) TearDownTest() {

}

// Running after all tests are completed
func (suite *TestSuiteEnv) TearDownSuite() {

}

// This gets run automatically by `go test` so we call `suite.Run` inside it
func TestAPIs(t *testing.T) {
	suite.Run(t, new(TestSuiteEnv))
}

func (suite *TestSuiteEnv) Test_API_Service() {

	apiURL := suite.keptnAPIURL + "/v1"

	apitest.New("Test api-service auth").EnableNetworking(getClient(1)).Post(apiURL + "/auth").Headers(map[string]string{"x-token": suite.token}).
		Expect(suite.T()).Status(http.StatusOK).End()

	apitest.New("Test api-service metadata").EnableNetworking(getClient(1)).Get(apiURL + "/metadata").Headers(map[string]string{"x-token": suite.token}).
		Expect(suite.T()).Status(http.StatusOK).End()

}

func (suite *TestSuiteEnv) Test_Configuration_Service() {

	apiURL := suite.keptnAPIURL + "/configuration-service/v1"

	req := apitest.New("Test configuration-service: not existing project").EnableNetworking(getClient(1)).Debug().
		Get(apiURL + "/project/unexisting-project/resource").
		Headers(map[string]string{"x-token": suite.token}).Expect(suite.T()).Status(http.StatusNotFound)

	if res, err := testutils.CompareServiceNameWithDeploymentName("configuration-service", "configuration-service"); err == nil && res {
		req.Assert(jsonpath.Equal(`$.message`, "Project does not exist")).End()
	} else {
		req.Assert(jsonpath.Equal(`$.message`, "Could not find credentials for upstream repository")).End()
	}
}

func getClient(sec time.Duration) *http.Client {
	return &http.Client{
		Timeout: time.Second * sec,
	}
}

func (suite *TestSuiteEnv) Test_ControlPlane() {

	apiURL := suite.keptnAPIURL + "/controlPlane/v1"

	// can add Body(`{"name": "jon", "id": "1234"}`) or .Assert(jsonpath.Equal(`$.key`, value )) Debug().
	apitest.New("Test control-plane: check uniform").EnableNetworking(getClient(1)).Debug().
		Get(apiURL+"/uniform/registration").Query("name", "lighthouse-service").
		Headers(map[string]string{"x-token": suite.token}).
		Expect(suite.T()).Status(http.StatusOK).Assert(jsonpath.Equal(`$[0].name`, "lighthouse-service")).End()

}

func (suite *TestSuiteEnv) Test_MongoDB() {

	apiURL := suite.keptnAPIURL + "/mongodb-datastore"

	apitest.New("Test mongo-datastore: not existing project").EnableNetworking(getClient(1)).Debug().
		Get(apiURL+"/event").Query("project", "keptn").Query("pageSize", "20").
		Headers(map[string]string{"x-token": suite.token}).
		Expect(suite.T()).Status(http.StatusOK).Body(`{"events":[], "pageSize":20}`).End()

}
