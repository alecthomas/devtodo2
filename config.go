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

func NewConfig() *config {
	currentUser, err := user.Current()
	configFile := ".todorc"
	if err == nil {
		configFile = strings.Join([]string{currentUser.HomeDir, configFile}, "/")
	}
	return &config{
		FGColors: map[Priority]string{
			VERYLOW:  BLUE,
			LOW:      CYAN,
			MEDIUM:   WHITE,
			HIGH:     BRIGHTYELLOW,
			VERYHIGH: BRIGHTRED,
		},
		BGColors: map[Priority]string{
			VERYLOW:  NOCOLOR,
			LOW:      NOCOLOR,
			MEDIUM:   NOCOLOR,
			HIGH:     NOCOLOR,
			VERYHIGH: NOCOLOR,
		},
		Priority:   medium,
		Graft:      "root",
		File:       ".todo2",
		LegacyFile: ".todo",
		Order:      "priority",
		ConfigFile: configFile,
	}
}

func (priority *Priority) UnmarshalText(data []byte) error {
	*priority = PriorityFromString(string(data))
	return nil
}

func loadConfigurationFile(config *config) (err error) {
	if file, err := os.Open(config.ConfigFile); err == nil {
		defer file.Close()
		decoder := json.NewDecoder(file)
		decodeErr := decoder.Decode(&config)
		if decodeErr != nil {
			return decodeErr
		}
	}
	return
}

func loadConfigCMD(config *config) {
	colors := []string{WHITE, BLUE, RED, CYAN, GREEN, YELLOW, BLACK, MAGENTA, BRIGHTWHITE, BRIGHTBLUE, BRIGHTRED, BRIGHTCYAN, BRIGHTGREEN, BRIGHTYELLOW, BRIGHTBLACK, BRIGHTMAGENTA, NOCOLOR}
	kingpin.Flag("priority", "Priority of newly created tasks (veryhigh,high,medium,low,verylow).").Short('p').EnumVar(&config.Priority, veryhigh, high, medium, low, verylow)
	kingpin.Flag("graft", "Task to graft new tasks to.").Short('g').Default("root").StringVar(&config.Graft)
	kingpin.Flag("file", "File to load task lists from.").StringVar(&config.File)
	kingpin.Flag("legacy-file", "File to load legacy task lists from.").StringVar(&config.LegacyFile)
	kingpin.Flag("order", "Specify display order of tasks (index,created,completed,text,priority,duration,done).").EnumVar(&config.Order, "index", "created", "completed", "text", "priority", "duration", "done")

	fgColorVeryLow := kingpin.Flag("verylowfgcolor", "Very low priority task texts foreground color.").Default(config.FGColors[priorityMapFromString[verylow]]).Enum(colors...)
	fgColorLow := kingpin.Flag("lowfgcolor", "Low priority task texts foreground color.").Default(config.FGColors[priorityMapFromString[low]]).Enum(colors...)
	fgColorMedium := kingpin.Flag("mediumfgcolor", "Medium priority task texts foreground color.").Default(config.FGColors[priorityMapFromString[medium]]).Enum(colors...)
	fgColorHigh := kingpin.Flag("highfgcolor", "High priority task texts foreground color.").Default(config.FGColors[priorityMapFromString[high]]).Enum(colors...)
	fgColorVeryHigh := kingpin.Flag("veryhighfgcolor", "Very high task texts foreground color.").Default(config.FGColors[priorityMapFromString[veryhigh]]).Enum(colors...)
	bgColorVeryLow := kingpin.Flag("verylowbgcolor", "Very low priority task texts background color.").Default(config.BGColors[priorityMapFromString[verylow]]).Enum(colors...)
	bgColorLow := kingpin.Flag("lowbgcolor", "Low priority task texts background color.").Default(config.BGColors[priorityMapFromString[low]]).Enum(colors...)
	bgColorMedium := kingpin.Flag("mediumbgcolor", "Medium priority task texts background color.").Default(config.BGColors[priorityMapFromString[medium]]).Enum(colors...)
	bgColorHigh := kingpin.Flag("highbgcolor", "High priority task texts background color.").Default(config.BGColors[priorityMapFromString[high]]).Enum(colors...)
	bgColorVeryHigh := kingpin.Flag("veryhighbgcolor", "Very high priority task texts background color.").Default(config.BGColors[priorityMapFromString[veryhigh]]).Enum(colors...)
	kingpin.Parse()
	config.FGColors[priorityMapFromString[verylow]] = *fgColorVeryLow
	config.FGColors[priorityMapFromString[low]] = *fgColorLow
	config.FGColors[priorityMapFromString[medium]] = *fgColorMedium
	config.FGColors[priorityMapFromString[high]] = *fgColorHigh
	config.FGColors[priorityMapFromString[veryhigh]] = *fgColorVeryHigh
	config.BGColors[priorityMapFromString[verylow]] = *bgColorVeryLow
	config.BGColors[priorityMapFromString[low]] = *bgColorLow
	config.BGColors[priorityMapFromString[medium]] = *bgColorMedium
	config.BGColors[priorityMapFromString[high]] = *bgColorHigh
	config.BGColors[priorityMapFromString[veryhigh]] = *bgColorVeryHigh
}
