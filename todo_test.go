/*
  Copyright 2011 Alec Thomas

  Licensed under the Apache License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.
*/

package main

import (
	"testing"
)

func TestFind(t *testing.T) {
	tasks := NewTaskList()
	tasks.AddTask("do A", MEDIUM)
	b := tasks.AddTask("do B", MEDIUM)
	b.AddTask("do C", MEDIUM)
	b.AddTask("do D", MEDIUM)
	
	task := tasks.Find([]int{1, 1})
	if task == nil || task.Text() != "do D" {
		t.Fail()
	}
}

func TestIndexFromString(t *testing.T) {
	index := IndexFromString("1.2.3")
	if len(index) != 3 || index[0] != 0 || index[1] != 1 || index[2] != 2 {
		t.Fail()
	}
}
