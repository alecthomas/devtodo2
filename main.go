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
	"strconv"
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
			fatal("invalid graft index '%s'", *graftFlag)
		}
	}

	switch {
	case *addFlag:
		if len(goopt.Args) == 0 {
			fatal("expected text for new task")
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
			fatal("expected <task> <new-parent> for reparenting")
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
		fatal("invalid task index %s", index)
	}
	return task
}

func expandRange(indexRange string) []string {
	// This whole function makes me sad. This kind of manipulation of strings and
	// arrays just should not be this verbose.
	//
	// For constrast, in Python:
	//
	// def expand_range(index):
	//   start_index, end = index.split('-')
	//   start_index, start = start_index.rsplit('.', 1)
	//   for i in range(int(start), int(end) + 1):
	//     yield '%s.%s' % (start_index, str(i))
	ranges := strings.Split(indexRange, "-", -1)
	if len(ranges) != 2 {
		return nil
	}
	startIndex := strings.Split(ranges[0], ".", -1)
	start, err := strconv.Atoi(startIndex[len(startIndex) - 1])
	if err != nil { return nil }
	end, err := strconv.Atoi(ranges[1])
	if err != nil { return nil }
	rangeIndexes := make([]string, 0)
	for i := start; i <= end; i++ {
		index := startIndex[:len(startIndex) - 1]
		index = append(index, fmt.Sprintf("%d", i))
		rangeIndexes = append(rangeIndexes, strings.Join(index, "."))
	}
	return rangeIndexes
}

func resolveTaskReferences(tasks TaskList, indices []string) []Task {
	references := make([]Task, 0, len(indices))
	for _, index := range indices {
		if strings.Index(index, "-") == -1 {
			task := resolveTaskReference(tasks, index)
			references = append(references, task)
		} else {
			// Expand ranges. eg. 1.2-5 expands to 1.2 1.3 1.4 1.5
			indexes := expandRange(index)
			if indexes == nil {
				fatal("invalid task range %s", index)
			}
			for _, rangeIndex := range indexes {
				task := resolveTaskReference(tasks, rangeIndex)
				if task != nil {
					references = append(references, task)
				}
			}
		}
	}
	if len(references) == 0 {
		fatal("no tasks provided to mark done")
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
	temp := filepath.Join(dir, file + "~")
	if file, err := os.Create(temp); err == nil {
		defer func () {
			if err = file.Close(); err != nil {
				os.Remove(temp)
			} else {
				if err = os.Rename(temp, path); err != nil {
					fatal("unable to rename %s to %s", temp, path)
				}
			}
		}()
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
		fatal("%s", err)
	}
	processAction(tasks)
}