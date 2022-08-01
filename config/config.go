package config

import (
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

type MapConfig struct {
	Key string  `yaml:"key"`
	Lat float64 `yaml:"lat"`
	Lng float64 `yaml:"lng"`
}

type TimeConfig struct {
	Timezone string `yaml:"timezone"`
}

type Config struct {
	Token string     `yaml:"token"`
	Map   MapConfig  `yaml:"map"`
	Time  TimeConfig `yaml:"time"`
}

var Conf Config

func Init(path string) (err error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDONLY, 0644)
	if err != nil {
		return
	}

	content, err := io.ReadAll(f)
	if err != nil {
		return
	}

	return yaml.Unmarshal(content, &Conf)
}
