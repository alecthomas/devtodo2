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

func loadConfigurationFile(config *config) {
	if file, err := os.Open(config.ConfigFile); err == nil {
		defer file.Close()
		decoder := json.NewDecoder(file)
		decoder.Decode(config)
	}
}
