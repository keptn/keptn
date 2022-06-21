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
	"testing"
)

// TestUpgradeUnknownCommand
func TestUpgradeUnknownCommand(t *testing.T) {
	testInvalidInputHelper("upgrade someUnknownCommand", "unknown command \"someUnknownCommand\" for \"keptn upgrade\"", t)
}

// TestUpgradeUnknownParameter
func TestUpgradeUnknownParameter(t *testing.T) {
	testInvalidInputHelper("upgrade --project=sockshop", "unknown flag: --project", t)
}

// TestUpgradeDeprecated
func TestUpgradeDeprecated(t *testing.T) {
	testInvalidInputHelper("upgrade", "this command is deprecated", t)
}
