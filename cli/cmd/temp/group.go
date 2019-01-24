// Copyright © 2018 NAME HERE <jbonds@jbvm.io>
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

var (
	groupList []string
)

var groupCmd = &cobra.Command{
	Use:     "group",
	Short:   "Search and Retrieve Group Info",
	Long:    `Example kafkactl group group1 group2`,
	Aliases: []string{"groups"},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			args = []string{""}
		}
		if targetTopic != "" {
			desc := []string{"group"}
			desc = append(desc, args...)
			describeCmd.Run(cmd, desc)
			return
		}
		if cmd.Flags().Changed("meta") || cmd.Flags().Changed("lag") || cmd.Flags().Changed("clientid") {
			desc := []string{"group"}
			desc = append(desc, args...)
			describeCmd.Run(cmd, desc)
			return
		}
		grps := searchGroupListMeta(args...)
		if len(grps) < 1 {
			closeFatal("no results found.\n")
		}
		printOutput(grps)
	},
}

func init() {
	rootCmd.AddCommand(groupCmd)
	groupCmd.Flags().BoolVarP(&exact, "exact", "x", false, "Find exact match")
	groupCmd.Flags().BoolVarP(&meta, "meta", "m", false, "Show extra/metadata details")
	groupCmd.Flags().BoolVarP(&showLag, "lag", "l", false, "Display Offset and Lag details (auto passes to --meta)")
	groupCmd.Flags().StringVarP(&clientID, "clientid", "i", "", "Find groups by ClientID")
}
