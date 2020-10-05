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

func TestVersionCmd(t *testing.T) {
	credentialmanager.MockAuthCreds = true
	mocking = true

	cmd := fmt.Sprintf("version")
	Version = "0.6.1"

	r := newRedirector()
	r.redirectStdOut()

	_, err := executeActionCommandC(cmd)
	if err != nil {
		t.Errorf(unexpectedErrMsg, err)
	}

	out := r.revertStdOut()
	if !strings.Contains(out, "CLI version: 0.6.1") {
		t.Errorf("unexpected used version: %s", out)
	}
	if !strings.Contains(out, "cluster version") {
		t.Error("expected cluster version")
	}
}
