//go:build !nokubectl
// +build !nokubectl

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

// installCmd represents the version command
var upgraderCmd = NewUpgraderCommand()

func NewUpgraderCommand() *cobra.Command {
	upgradeCmd := &cobra.Command{
		Use:          "upgrade",
		Deprecated:   MsgDeprecatedUseHelm,
		Args:         cobra.NoArgs,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("this command is deprecated, "+MsgDeprecatedUseHelm, Version)
		},
	}
	return upgradeCmd
}

func init() {
	rootCmd.AddCommand(upgraderCmd)
}
