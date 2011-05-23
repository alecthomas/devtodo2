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
	ACTION_VIEW = Action(0)
	ACTION_MARK_DONE = Action(1)
	ACTION_MARK_NOT_DONE = Action(2)
)

func parseCommandLine() Action {
	optarg.Header("Actions")
	optarg.Add("d", "done", "mark the given tasks as done", false)
	optarg.Add("D", "not-done", "mark the given tasks as not done", false)

	optarg.Header("Options")
	optarg.Add("p", "priority", "priority of newly created tasks", "medium")

	action := ACTION_VIEW

	for opt := range optarg.Parse() {
		switch opt.ShortName {
		// Actions
		case "d":
			action = ACTION_MARK_DONE
		case "D":
			action = ACTION_MARK_NOT_DONE

		// Options
		case "p":
		}
	}	
	return action
}

func main() {
	action := parseCommandLine()

	todoFile, err := os.Open(".todo")
	if err != nil {
		log.Fatal(err)
	}
	taskList := LoadLegacyTaskList(todoFile)

	switch action {
	case ACTION_VIEW:
		ConsoleView(taskList)
	case ACTION_MARK_DONE:
	case ACTION_MARK_NOT_DONE:
	}
}
