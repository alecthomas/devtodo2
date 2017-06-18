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
	"reflect"
	"testing"
)

const (
	expectedAndActualFormat = "Expected: %s, Got: %s \n"
)

func FileAndTestName(t *testing.T) string {
	var basisErrMessage bytes.Buffer
	privateFields := reflect.ValueOf(*t)
	testName := privateFields.FieldByName("name")
	fmt.Fprintf(&basisErrMessage, "Error at config_test.go:%s. ", testName)
	return basisErrMessage.String()
}

func TestDefaultLegacyFile(t *testing.T) {
	config := NewConfig()
	expected := ".todo"
	actual := config.LegacyFile
	if actual != expected {
		fmt.Printf(FileAndTestName(t)+expectedAndActualFormat, expected, actual)
		t.Fail()
	}
}

func TestOverrideLegacyFile(t *testing.T) {
	config := NewConfig()
	expected := ".testlegacyfile"
	os.Args = []string{"devtodo2", "--legacy-file", expected, "-A"}
	loadCLIConfig(config)
	actual := config.LegacyFile
	if actual != expected {
		fmt.Printf(FileAndTestName(t)+expectedAndActualFormat, expected, actual)
		t.Fail()
	}
}

func TestDefaultOrder(t *testing.T) {
	config := NewConfig()
	expected := "priority"
	actual := config.Order
	if actual != expected {
		fmt.Printf(FileAndTestName(t)+expectedAndActualFormat, expected, actual)
		t.Fail()
	}
}

func TestOverrideOrder(t *testing.T) {
	config := NewConfig()
	expected := orderToString[CREATED]
	os.Args = []string{"devtodo2", "--order", expected, "-A"}
	loadCLIConfig(config)
	actual := config.Order
	if actual != expected {
		fmt.Printf(FileAndTestName(t)+expectedAndActualFormat, expected, actual)
		t.Fail()
	}
}

func TestDefaultPriority(t *testing.T) {
	config := NewConfig()
	expected := medium
	actual := config.Priority
	if actual != expected {
		fmt.Printf(FileAndTestName(t)+expectedAndActualFormat, expected, actual)
		t.Fail()
	}
}

func TestOverridePriority(t *testing.T) {
	config := NewConfig()
	expected := high
	os.Args = []string{"devtodo2", "-p", expected, "-a", "Test task text"}
	loadCLIConfig(config)
	actual := config.Priority
	if actual != expected {
		fmt.Printf(FileAndTestName(t)+expectedAndActualFormat, expected, actual)
		t.Fail()
	}
}

func TestDefaultGraft(t *testing.T) {
	config := NewConfig()
	expected := "root"
	actual := config.Graft
	if actual != expected {
		fmt.Printf(FileAndTestName(t)+expectedAndActualFormat, expected, actual)
		t.Fail()
	}
}

func TestOverrideGraft(t *testing.T) {
	config := NewConfig()
	expected := "1"
	os.Args = []string{"devtodo2", "-g", expected, "-a", "Test task test"}
	loadCLIConfig(config)
	actual := config.Graft
	if actual != expected {
		fmt.Printf(FileAndTestName(t)+expectedAndActualFormat, expected, actual)
		t.Fail()
	}
}

func TestDefaultFile(t *testing.T) {
	config := NewConfig()
	expected := ".todo2"
	actual := config.File
	if actual != expected {
		fmt.Printf(FileAndTestName(t)+expectedAndActualFormat, expected, actual)
		t.Fail()
	}
}

func TestOverrideFile(t *testing.T) {
	config := NewConfig()
	expected := ".testtodo"
	os.Args = []string{"devtodo2", "--file", expected, "-A"}
	loadCLIConfig(config)
	actual := config.File
	if actual != expected {
		fmt.Printf(FileAndTestName(t)+expectedAndActualFormat, expected, actual)
		t.Fail()
	}
}

func TestDefaultVeryLowPriorityFGColor(t *testing.T) {
	config := NewConfig()
	expected := BLUE
	actual := config.FGColors[VERYLOW]
	if actual != expected {
		fmt.Printf(FileAndTestName(t)+expectedAndActualFormat, expected, actual)
		t.Fail()
	}
}

func TestDefaultLowPriorityFGColor(t *testing.T) {
	config := NewConfig()
	expected := CYAN
	actual := config.FGColors[LOW]
	if actual != expected {
		fmt.Printf(FileAndTestName(t)+expectedAndActualFormat, expected, actual)
		t.Fail()
	}
}

func TestDefaultMediumPriorityFGColor(t *testing.T) {
	config := NewConfig()
	expected := WHITE
	actual := config.FGColors[MEDIUM]
	if actual != expected {
		fmt.Printf(FileAndTestName(t)+expectedAndActualFormat, expected, actual)
		t.Fail()
	}
}

func TestDefaultHighPriorityFGColor(t *testing.T) {
	config := NewConfig()
	expected := BRIGHTYELLOW
	actual := config.FGColors[HIGH]
	if actual != expected {
		fmt.Printf(FileAndTestName(t)+expectedAndActualFormat, expected, actual)
		t.Fail()
	}
}

func TestDefaultVeryHighPriorityFGColor(t *testing.T) {
	config := NewConfig()
	expected := BRIGHTRED
	actual := config.FGColors[VERYHIGH]
	if actual != expected {
		fmt.Printf(FileAndTestName(t)+expectedAndActualFormat, expected, actual)
		t.Fail()
	}
}

func TestDefaultVeryHighPriorityBGColor(t *testing.T) {
	config := NewConfig()
	expected := NOCOLOR
	actual := config.BGColors[VERYHIGH]
	if actual != expected {
		fmt.Printf(FileAndTestName(t)+expectedAndActualFormat, expected, actual)
		t.Fail()
	}
}

func TestDefaultHighPriorityBGColor(t *testing.T) {
	config := NewConfig()
	expected := NOCOLOR
	actual := config.BGColors[HIGH]
	if actual != expected {
		fmt.Printf(FileAndTestName(t)+expectedAndActualFormat, expected, actual)
		t.Fail()
	}
}

func TestDefaultMediumPriorityBGColor(t *testing.T) {
	config := NewConfig()
	expected := NOCOLOR
	actual := config.BGColors[MEDIUM]
	if actual != expected {
		fmt.Printf(FileAndTestName(t)+expectedAndActualFormat, expected, actual)
		t.Fail()
	}
}

func TestDefaultLowPriorityBGColor(t *testing.T) {
	config := NewConfig()
	expected := NOCOLOR
	actual := config.BGColors[LOW]
	if actual != expected {
		fmt.Printf(FileAndTestName(t)+expectedAndActualFormat, expected, actual)
		t.Fail()
	}
}

func TestDefaultVeryLowPriorityBGColor(t *testing.T) {
	config := NewConfig()
	expected := NOCOLOR
	actual := config.BGColors[VERYLOW]
	if actual != expected {
		fmt.Printf(FileAndTestName(t)+expectedAndActualFormat, expected, actual)
		t.Fail()
	}
}
