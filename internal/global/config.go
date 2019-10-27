package global

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

var Config struct {
	AreaSize int `yaml:"area_size"`
	Match    struct {
		PartsLimit int `yaml:"parts_limit"`
		PartSize   int `yaml:"art_size"`
	} `yaml:"match"`
	BasePath string `yaml:"base_path"`
}

func InitConfig() {
	Config.AreaSize = 100
	Config.Match.PartsLimit = 200
	Config.Match.PartSize = 100
	Config.BasePath = os.ExpandEnv("$GOPATH/src/ants")

	f, err := os.Open("config.yml")
	if err != nil {
		log.Fatalf("Fatal error config file: %s \n", err)
	}

	d := yaml.NewDecoder(f)
	err = d.Decode(&Config)
	if err != nil {
		log.Fatalf("Fatal error config file: %s \n", err)
	}
}
