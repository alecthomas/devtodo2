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
	"fmt"
	goopt "github.com/droundy/goopt"
	"os"
	"path/filepath"
	"strings"
)

// Actions
var addFlag = goopt.Flag([]string{"-a", "--add"}, nil, "add a task", "")
var markDoneFlag = goopt.Flag([]string{"-d", "--done"}, nil, "mark the given tasks as done", "")
var markNotDoneFlag = goopt.Flag([]string{"-D", "--not-done"}, nil, "mark the given tasks as not done", "")
var removeFlag = goopt.Flag([]string{"--remove"}, nil, "remove the given tasks", "")
var reparentFlag = goopt.Flag([]string{"-R", "--reparent"}, nil, "reparent task A below task B", "")
// Options
var priorityFlag = goopt.String([]string{"-p", "--priority"}, "medium", "priority of newly created tasks (veryhigh,high,medium,low,verylow)")
var graftFlag = goopt.String([]string{"-g", "--graft"}, "root", "task to graft new tasks to")
var fileFlag = goopt.String([]string{"--file"}, ".todo2", "file to load task lists from")
var legacyFileFlag = goopt.String([]string{"--legacy-file"}, ".todo", "file to load legacy task lists from")
var allFlag = goopt.Flag([]string{"-A", "--all"}, nil, "show all tasks, even completed ones", "")

func doView(tasks TaskList) {
	ConsoleView(tasks, *allFlag)
}

func doAdd(tasks TaskList, graft TaskNode, priority Priority, text string) {
	graft.Create(text, priority)
	saveTaskList(tasks)
}

func doMarkDone(tasks TaskList, references []Task) {
	for _, task := range references {
		task.SetCompleted()
	}
	saveTaskList(tasks)
}

func doMarkNotDone(tasks TaskList, references []Task) {
	for _, task := range references {
		task.SetCompletionTime(nil)
	}
	saveTaskList(tasks)
}

func doReparent(tasks TaskList, task Task, below Task) {
	task.Reparent(below)
	saveTaskList(tasks)
}

func doRemove(tasks TaskList, references []Task) {
	for _, task := range references {
		task.Delete()
	}
	saveTaskList(tasks)
}

func processAction(tasks TaskList) {
	priority := PriorityFromString(*priorityFlag)
	var graft TaskNode = tasks
	if *graftFlag != "root" {
		if graft = tasks.Find(*graftFlag); graft == nil {
			printFatal("invalid graft index '%s'", *graftFlag)
		}
	}

	switch {
	case *addFlag:
		if len(goopt.Args) == 0 {
			printFatal("expected text for new task")
		}
		text := strings.Join(goopt.Args, " ")
		doAdd(tasks, graft, priority, text)
	case *markDoneFlag:
		doMarkDone(tasks, resolveTaskReferences(tasks, goopt.Args))
	case *markNotDoneFlag:
		doMarkNotDone(tasks, resolveTaskReferences(tasks, goopt.Args))
	case *removeFlag:
		doRemove(tasks, resolveTaskReferences(tasks, goopt.Args))
	case *reparentFlag:
		if len(goopt.Args) != 2 {
			printFatal("expected <task> <new-parent> for reparenting")
		}
		doReparent(tasks, resolveTaskReference(tasks, goopt.Args[0]),
							 resolveTaskReference(tasks, goopt.Args[1]))
	default:
		doView(tasks)
	}
}

func resolveTaskReference(tasks TaskList, index string) Task {
	task := tasks.Find(index)
	if task == nil {
		printFatal("invalid task index %s", index)
	}
	return task
}

func resolveTaskReferences(tasks TaskList, indices []string) []Task {
	references := make([]Task, len(indices))
	for i, index := range indices {
		task := resolveTaskReference(tasks, index)
		references[i] = task
	}
	if len(references) == 0 {
		printFatal("no tasks provided to mark done")
	}
	return references
}

func loadTaskList() (TaskList, os.Error) {
	// Try loading new-style task file
	if file, err := os.Open(*fileFlag); err == nil {
		defer file.Close()
		loader := NewJsonIO()
		return loader.Deserialize(file)
	}
	// Try loading legacy task file
	if file, err := os.Open(*legacyFileFlag); err == nil {
		defer file.Close()
		loader := NewLegacyIO()
		return loader.Deserialize(file)
	}
	return nil, nil
}

func saveTaskList(tasks TaskList) (err os.Error) {
	path := *fileFlag
	if path, err = filepath.Abs(path); err != nil {
		return err
	}
	dir, file := filepath.Split(path)
	println(dir, file)
	// TODO Use a temporary file
	if file, err := os.Create(*fileFlag); err == nil {
		defer file.Close()
		writer := NewJsonIO()
		return writer.Serialize(file, tasks)
	}
	return nil
}

func main() {
	goopt.Suite = "DevTodo2"
	goopt.Version = "2.0"
	goopt.Author = "Alec Thomas <alec@swapoff.org>"
	goopt.Description = func() string {
		return `DevTodo is a program aimed specifically at programmers (but usable by anybody
at the terminal) to aid in day-to-day development.

It maintains a list of items that have yet to be completed, one for each
project directory. This allows the programmer to track outstanding bugs or
items that need to be completed with very little effort.

Items can be prioritised and are displayed in a hierarchy, so that one item may
depend on another.


todo2 [-A]
  Display (all) tasks.

todo2 [-p <priority>] -a <text>
  Create a new task.

todo2 -d <index>
  Mark a task as complete.`
	}
	goopt.Summary = "DevTodo2 - a hierarchical command-line task manager"
	goopt.Usage = func () string {
		return fmt.Sprintf("usage: %s [<options>] ...\n\n%s\n\n%s",
						   os.Args[0], goopt.Summary, goopt.Help())
	}
	goopt.Parse(nil)

	tasks, err := loadTaskList()
	if err != nil {
		printFatal("%s", err)
	}
	processAction(tasks)
}
