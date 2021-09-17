package cmd

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/keptn/keptn/cli/pkg/credentialmanager"
	"github.com/keptn/keptn/cli/pkg/logging"
)

func init() {
	logging.InitLoggers(os.Stdout, os.Stdout, os.Stderr)
}

// TestEvaluationStart tests the start-evaluation command.
func TestEvaluationStart(t *testing.T) {

	credentialmanager.MockAuthCreds = true
	checkEndPointStatusMock = true

	*evaluationStart.Timeframe = ""
	*evaluationStart.Start = ""
	*evaluationStart.End = ""

	cmd := fmt.Sprintf("send event start-evaluation --project=%s --stage=%s --service=%s "+
		"--timeframe=%s --labels=foo=bar,bar=foo --mock", "sockshop", "hardening", "carts", "5m")
	_, err := executeActionCommandC(cmd)
	if err != nil {
		t.Errorf(unexpectedErrMsg, err)
	}
}

func TestEvaluationStartTimeSpecified(t *testing.T) {

	credentialmanager.MockAuthCreds = true
	checkEndPointStatusMock = true

	*evaluationStart.Timeframe = ""
	*evaluationStart.Start = ""
	*evaluationStart.End = ""

	cmd := fmt.Sprintf("send event start-evaluation --project=%s --stage=%s --service=%s "+
		"--timeframe=%s --start=%s --mock", "sockshop", "hardening", "carts", "5m", "2019-07-24T10:17:12.000Z")
	_, err := executeActionCommandC(cmd)

	if err != nil {
		t.Errorf(unexpectedErrMsg, err)
	}
}

func TestEvaluationStartAndEndTimeSpecified(t *testing.T) {

	credentialmanager.MockAuthCreds = true
	checkEndPointStatusMock = true

	*evaluationStart.Timeframe = ""
	*evaluationStart.Start = ""
	*evaluationStart.End = ""

	cmd := fmt.Sprintf("send event start-evaluation --project=%s --stage=%s --service=%s "+
		"--start=%s --end=%s --mock", "sockshop", "hardening", "carts", "2019-07-24T10:17:12.000Z", "2019-07-24T10:20:12.000Z")
	_, err := executeActionCommandC(cmd)

	if err != nil {
		t.Errorf(unexpectedErrMsg, err)
	}
}

func TestEvaluationStartAndEndTimeAndTimeframeSpecified(t *testing.T) {

	credentialmanager.MockAuthCreds = true
	checkEndPointStatusMock = true

	*evaluationStart.Timeframe = ""
	*evaluationStart.Start = ""
	*evaluationStart.End = ""

	cmd := fmt.Sprintf("send event start-evaluation --project=%s --stage=%s --service=%s "+
		"--start=%s --end=%s --timeframe=%s --mock", "sockshop", "hardening", "carts", "2019-07-24T10:17:12.000Z",
		"2019-07-24T10:20:12.000Z", "5m")
	_, err := executeActionCommandC(cmd)

	if err == nil {
		t.Error("An error occurred: expect an error due to too many parameters (start, end and timeframe) used at the same time")
	}

	assert.EqualValues(t, "You can not use --end together with --timeframe", err.Error())
}

func TestEvaluationStartAndEndTimeWrongOrder(t *testing.T) {

	credentialmanager.MockAuthCreds = true
	checkEndPointStatusMock = true

	*evaluationStart.Timeframe = ""
	*evaluationStart.Start = ""
	*evaluationStart.End = ""

	cmd := fmt.Sprintf("send event start-evaluation --project=%s --stage=%s --service=%s "+
		"--start=%s --end=%s  --mock", "sockshop", "hardening", "carts", "2019-07-24T10:17:12.000Z", "2019-07-24T10:10:12.000Z")
	_, err := executeActionCommandC(cmd)

	if err == nil {
		t.Error("An error occurred: expect an error as end time is before start time")
	}

	assert.EqualValues(t, "Start and end time of evaluation time frame not set: end date must be at least 1 minute after start date", err.Error())
}
