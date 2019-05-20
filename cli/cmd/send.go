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

	"github.com/spf13/cobra"
)

var eventFilePath *string

// sendCmd represents the send command
var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "Sends a keptn event.",
	Long: `Allows to send arbitrary keptn events, which are defined in the passed file.

Example:
	keptn send --event=new_artifact.json`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("send called")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(sendCmd)
	eventFilePath = serviceCmd.Flags().StringP("event", "e", "", "The file containing the event as Cloud Event in JSON.")
}
