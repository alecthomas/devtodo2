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
	"os"
	"strconv"
	"strings"
	"time"

	"gopkg.in/alecthomas/kingpin.v2"
)

const usage = `DevTodo2 - a hierarchical command-line task manager


DevTodo is a program aimed specifically at programmers (but usable by anybody
at the terminal) to aid in day-to-day development.

It maintains a list of items that have yet to be completed, one list for each
project directory. This allows the programmer to track outstanding bugs or
items that need to be completed with very little effort.

Items can be prioritised and are displayed in a hierarchy, so that one item may
depend on another.


  todo2 [-A]
    Display (all) tasks.

  todo2 [-p <priority>] -a <text>
    Create a new task.

  todo2 -d <index>
    Mark a task as complete.

  todo2 [-p <priority>] -e <task> [<text>]
    Edit an existing task.
`

var addFlag = kingpin.Flag("add", "Add a task.").Short('a').Bool()
var editFlag = kingpin.Flag("edit", "Edit a task, replacing its text.").Short('e').Bool()
var markDoneFlag = kingpin.Flag("done", "Mark the given tasks as done.").Short('d').Bool()
var markNotDoneFlag = kingpin.Flag("not-done", "Mark the given tasks as not done.").Short('D').Bool()
var removeFlag = kingpin.Flag("remove", "Remove the given tasks.").Bool()
var reparentFlag = kingpin.Flag("reparent", "Reparent task A below task B.").Bool()
var titleFlag = kingpin.Flag("title", "Set the task list title.").Bool()
var infoFlag = kingpin.Flag("info", "Show information on a task.").Bool()
var importFlag = kingpin.Flag("import", "Import and synchronise TODO items from source code.").Bool()
var purgeFlag = kingpin.Flag("purge", "Purge completed tasks older than this.").Default("-1s").PlaceHolder("0s").Duration()

var summaryFlag = kingpin.Flag("summary", "Summarise tasks to one line.").Short('s').Bool()
var allFlag = kingpin.Flag("all", "Show all tasks, even completed ones.").Short('A').Bool()
var taskText = kingpin.Arg("arg", "Task text or index.").Strings()

func doView(tasks TaskList, options *ViewOptions) {
	view := NewConsoleView()
	view.ShowTree(tasks, options)
}

func doAdd(tasks TaskList, graft TaskNode, priority Priority, text string, config *config) {
	graft.Create(text, priority)
	saveTaskList(tasks, config)
}

func doEditTask(tasks TaskList, task Task, priority Priority, text string, config *config) {
	if text != "" {
		task.SetText(text)
	}
	if priority != -1 {
		task.SetPriority(priority)
	}
	saveTaskList(tasks, config)
}

func doMarkDone(tasks TaskList, references []Task, config *config) {
	for _, task := range references {
		task.SetCompleted()
	}
	saveTaskList(tasks, config)
}

func doMarkNotDone(tasks TaskList, references []Task, config *config) {
	for _, task := range references {
		task.SetCompletionTime(time.Time{})
	}
	saveTaskList(tasks, config)
}

func doReparent(tasks TaskList, task TaskNode, below TaskNode, config *config) {
	ReparentTask(task, below)
	saveTaskList(tasks, config)
}

func doRemove(tasks TaskList, references []Task, config *config) {
	for _, task := range references {
		task.Delete()
	}
	saveTaskList(tasks, config)
}

func doPurge(tasks TaskList, age time.Duration, config *config) {
	cutoff := time.Now().Add(-age)
	matches := tasks.FindAll(func(task Task) bool {
		return !task.CompletionTime().IsZero() && task.CompletionTime().Before(cutoff)
	})
	for _, m := range matches {
		m.Delete()
	}
	saveTaskList(tasks, config)
}

func doSetTitle(tasks TaskList, args []string, config *config) {
	title := strings.Join(args, " ")
	tasks.SetTitle(title)
	saveTaskList(tasks, config)
}

func doShowInfo(tasks TaskList, index string, options *ViewOptions) {
	task := tasks.Find(index)
	if task == nil {
		fatalf("no such task %s", index)
	}
	view := NewConsoleView()
	view.ShowTaskInfo(task, options)
}

func processAction(tasks TaskList, config *config) {
	priority := PriorityFromString(config.Priority)
	var graft TaskNode = tasks // -golint
	if config.Graft != "root" {
		if graft = tasks.Find(config.Graft); graft == nil {
			fatalf("invalid graft index '%s'", config.Graft)
		}
	}
	order, reversed := OrderFromString(config.Order)
	options := NewViewOptions(*allFlag, *summaryFlag, order, reversed, config.FGColors, config.BGColors)

	switch {
	case *addFlag:
		if len(*taskText) == 0 {
			fatalf("expected text for new task")
		}
		text := strings.Join(*taskText, " ")
		doAdd(tasks, graft, priority, text, config)
	case *markDoneFlag:
		doMarkDone(tasks, resolveTaskReferences(tasks, *taskText), config)
	case *markNotDoneFlag:
		doMarkNotDone(tasks, resolveTaskReferences(tasks, *taskText), config)
	case *removeFlag:
		doRemove(tasks, resolveTaskReferences(tasks, *taskText), config)
	case *reparentFlag:
		if len(*taskText) < 1 {
			fatalf("expected <task> [<new-parent>] for reparenting")
		}
		var below TaskNode
		if len(*taskText) == 2 {
			below = resolveTaskReference(tasks, (*taskText)[1])
		} else {
			below = tasks
		}
		doReparent(tasks, resolveTaskReference(tasks, (*taskText)[0]), below, config)
	case *titleFlag:
		doSetTitle(tasks, *taskText, config)
	case *infoFlag:
		if len(*taskText) < 1 {
			fatalf("expected <task> for info")
		}
		doShowInfo(tasks, (*taskText)[0], options)
	case *importFlag:
		if len(*taskText) < 1 {
			fatalf("expected list of files to import")
		}
		doImport(tasks, *taskText)
	case *editFlag:
		if len(*taskText) < 1 {
			fatalf("expected [-p <priority>] <task> [<text>]")
		}
		task := tasks.Find((*taskText)[0])
		if task == nil {
			fatalf("invalid task %s", (*taskText)[0])
		}
		text := strings.Join(*taskText, " ")
		if config.Priority == "" {
			priority = -1
		}
		doEditTask(tasks, task, priority, text, config)
	case *purgeFlag != -1*time.Second:
		doPurge(tasks, *purgeFlag, config)
	default:
		doView(tasks, options)
	}
}

func resolveTaskReference(tasks TaskList, index string) Task {
	task := tasks.Find(index)
	if task == nil {
		fatalf("invalid task index %s", index)
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
	ranges := strings.Split(indexRange, "-")
	if len(ranges) != 2 {
		return nil
	}
	startIndex := strings.Split(ranges[0], ".")
	start, err := strconv.Atoi(startIndex[len(startIndex)-1])
	if err != nil {
		return nil
	}
	end, err := strconv.Atoi(ranges[1])
	if err != nil {
		return nil
	}
	rangeIndexes := []string{}
	for i := start; i <= end; i++ {
		index := startIndex[:len(startIndex)-1]
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
				fatalf("invalid task range %s", index)
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
		fatalf("no tasks provided to mark done")
	}
	return references
}

func loadTaskList(config *config) (tasks TaskList, err error) {
	// Try loading new-style task file
	if file, err := os.Open(config.File); err == nil {
		defer file.Close()
		loader := NewJSONIO()
		return loader.Deserialize(file)
	}
	// Try loading legacy task file
	if file, err := os.Open(config.LegacyFile); err == nil {
		defer file.Close()
		loader := NewLegacyIO()
		return loader.Deserialize(file)
	}
	return nil, nil
}

func saveTaskList(tasks TaskList, config *config) {
	path := config.File
	previous := path + "~"
	temp := path + "~~"
	var serializeError error
	if file, err := os.Create(temp); err == nil {
		defer func() {
			if err = file.Close(); err != nil {
				os.Remove(temp)
			} else {
				if serializeError != nil {
					return
				}
				if _, err = os.Stat(path); err == nil {
					if err = os.Rename(path, previous); err != nil {
						fatalf("unable to rename %s to %s", path, previous)
					}
				}
				if err = os.Rename(temp, path); err != nil {
					fatalf("unable to rename %s to %s", temp, path)
				}
			}
		}()
		writer := NewJSONIO()
		if serializeError = writer.Serialize(file, tasks); serializeError != nil {
			fatalf(serializeError.Error())
		}
	}
}

func init() {
	kingpin.CommandLine.Help = usage
	kingpin.Version("2.2.0").Author("Alec Thomas <alec@swapoff.org>")
}

func loadConfigCMD(config *config) {
	kingpin.Flag("priority", "Priority of newly created tasks (veryhigh,high,medium,low,verylow).").Short('p').EnumVar(&config.Priority, veryhigh, high, medium, low, verylow)
	kingpin.Flag("graft", "Task to graft new tasks to.").Short('g').Default("root").StringVar(&config.Graft)
	kingpin.Flag("file", "File to load task lists from.").StringVar(&config.File)
	kingpin.Flag("legacy-file", "File to load legacy task lists from.").StringVar(&config.LegacyFile)
	kingpin.Flag("order", "Specify display order of tasks (index,created,completed,text,priority,duration,done).").EnumVar(&config.Order, "index", "created", "completed", "text", "priority", "duration", "done")

	kingpin.Parse()

}

func main() {
	config := NewConfig()
	loadConfigurationFile(config)
	loadConfigCMD(config)

	tasks, err := loadTaskList(config)
	if err != nil {
		fatalf("%s", err)
	}
	if tasks == nil {
		tasks = NewTaskList()
	}
	processAction(tasks, config)
}
