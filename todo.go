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
	"strconv"
	"strings"
	"time"
)

type Priority int

// Priority constants.
const (
	VERYLOW = Priority(iota)
	LOW
	MEDIUM
	HIGH
	VERYHIGH
)

type Order int

// Order constants.
const (
	CREATED = Order(iota)
	COMPLETED
	TEXT
	PRIORITY
	DURATION
	DONE
	INDEX
)

type TaskListIO interface {
	Deserialize(reader io.Reader) (TaskList, error)
	Serialize(writer io.Writer, tasks TaskList) error
}

type TaskNode interface {
	ID() int
	At(index int) Task
	Len() int
	Equal(other TaskNode) bool

	Parent() TaskNode
	SetParent(parent TaskNode)

	Append(child TaskNode)
	Create(title string, priority Priority, comment string) Task
	Delete()
}

type Task interface {
	TaskNode

	Text() string
	SetText(text string)

	Priority() Priority
	SetPriority(priority Priority)

	SetCreationTime(time time.Time)
	CreationTime() time.Time

	SetCompleted()
	SetCompletionTime(time time.Time)
	CompletionTime() time.Time

	// Extra attributes usable by extensions
	Attributes() map[string]string
	SetComment(comment string)
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

var priorityMapFromString = map[string]Priority{
	"veryhigh": VERYHIGH,
	"high":     HIGH,
	"medium":   MEDIUM,
	"low":      LOW,
	"verylow":  VERYLOW,
}

var priorityToString = map[Priority]string{
	VERYHIGH: "veryhigh",
	HIGH:     "high",
	MEDIUM:   "medium",
	VERYLOW:  "verylow",
	LOW:      "low",
}

func (p Priority) String() string {
	return priorityToString[p]
}

func PriorityFromString(priority string) Priority {
	r := strings.NewReplacer(" ", "")
	priority = r.Replace(priority)
	priority = strings.ToLower(priority)
	if p, ok := priorityMapFromString[priority]; ok {
		return p
	}
	return MEDIUM
}

var orderFromString = map[string]Order{
	"started":    CREATED,
	"start":      CREATED,
	"creation":   CREATED,
	"created":    CREATED,
	"finish":     COMPLETED,
	"finished":   COMPLETED,
	"completion": COMPLETED,
	"completed":  COMPLETED,
	"text":       TEXT,
	"priority":   PRIORITY,
	"length":     DURATION,
	"lifetime":   DURATION,
	"duration":   DURATION,
	"done":       DONE,
	"index":      INDEX,
}

var orderToString = map[Order]string{
	CREATED:   "created",
	COMPLETED: "completed",
	TEXT:      "text",
	PRIORITY:  "priority",
	DURATION:  "duration",
	DONE:      "done",
	INDEX:     "index",
}

func (self Order) String() string {
	return orderToString[self]
}

func OrderFromString(order string) (Order, bool) {
	reversed := false
	if len(order) >= 1 && order[0] == '-' {
		reversed = true
		order = order[1:]
	}
	if o, ok := orderFromString[order]; ok {
		return o, reversed
	}
	return PRIORITY, false
}

type taskNodeImpl struct {
	id         int
	tasks      []TaskNode
	parent     TaskNode
	attributes map[string]string
}

func newTaskNode(id int) *taskNodeImpl {
	return &taskNodeImpl{
		id:         id,
		parent:     nil,
		attributes: make(map[string]string),
	}
}

func (self *taskNodeImpl) ID() int {
	return self.id
}

func (self *taskNodeImpl) Equal(other TaskNode) bool {
	return self == other
}

func (self *taskNodeImpl) Len() int {
	return len(self.tasks)
}

func (self *taskNodeImpl) At(index int) Task {
	if index >= len(self.tasks) {
		return nil
	}
	return self.tasks[index].(Task)
}

func (self *taskNodeImpl) Parent() TaskNode {
	return self.parent
}

func (self *taskNodeImpl) SetParent(parent TaskNode) {
	self.parent = parent
}

func (self *taskNodeImpl) Append(child TaskNode) {
	child.SetParent(self)
	self.tasks = append(self.tasks, child)
}

func (self *taskNodeImpl) Create(title string, priority Priority, comment string) Task {
	task := newTask(self.Len(), title, priority, comment)
	self.Append(task)
	return task
}

func (self *taskNodeImpl) Delete() {
	parent := self.Parent().(*taskNodeImpl)
	if parent == nil {
		panic("can not delete root node")
	}
	for i := 0; i < parent.Len(); i++ {
		if parent.At(i).Equal(self) {
			parent.tasks = append(parent.tasks[:i], parent.tasks[i+1:]...)
			self.parent = nil
			return
		}
	}
	panic("couldn't find self in parent in order to delete")
}

type taskImpl struct {
	*taskNodeImpl
	text               string
	priority           Priority
	created, completed time.Time
	attributes         map[string]string
}

func newTask(id int, text string, priority Priority, comment string) Task {
	return &taskImpl{
		taskNodeImpl: newTaskNode(id),
		text:         text,
		priority:     priority,
		created:      time.Now().UTC(),
		completed:    time.Time{},
		attributes:   map[string]string{"comment": comment},
	}
}

func (self *taskImpl) ID() int {
	return self.id
}

func (self *taskImpl) SetCreationTime(time time.Time) {
	self.created = time
}

func (self *taskImpl) CreationTime() time.Time {
	return self.created
}

func (self *taskImpl) SetCompleted() {
	self.SetCompletionTime(time.Now().UTC())
}

func (self *taskImpl) SetCompletionTime(time time.Time) {
	self.completed = time
}

func (self *taskImpl) CompletionTime() time.Time {
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

func (self *taskImpl) Attributes() map[string]string {
	return self.attributes
}

func (self *taskImpl) SetComment(comment string) {
	self.attributes["comment"] = comment
}

type taskListImpl struct {
	*taskNodeImpl
	title string
}

func NewTaskList() TaskList {
	return &taskListImpl{
		taskNodeImpl: newTaskNode(-1),
		title:        "",
	}
}

// Convert "1.2.3" to int[]{0, 1, 2} ready for indexing into TaskNodes
func indexFromString(index string) Index {
	tokens := strings.Split(index, ".")
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

func (self *taskListImpl) ID() int {
	return -1
}

func (self *taskListImpl) Find(index string) Task {
	numericIndex := indexFromString(index)
	if numericIndex == nil {
		return nil
	}
	var node TaskNode = self // -golint
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

func ReparentTask(node TaskNode, below TaskNode) {
	node.Delete()
	below.Append(node)
}
