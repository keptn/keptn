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

// getEventCmd represents the get command
var getEventCmd = &cobra.Command{
	Use:     "event [eventType]",
	Aliases: []string{"events"},
	Short:   `Get event in combination with "eventType" allows to retrieve a Keptn event`,
	Long: `Get event in combination with "eventType" allows to retrieve a Keptn event. Get event without subcommand cannot be used.

Example:
	keptn get event [eventType]`,
}

func init() {
	getCmd.AddCommand(getEventCmd)
}
