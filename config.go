package main

import (
	"encoding/json"
	"fmt"
	"io"
)

type ColorRGB struct {
	R int
	G int
	B int
}

type ColorConfigRGB struct {
	VeryLow  *ColorRGB
	Low      *ColorRGB
	Medium   *ColorRGB
	High     *ColorRGB
	VeryHigh *ColorRGB
}

type ColorConfigHex struct {
	VeryLow  string
	Low      string
	Medium   string
	High     string
	VeryHigh string
}

//@TODO: Figure out how bright white fits into this,
//	 it has the same RGB as regular 8 byte white.
//@TODO: Figure out how background colors should be implemented.
var rgbToAnsiEscapeSequence = map[ColorRGB]string{
	ColorRGB{0, 0, 0}:       FGBLACK,
	ColorRGB{170, 0, 0}:     FGRED,
	ColorRGB{0, 170, 0}:     FGGREEN,
	ColorRGB{170, 85, 0}:    FGYELLOW,
	ColorRGB{0, 0, 170}:     FGBLUE,
	ColorRGB{170, 0, 170}:   FGMAGENTA,
	ColorRGB{0, 170, 170}:   FGCYAN,
	ColorRGB{255, 255, 255}: FGWHITE,
	ColorRGB{255, 85, 85}:   BRIGHT + FGRED,
	ColorRGB{85, 255, 85}:   BRIGHT + FGGREEN,
	ColorRGB{255, 255, 85}:  BRIGHT + FGYELLOW,
	ColorRGB{85, 85, 255}:   BRIGHT + FGBLUE,
	ColorRGB{255, 85, 255}:  BRIGHT + FGMAGENTA,
	ColorRGB{85, 255, 255}:  BRIGHT + FGCYAN,
}

func findClosestAnsiEscapeSequenceToColor(color *ColorRGB) (string, error) {
	for k, v := range rgbToAnsiEscapeSequence {
		fmt.Printf("Key: %+v \n", k)
		fmt.Printf("Value: %s \n", v)
	}
	return "", nil
}

type UnmarshalledConfig struct {
	ColorConfigRGB *ColorConfigRGB
	ColorConfigHex *ColorConfigHex
}

type ConfigIO interface {
	Deserialize(reader io.Reader) (*UnmarshalledConfig, error)
}

type ConfigJSONIO struct {
}

type Config struct {
	VeryLow  string
	Low      string
	Medium   string
	High     string
	VeryHigh string
}

func NewConfig(unmarshalledConfig *UnmarshalledConfig) (*Config, error) {
	config := new(Config)
	medium, err := findClosestAnsiEscapeSequenceToColor(unmarshalledConfig.ColorConfigRGB.Medium)
	if err != nil {
		fmt.Print(err)
	}
	config.Medium = medium
	return config, nil
}

func NewConfigJSONIO() (configJSONIO *ConfigJSONIO) {
	return new(ConfigJSONIO)
}

func (configJSONIO *ConfigJSONIO) Deserialize(reader io.Reader) (*UnmarshalledConfig, error) {
	decoder := json.NewDecoder(reader)
	cfg := &UnmarshalledConfig{}
	if err := decoder.Decode(&cfg); err == nil {
		return cfg, nil
	}
	return nil, nil
}
