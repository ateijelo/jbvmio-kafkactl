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
	"fmt"

	"github.com/jbvmio/kafkactl"
	"github.com/spf13/cobra"
)

var (
	topicList   []string
	refreshMeta bool
)

var topicCmd = &cobra.Command{
	Use:   "topic",
	Short: "Search and Retrieve Available Topics",
	Long: `Provides a summary view of available topics.
  Example: kafkactl topic topic1 topic2 topic3
  
If no arguments are provided, all topics are retrieved.
To see detailed metadata information, use the meta command or the -m flag here.
  Example: kafkactl --broker kafkahost topic1 --exact`,
	Aliases: []string{"topics"},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			args = []string{""}
		}
		if preAllTopics {
			tom := chanGetTopicOffsetMap(findPRETopics(args...))
			if len(tom) < 1 {
				fmt.Println("\n No Topics Require PRE.\n")
				return
			}
			printOutput(tom)
			return
		}
		if meta || cmd.Flags().Changed("leaders") {
			desc := []string{"topic"}
			desc = append(desc, args...)
			describeCmd.Run(cmd, desc)
			return
		}
		if refreshMeta {
			refreshMetadata(args...)
			return
		}
		printOutput(kafkactl.GetTopicSummaries(searchTopicMeta(args...)))
	},
}

func init() {
	rootCmd.AddCommand(topicCmd)
	topicCmd.Flags().BoolVarP(&exact, "exact", "x", false, "Find exact match")
	topicCmd.Flags().BoolVarP(&meta, "meta", "m", false, "Show extra/metadata details")
	topicCmd.Flags().BoolVar(&preAllTopics, "needpre", false, "Show Topics that need a Preferred Leader Election")
	topicCmd.Flags().StringVar(&leaderList, "leaders", "", `Only show specified Leaders. (eg "1,3,7"; auto passes to --meta)`)
	topicCmd.Flags().BoolVarP(&refreshMeta, "refresh-metadata", "r", false, "Query the Cluster to Refresh the Available Metadata for the given Topic(s)")
}
