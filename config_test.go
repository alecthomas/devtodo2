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
	"bytes"
	"fmt"
	"os"
	"testing"
)

func TestLoadAllOptionsInCorrectOrderPrimarilyCMDLineOptions(t *testing.T) {
	var fail bytes.Buffer
	config := NewConfig()
	os.Args = []string{"devtodo2", "-A"}

	loadConfigCMD(config)

	if config.Priority != medium {
		fail.WriteString(fmt.Sprintf("\n Err at config_test.go: config.Priority Expected:%s , Got: %s \n", medium, config.Priority))
		t.Fail()
	}
	if config.Graft != "root" {
		fail.WriteString(fmt.Sprintf("\n Err at config_test.go: config.Graft Expected:%s , Got: %s \n", "root", config.Graft))
		t.Fail()
	}
	if config.File != ".todo2" {
		fail.WriteString(fmt.Sprintf("\n Err at config_test.go: config.File Expected:%s , Got: %s \n", ".todo2", config.File))
		t.Fail()
	}

	veryLowFGColorFromConfig := config.FGColors[PriorityFromString(verylow)]
	expectedVeryLowFGColorFromConfig := BLUE
	if veryLowFGColorFromConfig != expectedVeryLowFGColorFromConfig {
		fail.WriteString(fmt.Sprintf("\n Err at config_test.go: config.FGColors.verylow Expected:%s , Got: %s \n", expectedVeryLowFGColorFromConfig, veryLowFGColorFromConfig))
		t.Fail()
	}

	lowFGColorFromConfig := config.FGColors[PriorityFromString(low)]
	expectedLowFGColorFromConfig := CYAN
	if lowFGColorFromConfig != expectedLowFGColorFromConfig {
		fail.WriteString(fmt.Sprintf("\n hmm? Err at config_test.go: config.FGColors.low Expected:%s , Got: %s \n", expectedLowFGColorFromConfig, lowFGColorFromConfig))
		t.Fail()
	}

	mediumFGColorFromConfig := config.FGColors[PriorityFromString(medium)]
	expectedMediumFGColorFromConfig := WHITE
	if mediumFGColorFromConfig != expectedMediumFGColorFromConfig {
		fail.WriteString(fmt.Sprintf("\n Err at config_test.go: config.FGColors.medium Expected:%s , Got: %s \n", expectedMediumFGColorFromConfig, mediumFGColorFromConfig))
		t.Fail()
	}

	highFGColorFromConfig := config.FGColors[PriorityFromString(high)]
	expectedHighFGColorFromConfig := BRIGHTYELLOW
	if highFGColorFromConfig != expectedHighFGColorFromConfig {
		fail.WriteString(fmt.Sprintf("\n Err at config_test.go: config.FGColors.high  Expected:%s , Got: %s \n", expectedHighFGColorFromConfig, highFGColorFromConfig))
		t.Fail()
	}

	veryHighFGColorFromConfig := config.FGColors[PriorityFromString(veryhigh)]
	expectedVeryHighFGColorFromConfig := BRIGHTRED
	if veryHighFGColorFromConfig != expectedVeryHighFGColorFromConfig {
		fail.WriteString(fmt.Sprintf("\n Err at config_test.go: config.FGColors.veryhigh  Expected:%s , Got: %s \n", expectedVeryHighFGColorFromConfig, highFGColorFromConfig))
		t.Fail()
	}

	veryLowBGColorFromConfig := config.BGColors[PriorityFromString(verylow)]
	expectedVeryLowBGColorFromConfig := NOCOLOR
	if veryLowBGColorFromConfig != expectedVeryLowBGColorFromConfig {
		fail.WriteString(fmt.Sprintf("\n Err at config_test.go: config.BGColors.verylow  Expected:%s , Got: %s \n", expectedVeryLowBGColorFromConfig, veryLowBGColorFromConfig))
		t.Fail()
	}

	lowBGColorFromConfig := config.BGColors[PriorityFromString(low)]
	expectedLowBGColorFromConfig := NOCOLOR
	if lowBGColorFromConfig != expectedLowBGColorFromConfig {
		fail.WriteString(fmt.Sprintf("\n Err at config_test.go: config.BGColors.low  Expected:%s , Got: %s \n", expectedLowBGColorFromConfig, lowBGColorFromConfig))
		t.Fail()
	}

	mediumBGColorFromConfig := config.BGColors[PriorityFromString(medium)]
	expectedMediumBGColorFromConfig := NOCOLOR
	if mediumBGColorFromConfig != expectedMediumBGColorFromConfig {
		fail.WriteString(fmt.Sprintf("\n Err at config_test.go: config.BGColors.medium  Expected:%s , Got: %s \n", expectedMediumBGColorFromConfig, mediumBGColorFromConfig))
		t.Fail()
	}

	highBGColorFromConfig := config.BGColors[PriorityFromString(high)]
	expectedHighBGColorFromConfig := NOCOLOR
	if highBGColorFromConfig != expectedHighBGColorFromConfig {
		fail.WriteString(fmt.Sprintf("\n Err at config_test.go: config.BGColors.high  Expected:%s , Got: %s \n", expectedHighBGColorFromConfig, highBGColorFromConfig))
		t.Fail()
	}

	veryHighBGColorFromConfig := config.BGColors[PriorityFromString(veryhigh)]
	expectedVeryHighBGColorFromConfig := NOCOLOR
	if veryHighBGColorFromConfig != expectedVeryHighBGColorFromConfig {
		fail.WriteString(fmt.Sprintf("\n Err at config_test.go: config.BGColors.veryhigh  Expected:%s , Got: %s \n", expectedVeryHighBGColorFromConfig, veryHighBGColorFromConfig))
		t.Fail()
	}

	fmt.Println(fail.String())
}

func TestLoadAllOptionsInCorrectOrderPrimarilyConfigurationFileOptions(t *testing.T) {
	var fail bytes.Buffer
	config := NewConfig()

	priority := veryhigh
	os.Args = []string{"devtodo2", "-p", priority, "-a", "Test task name whatever."}

	loadConfigCMD(config)

	if config.Priority != priority {
		fail.WriteString(fmt.Sprintf("\n Err at config_test.go: config.Priority Expected:%s , Got: %s \n", priority, config.Priority))
		t.Fail()
	}
	if config.Graft != "root" {
		fail.WriteString(fmt.Sprintf("\n Err at config_test.go: config.Graft Expected:%s , Got: %s \n", "root", config.Graft))
		t.Fail()
	}
	if config.File != ".todo2" {
		fail.WriteString(fmt.Sprintf("\n Err at config_test.go: config.File Expected:%s , Got: %s \n", ".todo2", config.File))
		t.Fail()
	}

	veryLowFGColorFromConfig := config.FGColors[PriorityFromString(verylow)]
	expectedVeryLowFGColorFromConfig := BLUE
	if veryLowFGColorFromConfig != expectedVeryLowFGColorFromConfig {
		fail.WriteString(fmt.Sprintf("\n Err at config_test.go: config.FGColors.verylow Expected:%s , Got: %s \n", expectedVeryLowFGColorFromConfig, veryLowFGColorFromConfig))
		t.Fail()
	}

	lowFGColorFromConfig := config.FGColors[PriorityFromString(low)]
	expectedLowFGColorFromConfig := CYAN
	if lowFGColorFromConfig != expectedLowFGColorFromConfig {
		fail.WriteString(fmt.Sprintf("\n Err at config_test.go: config.FGColors.low Expected:%s , Got: %s \n", expectedLowFGColorFromConfig, lowFGColorFromConfig))
		t.Fail()
	}

	mediumFGColorFromConfig := config.FGColors[PriorityFromString(medium)]
	expectedMediumFGColorFromConfig := WHITE
	if mediumFGColorFromConfig != expectedMediumFGColorFromConfig {
		fail.WriteString(fmt.Sprintf("\n Err at config_test.go: config.FGColors.medium Expected:%s , Got: %s \n", expectedMediumFGColorFromConfig, mediumFGColorFromConfig))
		t.Fail()
	}

	highFGColorFromConfig := config.FGColors[PriorityFromString(high)]
	expectedHighFGColorFromConfig := BRIGHTYELLOW
	if highFGColorFromConfig != expectedHighFGColorFromConfig {
		fail.WriteString(fmt.Sprintf("\n Err at config_test.go: config.FGColors.high  Expected:%s , Got: %s \n", expectedHighFGColorFromConfig, highFGColorFromConfig))
		t.Fail()
	}

	veryHighFGColorFromConfig := config.FGColors[PriorityFromString(veryhigh)]
	expectedVeryHighFGColorFromConfig := BRIGHTRED
	if veryHighFGColorFromConfig != expectedVeryHighFGColorFromConfig {
		fail.WriteString(fmt.Sprintf("\n Err at config_test.go: config.FGColors.veryhigh  Expected:%s , Got: %s \n", expectedVeryHighFGColorFromConfig, highFGColorFromConfig))
		t.Fail()
	}

	veryLowBGColorFromConfig := config.BGColors[PriorityFromString(verylow)]
	expectedVeryLowBGColorFromConfig := NOCOLOR
	if veryLowBGColorFromConfig != expectedVeryLowBGColorFromConfig {
		fail.WriteString(fmt.Sprintf("\n Err at config_test.go: config.BGColors.verylow  Expected:%s , Got: %s \n", expectedVeryLowBGColorFromConfig, veryLowBGColorFromConfig))
		t.Fail()
	}

	lowBGColorFromConfig := config.BGColors[PriorityFromString(low)]
	expectedLowBGColorFromConfig := NOCOLOR
	if lowBGColorFromConfig != expectedLowBGColorFromConfig {
		fail.WriteString(fmt.Sprintf("\n Err at config_test.go: config.BGColors.low  Expected:%s , Got: %s \n", expectedLowBGColorFromConfig, lowBGColorFromConfig))
		t.Fail()
	}

	mediumBGColorFromConfig := config.BGColors[PriorityFromString(medium)]
	expectedMediumBGColorFromConfig := NOCOLOR
	if mediumBGColorFromConfig != expectedMediumBGColorFromConfig {
		fail.WriteString(fmt.Sprintf("\n Err at config_test.go: config.BGColors.medium  Expected:%s , Got: %s \n", expectedMediumBGColorFromConfig, mediumBGColorFromConfig))
		t.Fail()
	}

	highBGColorFromConfig := config.BGColors[PriorityFromString(high)]
	expectedHighBGColorFromConfig := NOCOLOR
	if highBGColorFromConfig != expectedHighBGColorFromConfig {
		fail.WriteString(fmt.Sprintf("\n herpderp Err at config_test.go: config.BGColors.high  Expected:%s , Got: %s \n", expectedHighBGColorFromConfig, highBGColorFromConfig))
		t.Fail()
	}

	veryHighBGColorFromConfig := config.BGColors[PriorityFromString(veryhigh)]
	expectedVeryHighBGColorFromConfig := NOCOLOR
	if veryHighBGColorFromConfig != expectedVeryHighBGColorFromConfig {
		fail.WriteString(fmt.Sprintf("\n Err at config_test.go: config.BGColors.veryhigh  Expected:%s , Got: %s \n", expectedVeryHighBGColorFromConfig, veryHighBGColorFromConfig))
		t.Fail()
	}

	fmt.Println(fail.String())
}
