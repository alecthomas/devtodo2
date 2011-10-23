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
	"sort"
)

type ViewOptions struct {
	Order Order
	Reversed bool
	Summarise bool
	ShowAll bool
}

type View interface {
	ShowTree(tasks TaskList, options *ViewOptions)
	ShowTaskInfo(task Task)
}

// A filtered, ordered view of a Tasks children.
type TaskView struct {
	tasks []Task
	options *ViewOptions
}

func CreateTaskView(node TaskNode, options *ViewOptions) *TaskView {
	view := &TaskView{
		tasks: make([]Task, node.Len()),
		options: options,
	}
	for i := 0; i < node.Len(); i++ {
		view.tasks[i] = node.At(i)
	}
	sort.Sort(view)
	return view
}

func (self *TaskView) At(i int) Task {
	return self.tasks[i]
}

func (self *TaskView) Len() int {
	return len(self.tasks)
}

func (self *TaskView) Less(i, j int) bool {
	left := self.tasks[i]
	right := self.tasks[j]
	less := false
	switch self.options.Order {
	case CREATED:
		less = left.CreationTime().Seconds() < right.CreationTime().Seconds()
	case COMPLETED:
		less = left.CompletionTime().Seconds() < right.CompletionTime().Seconds()
	case TEXT:
		less = left.Text() < right.Text()
	case PRIORITY:
		less = left.Priority() < right.Priority()
	case DURATION:
		var leftDuration, rightDuration int64
		leftCompletion := left.CompletionTime()
		rightCompletion := right.CompletionTime()
		if leftCompletion != nil {
			leftDuration = leftCompletion.Seconds() - left.CreationTime().Seconds()
		} else {
			leftDuration = 0
		}
		if rightCompletion != nil {
			rightDuration = rightCompletion.Seconds() - right.CreationTime().Seconds()
		} else {
			rightDuration = 0
		}
		less = leftDuration < rightDuration
	case DONE:
		less = left.CompletionTime() != nil && right.CompletionTime() == nil
	default:
		panic("invalid ordering")
	}
	if !self.options.Reversed {
		less = !less
	}
	return less
}

func (self *TaskView) Swap(i, j int) {
	self.tasks[j], self.tasks[i] = self.tasks[i], self.tasks[j]
}


