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
	"container/vector"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

type Priority int

const (
	VERYLOW = Priority(iota)
	LOW
	MEDIUM
	HIGH
	VERYHIGH
)

type TaskListIO interface {
	Deserialize(reader io.Reader) (TaskList, os.Error)
	Serialize(writer io.Writer, tasks TaskList) os.Error
}

type TaskNode interface {
	// Return an iterator over child tasks. nil if no children.
	At(index int) Task
	Len() int
	Equal(other TaskNode) bool

	Parent() TaskNode
	SetParent(parent TaskNode)

	Append(child TaskNode)
	Create(title string, priority Priority) Task
	Reparent(below TaskNode)
	Delete()
}

type Task interface {
	TaskNode

	Text() string
	SetText(text string)

	Priority() Priority
	SetPriority(priority Priority)

	SetCreationTime(time *time.Time)
	CreationTime() *time.Time

	SetCompleted()
	SetCompletionTime(time *time.Time)
	CompletionTime() *time.Time
}

type TaskList interface {
	TaskNode

	Title() string
	SetTitle(title string)

	Find(index string) Task
}

// Index referencing a task
type Index []int

// Implementation

var priorityMapFromString map[string]Priority = map[string]Priority {
	"veryhigh": VERYHIGH,
	"high": HIGH,
	"medium": MEDIUM,
	"low": LOW,
	"verylow": VERYLOW,
}

var priorityToString map[Priority]string = map[Priority]string {
	VERYHIGH: "veryhigh",
	HIGH: "high",
	MEDIUM: "medium",
	VERYLOW: "verylow",
	LOW: "low",
}

func (p Priority) String() string {
	return priorityToString[p]
}

func PriorityFromString(priority string) Priority {
	if p, ok := priorityMapFromString[priority]; ok {
		return p
	}
	return MEDIUM
}

type taskNodeImpl struct {
	tasks vector.Vector
	parent TaskNode
}

func newTaskNode() *taskNodeImpl {
	return &taskNodeImpl{
		parent: nil,
	}
}

func (self *taskNodeImpl) Equal(other TaskNode) bool {
	return self == other
}

func (self *taskNodeImpl) Len() int {
	return self.tasks.Len()
}

func (self *taskNodeImpl) At(index int) Task {
	if index >= len(self.tasks) {
		return nil
	}
	return self.tasks.At(index).(Task)
}

func (self *taskNodeImpl) Parent() TaskNode {
	return self.parent
}

func (self *taskNodeImpl) SetParent(parent TaskNode) {
	self.parent = parent
}

func (self *taskNodeImpl) Append(child TaskNode) {
	child.SetParent(self)
	self.tasks.Push(child)
}

func (self *taskNodeImpl) Create(title string, priority Priority) Task {
	task := newTask(title, priority)
	self.Append(task)
	return task
}

func (self *taskNodeImpl) Reparent(below TaskNode) {
	self.Delete()
	below.Append(self)
}

func (self *taskNodeImpl) Delete() {
	parent := self.Parent().(*taskNodeImpl)
	if parent == nil {
		panic("can not delete root node")
	}
	for i := 0; i < parent.Len(); i++ {
		if parent.At(i).Equal(self) {
			parent.tasks.Delete(i)
			self.parent = nil
			return
		}
	}
	panic("couldn't find self in parent in order to delete")
}

type taskImpl struct {
	*taskNodeImpl
	text string
	priority Priority
	created, completed *time.Time
}

func newTask(text string, priority Priority) Task {
	return &taskImpl{
		taskNodeImpl: newTaskNode(),
		text: text,
		priority: priority,
		created: time.UTC(),
		completed: nil,
	}
}

func (self *taskImpl) SetCreationTime(time *time.Time) {
	self.created = time
}

func (self *taskImpl) CreationTime() *time.Time {
	return self.created
}

func (self *taskImpl) SetCompleted() {
	self.SetCompletionTime(time.UTC())
}

func (self *taskImpl) SetCompletionTime(time *time.Time) {
	self.completed = time
}

func (self *taskImpl) CompletionTime() *time.Time {
	return self.completed
}

func (self *taskImpl) Text() string {
	return self.text
}

func (self *taskImpl) SetText(text string) {
	self.text = text
}

func (self *taskImpl) Priority() Priority {
	return self.priority
}

func (self *taskImpl) SetPriority(priority Priority) {
	self.priority = priority
}

type taskListImpl struct {
	*taskNodeImpl
	title string
}

func NewTaskList() TaskList {
	return &taskListImpl{
		taskNodeImpl: newTaskNode(),
		title: "",
	}
}

// Convert "1.2.3" to int[]{0, 1, 2} ready for indexing into TaskNodes
func indexFromString(index string) Index {
	tokens := strings.Split(index, ".", -1)
	numericIndex := make(Index, len(tokens))
	for i, token := range tokens {
		value, err := strconv.Atoi(token)
		if err != nil || value < 1 {
			return nil
		}
		numericIndex[i] = value - 1
	}
	return numericIndex
}

func (self *taskListImpl) Find(index string) Task {
	numericIndex := indexFromString(index)
	if numericIndex == nil {
		return nil
	}
	var node TaskNode = self
	for _, i := range numericIndex {
		if node = node.At(i); node == nil {
			return nil
		}
	}
	return node.(Task)
}

func (self *taskListImpl) Title() string {
	return self.title
}

func (self *taskListImpl) SetTitle(title string) {
	self.title = title
}
