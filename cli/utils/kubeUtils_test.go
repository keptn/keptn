package utils

import (
	"errors"
	"testing"

	keptnutils "github.com/keptn/go-utils/pkg/utils"
	"github.com/stretchr/testify/require"
)

const v115 = `Client Version: version.Info{Major:"1", Minor:"17", GitVersion:"v1.17.2", GitCommit:"59603c6e503c87169aea6106f57b9f242f64df89", GitTreeState:"clean", BuildDate:"2020-01-23T14:21:54Z", GoVersion:"go1.13.6", Compiler:"gc", Platform:"darwin/amd64"}
Server Version: version.Info{Major:"1", Minor:"15", GitVersion:"v1.15.5", GitCommit:"20c265fef0741dd71a66480e35bd69f18351daea", GitTreeState:"clean", BuildDate:"2019-10-15T19:07:57Z", GoVersion:"go1.12.10", Compiler:"gc", Platform:"linux/amd64"}`
const v9999 = `Client Version: version.Info{Major:"1", Minor:"17", GitVersion:"v1.17.2", GitCommit:"59603c6e503c87169aea6106f57b9f242f64df89", GitTreeState:"clean", BuildDate:"2020-01-23T14:21:54Z", GoVersion:"go1.13.6", Compiler:"gc", Platform:"darwin/amd64"}
Server Version: version.Info{Major:"99", Minor:"99", GitVersion:"v1.15.5", GitCommit:"20c265fef0741dd71a66480e35bd69f18351daea", GitTreeState:"clean", BuildDate:"2019-10-15T19:07:57Z", GoVersion:"go1.12.10", Compiler:"gc", Platform:"linux/amd64"}`

var checkSplitTests = []struct {
	constraints   string
	err           string
	executeOutput string
	executeError  error
}{
	{">= 1.13, <= 1.15", "", v115, nil},
	{">= 1.13, <= 1.15", "The Kubernetes Server Version '99.99' doesn't satisfy constraints '>= 1.13, <= 1.15'", v9999, nil},
	{"< 1.13", "The Kubernetes Server Version '1.15' doesn't satisfy constraints '< 1.13'", v115, nil},
	{"wrong constraints", "Malformed constraint: wrong constraints", v115, nil},
	{">= 1.13, <= 1.15", "execute error", v115, errors.New("execute error")},
	{">= 1.13, <= 1.15", "Server Version not found: no version", "no version", nil},
}

func TestSplitCheckKubeServerVersion(t *testing.T) {
	defer func() {
		executeCommandFunc = keptnutils.ExecuteCommand
	}()
	var executeOutput string
	var executeError error
	executeCommandFunc = func(string, []string) (string, error) {
		return executeOutput, executeError
	}
	for _, tt := range checkSplitTests {
		t.Run(tt.constraints, func(t *testing.T) {
			executeOutput = tt.executeOutput
			executeError = tt.executeError
			err := CheckKubeServerVersion(tt.constraints)
			if tt.err == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tt.err)
			}
		})
	}
}

func TestIsKubectlAvailable(t *testing.T) {
	defer func() {
		executeCommandFunc = keptnutils.ExecuteCommand
	}()
	executeCommandFunc = func(string, []string) (string, error) {
		return "true", nil
	}
	res, err := IsKubectlAvailable()
	require.True(t, res)
	require.NoError(t, err)

	executeCommandFunc = func(string, []string) (string, error) {
		return "false", errors.New("fake error")
	}
	res, err = IsKubectlAvailable()
	require.False(t, res)
	require.EqualError(t, err, "fake error")
}
