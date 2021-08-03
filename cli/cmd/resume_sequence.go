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

var resumeSequenceParams sequenceControlStruct

var resumeSequenceCmd = &cobra.Command{
	Use:          "sequence",
	Short:        "Resumes the execution of a sequence",
	Example:      `keptn resume sequence --project <my-project> --keptn-context <keptn-context>`,
	SilenceUsage: true,
	Args: func(cmd *cobra.Command, args []string) error {
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return ResumeSequence(resumeSequenceParams)
	},
}

func init() {
	resumeCmd.AddCommand(resumeSequenceCmd)
	resumeSequenceParams.keptnContext = resumeSequenceCmd.Flags().StringP("keptn-context", "c", "",
		"The Keptn context the sequence execution is bound to")
	resumeSequenceParams.project = resumeSequenceCmd.Flags().StringP("project", "p", "",
		"The Keptn project the sequence belongs to")
	resumeSequenceParams.stage = resumeSequenceCmd.Flags().StringP("stage", "s", "",
		"The Keptn stage in which the sequence shall be resumed")
	resumeSequenceCmd.MarkFlagRequired("keptn-context")
	resumeSequenceCmd.MarkFlagRequired("project")
}
