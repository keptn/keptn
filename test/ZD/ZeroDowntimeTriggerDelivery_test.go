package ZD

import (
	testutils "github.com/keptn/keptn/test/go-tests"
	"github.com/stretchr/testify/suite"
	"testing"
)

type TestTriggerKeptn struct {
	suite.Suite
	token       string
	keptnAPIURL string
}

// Tests are run before they start
func (suite *TestTriggerKeptn) SetupSuite() {
	var err error
	suite.token, suite.keptnAPIURL, err = testutils.GetApiCredentials()
	suite.Assert().Nil(err)
}

// Running after each test
func (suite *TestTriggerKeptn) TearDownTest() {

}

// Running after all tests are completed
func (suite *TestTriggerKeptn) TearDownSuite() {

}

// This gets run automatically by `go test` so we call `suite.Run` inside it
func TestDelivery(t *testing.T) {
	suite.Run(t, new(TestTriggerKeptn))
}

func (suite *TestTriggerKeptn) Test_Delivery() {

}

func (suite *TestTriggerKeptn) Test_Evaluation() {

}
