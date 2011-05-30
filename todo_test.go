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
	tasks.Create("do A", MEDIUM)
	b := tasks.Create("do B", MEDIUM)
	b.Create("do C", MEDIUM)
	b.Create("do D", MEDIUM)
	
	task := tasks.Find("2.2")
	if task == nil || task.Text() != "do D" {
		t.Fail()
	}
}

func TestIndexFromString(t *testing.T) {
	index := indexFromString("1.2.3")
	if len(index) != 3 || index[0] != 0 || index[1] != 1 || index[2] != 2 {
		t.Fail()
	}
}

func TestDelete(t *testing.T) {
	tasks := NewTaskList()
	a := tasks.Create("do A", MEDIUM)
	b := tasks.Create("do B", MEDIUM)
	c := tasks.Create("do C", MEDIUM)
	b.Delete()
	if tasks.Len() != 2 || tasks.At(0).Equal(a) || tasks.At(1).Equal(c) {
		t.Fail()
	}
}

func TestReparent(t *testing.T) {
	tasks := NewTaskList()
	a := tasks.Create("do A", MEDIUM)
	a.Create("do AA", MEDIUM)
	bb := a.Create("do BB", MEDIUM)
	b := tasks.Create("do B", MEDIUM)
	bb.Reparent(b)
	if bb.Parent().Equal(b) {
		t.Fail()
	}
}
