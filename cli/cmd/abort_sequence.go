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

var abortSequenceParams sequenceControlStruct

var abortSequenceCmd = &cobra.Command{
	Use:          "sequence",
	Short:        "Aborts the execution of a sequence",
	Long:         `Aborts the execution of a sequence. Currently running task(s) will not be aborted.`,
	Example:      `keptn abort sequence --project <my-project> --keptn-context <keptn-context>`,
	SilenceUsage: true,
	Args: func(cmd *cobra.Command, args []string) error {
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := AbortSequence(abortSequenceParams); err != nil {
			return err
		}
		fmt.Println("Successfully aborted sequence")
		return nil
	},
}

func init() {
	abortCmd.AddCommand(abortSequenceCmd)
	abortSequenceParams.keptnContext = abortSequenceCmd.Flags().StringP("keptn-context", "c", "",
		"The Keptn context the sequence execution is bound to")
	abortSequenceParams.project = abortSequenceCmd.Flags().StringP("project", "p", "",
		"The Keptn project the sequence belongs to")
	abortSequenceParams.stage = abortSequenceCmd.Flags().StringP("stage", "s", "",
		"The Keptn stage in which the sequence shall be aborted")
	abortSequenceCmd.MarkFlagRequired("keptn-context")
	abortSequenceCmd.MarkFlagRequired("project")
}
