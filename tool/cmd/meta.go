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
	"strings"

	"github.com/jbvmio/kafkactl"

	"github.com/spf13/cobra"
)

var metaCmd = &cobra.Command{
	Use:     "meta",
	Short:   "Return Metadata",
	Aliases: []string{"metadata"},
	Run: func(cmd *cobra.Command, args []string) {
		meta, err := client.GetClusterMeta()
		if err != nil {
			closeFatal("Error getting cluster metadata: %v\n", err)
		}
		if kafkaVer == "" {
			ver, _ := kafkactl.MatchKafkaVersion(getKafkaVersion(meta.APIMaxVersions))
			kafkaVer = ver.String()
		}
		c, err := client.Controller()
		if err != nil {
			closeFatal("Error obtaining controller: %v\n", err)
		}
		if len(meta.ErrorStack) > 0 {
			fmt.Println("ERRORs:")
			for _, e := range meta.ErrorStack {
				fmt.Printf(" %v\n", e)
			}
		}
		fmt.Println("\nBrokers: ", meta.BrokerCount())
		fmt.Println(" Topics: ", meta.TopicCount())
		fmt.Println(" Groups: ", meta.GroupCount())
		fmt.Printf("\nCluster: (Kafka: %v)\n", kafkaVer)
		for _, b := range meta.Brokers {
			if strings.Contains(b, c.Addr()) {
				fmt.Println("*", b)
			} else {
				fmt.Println(" ", b)
			}
		}
		fmt.Printf("\n(*)Controller\n\n")
	},
}

func init() {
	rootCmd.AddCommand(metaCmd)
	metaCmd.Flags().BoolVarP(&exact, "exact", "x", false, "Find exact match")
}
