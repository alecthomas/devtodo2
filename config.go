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
	"sync"
)

var (
	instance *config
	once     sync.Once
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

func GetConfigInstance() *config {
	once.Do(func() {
		fgColors := map[Priority]string{
			VERYLOW:  BLUE,
			LOW:      CYAN,
			MEDIUM:   WHITE,
			HIGH:     BRIGHTYELLOW,
			VERYHIGH: BRIGHTRED,
		}
		bgColors := map[Priority]string{
			VERYLOW:  "BLACK",
			LOW:      "BLACK",
			MEDIUM:   "BLACK",
			HIGH:     "BLACK",
			VERYHIGH: "BLACK",
		}
		priority := "medium"
		graft := "root"
		file := ".todo2"
		legacyFile := ".todo"
		order := "priority"
		configFile := ".todorc"
		instance = &config{
			BGColors:   bgColors,
			FGColors:   fgColors,
			Priority:   priority,
			Graft:      graft,
			File:       file,
			LegacyFile: legacyFile,
			Order:      order,
			ConfigFile: configFile,
		}
	})
	return instance
}

type ConfigIO interface {
	Deserialize(reader io.Reader) (*MarshalableConfig, error)
}

type ConfigIOImpl struct {
}

func NewConfigIO() ConfigIO {
	return ConfigIOImpl{}
}

type MarshalableConfig struct {
	BGColors   map[string]string
	FGColors   map[string]string
	Priority   string
	Graft      string
	File       string
	LegacyFile string
	Order      string
}

func copyToConfigFromMarshalableConfig(marshalableConfig *MarshalableConfig, config *config) {
	if marshalableConfig.Priority != "" {
		config.Priority = marshalableConfig.Priority
	}
	if marshalableConfig.File != "" {
		config.File = marshalableConfig.File
	}
	if marshalableConfig.Graft != "" {
		config.Graft = marshalableConfig.Graft
	}
	if marshalableConfig.LegacyFile != "" {
		config.LegacyFile = marshalableConfig.LegacyFile
	}
	if marshalableConfig.Order != "" {
		config.Order = marshalableConfig.Order
	}
	if marshalableConfig.FGColors != nil {
		fgColors := map[Priority]string{}
		for key, value := range marshalableConfig.FGColors {
			fgColors[priorityMapFromString[key]] = value
		}
		config.FGColors = fgColors
	}
	if marshalableConfig.BGColors != nil {
		bgColors := map[Priority]string{}
		for key, value := range marshalableConfig.BGColors {
			bgColors[priorityMapFromString[key]] = value
		}
		config.BGColors = bgColors
	}
}

func (configIOImpl ConfigIOImpl) Deserialize(reader io.Reader) (marshalableConfig *MarshalableConfig, err error) {
	decoder := json.NewDecoder(reader)
	marshalableConfig = &MarshalableConfig{}
	err = decoder.Decode(marshalableConfig)
	return marshalableConfig, err
}

func loadConfigurationFile(config *config) (marshalableConfig *MarshalableConfig, err error) {
	if file, err := os.Open(config.ConfigFile); err == nil {
		defer file.Close()
		configIO := NewConfigIO()
		return configIO.Deserialize(file)
	}
	return nil, err
}

func copyToConfigFromCMDOptions(config *config) {
	kingpin.Flag("priority", "Priority of newly created tasks (veryhigh,high,medium,low,verylow).").Short('p').EnumVar(&config.Priority, "veryhigh", "high", "medium", "low", "verylow")
	kingpin.Flag("graft", "Task to graft new tasks to.").Short('g').Default("root").StringVar(&config.Graft)
	kingpin.Flag("file", "File to load task lists from.").StringVar(&config.File)
	kingpin.Flag("legacy-file", "File to load legacy task lists from.").StringVar(&config.LegacyFile)
	kingpin.Flag("order", "Specify display order of tasks (index,created,completed,text,priority,duration,done).").EnumVar(&config.Order, "index", "created", "completed", "text", "priority", "duration", "done")

	colors := []string{WHITE, BLUE, RED, CYAN, GREEN, YELLOW, BLACK, MAGENTA, BRIGHTWHITE, BRIGHTBLUE, BRIGHTRED, BRIGHTCYAN, BRIGHTGREEN, BRIGHTYELLOW, BRIGHTBLACK, BRIGHTMAGENTA, NOCOLOR}

	var veryLowFGColor = kingpin.Flag("verylowfgcolor", "Very low priority task texts foreground color.").Default(NOCOLOR).Enum(colors...)
	if *veryLowFGColor != NOCOLOR {
		config.FGColors[priorityMapFromString[verylow]] = *veryLowFGColor
	}

	var lowFGColor = kingpin.Flag("lowfgcolor", "Low priority task texts foreground color.").Default(NOCOLOR).Enum(colors...)
	if *lowFGColor != NOCOLOR {
		config.FGColors[priorityMapFromString[low]] = *lowFGColor
	}

	var mediumFGColor = kingpin.Flag("mediumfgcolor", "Medium priority task texts foreground color.").Default(NOCOLOR).Enum(colors...)
	if *mediumFGColor != NOCOLOR {
		config.FGColors[priorityMapFromString[medium]] = *mediumFGColor
	}

	var highFGColor = kingpin.Flag("highfgcolor", "High priority task texts foreground color.").Default(NOCOLOR).Enum(colors...)
	if *highFGColor != NOCOLOR {
		config.FGColors[priorityMapFromString[high]] = *highFGColor
	}

	var veryHighFGColor = kingpin.Flag("veryhighfgcolor", "Very high task texts foreground color.").Default(NOCOLOR).Enum(colors...)
	if *veryHighFGColor != NOCOLOR {
		config.FGColors[priorityMapFromString[veryhigh]] = *veryHighFGColor
	}

	var veryLowBGColor = kingpin.Flag("verylowbgcolor", "Very low priority task texts background color.").Default(NOCOLOR).Enum(colors...)
	if *veryLowBGColor != NOCOLOR {
		config.BGColors[priorityMapFromString[verylow]] = *veryLowFGColor
	}

	var lowBGColor = kingpin.Flag("lowbgcolor", "Low priority task texts background color.").Default(NOCOLOR).Enum(colors...)
	if *lowBGColor != NOCOLOR {
		config.BGColors[priorityMapFromString[low]] = *lowBGColor
	}

	var mediumBGColor = kingpin.Flag("mediumbgcolor", "Medium priority task texts background color.").Default(NOCOLOR).Enum(colors...)
	if *mediumBGColor != NOCOLOR {
		config.BGColors[priorityMapFromString[medium]] = *mediumBGColor
	}

	var highBGColor = kingpin.Flag("highbgcolor", "High priority task texts background color.").Default(NOCOLOR).Enum(colors...)
	if *highBGColor != NOCOLOR {
		config.BGColors[priorityMapFromString[high]] = *highBGColor
	}

	var veryHighBGColor = kingpin.Flag("veryhighbgcolor", "Very high priority task texts background color.").Default(NOCOLOR).Enum(colors...)
	if *veryHighBGColor != NOCOLOR {
		config.BGColors[priorityMapFromString[veryhigh]] = *veryHighBGColor
	}
}
