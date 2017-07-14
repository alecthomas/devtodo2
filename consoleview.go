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
	"strings"
	"syscall"
	"unsafe"
)

// ANSI color constants.
const (
	RESET      = "\x1b[0m"
	BRIGHT     = "\x1b[1m"
	DIM        = "\x1b[2m"
	UNDERSCORE = "\x1b[4m"
	BLINK      = "\x1b[5m"
	REVERSE    = "\x1b[7m"
	HIDDEN     = "\x1b[8m"
	FGBLACK    = "\x1b[30m"
	FGRED      = "\x1b[31m"
	FGGREEN    = "\x1b[32m"
	FGYELLOW   = "\x1b[33m"
	FGBLUE     = "\x1b[34m"
	FGMAGENTA  = "\x1b[35m"
	FGCYAN     = "\x1b[36m"
	FGWHITE    = "\x1b[37m"
	BGBLACK    = "\x1b[40m"
	BGRED      = "\x1b[41m"
	BGGREEN    = "\x1b[42m"
	BGYELLOW   = "\x1b[43m"
	BGBLUE     = "\x1b[44m"
	BGMAGENTA  = "\x1b[45m"
	BGCYAN     = "\x1b[46m"
	BGWHITE    = "\x1b[47m"

	TITLE_COLOUR = BRIGHT + FGGREEN
	NUMBER_COLOR = FGGREEN
)

// Map for priority level to ANSI colour
var colourPriorityMap = map[Priority]string{
	VERYHIGH: BRIGHT + FGRED,
	HIGH:     BRIGHT + FGYELLOW,
	MEDIUM:   FGWHITE,
	LOW:      FGCYAN,
	VERYLOW:  FGBLUE,
}

func getTerminalWidth() int {
	type winsize struct {
		wsRow, wsCol       uint16
		wsXPixel, wsYPixel uint16
	}

	ws := winsize{}
	syscall.Syscall(syscall.SYS_IOCTL,
		uintptr(0), syscall.TIOCGWINSZ,
		uintptr(unsafe.Pointer(&ws)))
	return int(ws.wsCol)
}

func taskState(task Task) int {
	if task.Len() != 0 {
		return '+'
	}
	if !task.CompletionTime().IsZero() {
		return '-'
	}
	return ' '
}

func fatalf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "error: %s\n", fmt.Sprintf(format, args...))
	os.Exit(1)
}

func printWrappedText(text string, width, subsequentIndent int) {
	tokens := strings.Split(text, " ")
	offset := 0
	for i, token := range tokens {
		if i > 0 && offset+len(token) > width {
			fmt.Printf("\n%s", strings.Repeat(" ", subsequentIndent))
			offset = 0
		}
		fmt.Printf("%s", token)
		offset += len(token)
		if offset < width && i != len(tokens)-1 {
			fmt.Print(" ")
			offset++
		}
	}
}

func formatTask(width, depth int, task Task, options *ViewOptions) {
	indent := depth*4 + 4
	width -= indent
	state := taskState(task)
	fmt.Printf("%s%s%c%2d.%s%s", strings.Repeat("    ", depth), NUMBER_COLOR, state,
		task.ID()+1, RESET, colourPriorityMap[task.Priority()])
	text := task.Text()
	trimmed := false
	if options.Summarise {
		if len(text) > width {
			text = strings.TrimSpace(text[:width-1])
			trimmed = true
		}
	}
	printWrappedText(text, width, indent)
	if trimmed {
		fmt.Printf("%s+%s\n", TITLE_COLOUR, RESET)
	} else {
		fmt.Printf("%s\n", RESET)
	}
}

func consoleDisplayTask(width, depth int, task Task, options *ViewOptions) {
	if depth >= 0 && (!options.ShowAll && !task.CompletionTime().IsZero()) {
		return
	}
	if depth >= 0 {
		formatTask(width, depth, task, options)
	}
	if !options.Summarise {
		view := CreateTaskView(task, options)
		for i := 0; i < view.Len(); i++ {
			consoleDisplayTask(width, depth+1, view.At(i), options)
		}
	}
}

type ConsoleView struct{}

func NewConsoleView() *ConsoleView {
	return &ConsoleView{}
}

func (c *ConsoleView) ShowTree(tasks TaskList, options *ViewOptions) {
	width := getTerminalWidth()
	if tasks.Title() != "" {
		fmt.Print(TITLE_COLOUR)
		printWrappedText("    "+tasks.Title(), width, 4)
		fmt.Printf("%s\n", RESET)
	}
	view := CreateTaskView(tasks, options)
	for i := 0; i < view.Len(); i++ {
		consoleDisplayTask(width, 0, view.At(i), options)
	}
}

func (c *ConsoleView) ShowTaskInfo(task Task) {
	width := getTerminalWidth()
	fmt.Print(colourPriorityMap[task.Priority()])
	printWrappedText(task.Text(), width, 0)
	fmt.Printf("%s\n\n", RESET)
	fmt.Printf("%sPriority%s %s%s%s\n", BRIGHT, RESET, colourPriorityMap[task.Priority()], task.Priority().String(), RESET)
	fmt.Printf("%sCreated:%s %s\n", BRIGHT, RESET, task.CreationTime().Local().String())
	completed := "incomplete"
	if !task.CompletionTime().IsZero() {
		completed = task.CompletionTime().Local().String()
	}
	fmt.Printf("%sCompleted:%s %s\n", BRIGHT, RESET, completed)
}
