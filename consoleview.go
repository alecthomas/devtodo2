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

const (
	RESET = "\x1b[0m"
	BRIGHT = "\x1b[1m"
	DIM = "\x1b[2m"
	UNDERSCORE = "\x1b[4m"
	BLINK = "\x1b[5m"
	REVERSE = "\x1b[7m"
	HIDDEN = "\x1b[8m"
	FGBLACK = "\x1b[30m"
	FGRED = "\x1b[31m"
	FGGREEN = "\x1b[32m"
	FGYELLOW = "\x1b[33m"
	FGBLUE = "\x1b[34m"
	FGMAGENTA = "\x1b[35m"
	FGCYAN = "\x1b[36m"
	FGWHITE = "\x1b[37m"
	BGBLACK = "\x1b[40m"
	BGRED = "\x1b[41m"
	BGGREEN = "\x1b[42m"
	BGYELLOW = "\x1b[43m"
	BGBLUE = "\x1b[44m"
	BGMAGENTA = "\x1b[45m"
	BGCYAN = "\x1b[46m"
	BGWHITE = "\x1b[47m"
)

const TITLE_COLOUR = BRIGHT + FGGREEN
const NUMBER_COLOR = FGGREEN

// Map for priority level to ANSI colour
var colourPriorityMap map[Priority]string = map[Priority]string {
	VERYHIGH: BRIGHT + FGRED,
	HIGH: BRIGHT + FGYELLOW,
	MEDIUM: FGWHITE,
	LOW: FGCYAN,
	VERYLOW: FGBLUE,
}

func getTerminalWidth() int {
	type winsize struct {
		ws_row, ws_col uint16
		ws_xpixel, ws_ypixel uint16
	}

	ws := winsize{}
	syscall.Syscall(syscall.SYS_IOCTL,
		uintptr(0), uintptr(TIOCGWINSZ),
		uintptr(unsafe.Pointer(&ws)))
	return int(ws.ws_col)
}


func taskState(task Task) int {
	if task.Len() != 0 {
		return '+'
	}
	if task.CompletionTime() != nil {
		return '-'
	}
	return ' '
}

func fatal(format string, args ...interface{}) {
	fmt.Printf("error: %s\n", fmt.Sprintf(format, args...))
	os.Exit(1)
}

func printWrappedText(text string, width, subsequentIndent int) {
	tokens := strings.Split(text, " ", -1)
	offset := 0
	for i, token := range tokens {
		if i > 0 && offset + len(token) > width {
			fmt.Printf("\n%s", strings.Repeat(" ", subsequentIndent))
			offset = 0
		}
		fmt.Printf("%s", token)
		offset += len(token)
		if offset < width && i != len(tokens) - 1 {
			fmt.Print(" ")
			offset += 1
		}
	}
}

func formatTask(width, depth, index int, task Task) {
	indent := depth * 4 + 4
	width -= indent
	state := taskState(task)
	fmt.Printf("%s%s%c%2d.%s%s", strings.Repeat("    ", depth), NUMBER_COLOR, state,
			   index + 1, RESET, colourPriorityMap[task.Priority()])
	printWrappedText(task.Text(), width, indent)
	fmt.Printf("%s\n", RESET)
}

func consoleDisplayTask(width, depth, index int, task Task, showAll bool) {
	if depth >= 0 && (!showAll && task.CompletionTime() != nil) {
		return
	}
	if depth >= 0 {
		formatTask(width, depth, index, task)
	}
	for i := 0; i < task.Len(); i++ {
		consoleDisplayTask(width, depth + 1, i, task.At(i), showAll)
	}
}

func ConsoleView(tasks TaskList, showAll bool) {
	width := getTerminalWidth()
	if tasks.Title() != "" {
		fmt.Print(TITLE_COLOUR)
		printWrappedText("    " + tasks.Title(), width, 4)
		fmt.Printf("%s\n", RESET)
	}
	for i := 0; i < tasks.Len(); i++ {
		consoleDisplayTask(width, 0, i, tasks.At(i), showAll)
	}
}
