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
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/spf13/cobra"

	apiutils "github.com/keptn/go-utils/pkg/api/utils"

	"github.com/keptn/keptn/cli/pkg/config"
	"github.com/keptn/keptn/cli/pkg/credentialmanager"
	"github.com/keptn/keptn/cli/pkg/logging"
	"github.com/keptn/keptn/cli/pkg/version"
)

var (
	// Version information which is passed by ldflags
	Version string
)

const keptnReleaseDocsURL = "0.8.x" // ToDo: Can we automate this?

const versionCheckInfo = `Daily version check is %s. 
Keptn will%s collect statistical data and will%s notify about new versions and security patches for Keptn. Details can be found at: https://keptn.sh/docs/` + keptnReleaseDocsURL + `/reference/version_check
---------------------------------------------------
`
const setVersionCheckMsg = `* To %s the daily version check, please execute:
 - keptn set config AutomaticVersionCheck %s

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

		var cliChecked, keptnChecked bool

		// Keptn CLI
		fmt.Println("\nKeptn CLI version: " + Version)
		if isLastCheckStale {
			vChecker := version.NewVersionChecker()
			cliChecked, _ = vChecker.CheckCLIVersion(Version, false)
		}

		// Keptn
		fmt.Print("Keptn cluster version: ")
		keptnVersion, err := getKeptnServerVersion()
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(keptnVersion)
			if isLastCheckStale {
				kvChecker := version.NewKeptnVersionChecker()
				keptnChecked, _ = kvChecker.CheckKeptnVersion(Version, keptnVersion, false)
			}
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

func getKeptnServerVersion() (string, error) {
	var endPoint url.URL
	var apiToken string
	var err error
	if !mocking {
		endPoint, apiToken, err = credentialmanager.NewCredentialManager(false).GetCreds(namespace)
	} else {
		endPointPtr, _ := url.Parse(os.Getenv("MOCK_SERVER"))
		endPoint = *endPointPtr
		apiToken = ""
	}

	if err != nil {
		return "", errors.New(authErrorMsg)
	}
	if endPointErr := checkEndPointStatus(endPoint.String()); endPointErr != nil {
		return "", fmt.Errorf("Error connecting to server: %s"+endPointErrorReasons, endPointErr)
	}
	apiHandler := apiutils.NewAuthenticatedAPIHandler(endPoint.String(), apiToken, "x-token", nil, endPoint.Scheme)
	metadataData, errMetadata := apiHandler.GetMetadata()
	if errMetadata != nil {
		if errMetadata.Message != nil {
			return "", errors.New("Error occurred with response code " + strconv.FormatInt(errMetadata.Code, 10) + " with message " + *errMetadata.Message)
		}
		return "", errors.New("received invalid response from Keptn API")
	}
	return metadataData.Keptnversion, nil
}
