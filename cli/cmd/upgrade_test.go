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
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/keptn/keptn/cli/pkg/common"
	commonfake "github.com/keptn/keptn/cli/pkg/common/fake"
	"github.com/keptn/keptn/cli/pkg/credentialmanager"
	credmanagerfake "github.com/keptn/keptn/cli/pkg/credentialmanager/fake"
	"github.com/keptn/keptn/cli/pkg/helm"
	helmfake "github.com/keptn/keptn/cli/pkg/helm/fake"
	"github.com/keptn/keptn/cli/pkg/kube"
	kubefake "github.com/keptn/keptn/cli/pkg/kube/fake"
	"github.com/keptn/keptn/cli/pkg/version"
	"github.com/stretchr/testify/assert"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/release"
)

func TestUpgradeCmdHandler_doUpgradePreRunCheck(t *testing.T) {
	type fields struct {
		vChecker          *version.KeptnVersionChecker
		helmHelper        *helmfake.IHelperMock
		namespaceHandler  *kubefake.IKeptnNamespaceHandlerMock
		userInput         *commonfake.IUserInputMock
		credentialManager *credmanagerfake.CredentialManagerInterfaceMock
	}
	tests := []struct {
		name              string
		fields            fields
		args              installUpgradeParams
		chartsToBeApplied []*chart.Chart
		wantErr           bool
	}{
		{
			name: "upgrade pre-run check: namespace exists, end check",
			fields: fields{
				vChecker: version.NewKeptnVersionChecker(),
				helmHelper: &helmfake.IHelperMock{
					DownloadChartFunc: func(chartRepoURL string) (*chart.Chart, error) {
						return nil, errors.New("DownloadChart should not be called")
					},
					GetHistoryFunc: func(releaseName string, namespace string) ([]*release.Release, error) {
						return nil, errors.New("GetHistory should not be called")
					},
				},
				namespaceHandler:  &kubefake.IKeptnNamespaceHandlerMock{},
				userInput:         &commonfake.IUserInputMock{},
				credentialManager: &credmanagerfake.CredentialManagerInterfaceMock{},
			},
			args: installUpgradeParams{
				PatchNamespace: boolp(true),
			},
			chartsToBeApplied: []*chart.Chart{},
			wantErr:           false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UpgradeCmdHandler{
				helmHelper:        tt.fields.helmHelper,
				namespaceHandler:  tt.fields.namespaceHandler,
				userInput:         tt.fields.userInput,
				credentialManager: tt.fields.credentialManager,
			}

			if err := u.doUpgradePreRunCheck(tt.fields.vChecker, tt.args); (err != nil) != tt.wantErr {
				t.Errorf("doUpgradePreRunCheck error = %v, wantErr %v", err, tt.wantErr)
			}

			for index, upgradeChartCall := range tt.fields.helmHelper.UpgradeChartCalls() {
				assert.Equal(t, tt.chartsToBeApplied[index], upgradeChartCall.Ch)
			}
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
	cm := credentialmanager.NewCredentialManager(assumeYes)

	testUpgradeCmd := NewUpgradeCommand(vChecker, helmHelper, namespaceHandler, userInput, cm)

	upgradeCmd.RunE = testUpgradeCmd.RunE
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
