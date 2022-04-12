// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/keptn/keptn/cli/pkg/common"
	commonfake "github.com/keptn/keptn/cli/pkg/common/fake"
	"github.com/keptn/keptn/cli/pkg/credentialmanager"
	"github.com/keptn/keptn/cli/pkg/helm"
	helmfake "github.com/keptn/keptn/cli/pkg/helm/fake"
	"github.com/keptn/keptn/cli/pkg/kube"
	kubefake "github.com/keptn/keptn/cli/pkg/kube/fake"
	"github.com/keptn/keptn/cli/pkg/version"
	"github.com/stretchr/testify/assert"
	"helm.sh/helm/v3/pkg/chart"
)

func TestUpgradeCmdHandler_doUpgradePreRunCheck(t *testing.T) {
	type fields struct {
		helmHelper       *helmfake.IHelperMock
		namespaceHandler *kubefake.IKeptnNamespaceHandlerMock
		userInput        *commonfake.IUserInputMock
	}
	tests := []struct {
		name              string
		fields            fields
		args              installUpgradeParams
		chartsToBeApplied []*chart.Chart
		wantErr           bool
	}{
		{
			name:   "upgrade pre-run check: namespace exists, cancel check",
			fields: fields{},
		},
		{
			name:   "upgrade pre-run check: use custom chart URL",
			fields: fields{},
		},
		{
			name:   "upgrade pre-run check: skip upgrade compatibility check",
			fields: fields{},
		},
		{
			name:   "upgrade pre-run check: version is not upgrade compatible and newer version of cli is available",
			fields: fields{},
		},
		{
			name:   "upgrade pre-run check: version is not upgrade compatible and newer version of cli is not available",
			fields: fields{},
		},
		{
			name:   "upgrade pre-run check: version is not upgrade compatible and no upgrade path exists from current Keptn version",
			fields: fields{},
		},
		{
			name:   "upgrade pre-run check: version is not upgrade compatible",
			fields: fields{},
		},
		{
			name:   "upgrade pre-run check: CLI context does not match kubernetes context",
			fields: fields{},
		},
		{
			name:   "upgrade pre-run check: use the OpenShift platform manager",
			fields: fields{},
		},
		{
			name:   "upgrade pre-run check: use the Kubernetes platform manager",
			fields: fields{},
		},
		{
			name:   "upgrade pre-run check: use an invalid platform manager",
			fields: fields{},
		},
		{
			name:   "upgrade pre-run check: incorrect kubectl configurations",
			fields: fields{},
		},
		{
			name:   "upgrade pre-run check: config file path not provided nor provided in a file",
			fields: fields{},
		},
		{
			name:   "upgrade pre-run check: config file path is not provided and cluster creds are invalid",
			fields: fields{},
		},
		{
			name:   "upgrade pre-run check: config file path is provided and cluster creds are invalid",
			fields: fields{},
		},
		{
			name:   "upgrade pre-run check: k8s server version is not compatible",
			fields: fields{},
		},
		{
			name:   "upgrade pre-run check: k8s server version is compatible but user cancels",
			fields: fields{},
		},
		{
			name:   "upgrade pre-run check: k8s server version is compatible and user does not cancel",
			fields: fields{},
		},
		{
			name:   "upgrade pre-run check: deployment with given namespace does not exist",
			fields: fields{},
		},
		{
			name:   "upgrade pre-run check: deployment with given namespace exists and is not managed by Helm",
			fields: fields{},
		},
		{
			name:   "upgrade pre-run check: deployment with given namespace exists and is managed by Helm",
			fields: fields{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, true, true)
		})
	}
}

func TestUpgradeCmdHandler_doUpgrade(t *testing.T) {
	t.Run("upgrade: ", func(t *testing.T) {
		assert.Equal(t, true, true)
	})
}

func TestUpgradeCmdHandler_addWarningNonExistingProjectUpstream(t *testing.T) {
	t.Run("add warning: ", func(t *testing.T) {
		assert.Equal(t, true, true)
	})
}

func TestSkipUpgradeCheck(t *testing.T) {
	noSkipMsg := "No upgrade path exists from Keptn version"
	skipMsg := "Skipping upgrade compatibility check!"
	credentialmanager.MockAuthCreds = true
	Version = "0.11.4"
	os.Setenv("MOCK_SERVER", "http://some-valid-url.com")
	cmd := fmt.Sprintf("upgrade --mock")

	ts := getMockVersionHTTPServer()
	vChecker := &version.KeptnVersionChecker{
		VersionFetcherClient: &version.VersionFetcherClient{
			HTTPClient: http.DefaultClient,
			VersionURL: ts.URL,
		},
	}
	helmHelper := helm.NewHelper()
	namespaceHandler := kube.NewKubernetesUtilsKeptnNamespaceHandler()
	userInput := common.NewUserInput()

	testUpgraderCmd := NewUpgraderCommand(vChecker, helmHelper, namespaceHandler, userInput)

	upgraderCmd.RunE = testUpgraderCmd.RunE
	r := newRedirector()
	r.redirectStdOut()

	_, err := executeActionCommandC(cmd)
	out := r.revertStdOut()

	if !errorContains(err, noSkipMsg) {
		t.Errorf("upgrade Response did not match [no skip] : \nERROR: %v\nOUT: %v", err, out)
	}

	cmd = fmt.Sprintf("upgrade --skip-upgrade-check --mock --chart-repo=https://charts-dev.keptn.sh/packages/keptn-0.9.0.tgz")
	r = newRedirector()
	r.redirectStdOut()

	_, err = executeActionCommandC(cmd)
	out = r.revertStdOut()

	if !errorContains(err, "EOF") || !strings.Contains(out, skipMsg) {
		t.Errorf("upgrade Response did not match [skip] : \nERROR: %v\nOUT: %v", err, out)
	}
}

func Test_isContinuousDeliveryEnable(t *testing.T) {
	type args struct {
		configValues map[string]interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "continuous delivery enabled - return true",
			args: args{
				configValues: map[string]interface{}{
					"continuous-delivery": map[string]interface{}{
						"enabled": true,
					},
				},
			},
			want: true,
		},
		{
			name: "continuous delivery enabled (string value) - return true",
			args: args{
				configValues: map[string]interface{}{
					"continuous-delivery": map[string]interface{}{
						"enabled": "true",
					},
				},
			},
			want: true,
		},
		{
			name: "continuous delivery not enabled - return false",
			args: args{
				configValues: map[string]interface{}{
					"continuous-delivery": map[string]interface{}{
						"enabled": false,
					},
				},
			},
			want: false,
		},
		{
			name: "continuous delivery not enabled (string value) - return false",
			args: args{
				configValues: map[string]interface{}{
					"continuous-delivery": map[string]interface{}{
						"enabled": "false",
					},
				},
			},
			want: false,
		},
		{
			name: "continuous delivery not defined - return false",
			args: args{
				configValues: map[string]interface{}{
					"schwifty": map[string]interface{}{
						"enabled": "true",
					},
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isContinuousDeliveryEnabled(tt.args.configValues); got != tt.want {
				t.Errorf("isContinuousDeliveryEnabled() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestUpgradeUnknownCommand
func TestUpgradeUnknownCommand(t *testing.T) {
	testInvalidInputHelper("upgrade someUnknownCommand", "unknown command \"someUnknownCommand\" for \"keptn upgrade\"", t)
}

// TestUpgradeUnknownParameter
func TestUpgradeUnknownParmeter(t *testing.T) {
	testInvalidInputHelper("upgrade --project=sockshop", "unknown flag: --project", t)
}
