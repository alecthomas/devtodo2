package main

import (
	"encoding/json"
	"io"
)

type Color struct {
	R int
	G int
	B int
}

type ColorConfig struct {
	VeryLow  Color
	Low      Color
	Medium   Color
	High     Color
	VeryHigh Color
}

type Config struct {
	Color ColorConfig
}

type ConfigIO interface {
	Deserialize(reader io.Reader) (*Config, error)
}

type ConfigJSONIO struct {
}

func NewConfigJSONIO() (configJSONIO *ConfigJSONIO) {
	return new(ConfigJSONIO)
}

func (configJSONIO *ConfigJSONIO) Deserialize(reader io.Reader) (config *Config, err error) {
	decoder := json.NewDecoder(reader)
	cfg := &Config{}
	if err = decoder.Decode(&cfg); err == nil {
		return cfg, nil
	}
	return nil, err
}
