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
	"strings"
	"testing"

	"github.com/keptn/keptn/cli/pkg/credentialmanager"
)

func TestSkipUpgradeCheck(t *testing.T) {
	noSkipMsg := "No upgrade path exists from Keptn version"
	skipMsg := "Skipping upgrade compatibility check!"
	credentialmanager.MockAuthCreds = true
	checkEndPointStatusMock = true
	mocking = true
	cmd := fmt.Sprintf("upgrade --mock")

	r := newRedirector()
	r.redirectStdOut()

	_, err := executeActionCommandC(cmd)
	out := r.revertStdOut()

	if !errorContains(err, noSkipMsg) {
		t.Errorf("upgrade Response did not match [no skip] : \nERROR: %v\nOUT: %v", err, out)
	}

	cmd = fmt.Sprintf("upgrade --skip-upgrade-check --mock")
	r = newRedirector()
	r.redirectStdOut()
	_, err = executeActionCommandC(cmd)
	out = r.revertStdOut()

	if !errorContains(err, "EOF") || !strings.Contains(out, skipMsg) {
		t.Errorf("upgrade Response did not match [skip] : \nERROR: %v\nOUT: %v", err, out)
	}
}
