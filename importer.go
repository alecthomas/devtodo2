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
	"bufio"
	"io"
	"os"
	"strings"
)

func importFile(file string) {
	f, e := os.Open(file)
	if e != nil {
		fatal("failed to open %s: %s", file, e.Error())
	}
	defer f.Close()
	reader := bufio.NewReader(f)
	for {
		line, _, e := reader.ReadLine()
		if e == io.EOF {
			break
		}
		if e != nil {
			fatal("error reading %s: %s", file, e.Error())
		}
		text := string(line)
		if strings.Contains(text, "TODO") {
			println(text)
		}
	}
}

func doImport(tasks TaskList, files []string) {
	for _, file := range files {
		importFile(file)
	}
}
