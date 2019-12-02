package cmd

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"

	"github.com/keptn/keptn/cli/pkg/logging"
	"github.com/keptn/keptn/cli/utils/credentialmanager"
)

func init() {
	logging.InitLoggers(os.Stdout, os.Stdout, os.Stderr)
}

// TestEvaluationStart tests the start-evaluation command.
func TestEvaluationStart(t *testing.T) {

	credentialmanager.MockAuthCreds = true
	buf := new(bytes.Buffer)
	rootCmd.SetOutput(buf)

	*evaluationStart.Timeframe = ""
	*evaluationStart.Start = ""
	*evaluationStart.End = ""

	args := []string{
		"send",
		"event",
		"start-evaluation",
		fmt.Sprintf("--project=%s", "sockshop"),
		fmt.Sprintf("--stage=%s", "hardening"),
		fmt.Sprintf("--service=%s", "carts"),
		fmt.Sprintf("--timeframe=%s", "5m"),
		fmt.Sprintf("--labels=foo=bar,bar=foo"),
		"--mock",
	}
	rootCmd.SetArgs(args)
	err := rootCmd.Execute()

	if err != nil {
		t.Errorf("An error occured: %v", err)
	}
}

func TestEvaluationStartWrongFormat(t *testing.T) {

	credentialmanager.MockAuthCreds = true
	buf := new(bytes.Buffer)
	rootCmd.SetOutput(buf)

	*evaluationStart.Timeframe = ""
	*evaluationStart.Start = ""
	*evaluationStart.End = ""

	args := []string{
		"send",
		"event",
		"start-evaluation",
		fmt.Sprintf("--project=%s", "sockshop"),
		fmt.Sprintf("--stage=%s", "hardening"),
		fmt.Sprintf("--service=%s", "carts"),
		fmt.Sprintf("--timeframe=%s", "5h"),
		"--mock",
	}
	rootCmd.SetArgs(args)
	err := rootCmd.Execute()

	if err == nil {
		t.Error("An error occured: expect an error due to wrong time frame format")
	}
}

func TestEvaluationStartTimeSpecified(t *testing.T) {

	credentialmanager.MockAuthCreds = true
	buf := new(bytes.Buffer)
	rootCmd.SetOutput(buf)

	*evaluationStart.Timeframe = ""
	*evaluationStart.Start = ""
	*evaluationStart.End = ""

	args := []string{
		"send",
		"event",
		"start-evaluation",
		fmt.Sprintf("--project=%s", "sockshop"),
		fmt.Sprintf("--stage=%s", "hardening"),
		fmt.Sprintf("--service=%s", "carts"),
		fmt.Sprintf("--timeframe=%s", "5m"),
		fmt.Sprintf("--start=%s", "2019-07-24T10:17:12"),
		"--mock",
	}
	rootCmd.SetArgs(args)
	err := rootCmd.Execute()

	if err != nil {
		t.Errorf("An error occured: %v", err)
	}
}

func TestEvaluationStartAndEndTimeSpecified(t *testing.T) {

	credentialmanager.MockAuthCreds = true
	buf := new(bytes.Buffer)
	rootCmd.SetOutput(buf)

	*evaluationStart.Timeframe = ""
	*evaluationStart.Start = ""
	*evaluationStart.End = ""

	args := []string{
		"send",
		"event",
		"start-evaluation",
		fmt.Sprintf("--project=%s", "sockshop"),
		fmt.Sprintf("--stage=%s", "hardening"),
		fmt.Sprintf("--service=%s", "carts"),
		fmt.Sprintf("--start=%s", "2019-07-24T10:17:12"),
		fmt.Sprintf("--end=%s", "2019-07-24T10:20:12"),
		"--mock",
	}

	rootCmd.SetArgs(args)
	err := rootCmd.Execute()

	if err != nil {
		t.Errorf("An error occured: %v", err)
	}
}

func TestEvaluationStartAndEndTimeAndTimeframeSpecified(t *testing.T) {

	credentialmanager.MockAuthCreds = true
	buf := new(bytes.Buffer)
	rootCmd.SetOutput(buf)

	*evaluationStart.Timeframe = ""
	*evaluationStart.Start = ""
	*evaluationStart.End = ""

	args := []string{
		"send",
		"event",
		"start-evaluation",
		fmt.Sprintf("--project=%s", "sockshop"),
		fmt.Sprintf("--stage=%s", "hardening"),
		fmt.Sprintf("--service=%s", "carts"),
		fmt.Sprintf("--start=%s", "2019-07-24T10:17:12"),
		fmt.Sprintf("--end=%s", "2019-07-24T10:20:12"),
		fmt.Sprintf("--timeframe=%s", "5m"),
		"--mock",
	}
	rootCmd.SetArgs(args)
	err := rootCmd.Execute()

	if err == nil {
		t.Error("An error occured: expect an error due to too many parameters (start, end and timeframe) used at the same time")
	}

	assert.EqualValues(t, "Start and end time of evaluation time frame not set: You can not use --end together with --timeframe", err.Error())
}

func TestEvaluationStartAndEndTimeWrongOrder(t *testing.T) {

	credentialmanager.MockAuthCreds = true
	buf := new(bytes.Buffer)
	rootCmd.SetOutput(buf)

	*evaluationStart.Timeframe = ""
	*evaluationStart.Start = ""
	*evaluationStart.End = ""

	args := []string{
		"send",
		"event",
		"start-evaluation",
		fmt.Sprintf("--project=%s", "sockshop"),
		fmt.Sprintf("--stage=%s", "hardening"),
		fmt.Sprintf("--service=%s", "carts"),
		fmt.Sprintf("--start=%s", "2019-07-24T10:17:12"),
		fmt.Sprintf("--end=%s", "2019-07-24T10:10:12"),
		"--mock",
	}
	rootCmd.SetArgs(args)
	err := rootCmd.Execute()
	if err == nil {
		t.Error("An error occured: expect an error as end time is before start time")
	}

	assert.EqualValues(t, "Start and end time of evaluation time frame not set: end date must be at least 1 minute after start date", err.Error())
}
