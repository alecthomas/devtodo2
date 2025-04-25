/*
  Author 2025 Konstantin Volokh

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
    "time"
)

import "github.com/dispatchrun/coroutine"

type Exporter struct{
    databaseFile string;
    coro coroutine.Coroutine[Task, any];
}

func NewExporter(databaseFile string) *Exporter {
	return &Exporter{
		databaseFile: databaseFile,
	}
}

func writeLine(file *os.File, line string) error {
    _, err := file.WriteString(line + "\n")
    if err != nil {
        return fmt.Errorf("failed to write to file: %w", err)
    }

    return nil;
}

func (c *Exporter) writeStringsToFile(filename string) error {
    // Open file for writing (create if not exists, truncate if exists)
    file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
    if err != nil {
        return fmt.Errorf("failed to open file: %w", err)
    }
    defer file.Close()

    for c.coro.Next() {
        task := c.coro.Recv()
        if !task.CompletionTime().IsZero() {
		continue;
	}
        line := fmt.Sprintf("- %s\n  (added %s, priority %s)\n", task.Text(), task.CreationTime().Format(time.ANSIC), task.Priority())
	_ = writeLine(file, line)
    }

    return nil
}

func SaveTask(task Task, level int) {
	for i := 0; i < task.Len(); i++ {
		task_ := task.At(i)
		coroutine.Yield[Task, any](task_)
		for y := 0; y < task_.Len(); y++ {
			SaveTask(task_.At(y), level + 1)
		}
	}
}

func (c *Exporter) SaveTodo(tasks TaskList, options *ViewOptions) {
	view := CreateTaskView(tasks, options)
	c.coro = coroutine.New[Task, any](func() {
	//	SaveTask(view)
		for i := 0; i < view.Len(); i++ {
			task_ := view.At(i)
			coroutine.Yield[Task, any](task_)
			SaveTask(task_, 1)
		}
	})

        c.writeStringsToFile("TODO")
}
