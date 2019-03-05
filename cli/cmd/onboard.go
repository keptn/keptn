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

	"github.com/keptn/keptn/cli/utils"
	"github.com/spf13/cobra"
)

// onboardCmd represents the onboard command
var onboardCmd = &cobra.Command{
	Use:   "onboard [service]",
	Short: "onboard allows to onbard a new service",
	Long:  `onbaord currently allows to onboard a new service with \"onboard service\". onboard without subcommand cannot be used.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		utils.Info.Println("onboard called")
		return errors.New("onboard can only be called in combination with \"service\"")
	},
}

func init() {
	rootCmd.AddCommand(onboardCmd)
}
