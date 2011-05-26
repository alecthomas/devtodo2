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
	"os"
	optarg "github.com/jteeuwen/go-pkg-optarg"
	"log"
)

type Action int

const (
	ACTION_VIEW = Action(iota)
	ACTION_MARK_DONE
	ACTION_MARK_NOT_DONE
	ACTION_ADD_TASK
)

func doView(tasks TaskList) {
	ConsoleView(tasks)
}

func processAction(tasks TaskList) {
	optarg.Header("Actions")
	optarg.Add("a", "add", "add a task", false)
	optarg.Add("d", "done", "mark the given tasks as done", false)
	optarg.Add("D", "not-done", "mark the given tasks as not done", false)

	optarg.Header("Task creation options")
	optarg.Add("p", "priority", "priority of newly created tasks", "medium")
	optarg.Add("g", "graft", "index to graft new task to", "")

	action := ACTION_VIEW
	priority := MEDIUM
	var graft Task = nil

	// First pass, collect options.
	for opt := range optarg.Parse() {
		switch opt.ShortName {
		// Actions
		case "d":
			action = ACTION_MARK_DONE
		case "D":
			action = ACTION_MARK_NOT_DONE
		case "a":
			action = ACTION_ADD_TASK

		// Options
		case "p":
			priority = PriorityFromString(opt.String())
		case "g":
			if graft = tasks.Find(opt.String()); graft == nil {
				printFatal("invalid graft index '%s'", opt.String())
			}
		}
	}

	switch action {
		case ACTION_VIEW:
			doView(tasks)
		case ACTION_MARK_DONE:
		case ACTION_MARK_NOT_DONE:
		case ACTION_ADD_TASK:
	}

	println(priority)
	println(graft)
	println(optarg.Remainder)
}

func main() {
	todoFile, err := os.Open(".todo")
	if err != nil {
		log.Fatal(err)
	}
	tasks := LoadLegacyTaskList(todoFile)
	processAction(tasks)
}
