package structs

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Token   string `yaml:"token"`
	GuildID string `yaml:"guildID"`
	Status  string `yaml:"status"`
}

func (c *Config) Load() {
	var filename string

	if os.Getenv("PROD") == "1" {
		filename = "config.yaml"
	} else {
		filename = "dev-config.yaml"
	}

	file, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("Failed to load %s\n", filename)
	}

	err = yaml.Unmarshal(file, &c)
	if err != nil {
		log.Fatal("Error parsing yaml")
	}
}
