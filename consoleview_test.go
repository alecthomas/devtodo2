package main

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func TestFormatTaskForCorrectVeryLowFGColor(t *testing.T) {
	var fail bytes.Buffer
	id := 1
	text := "Test text of task"
	priority := PriorityFromString(verylow)

	width := 5
	depth := 0
	task := newTask(id, text, priority)
	order, reversed := OrderFromString("priority")
	fgColors := map[Priority]string{
		VERYLOW:  BLUE,
		LOW:      CYAN,
		MEDIUM:   GREEN,
		HIGH:     BRIGHTYELLOW,
		VERYHIGH: BRIGHTRED,
	}
	bgColors := map[Priority]string{
		VERYLOW:  NOCOLOR,
		LOW:      NOCOLOR,
		MEDIUM:   NOCOLOR,
		HIGH:     NOCOLOR,
		VERYHIGH: NOCOLOR,
	}

	options := NewViewOptions(false, false, order, reversed, fgColors, bgColors)
	formattedTask := formatTask(width, depth, task, options)
	expected := FGBLUE
	if !strings.Contains(formattedTask, expected) {
		fail.WriteString(fmt.Sprintf("formattedTask string did not contain the color we expected. Expected to contain: %s . Got formatted task: %s", expected, formattedTask))
		t.Fail()
	}
	fmt.Println(fail.String())
}

func TestFormatTaskForCorrectLowFGColor(t *testing.T) {
	var fail bytes.Buffer
	id := 1
	text := "Test text of task"
	priority := PriorityFromString(low)

	width := 5
	depth := 0
	task := newTask(id, text, priority)
	order, reversed := OrderFromString("priority")
	fgColors := map[Priority]string{
		VERYLOW:  BLUE,
		LOW:      CYAN,
		MEDIUM:   GREEN,
		HIGH:     BRIGHTYELLOW,
		VERYHIGH: BRIGHTRED,
	}
	bgColors := map[Priority]string{
		VERYLOW:  BLACK,
		LOW:      BLACK,
		MEDIUM:   BLACK,
		HIGH:     BLACK,
		VERYHIGH: BLACK,
	}

	options := NewViewOptions(false, false, order, reversed, fgColors, bgColors)
	formattedTask := formatTask(width, depth, task, options)
	expected := FGCYAN
	if !strings.Contains(formattedTask, expected) {
		fail.WriteString(fmt.Sprintf("formattedTask string did not contain the color we expected. Expected to contain: %s . Got formatted task: %s", expected, formattedTask))
		t.Fail()
	}
	fmt.Println(fail.String())
}

func TestFormatTaskForCorrectMediumFGColor(t *testing.T) {
	var fail bytes.Buffer
	id := 1
	text := "Test text of task"
	priority := PriorityFromString(medium)

	width := 5
	depth := 0
	task := newTask(id, text, priority)
	order, reversed := OrderFromString("priority")
	fgColors := map[Priority]string{
		VERYLOW:  BLUE,
		LOW:      CYAN,
		MEDIUM:   GREEN,
		HIGH:     BRIGHTYELLOW,
		VERYHIGH: BRIGHTRED,
	}
	bgColors := map[Priority]string{
		VERYLOW:  BLACK,
		LOW:      BLACK,
		MEDIUM:   BLACK,
		HIGH:     BLACK,
		VERYHIGH: BLACK,
	}

	options := NewViewOptions(false, false, order, reversed, fgColors, bgColors)
	formattedTask := formatTask(width, depth, task, options)
	expected := FGGREEN
	if !strings.Contains(formattedTask, expected) {
		fail.WriteString(fmt.Sprintf("formattedTask string did not contain the color we expected. Expected to contain: %s . Got formatted task: %s", expected, formattedTask))
		t.Fail()
	}
	fmt.Println(fail.String())
}

func TestFormatTaskForCorrectHighFGColor(t *testing.T) {
	var fail bytes.Buffer
	id := 1
	text := "Test text of task"
	priority := PriorityFromString(high)

	width := 5
	depth := 0
	task := newTask(id, text, priority)
	order, reversed := OrderFromString("priority")
	fgColors := map[Priority]string{
		VERYLOW:  BLUE,
		LOW:      CYAN,
		MEDIUM:   GREEN,
		HIGH:     BRIGHTYELLOW,
		VERYHIGH: BRIGHTRED,
	}
	bgColors := map[Priority]string{
		VERYLOW:  BLACK,
		LOW:      BLACK,
		MEDIUM:   BLACK,
		HIGH:     BLACK,
		VERYHIGH: BLACK,
	}

	options := NewViewOptions(false, false, order, reversed, fgColors, bgColors)
	formattedTask := formatTask(width, depth, task, options)
	expected := BRIGHT + FGYELLOW
	if !strings.Contains(formattedTask, expected) {
		fail.WriteString(fmt.Sprintf("formattedTask string did not contain the color we expected. Expected to contain: %s . Got formatted task: %s", expected, formattedTask))
		t.Fail()
	}
	fmt.Println(fail.String())
}

func TestFormatTaskForCorrectVeryHighFGColor(t *testing.T) {
	var fail bytes.Buffer
	id := 1
	text := "Test text of task"
	priority := PriorityFromString(high)

	width := 5
	depth := 0
	task := newTask(id, text, priority)
	order, reversed := OrderFromString("priority")
	fgColors := map[Priority]string{
		VERYLOW:  BLUE,
		LOW:      CYAN,
		MEDIUM:   GREEN,
		HIGH:     BRIGHTYELLOW,
		VERYHIGH: BRIGHTRED,
	}
	bgColors := map[Priority]string{
		VERYLOW:  BLACK,
		LOW:      BLACK,
		MEDIUM:   BLACK,
		HIGH:     BLACK,
		VERYHIGH: BLACK,
	}

	options := NewViewOptions(false, false, order, reversed, fgColors, bgColors)
	formattedTask := formatTask(width, depth, task, options)
	expected := BRIGHT + FGYELLOW
	if !strings.Contains(formattedTask, expected) {
		fail.WriteString(fmt.Sprintf("formattedTask string did not contain the color we expected. Expected to contain: %s . Got formatted task: %s", expected, formattedTask))
		t.Fail()
	}
	fmt.Println(fail.String())
}

func TestFormatTaskForCorrectVeryLowBGColor(t *testing.T) {
	var fail bytes.Buffer
	id := 1
	text := "Test text of task"
	priority := PriorityFromString(verylow)

	width := 5
	depth := 0
	task := newTask(id, text, priority)
	order, reversed := OrderFromString("priority")
	fgColors := map[Priority]string{
		VERYLOW:  BLUE,
		LOW:      CYAN,
		MEDIUM:   GREEN,
		HIGH:     BRIGHTYELLOW,
		VERYHIGH: BRIGHTRED,
	}
	bgColors := map[Priority]string{
		VERYLOW:  BLUE,
		LOW:      YELLOW,
		MEDIUM:   WHITE,
		HIGH:     GREEN,
		VERYHIGH: RED,
	}

	options := NewViewOptions(false, false, order, reversed, fgColors, bgColors)
	formattedTask := formatTask(width, depth, task, options)
	expected := BGBLUE
	if !strings.Contains(formattedTask, expected) {
		fail.WriteString(fmt.Sprintf("formattedTask string did not contain the bg color we expected. Expected to contain: %s . Got formatted task: %s", expected, formattedTask))
		t.Fail()
	}
	fmt.Println(fail.String())
}

func TestFormatTaskForCorrectLowBGColor(t *testing.T) {
	var fail bytes.Buffer
	id := 1
	text := "Test text of task"
	priority := PriorityFromString(low)

	width := 5
	depth := 0
	task := newTask(id, text, priority)
	order, reversed := OrderFromString("priority")
	fgColors := map[Priority]string{
		VERYLOW:  BLUE,
		LOW:      CYAN,
		MEDIUM:   GREEN,
		HIGH:     BRIGHTYELLOW,
		VERYHIGH: BRIGHTRED,
	}
	bgColors := map[Priority]string{
		VERYLOW:  BLUE,
		LOW:      YELLOW,
		MEDIUM:   WHITE,
		HIGH:     GREEN,
		VERYHIGH: RED,
	}

	options := NewViewOptions(false, false, order, reversed, fgColors, bgColors)
	formattedTask := formatTask(width, depth, task, options)
	expected := BGYELLOW
	if !strings.Contains(formattedTask, expected) {
		fail.WriteString(fmt.Sprintf("formattedTask string did not contain the bg color we expected. Expected to contain: %s . Got formatted task: %s", expected, formattedTask))
		t.Fail()
	}
	fmt.Println(fail.String())
}

func TestFormatTaskForCorrectMediumBGColor(t *testing.T) {
	var fail bytes.Buffer
	id := 1
	text := "Test text of task"
	priority := PriorityFromString(medium)

	width := 5
	depth := 0
	task := newTask(id, text, priority)
	order, reversed := OrderFromString("priority")
	fgColors := map[Priority]string{
		VERYLOW:  BLUE,
		LOW:      CYAN,
		MEDIUM:   GREEN,
		HIGH:     BRIGHTYELLOW,
		VERYHIGH: BRIGHTRED,
	}
	bgColors := map[Priority]string{
		VERYLOW:  BLUE,
		LOW:      YELLOW,
		MEDIUM:   WHITE,
		HIGH:     GREEN,
		VERYHIGH: RED,
	}

	options := NewViewOptions(false, false, order, reversed, fgColors, bgColors)
	formattedTask := formatTask(width, depth, task, options)
	expected := BGWHITE
	if !strings.Contains(formattedTask, expected) {
		fail.WriteString(fmt.Sprintf("formattedTask string did not contain the bg color we expected. Expected to contain: %s . Got formatted task: %s", expected, formattedTask))
		t.Fail()
	}
	fmt.Println(fail.String())
}

func TestFormatTaskForCorrectHighBGColor(t *testing.T) {
	var fail bytes.Buffer
	id := 1
	text := "Test text of task"
	priority := PriorityFromString(high)

	width := 5
	depth := 0
	task := newTask(id, text, priority)
	order, reversed := OrderFromString("priority")
	fgColors := map[Priority]string{
		VERYLOW:  BLUE,
		LOW:      CYAN,
		MEDIUM:   GREEN,
		HIGH:     BRIGHTYELLOW,
		VERYHIGH: BRIGHTRED,
	}
	bgColors := map[Priority]string{
		VERYLOW:  BLUE,
		LOW:      YELLOW,
		MEDIUM:   WHITE,
		HIGH:     GREEN,
		VERYHIGH: RED,
	}

	options := NewViewOptions(false, false, order, reversed, fgColors, bgColors)
	formattedTask := formatTask(width, depth, task, options)
	expected := BGGREEN
	if !strings.Contains(formattedTask, expected) {
		fail.WriteString(fmt.Sprintf("formattedTask string did not contain the bg color we expected. Expected to contain: %s . Got formatted task: %s", expected, formattedTask))
		t.Fail()
	}
	fmt.Println(fail.String())
}

func TestFormatTaskForCorrectVeryHighBGColor(t *testing.T) {
	var fail bytes.Buffer
	id := 1
	text := "Test text of task"
	priority := PriorityFromString(veryhigh)

	width := 5
	depth := 0
	task := newTask(id, text, priority)
	order, reversed := OrderFromString("priority")
	fgColors := map[Priority]string{
		VERYLOW:  BLUE,
		LOW:      CYAN,
		MEDIUM:   GREEN,
		HIGH:     BRIGHTYELLOW,
		VERYHIGH: BRIGHTRED,
	}
	bgColors := map[Priority]string{
		VERYLOW:  BLUE,
		LOW:      YELLOW,
		MEDIUM:   WHITE,
		HIGH:     GREEN,
		VERYHIGH: RED,
	}

	options := NewViewOptions(false, false, order, reversed, fgColors, bgColors)
	formattedTask := formatTask(width, depth, task, options)
	expected := BGRED
	if !strings.Contains(formattedTask, expected) {
		fail.WriteString(fmt.Sprintf("formattedTask string did not contain the bg color we expected. Expected to contain: %s . Got formatted task: %s", expected, formattedTask))
		t.Fail()
	}
	fmt.Println(fail.String())
}
