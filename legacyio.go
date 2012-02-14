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

// Loads legacy devtodo XML files.

package main

import (
	"encoding/xml"
	"errors"

	"io"
	"strconv"
	"strings"
	"time"
)

type legacyIO struct {
	TaskListIO
}

type xmlNote struct {
	Priority string `xml:"attr"`
	Time     string `xml:"attr"`
	Done     string `xml:"attr"`
	Text     string `xml:"chardata"`
	Note     []xmlNote
}

type xmlTodo struct {
	Title string
	Note  []xmlNote
}

func parseXmlNote(parent TaskNode, from []xmlNote) {
	if from == nil {
		return
	}
	for _, note := range from {
		text := strings.TrimSpace(note.Text)
		priority := PriorityFromString(note.Priority)

		task := parent.Create(text, priority)

		created, _ := strconv.ParseInt(note.Time, 10, 64)
		completed, _ := strconv.ParseInt(note.Done, 10, 64)
		task.SetCreationTime(time.Unix(created, 0).UTC())
		if completed != 0 {
			task.SetCompletionTime(time.Unix(completed, 0).UTC())
		}
		parseXmlNote(task, note.Note)
	}
}

func NewLegacyIO() TaskListIO {
	return &legacyIO{}
}

func (self *legacyIO) Deserialize(reader io.Reader) (tasks TaskList, err error) {
	todoXml := &xmlTodo{}
	if err = xml.NewDecoder(reader).Decode(&todoXml); err != nil {
		return
	}

	tasks = NewTaskList()
	tasks.SetTitle(strings.TrimSpace(todoXml.Title))
	parseXmlNote(tasks, todoXml.Note)
	return
}

func (self *legacyIO) Serialize(writer io.Writer, tasks TaskList) error {
	return errors.New("serialization to legacy format not supported")
}
