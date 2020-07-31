package internal

import (
	util "github.com/Floor-Gang/utilpkg/config"
	"log"
)

type Config struct {
	Auth     string            `yaml:"auth_server"`
	Token    string            `yaml:"token"`
	Prefix   string            `yaml:"prefix"`
	Channels map[string]string `yaml:"linked_channels"`
}

const configPath = "./config.yml"

func GetConfig() Config {
	config := Config{
		Prefix:   ".link",
		Channels: make(map[string]string),
	}

	err := util.GetConfig(configPath, &config)

	if err != nil {
		log.Fatalln(err)
	}

	return config
}

func (config *Config) Save() {
	if err := util.Save(configPath, config); err != nil {
		log.Fatalln(err)
	}
}
