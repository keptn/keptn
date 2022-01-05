package go_tests

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {
	if err := setup(); err != nil {
		fmt.Printf("TestMain: error while setting up the tests: %v", err)
		os.Exit(-1)
	}
	code := m.Run()

	events, err := GetOOMEvents()
	if code == 0 && (len(events.Items) != 0 || err != nil) {
		println("There were some out of memory Errors!")
		os.Exit(-1)
	}

	os.Exit(code)
}

func setup() error {
	// before executing the tests, we check whether the context of the Keptn CLI matches the one of the kubectl CLI
	// i.e. The kubectl CLI should be connected to the cluster the Keptn CLI is currently authenticated against.
	// this prevents unintended kubectl commands from being executed against a different cluster than the one containing the Keptn instance that should be tested
	match, err := endpointsMatch()
	if err != nil {
		return fmt.Errorf("could not compare endpoints of kubectl context and keptn CLI: %s", err.Error())
	}
	if !match {
		return errors.New("endpoint mismatch between CLI and kubectl detected")
	}
	return nil
}

func endpointsMatch() (bool, error) {
	_, keptnAPIURL, err := GetApiCredentials()
	if err != nil {
		return false, err
	}
	statusCmdOutput, err := ExecuteCommand("keptn status")
	if err != nil {
		return false, err
	}
	statusOutputLines := strings.Split(statusCmdOutput, "\n")

	var apiURLFromStatusCommand string
	for _, line := range statusOutputLines {
		if strings.Contains(line, "Successfully authenticated") {
			endpointLineSplit := strings.Split(line, " ")
			apiURLFromStatusCommand = endpointLineSplit[len(endpointLineSplit)-1]
			break
		}
	}

	if apiURLFromStatusCommand != keptnAPIURL {
		return false, nil
	}
	return true, nil
}
