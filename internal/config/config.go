package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	AreaSize int `yaml:"area_size"`
	Match    struct {
		PartsLimit int `yaml:"parts_limit"`
		PartSize   int `yaml:"art_size"`
	} `yaml:"match"`
	BasePath string `yaml:"base_path"`
}

func NewConfig() *Config {
	config := new(Config)
	config.AreaSize = 100
	config.Match.PartsLimit = 200
	config.Match.PartSize = 100
	config.BasePath = os.ExpandEnv("$GOPATH/src/ants")

	f, err := os.Open("config.yml")
	if err != nil {
		log.Fatalf("Fatal error config file: %s \n", err)
	}

	d := yaml.NewDecoder(f)
	err = d.Decode(&config)
	if err != nil {
		log.Fatalf("Fatal error config file: %s \n", err)
	}

	return config
}