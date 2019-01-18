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
	"bytes"
	"sort"
	"strings"

	"github.com/jbvmio/kafkactl"
	"github.com/spf13/cast"
)

func searchTopicMeta(topics ...string) []kafkactl.TopicMeta {
	tMeta, err := client.GetTopicMeta()
	if err != nil {
		closeFatal("Error getting topic metadata: %s\n", err)
	}
	var topicMeta []kafkactl.TopicMeta
	if len(topics) >= 1 {
		if topics[0] != "" {
			for _, t := range topics {
				for _, m := range tMeta {
					if exact {
						if m.Topic == t {
							topicMeta = append(topicMeta, m)
						}
					} else {
						if strings.Contains(m.Topic, t) {
							topicMeta = append(topicMeta, m)
						}
					}
				}
			}
		} else {
			topicMeta = tMeta
		}
	}
	sort.Slice(topicMeta, func(i, j int) bool {
		if topicMeta[i].Topic < topicMeta[j].Topic {
			return true
		}
		if topicMeta[i].Topic > topicMeta[j].Topic {
			return false
		}
		return topicMeta[i].Partition < topicMeta[j].Partition
	})
	return topicMeta
}

func getTopicOffsetMap(tm []kafkactl.TopicMeta) []kafkactl.TopicOffsetMap {
	return client.MakeTopicOffsetMap(tm)
}

func filterTOMByLeader(tom []kafkactl.TopicOffsetMap, leaders []int32) []kafkactl.TopicOffsetMap {
	var TOM []kafkactl.TopicOffsetMap
	done := make(map[string]bool)
	for _, t := range tom {
		var topicMeta []kafkactl.TopicMeta
		if !done[t.Topic] {
			done[t.Topic] = true
			for _, tm := range t.TopicMeta {
				for _, leader := range leaders {
					if t.PartitionLeaders[tm.Partition] == leader {
						topicMeta = append(topicMeta, tm)
					}
				}
			}
		}
		tom := kafkactl.TopicOffsetMap{
			Topic:            t.Topic,
			TopicMeta:        topicMeta,
			PartitionOffsets: t.PartitionOffsets,
			PartitionLeaders: t.PartitionLeaders,
		}
		TOM = append(TOM, tom)
	}
	return TOM
}

func chanGetTopicOffsetMap(t []kafkactl.TopicMeta) []kafkactl.TopicOffsetMap {
	var TOM []kafkactl.TopicOffsetMap
	var count int
	tmMap := make(map[string][]kafkactl.TopicMeta)
	for _, tm := range t {
		tmMap[tm.Topic] = append(tmMap[tm.Topic], tm)
	}
	done := make(map[string]bool, len(tmMap))
	tomChan := make(chan []kafkactl.TopicOffsetMap, 100)
	for t, meta := range tmMap {
		if !done[t] {
			count++
			done[t] = true
		}
		go chanMakeTOM(client, meta, tomChan)
	}
	for i := 0; i < count; i++ {
		tom := <-tomChan
		TOM = append(TOM, tom...)
	}
	return TOM
}

func chanMakeTOM(client *kafkactl.KClient, tMeta []kafkactl.TopicMeta, tomChan chan []kafkactl.TopicOffsetMap) {
	tom := client.MakeTopicOffsetMap(tMeta)
	tomChan <- tom
	return
}

func refreshMetadata(topics ...string) {
	errd = client.RefreshMetadata(topics...)
	if errd != nil {
		closeFatal("Error refreshing topic metadata: %v\n", errd)
	}
}

func findPRETopics(topics ...string) []kafkactl.TopicMeta {
	var preMeta []kafkactl.TopicMeta
	tMeta := searchTopicMeta(topics...)
	for _, tm := range tMeta {
		if len(tm.Replicas) > 0 {
			if tm.Leader != tm.Replicas[0] {
				preMeta = append(preMeta, tm)
			}
		}
	}
	return preMeta
}

func validateLeaders(leaders []int32) {
	pMap := make(map[int32]bool, len(leaders))
	for _, p := range leaders {
		if pMap[p] {
			closeFatal("Error: invalid leader/brokerID entered or duplicate.")
		} else {
			pMap[p] = true
		}
	}
}

// Parses Stdin passed from kafkactl topic metadata
func parseTopicStdin(b []byte) []topicStdinData {
	bits := bytes.TrimSpace(b)
	lines := string(bits)

	var topicData []topicStdinData
	a := strings.Split(lines, "\n")
	headers := strings.Fields(strings.TrimSpace(a[0]))
	if len(headers) < 3 {
		closeFatal("Invalid Stdin Passed\n")
	}
	if headers[0] != "TOPIC" || headers[1] != "PART" || headers[2] != "OFFSET" {
		closeFatal("Best to pass stdin through kafkactl itself.\n")
	}
	for _, b := range a[1:] {
		td := topicStdinData{}
		b := strings.TrimSpace(b)
		td.topic = cutField(b, 1)
		td.partition = cast.ToInt32(cutField(b, 2))
		topicData = append(topicData, td)
	}
	return topicData
}

type topicStdinData struct {
	topic     string
	partition int32
}
