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
)

func deleteTopic(topic string) {
	errd = client.RemoveTopic(topic)
	if errd != nil {
		closeFatal("Error deleting topic: %v\n", errd)
	}
	fmt.Println("\nSuccessfully deleted topic", topic, "\n")
}

func deleteGroup(group string) {
	errd = client.RemoveGroup(group)
	if errd != nil {
		closeFatal("Error removing group: %v\n", errd)
	}
	fmt.Println("\nSuccessfully removed group", group, "\n")
}

func deleteToOffset(topic string, partition int32, offset int64) {
	errd = client.DeleteToOffset(topic, partition, offset)
	if errd != nil {
		closeFatal("Error deleting to offset: %v\n", errd)
	}
	fmt.Println("\nSuccessfully deleted to offset", offset, "\n")
}
