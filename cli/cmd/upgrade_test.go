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

	"github.com/keptn/keptn/cli/pkg/credentialmanager"
	"github.com/keptn/keptn/cli/pkg/version"
)

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

	testUpgraderCmd := NewUpgraderCommand(vChecker)

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
					"continuousDelivery": map[string]interface{}{
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
					"continuousDelivery": map[string]interface{}{
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
					"continuousDelivery": map[string]interface{}{
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
					"continuousDelivery": map[string]interface{}{
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
