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
	"time"

	"github.com/keptn/keptn/cli/pkg/logging"

	"github.com/keptn/keptn/cli/pkg/config"
	"github.com/keptn/keptn/cli/pkg/version"
	"github.com/spf13/cobra"
)

var (
	// Version information which is passed by ldflags
	Version string
)

const versionCheckInfo = "Daily version check is %s. Keptn will%s collect statistical data and will%s notify about new versions and security patches for Keptn. Details can be found at https://keptn.sh/docs/0.7.x/reference/version_check\n"
const enableVersionCheckMsg = "To %s the daily version check, please execute: \nkeptn set config AutomaticVersionCheck %s\n"

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:     "version",
	Short:   "Shows the CLI version of Keptn.",
	Long:    `Shows the CLI version of Keptn, as well as an indication whether a new version is available.`,
	Example: `keptn version`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("CLI version: " + Version)

		configMng := config.NewCLIConfigManager()
		cliConfig, err := configMng.LoadCLIConfig()
		if err != nil {
			logging.PrintLog(err.Error(), logging.InfoLevel)
			return
		}
		checkTime := time.Now()
		if cliConfig.LastVersionCheck == nil || checkTime.Sub(*cliConfig.LastVersionCheck) >= time.Second {
			vChecker := version.NewVersionChecker()
			vChecker.CheckCLIVersion(Version, false)
		}
		if cliConfig.AutomaticVersionCheck {
			fmt.Printf(versionCheckInfo, "enabled", "", "")
			fmt.Printf(enableVersionCheckMsg, "disable", "false")
		} else {
			fmt.Printf(versionCheckInfo, "disabled", " not", " not")
			fmt.Printf(enableVersionCheckMsg, "enable", "true")
		}
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
