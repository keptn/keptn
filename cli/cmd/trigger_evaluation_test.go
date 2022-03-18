package cmd

import (
	"fmt"
	"os"
	"testing"

	"github.com/keptn/keptn/cli/pkg/credentialmanager"
	"github.com/keptn/keptn/cli/pkg/logging"
	"github.com/stretchr/testify/assert"
)

func init() {
	logging.InitLoggers(os.Stdout, os.Stdout, os.Stderr)
}

func TestTriggerEvaluation(t *testing.T) {

	credentialmanager.MockAuthCreds = true

	*triggerEvaluation.Timeframe = ""
	*triggerEvaluation.Start = ""
	*triggerEvaluation.End = ""

	cmd := fmt.Sprintf("trigger evaluation --project=%s --stage=%s --service=%s "+
		"--timeframe=%s --labels=foo=bar,bar=foo --mock", "sockshop", "hardening", "carts", "5m")
	_, err := executeActionCommandC(cmd)
	if err != nil {
		t.Errorf(unexpectedErrMsg, err)
	}
}

func TestTriggerEvaluationWrongFormat(t *testing.T) {

	credentialmanager.MockAuthCreds = true

	*triggerEvaluation.Timeframe = ""
	*triggerEvaluation.Start = ""
	*triggerEvaluation.End = ""

	cmd := fmt.Sprintf("send event start-evaluation --project=%s --stage=%s --service=%s "+
		"--timeframe=%s --mock", "sockshop", "hardening", "carts", "5h")
	_, err := executeActionCommandC(cmd)

	if err == nil {
		t.Error("An error occurred: expect an error due to wrong time frame format")
	}
}

func TestTriggerEvaluationVariousFormats(t *testing.T) {

	credentialmanager.MockAuthCreds = true

	times := []string{
		"2020-01-02T15:00:00",
		"2020-01-02T15:00:00Z",
		"2020-01-02T15:00:00+10:00",
		"2020-01-02T15:00:00.000Z",
		"2020-01-02T15:00:00.000000000Z",
	}

	*triggerEvaluation.Timeframe = ""
	*triggerEvaluation.Start = ""
	*triggerEvaluation.End = ""

	for _, time := range times {

		cmd := fmt.Sprintf("trigger evaluation --project=%s --stage=%s --service=%s "+
			"--start=%s --end=%s --mock", "sockshop", "hardening", "carts", time, "2020-01-02T15:10:12.000Z")
		_, err := executeActionCommandC(cmd)

		if err != nil {
			t.Errorf(unexpectedErrMsg, err)
		}

	}
}

func TestTriggerEvaluationStartTimeSpecified(t *testing.T) {

	credentialmanager.MockAuthCreds = true

	*triggerEvaluation.Timeframe = ""
	*triggerEvaluation.Start = ""
	*triggerEvaluation.End = ""

	cmd := fmt.Sprintf("trigger evaluation --project=%s --stage=%s --service=%s "+
		"--timeframe=%s --start=%s --mock", "sockshop", "hardening", "carts", "5m", "2019-07-24T10:17:12.000Z")
	_, err := executeActionCommandC(cmd)

	if err != nil {
		t.Errorf(unexpectedErrMsg, err)
	}
}

func TestTriggerEvaluationStartTimeNotSpecified(t *testing.T) {

	credentialmanager.MockAuthCreds = true

	*triggerEvaluation.Timeframe = ""
	*triggerEvaluation.Start = ""
	*triggerEvaluation.End = ""

	cmd := fmt.Sprintf("trigger evaluation --project=%s --stage=%s --service=%s "+
		"--timeframe=%s --mock", "sockshop", "hardening", "carts", "5m")
	_, err := executeActionCommandC(cmd)

	if err != nil {
		assert.EqualValues(t, "required flag(s) \"start\" not set", err.Error())
	}
}

func TestTriggerEvaluationProjectNotSpecified(t *testing.T) {

	credentialmanager.MockAuthCreds = true

	*triggerEvaluation.Timeframe = ""
	*triggerEvaluation.Start = ""
	*triggerEvaluation.End = ""

	cmd := fmt.Sprintf("trigger evaluation --stage=%s --service=%s "+
		"--timeframe=%s --start=%s --mock", "hardening", "carts", "5m", "2019-07-24T10:17:12.000Z")
	_, err := executeActionCommandC(cmd)

	if err != nil {
		assert.EqualValues(t, "required flag(s) \"project\" not set", err.Error())
	}
}

func TestTriggerEvaluationStageNotSpecified(t *testing.T) {

	credentialmanager.MockAuthCreds = true

	*triggerEvaluation.Timeframe = ""
	*triggerEvaluation.Start = ""
	*triggerEvaluation.End = ""

	cmd := fmt.Sprintf("trigger evaluation --project=%s --service=%s "+
		"--timeframe=%s --start=%s --mock", "sockshop", "carts", "5m", "2019-07-24T10:17:12.000Z")
	_, err := executeActionCommandC(cmd)

	if err != nil {
		assert.EqualValues(t, "required flag(s) \"stage\" not set", err.Error())
	}
}

func TestTriggerEvaluationServiceNotSpecified(t *testing.T) {

	credentialmanager.MockAuthCreds = true

	*triggerEvaluation.Timeframe = ""
	*triggerEvaluation.Start = ""
	*triggerEvaluation.End = ""

	cmd := fmt.Sprintf("trigger evaluation --project=%s --stage=%s "+
		"--timeframe=%s --start=%s --mock", "sockshop", "hardening", "5m", "2019-07-24T10:17:12.000Z")
	_, err := executeActionCommandC(cmd)

	if err != nil {
		assert.EqualValues(t, "required flag(s) \"service\" not set", err.Error())
	}
}

func TestTriggerEvaluationStartAndEndTimeSpecified(t *testing.T) {

	credentialmanager.MockAuthCreds = true

	*triggerEvaluation.Timeframe = ""
	*triggerEvaluation.Start = ""
	*triggerEvaluation.End = ""

	cmd := fmt.Sprintf("trigger evaluation --project=%s --stage=%s --service=%s "+
		"--start=%s --end=%s --mock", "sockshop", "hardening", "carts", "2019-07-24T10:17:12.000Z", "2019-07-24T10:20:12.000Z")
	_, err := executeActionCommandC(cmd)

	if err != nil {
		t.Errorf(unexpectedErrMsg, err)
	}
}

func TestTriggerEvaluationStartAndEndTimeAndTimeframeSpecified(t *testing.T) {

	credentialmanager.MockAuthCreds = true

	*triggerEvaluation.Timeframe = ""
	*triggerEvaluation.Start = ""
	*triggerEvaluation.End = ""

	cmd := fmt.Sprintf("trigger evaluation --project=%s --stage=%s --service=%s "+
		"--start=%s --end=%s --timeframe=%s --mock", "sockshop", "hardening", "carts", "2019-07-24T10:17:12.000Z",
		"2019-07-24T10:20:12.000Z", "5m")
	_, err := executeActionCommandC(cmd)

	if err == nil {
		t.Error("An error occurred: expect an error due to too many parameters (start, end and timeframe) used at the same time")
	}

	assert.EqualValues(t, "You can not use --end together with --timeframe", err.Error())
}

func TestTriggerEvaluationStartAndEndTimeWrongOrder(t *testing.T) {

	credentialmanager.MockAuthCreds = true

	*triggerEvaluation.Timeframe = ""
	*triggerEvaluation.Start = ""
	*triggerEvaluation.End = ""

	cmd := fmt.Sprintf("trigger evaluation --project=%s --stage=%s --service=%s "+
		"--start=%s --end=%s  --mock", "sockshop", "hardening", "carts", "2019-07-24T10:17:12.000Z", "2019-07-24T10:10:12.000Z")
	_, err := executeActionCommandC(cmd)

	if err == nil {
		t.Error("An error occurred: expect an error as end time is before start time")
	}

	assert.EqualValues(t, "Start and end time of evaluation time frame not set: end date must be at least 1 minute after start date", err.Error())
}

// TestTriggerEvaluationUnknownCommand
func TestTriggerEvaluationUnknownCommand(t *testing.T) {
	testInvalidInputHelper("trigger evaluation someUnknownCommand --project=sockshop --service=service --timeframe=5m --start=2019-10-31T11:59:59", "unknown command \"someUnknownCommand\" for \"keptn trigger evaluation\"", t)
}

// TestTriggerEvaluationUnknownParameter
func TestTriggerEvaluationUnknownParmeter(t *testing.T) {
	testInvalidInputHelper("trigger evaluation --projectt=sockshop --service=service --timeframe=5m --start=2019-10-31T11:59:59", "unknown flag: --projectt", t)
}
