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
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
	"testing"
)

func TestLoadAllOptionsInCorrectOrderEspeciallyNonColorOptions(t *testing.T) {
	var fail bytes.Buffer
	config := GetConfigInstance()
	fgColors := map[string]string{
		verylow:  BLUE,
		low:      CYAN,
		medium:   WHITE,
		high:     BRIGHTYELLOW,
		veryhigh: BRIGHTRED,
	}
	bgColors := map[string]string{
		verylow:  BLACK,
		low:      YELLOW,
		medium:   WHITE,
		high:     BLUE,
		veryhigh: CYAN,
	}
	defaultPriority := "medium"
	graft := "root"
	file := ".todo2"
	legacyFile := ".todo"
	order := "priority"
	marshalableConfig := &MarshalableConfig{
		FGColors:   fgColors,
		BGColors:   bgColors,
		Priority:   defaultPriority,
		Graft:      graft,
		File:       file,
		LegacyFile: legacyFile,
		Order:      order,
	}
	priority := "veryhigh"
	os.Args = []string{"devtodo2", "-A", "--lowfgcolor", BLACK}
	copyToConfigFromMarshalableConfig(marshalableConfig, config)

	copyToConfigFromCMDOptions(config)
	kingpin.Parse()

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
	expectedVeryLowFGColorFromConfig := BLACK
	if veryLowFGColorFromConfig != expectedVeryLowFGColorFromConfig {
		fail.WriteString(fmt.Sprintf("\n Err at config_test.go: config.FGColors.verylow Expected:%s , Got: %s \n", veryLowFGColorFromConfig, expectedVeryLowFGColorFromConfig))
		t.Fail()
	}

	lowFGColorFromConfig := config.FGColors[PriorityFromString(low)]
	expectedLowFGColorFromConfig := CYAN
	if lowFGColorFromConfig != expectedLowFGColorFromConfig {
		fail.WriteString(fmt.Sprintf("\n Err at config_test.go: config.FGColors.low Expected:%s , Got: %s \n", lowFGColorFromConfig, expectedLowFGColorFromConfig))
		t.Fail()
	}

	mediumFGColorFromConfig := config.FGColors[PriorityFromString(medium)]
	expectedMediumFGColorFromConfig := WHITE
	if mediumFGColorFromConfig != expectedMediumFGColorFromConfig {
		fail.WriteString(fmt.Sprintf("\n Err at config_test.go: config.FGColors.medium Expected:%s , Got: %s \n", mediumFGColorFromConfig, expectedMediumFGColorFromConfig))
		t.Fail()
	}

	highFGColorFromConfig := config.FGColors[PriorityFromString(high)]
	expectedHighFGColorFromConfig := BRIGHTYELLOW
	if highFGColorFromConfig != expectedHighFGColorFromConfig {
		fail.WriteString(fmt.Sprintf("\n Err at config_test.go: config.FGColors.high  Expected:%s , Got: %s \n", highFGColorFromConfig, expectedHighFGColorFromConfig))
		t.Fail()
	}

	veryHighFGColorFromConfig := config.FGColors[PriorityFromString(veryhigh)]
	expectedVeryHighFGColorFromConfig := BRIGHTRED
	if veryHighFGColorFromConfig != expectedVeryHighFGColorFromConfig {
		fail.WriteString(fmt.Sprintf("\n Err at config_test.go: config.FGColors.veryhigh  Expected:%s , Got: %s \n", veryHighFGColorFromConfig, expectedHighFGColorFromConfig))
		t.Fail()
	}

	veryLowBGColorFromConfig := config.BGColors[PriorityFromString(verylow)]
	expectedVeryLowBGColorFromConfig := BLACK
	if veryLowBGColorFromConfig != expectedVeryLowBGColorFromConfig {
		fail.WriteString(fmt.Sprintf("\n Err at config_test.go: config.BGColors.verylow  Expected:%s , Got: %s \n", veryLowBGColorFromConfig, expectedVeryLowBGColorFromConfig))
		t.Fail()
	}

	lowBGColorFromConfig := config.BGColors[PriorityFromString(low)]
	expectedLowBGColorFromConfig := YELLOW
	if lowBGColorFromConfig != expectedLowBGColorFromConfig {
		fail.WriteString(fmt.Sprintf("\n Err at config_test.go: config.BGColors.low  Expected:%s , Got: %s \n", lowBGColorFromConfig, expectedLowBGColorFromConfig))
		t.Fail()
	}

	mediumBGColorFromConfig := config.BGColors[PriorityFromString(medium)]
	expectedMediumBGColorFromConfig := WHITE
	if mediumBGColorFromConfig != expectedMediumBGColorFromConfig {
		fail.WriteString(fmt.Sprintf("\n Err at config_test.go: config.BGColors.medium  Expected:%s , Got: %s \n", mediumBGColorFromConfig, expectedMediumBGColorFromConfig))
		t.Fail()
	}

	highBGColorFromConfig := config.BGColors[PriorityFromString(high)]
	expectedHighBGColorFromConfig := BLUE
	if highBGColorFromConfig != expectedHighBGColorFromConfig {
		fail.WriteString(fmt.Sprintf("\n Err at config_test.go: config.BGColors.high  Expected:%s , Got: %s \n", highBGColorFromConfig, expectedHighBGColorFromConfig))
		t.Fail()
	}

	veryHighBGColorFromConfig := config.BGColors[PriorityFromString(veryhigh)]
	expectedVeryHighBGColorFromConfig := CYAN
	if veryHighBGColorFromConfig != expectedVeryHighBGColorFromConfig {
		fail.WriteString(fmt.Sprintf("\n Err at config_test.go: config.BGColors.veryhigh  Expected:%s , Got: %s \n", veryHighBGColorFromConfig, expectedVeryHighBGColorFromConfig))
		t.Fail()
	}

	fmt.Println(fail.String())
}

func TestLoadAllOptionsInCorrectOrderEspeciallyColors(t *testing.T) {
	var fail bytes.Buffer
	config := GetConfigInstance()
	fgColors := map[string]string{
		verylow:  BLUE,
		low:      CYAN,
		medium:   WHITE,
		high:     BRIGHTYELLOW,
		veryhigh: BRIGHTRED,
	}
	bgColors := map[string]string{
		verylow:  BLACK,
		low:      YELLOW,
		medium:   WHITE,
		high:     BLUE,
		veryhigh: CYAN,
	}
	defaultPriority := "medium"
	graft := "root"
	file := ".todo2"
	legacyFile := ".todo"
	order := "priority"
	marshalableConfig := &MarshalableConfig{
		FGColors:   fgColors,
		BGColors:   bgColors,
		Priority:   defaultPriority,
		Graft:      graft,
		File:       file,
		LegacyFile: legacyFile,
		Order:      order,
	}
	priority := veryhigh
	os.Args = []string{"devtodo2", "-p", priority, "-a", "task name whatever"}
	copyToConfigFromMarshalableConfig(marshalableConfig, config)

	copyToConfigFromCMDOptions(config)
	kingpin.Parse()

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
		fail.WriteString(fmt.Sprintf("\n Err at config_test.go: config.FGColors.verylow Expected:%s , Got: %s \n", veryLowFGColorFromConfig, expectedVeryLowFGColorFromConfig))
		t.Fail()
	}

	lowFGColorFromConfig := config.FGColors[PriorityFromString(low)]
	expectedLowFGColorFromConfig := CYAN
	if lowFGColorFromConfig != expectedLowFGColorFromConfig {
		fail.WriteString(fmt.Sprintf("\n Err at config_test.go: config.FGColors.low Expected:%s , Got: %s \n", lowFGColorFromConfig, expectedLowFGColorFromConfig))
		t.Fail()
	}

	mediumFGColorFromConfig := config.FGColors[PriorityFromString(medium)]
	expectedMediumFGColorFromConfig := WHITE
	if mediumFGColorFromConfig != expectedMediumFGColorFromConfig {
		fail.WriteString(fmt.Sprintf("\n Err at config_test.go: config.FGColors.medium Expected:%s , Got: %s \n", mediumFGColorFromConfig, expectedMediumFGColorFromConfig))
		t.Fail()
	}

	highFGColorFromConfig := config.FGColors[PriorityFromString(high)]
	expectedHighFGColorFromConfig := BRIGHTYELLOW
	if highFGColorFromConfig != expectedHighFGColorFromConfig {
		fail.WriteString(fmt.Sprintf("\n Err at config_test.go: config.FGColors.high  Expected:%s , Got: %s \n", highFGColorFromConfig, expectedHighFGColorFromConfig))
		t.Fail()
	}

	veryHighFGColorFromConfig := config.FGColors[PriorityFromString(veryhigh)]
	expectedVeryHighFGColorFromConfig := BRIGHTRED
	if veryHighFGColorFromConfig != expectedVeryHighFGColorFromConfig {
		fail.WriteString(fmt.Sprintf("\n Err at config_test.go: config.FGColors.veryhigh  Expected:%s , Got: %s \n", veryHighFGColorFromConfig, expectedHighFGColorFromConfig))
		t.Fail()
	}

	veryLowBGColorFromConfig := config.BGColors[PriorityFromString(verylow)]
	expectedVeryLowBGColorFromConfig := BLACK
	if veryLowBGColorFromConfig != expectedVeryLowBGColorFromConfig {
		fail.WriteString(fmt.Sprintf("\n Err at config_test.go: config.BGColors.verylow  Expected:%s , Got: %s \n", veryLowBGColorFromConfig, expectedVeryLowBGColorFromConfig))
		t.Fail()
	}

	lowBGColorFromConfig := config.BGColors[PriorityFromString(low)]
	expectedLowBGColorFromConfig := YELLOW
	if lowBGColorFromConfig != expectedLowBGColorFromConfig {
		fail.WriteString(fmt.Sprintf("\n Err at config_test.go: config.BGColors.low  Expected:%s , Got: %s \n", lowBGColorFromConfig, expectedLowBGColorFromConfig))
		t.Fail()
	}

	mediumBGColorFromConfig := config.BGColors[PriorityFromString(medium)]
	expectedMediumBGColorFromConfig := WHITE
	if mediumBGColorFromConfig != expectedMediumBGColorFromConfig {
		fail.WriteString(fmt.Sprintf("\n Err at config_test.go: config.BGColors.medium  Expected:%s , Got: %s \n", mediumBGColorFromConfig, expectedMediumBGColorFromConfig))
		t.Fail()
	}

	highBGColorFromConfig := config.BGColors[PriorityFromString(high)]
	expectedHighBGColorFromConfig := BLUE
	if highBGColorFromConfig != expectedHighBGColorFromConfig {
		fail.WriteString(fmt.Sprintf("\n Err at config_test.go: config.BGColors.high  Expected:%s , Got: %s \n", highBGColorFromConfig, expectedHighBGColorFromConfig))
		t.Fail()
	}

	veryHighBGColorFromConfig := config.BGColors[PriorityFromString(veryhigh)]
	expectedVeryHighBGColorFromConfig := CYAN
	if veryHighBGColorFromConfig != expectedVeryHighBGColorFromConfig {
		fail.WriteString(fmt.Sprintf("\n Err at config_test.go: config.BGColors.veryhigh  Expected:%s , Got: %s \n", veryHighBGColorFromConfig, expectedVeryHighBGColorFromConfig))
		t.Fail()
	}

	fmt.Println(fail.String())
}
