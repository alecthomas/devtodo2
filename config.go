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
	"encoding/json"
	"gopkg.in/alecthomas/kingpin.v2"
	"io"
	"os"
	"os/user"
	"strings"
)

type config struct {
	FGColors   map[Priority]string
	BGColors   map[Priority]string
	Priority   string
	Graft      string
	File       string
	LegacyFile string
	Order      string
	ConfigFile string
}

// Map for priority level to ANSI colour
var colourPriorityMap = map[Priority]string{
	VERYHIGH: BRIGHT + FGRED,
	HIGH:     BRIGHT + FGYELLOW,
	MEDIUM:   FGWHITE,
	LOW:      FGCYAN,
	VERYLOW:  FGBLUE,
}

func NewConfig() *config {
	currentUser, err := user.Current()
	configFile := ".todorc"
	if err == nil {
		configFile = strings.Join([]string{currentUser.HomeDir, configFile}, "/")
	}
	return &config{
		BGColors: map[Priority]string{
			VERYLOW:  NOCOLOR,
			LOW:      NOCOLOR,
			MEDIUM:   NOCOLOR,
			HIGH:     NOCOLOR,
			VERYHIGH: NOCOLOR,
		},
		FGColors: map[Priority]string{
			VERYLOW:  BLUE,
			LOW:      CYAN,
			MEDIUM:   WHITE,
			HIGH:     BRIGHTYELLOW,
			VERYHIGH: BRIGHTRED,
		},
		Priority:   "medium",
		Graft:      "root",
		File:       ".todo2",
		LegacyFile: ".todo",
		Order:      "priority",
		ConfigFile: configFile,
	}
}

type ConfigIO interface {
	Deserialize(reader io.Reader) (*config, error)
}

type ConfigIOImpl struct{}

func NewConfigIO() ConfigIO {
	return ConfigIOImpl{}
}

func (priority *Priority) Unmarshaler(bytes []byte) (err error) {
	p := int(0)
	if err = json.Unmarshal(bytes, &p); err == nil {
		*priority = Priority(p)

	}
	return
}

func loadConfigurationFile(config *config) (err error) {
	if file, err := os.Open(".todorc"); err == nil {
		defer file.Close()
		configIO := NewConfigIO()
		decoder := json.NewDecoder(file)
		decodeErr = decoder.Decode(&config)
		if decodeErr != nil {
			return decodeErr
		}
	}
	return
}

func loadConfigCMD(config *config) {
	kingpin.Flag("priority", "Priority of newly created tasks (veryhigh,high,medium,low,verylow).").Short('p').EnumVar(&config.Priority, veryhigh, high, medium, low, verylow)
	kingpin.Flag("graft", "Task to graft new tasks to.").Short('g').Default("root").StringVar(&config.Graft)
	kingpin.Flag("file", "File to load task lists from.").StringVar(&config.File)
	kingpin.Flag("legacy-file", "File to load legacy task lists from.").StringVar(&config.LegacyFile)
	kingpin.Flag("order", "Specify display order of tasks (index,created,completed,text,priority,duration,done).").EnumVar(&config.Order, "index", "created", "completed", "text", "priority", "duration", "done")

	colors := []string{WHITE, BLUE, RED, CYAN, GREEN, YELLOW, BLACK, MAGENTA, BRIGHTWHITE, BRIGHTBLUE, BRIGHTRED, BRIGHTCYAN, BRIGHTGREEN, BRIGHTYELLOW, BRIGHTBLACK, BRIGHTMAGENTA, NOCOLOR}

	var veryLowFGColor = kingpin.Flag("verylowfgcolor", "Very low priority task texts foreground color.").Enum(colors...)
	var lowFGColor = kingpin.Flag("lowfgcolor", "Low priority task texts foreground color.").Enum(colors...)
	var mediumFGColor = kingpin.Flag("mediumfgcolor", "Medium priority task texts foreground color.").Enum(colors...)
	var highFGColor = kingpin.Flag("highfgcolor", "High priority task texts foreground color.").Enum(colors...)
	var veryHighFGColor = kingpin.Flag("veryhighfgcolor", "Very high task texts foreground color.").Enum(colors...)
	var veryLowBGColor = kingpin.Flag("verylowbgcolor", "Very low priority task texts background color.").Enum(colors...)
	var lowBGColor = kingpin.Flag("lowbgcolor", "Low priority task texts background color.").Enum(colors...)
	var mediumBGColor = kingpin.Flag("mediumbgcolor", "Medium priority task texts background color.").Enum(colors...)
	var highBGColor = kingpin.Flag("highbgcolor", "High priority task texts background color.").Enum(colors...)
	var veryHighBGColor = kingpin.Flag("veryhighbgcolor", "Very high priority task texts background color.").Enum(colors...)
	kingpin.Parse()

	if *veryLowFGColor != NOCOLOR {
		config.FGColors[priorityMapFromString[verylow]] = *veryLowFGColor
	}

	if *lowFGColor != NOCOLOR {
		config.FGColors[priorityMapFromString[low]] = *lowFGColor
	}

	if *mediumFGColor != NOCOLOR {
		config.FGColors[priorityMapFromString[medium]] = *mediumFGColor
	}

	if *highFGColor != NOCOLOR {
		config.FGColors[priorityMapFromString[high]] = *highFGColor
	}

	if *veryHighFGColor != NOCOLOR {
		config.FGColors[priorityMapFromString[veryhigh]] = *veryHighFGColor
	}

	if *veryLowBGColor != NOCOLOR {
		config.BGColors[priorityMapFromString[verylow]] = *veryLowFGColor
	}

	if *lowBGColor != NOCOLOR {
		config.BGColors[priorityMapFromString[low]] = *lowBGColor
	}

	if *mediumBGColor != NOCOLOR {
		config.BGColors[priorityMapFromString[medium]] = *mediumBGColor
	}

	if *highBGColor != NOCOLOR {
		config.BGColors[priorityMapFromString[high]] = *highBGColor
	}

	if *veryHighBGColor != NOCOLOR {
		config.BGColors[priorityMapFromString[veryhigh]] = *veryHighBGColor
	}

}
