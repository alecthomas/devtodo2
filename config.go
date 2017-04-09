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
	"io"
)

type ConfigIO interface {
	Deserialize(reader io.Reader) (*Config, error)
}

func NewConfigIO() ConfigIO {
	return &ConfigIO{}
}

type marshalableConfig struct {
	Colors     map[string]string
	Priority   string
	Graft      string
	File       string
	LegacyFile string
	Order      string
}

func fromMarshalableConfig(marshalableConfig *marshalableConfig) (config Config) {
	config := &Config{
		Colors:     marshalableConfig.Colors,
		Priority:   marshalableConfig.Priority,
		Graft:      marshalableConfig.Graft,
		File:       marshalableConfig.File,
		LegacyFile: marshalableConfig.LegacyFile,
		Order:      marshalableConfig.Order,
	}
	return
}

func (configIO *ConfigIO) Deserialize(reader io.Reader) (config Config, err error) {
	decoder := json.NewDecoder(reader)
	marshalableConfig := &marshalableConfig{}
	err = decoder.Decode(&marshalableConfig)
	if err == nil {
		config = fromMarshalableConfig(marshalableConfig)
	}
	return
}

type Config struct {
	Colors     map[string]string
	Priority   string
	Graft      string
	File       string
	LegacyFile string
	Order      string
}

var (
	colors = map[string]string{
		"VeryLow":  "",
		"Low":      "",
		"Medium":   "",
		"High":     "",
		"VeryHigh": "",
	}
	priority   = "medium"
	graft      = "root"
	file       = ".todo2"
	legacyFile = ".todo"
	order      = "priority"
	configFile = ".todorc"

	config = &Config{
		Colors:     colors,
		Priority:   priority,
		Graft:      graft,
		File:       file,
		LegacyFile: legacyFile,
		Order:      order,
		ConfigFile: configFile,
	}
)
