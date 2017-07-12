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
	Order     Order
	Reversed  bool
	Summarise bool
	ShowAll   bool
	FGColors  map[Priority]string
	BGColors  map[Priority]string
}

func NewViewOptions(showAll bool, summarise bool, order Order, reversed bool, fgColors map[Priority]string, bgColors map[Priority]string) *ViewOptions {
	return &ViewOptions{
		ShowAll:   showAll,
		Summarise: summarise,
		Order:     order,
		Reversed:  reversed,
		FGColors:  fgColors,
		BGColors:  bgColors,
	}
}

func (viewOptions *ViewOptions) GetBGColor(priority Priority) string {
	return bgColourEnumMap[viewOptions.BGColors[priority]]
}

func (viewOptions *ViewOptions) GetFGColor(priority Priority) string {
	return fgColourEnumMap[viewOptions.FGColors[priority]]
}

type View interface {
	ShowTree(tasks TaskList, options *ViewOptions)
	ShowTaskInfo(task Task, options *ViewOptions)
}

// TaskView is a filtered, ordered view of a Tasks children.
type TaskView struct {
	tasks   []Task
	options *ViewOptions
}

func CreateTaskView(node TaskNode, options *ViewOptions) *TaskView {
	view := &TaskView{
		tasks:   make([]Task, node.Len()),
		options: options,
	}
	for i := 0; i < node.Len(); i++ {
		view.tasks[i] = node.At(i)
	}
	sort.Sort(view)
	return view
}

func (t *TaskView) At(i int) Task {
	return t.tasks[i]
}

func (t *TaskView) Len() int {
	return len(t.tasks)
}

func (t *TaskView) Less(i, j int) bool {
	left := t.tasks[i]
	right := t.tasks[j]
	less := false
	switch t.options.Order {
	case INDEX:
		less = left.ID() < right.ID()
	case CREATED:
		less = left.CreationTime().Unix() < right.CreationTime().Unix()
	case COMPLETED:
		less = left.CompletionTime().Unix() < right.CompletionTime().Unix()
	case TEXT:
		less = left.Text() < right.Text()
	case PRIORITY:
		less = left.Priority() < right.Priority()
	case DURATION:
		var leftDuration, rightDuration int64
		leftCompletion := left.CompletionTime()
		rightCompletion := right.CompletionTime()
		if !leftCompletion.IsZero() {
			leftDuration = leftCompletion.Unix() - left.CreationTime().Unix()
		} else {
			leftDuration = 0
		}
		if !rightCompletion.IsZero() {
			rightDuration = rightCompletion.Unix() - right.CreationTime().Unix()
		} else {
			rightDuration = 0
		}
		less = leftDuration < rightDuration
	case DONE:
		less = !left.CompletionTime().IsZero() && !right.CompletionTime().IsZero()
	default:
		panic("invalid ordering")
	}
	if t.options.Reversed {
		less = !less
	}
	return less
}

func (t *TaskView) Swap(i, j int) {
	t.tasks[j], t.tasks[i] = t.tasks[i], t.tasks[j]
}
