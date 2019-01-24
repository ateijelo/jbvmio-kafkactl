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
	"os"
	"strings"

	"github.com/jbvmio/kafkactl"

	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

var (
	showLag    bool
	leaderList string
	clientID   string
)

const (
	targetMatch = true
)

var describeCmd = &cobra.Command{
	Use:   "describe",
	Short: "Return Topic or Group details",
	Long: `Examples:
  kafkactl describe group myConsumerGroup
  kafkactl describe topic myTopic`,
	Aliases: []string{"desc", "descr", "des", "get"},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			cmd.Help()
			return
		}
		target := args[0]
		args = args[1:]
		switch targetMatch {
		case strings.Contains(target, "gro"):
			grps := searchGroupMeta(args...)
			if cmd.Flags().Changed("topic") {
				grps = groupMetaByTopic(targetTopic, grps)
			}
			if cmd.Flags().Changed("clientid") {
				grps = groupMetaByMember(clientID, grps)
			}
			if len(grps) < 1 {
				closeFatal("no results for that group/topic combination\n")
			}
			if showLag {
				var partitionLag []PartitionLag
				if useFast {
					partitionLag = chanGetPartitionLag(grps)
				} else {
					partitionLag = getPartitionLag(grps)
				}
				printOutput(partitionLag)
				return
			}
			printOutput(grps)
			return
		case strings.Contains(target, "top"):
			var tom []kafkactl.TopicOffsetMap
			if useFast {
				tom = chanGetTopicOffsetMap(searchTopicMeta(args...))
			} else {
				tom = getTopicOffsetMap(searchTopicMeta(args...))
			}
			if len(tom) < 1 {
				closeFatal("no results for that group/topic combination\n")
			}
			if cmd.Flags().Changed("leaders") {
				var leaders []int32
				ldrs := strings.Split(leaderList, ",")
				for _, l := range ldrs {
					leaders = append(leaders, cast.ToInt32(l))
				}
				validateLeaders(leaders)
				tom = filterTOMByLeader(tom, leaders)
				printOutput(tom)
				return
			}
			printOutput(tom)
			return
		default:
			fmt.Println("no such resource to describe:", target)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(describeCmd)
	describeCmd.Flags().BoolVarP(&exact, "exact", "x", false, "Find exact match")
	describeCmd.Flags().BoolVarP(&showLag, "lag", "l", false, "Display Offset and Lag details")
	describeCmd.Flags().StringVar(&leaderList, "leaders", "", `Only show specified Leaders. (eg "1,3,7")`)
	describeCmd.Flags().StringVarP(&clientID, "clientid", "i", "", "Find groups by ClientID")
}
