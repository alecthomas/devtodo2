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
	"io"
	"json"
	"os"
	"time"
)

type jsonIO struct {
	TaskListIO
}

func NewJsonIO() TaskListIO {
	return &jsonIO{}
}

func (self *jsonIO) Deserialize(reader io.Reader) (tasks TaskList, err os.Error) {
	decoder := json.NewDecoder(reader)
	mtl := &marshalableTaskList{}
	if err = decoder.Decode(&mtl); err == nil {
		tasks = fromMarshalableTaskList(mtl)
	}
	return
}

func (self *jsonIO) Serialize(writer io.Writer, tasks TaskList) (err os.Error) {
	translated := toMarshalableTaskList(tasks)
	data, err := json.MarshalIndent(translated, "", "  ")
	if err != nil {
		return err
	}
	_, err = writer.Write(data)
	return
}

// Utilitie functions and structures for marshaling

type marshalableTask struct {
	Text string "text"
	Priority string "priority"
	Creation int64 "creation"
	Completion int64 "completion"
	Tasks []*marshalableTask "tasks"
}

type marshalableTaskList struct {
	Title string "title"
	Tasks []*marshalableTask "tasks"
}

func toMarshalableTaskList(t TaskList) *marshalableTaskList {
	return &marshalableTaskList{
		Title: t.Title(),
		Tasks: toMarshalableTask(t),
	}
}

func toMarshalableTask(n TaskNode) []*marshalableTask {
	children := make([]*marshalableTask, n.Len())
	for i := 0; i < n.Len(); i++ {
		t := n.At(i)
		var created, completed int64 = 0, 0
		if t.CreationTime() != nil {
			created = t.CreationTime().Seconds()
		}
		if t.CompletionTime() != nil {
			completed = t.CompletionTime().Seconds()
		}
		children[i] = &marshalableTask{
			Text: t.Text(),
			Priority: t.Priority().String(),
			Creation: created,
			Completion: completed,
			Tasks: toMarshalableTask(t),
		}
	}
	return children
}

func fromMarshalableTaskList(l *marshalableTaskList) TaskList {
	tasks := NewTaskList()
	tasks.SetTitle(l.Title)
	fromMarshalableTask(tasks, l.Tasks)
	return tasks
}

func fromMarshalableTask(node TaskNode, t []*marshalableTask) {
	for _, j := range t {
		task := node.Create(j.Text, PriorityFromString(j.Priority))
		task.SetCreationTime(time.SecondsToUTC(j.Creation))
		if j.Completion != 0 {
			task.SetCompletionTime(time.SecondsToUTC(j.Completion))
		}
		fromMarshalableTask(task, j.Tasks)
	}
}
