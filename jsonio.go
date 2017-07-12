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
	"encoding/json"
	"io"
	"time"
)

type jsonIO struct {
	TaskListIO
}

func NewJSONIO() TaskListIO {
	return &jsonIO{}
}

func (j *jsonIO) Deserialize(reader io.Reader) (tasks TaskList, err error) {
	decoder := json.NewDecoder(reader)
	mtl := &marshalableTaskList{}
	if err = decoder.Decode(&mtl); err == nil {
		tasks = fromMarshalableTaskList(mtl)
	}
	return
}

func (j *jsonIO) Serialize(writer io.Writer, tasks TaskList) (err error) {
	translated := toMarshalableTaskList(tasks)
	data, err := json.MarshalIndent(translated, "", "  ")
	if err != nil {
		return err
	}
	_, err = writer.Write(data)
	return
}

// Utility functions and structures for marshaling

type marshalableTask struct {
	Text       string             `json:"text"`
	Priority   string             `json:"priority"`
	Creation   int64              `json:"creation"`
	Completion int64              `json:"completion,omitempty"`
	Tasks      []*marshalableTask `json:"tasks,omitempty"`
}

type marshalableTaskList struct {
	Title string             `json:"title"`
	Tasks []*marshalableTask `json:"tasks"`
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
		if !t.CreationTime().IsZero() {
			created = t.CreationTime().Unix()
		}
		if !t.CompletionTime().IsZero() {
			completed = t.CompletionTime().Unix()
		}
		children[i] = &marshalableTask{
			Text:       t.Text(),
			Priority:   t.Priority().String(),
			Creation:   created,
			Completion: completed,
			Tasks:      toMarshalableTask(t),
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
		task.SetCreationTime(time.Unix(j.Creation, 0).UTC())
		if j.Completion != 0 {
			task.SetCompletionTime(time.Unix(j.Completion, 0).UTC())
		}
		fromMarshalableTask(task, j.Tasks)
	}
}
