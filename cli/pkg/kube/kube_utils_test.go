package kube

import (
	"errors"
	"testing"

	keptnutils "github.com/keptn/go-utils/pkg/lib"
	"github.com/stretchr/testify/require"
)

const v115 = `Client Version: version.Info{Major:"1", Minor:"17", GitVersion:"v1.17.2", GitCommit:"59603c6e503c87169aea6106f57b9f242f64df89", GitTreeState:"clean", BuildDate:"2020-01-23T14:21:54Z", GoVersion:"go1.13.6", Compiler:"gc", Platform:"darwin/amd64"}
Server Version: version.Info{Major:"1", Minor:"15", GitVersion:"v1.15.5", GitCommit:"20c265fef0741dd71a66480e35bd69f18351daea", GitTreeState:"clean", BuildDate:"2019-10-15T19:07:57Z", GoVersion:"go1.12.10", Compiler:"gc", Platform:"linux/amd64"}`
const v9999 = `Client Version: version.Info{Major:"1", Minor:"17", GitVersion:"v1.17.2", GitCommit:"59603c6e503c87169aea6106f57b9f242f64df89", GitTreeState:"clean", BuildDate:"2020-01-23T14:21:54Z", GoVersion:"go1.13.6", Compiler:"gc", Platform:"darwin/amd64"}
Server Version: version.Info{Major:"99", Minor:"99", GitVersion:"v1.15.5", GitCommit:"20c265fef0741dd71a66480e35bd69f18351daea", GitTreeState:"clean", BuildDate:"2019-10-15T19:07:57Z", GoVersion:"go1.12.10", Compiler:"gc", Platform:"linux/amd64"}`
const v113Plus = `Client Version: version.Info{Major:"1", Minor:"17", GitVersion:"v1.17.2", GitCommit:"59603c6e503c87169aea6106f57b9f242f64df89", GitTreeState:"clean", BuildDate:"2020-01-23T14:21:54Z", GoVersion:"go1.13.6", Compiler:"gc", Platform:"darwin/amd64"}
Server Version: version.Info{Major:"1", Minor:"13+d", GitVersion:"v1.15.5", GitCommit:"20c265fef0741dd71a66480e35bd69f18351daea", GitTreeState:"clean", BuildDate:"2019-10-15T19:07:57Z", GoVersion:"go1.12.10", Compiler:"gc", Platform:"linux/amd64"}`
const v113PlusBeta = `Client Version: version.Info{Major:"1", Minor:"17", GitVersion:"v1.17.2", GitCommit:"59603c6e503c87169aea6106f57b9f242f64df89", GitTreeState:"clean", BuildDate:"2020-01-23T14:21:54Z", GoVersion:"go1.13.6", Compiler:"gc", Platform:"darwin/amd64"}
Server Version: version.Info{Major:"1", Minor:"13+beta", GitVersion:"v1.15.5", GitCommit:"20c265fef0741dd71a66480e35bd69f18351daea", GitTreeState:"clean", BuildDate:"2019-10-15T19:07:57Z", GoVersion:"go1.12.10", Compiler:"gc", Platform:"linux/amd64"}`

var checkSplitTests = []struct {
	constraints    string
	isNewerVersion bool
	err            string
	executeOutput  string
	executeError   error
}{
	{">= 1.13, <= 1.15", false, "", v115, nil},
	{">= 1.13, <= 1.15", true, "", v9999, nil},
	{"< 1.13", true, "", v115, nil},
	{"1.13", true, "", v115, nil},
	{"wrong constraints", false, "Malformed constraint: wrong constraints", v115, nil},
	{">= 1.13, <= 1.15", false, "execute error", v115, errors.New("execute error")},
	{">= 1.13, <= 1.15", false, "Server Version not found: no version", "no version", nil},
	{">= 1.13, <= 1.15", false, "", v113Plus, nil},
	{">= 1.13, <= 1.15", false, "", v113PlusBeta, nil},
	{">= 1.14, <= 1.15", false, "The Kubernetes Server Version '1.13+beta' doesn't satisfy constraints '>= 1.14, <= 1.15'", v113PlusBeta, nil},
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
			isNewerVersion, err := CheckKubeServerVersion(tt.constraints)

			if tt.err == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tt.err)
			}
			require.Equal(t, isNewerVersion, tt.isNewerVersion)
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

func TestCheckDeploymentManagedByHelm(t *testing.T) {
	defer func() {
		executeCommandFunc = keptnutils.ExecuteCommand
	}()

	responseWithManagedByLabelValueHelm := `{"metadata":{"labels":{"app.kubernetes.io/managed-by":"Helm"}}}`
	responseWithManagedByLabelNoValue := `{"metadata":{"labels":{"app.kubernetes.io/managed-by":""}}}`
	responseWithoutManagedByLabel := `{"metadata":{"labels":{}}}`

	var tests = []struct {
		name           string
		deploymentName string
		executedCmd    func(string, []string) (string, error)
		result         bool
		resulterr      bool
	}{
		{"CheckDeploymentManagedByHelm - labelAndValueHelm", "my-deployment", func(string, []string) (string, error) { return responseWithManagedByLabelValueHelm, nil }, true, false},
		{"CheckDeploymentManagedByHelm - labelWithoutAValue", "my-deployment", func(string, []string) (string, error) { return responseWithManagedByLabelNoValue, nil }, false, false},
		{"CheckDeploymentManagedByHelm - noLabelPresent", "my-deployment", func(string, []string) (string, error) { return responseWithoutManagedByLabel, nil }, false, false},
		{"CheckDeploymentManagedByHelm - errorOnKubectlCmd", "my-deployment", func(string, []string) (string, error) { return "", errors.New("Whoops...") }, false, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			executeCommandFunc = tt.executedCmd
			result, err := CheckDeploymentManagedByHelm(tt.deploymentName, "ns")
			if tt.resulterr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tt.result, result)
		})
	}
}

func TestCheckDeploymentAvailable(t *testing.T) {
	defer func() {
		executeCommandFunc = keptnutils.ExecuteCommand
	}()

	responseWithDeployment := `{"items":[{"metadata":{"name":"my-deployment"}}]}`
	responseWithNoDeployment := `{"items":[{}]}`

	var tests = []struct {
		name           string
		deploymentName string
		executedCmd    func(string, []string) (string, error)
		result         bool
		resulterr      bool
	}{
		{"CheckDeploymentAvailable - responseWithDeploymentAvailable", "my-deployment", func(string, []string) (string, error) { return responseWithDeployment, nil }, true, false},
		{"CheckDeploymentAvailable - responseWithDeploymentAvailable", "my-deployment", func(string, []string) (string, error) { return responseWithNoDeployment, nil }, false, false},
		{"CheckDeploymentAvailable - responseWithDeploymentAvailable", "my-deployment", func(string, []string) (string, error) { return "", errors.New("Whoops...") }, false, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			executeCommandFunc = tt.executedCmd
			result, err := CheckDeploymentAvailable(tt.deploymentName, "ns")
			if tt.resulterr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tt.result, result)
		})
	}
}
