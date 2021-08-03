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
	"github.com/spf13/cobra"
)

var pauseSequenceParams sequenceControlStruct

var pauseSequenceCmd = &cobra.Command{
	Use:          "sequence",
	Short:        "Pauses the execution of a sequence",
	Long:         `Pauses the execution of a sequence. Currently running task(s) will not be paused.`,
	Example:      `keptn pause sequence --project <my-project> --keptn-context <keptn-context>`,
	SilenceUsage: true,
	Args: func(cmd *cobra.Command, args []string) error {
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return PauseSequence(pauseSequenceParams)
	},
}

func init() {
	pauseCmd.AddCommand(pauseSequenceCmd)
	pauseSequenceParams.keptnContext = pauseSequenceCmd.Flags().StringP("keptn-context", "c", "",
		"The Keptn context the sequence execution is bound to")
	pauseSequenceParams.project = pauseSequenceCmd.Flags().StringP("project", "p", "",
		"The Keptn project the sequence belongs to")
	pauseSequenceParams.stage = pauseSequenceCmd.Flags().StringP("stage", "s", "",
		"The Keptn stage in which the sequence shall be paused")
	pauseSequenceCmd.MarkFlagRequired("keptn-context")
	pauseSequenceCmd.MarkFlagRequired("project")
}
