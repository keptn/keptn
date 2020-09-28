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

const versionCheckInfo = `Daily version check is %s. 
Keptn will%s collect statistical data and will%s notify about new versions and security patches for Keptn. Details can be found at: https://keptn.sh/docs/0.7.x/reference/version_check
---------------------------------------------------
`
const setVersionCheckMsg = `* To %s the daily version check, please execute:
 - keptn set config AutomaticVersionCheck %s

`
const keptnReleaseDocsURL = "0.7.x"

const updateInfoMsg = `
Please visit https://keptn.sh for more information about updating.
`

const disableVersionCheckMsg = "To disable this notice, run: '%s set config AutomaticVersionCheck false'"

// KubeServerVersionConstraints the Kubernetes Cluster version's constraints is passed by ldflags
var KubeServerVersionConstraints string

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:     "version",
	Short:   "Shows the version of Keptn and Keptn CLI",
	Long:    `Shows the version of Keptn and Keptn CLI, and a note when a new version is available.`,
	Example: `keptn version`,
	Run: func(cmd *cobra.Command, args []string) {
		isLastCheckStale, err := isLastCheckStale()
		if err != nil {
			logging.PrintLog(err.Error(), logging.InfoLevel)
			return
		}

		var cliMsgPrinted, cliChecked, keptnMsgPrinted, keptnChecked bool

		// Keptn CLI
		fmt.Println("\nKeptn CLI version: " + Version)
		if isLastCheckStale {
			vChecker := version.NewVersionChecker()
			cliMsgPrinted, cliChecked = vChecker.CheckCLIVersion(Version, false)
		}

		// Keptn
		keptnVersion, err := getInstalledKeptnVersion()
		if err != nil {
			logging.PrintLog(err.Error(), logging.InfoLevel)
			return
		}
		fmt.Println("\nKeptn cluster version: " + keptnVersion)
		if isLastCheckStale {
			kvChecker := version.NewKeptnVersionChecker()
			keptnMsgPrinted, keptnChecked = kvChecker.CheckKeptnVersion(Version, keptnVersion, false)
		}

		if cliMsgPrinted || keptnMsgPrinted {
			fmt.Println(updateInfoMsg)
		}

		if cliChecked || keptnChecked {
			updateLastVersionCheck()
		}

		if err := printDailyVersionCheckInfo(); err != nil {
			logging.PrintLog(err.Error(), logging.InfoLevel)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

func isLastCheckStale() (bool, error) {
	configMng := config.NewCLIConfigManager()
	cliConfig, err := configMng.LoadCLIConfig()
	if err != nil {
		return false, err
	}
	return cliConfig.LastVersionCheck == nil || time.Now().Sub(*cliConfig.LastVersionCheck) >= time.Second, nil
}

func printDailyVersionCheckInfo() error {
	configMng := config.NewCLIConfigManager()
	cliConfig, err := configMng.LoadCLIConfig()
	if err != nil {
		return err
	}
	fmt.Println()
	if cliConfig.AutomaticVersionCheck {
		fmt.Printf(versionCheckInfo, "enabled", "", "")
		fmt.Printf(setVersionCheckMsg, "disable", "false")
	} else {
		fmt.Printf(versionCheckInfo, "disabled", " not", " not")
		fmt.Printf(setVersionCheckMsg, "enable", "true")
	}
	return nil
}

func updateLastVersionCheck() {
	configMng := config.NewCLIConfigManager()
	cliConfig, err := configMng.LoadCLIConfig()
	if err != nil {
		logging.PrintLog(err.Error(), logging.InfoLevel)
		return
	}
	currentTime := time.Now()
	cliConfig.LastVersionCheck = &currentTime
	if err := configMng.StoreCLIConfig(cliConfig); err != nil {
		logging.PrintLog(err.Error(), logging.InfoLevel)
	}
}
