package config

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	Server   ServerConfig   `json:"server"`
	Database PostgresConfig `json:"database"`
}

const (
	configFile = ".config"
)

func DefaultConfig() Config {
	return Config{
		Server:   DefaultServerConfig(),
		Database: DefaultPostgresConfig(),
	}
}

func LoadConfig(configRequired bool) Config {
	f, err := os.Open(configFile)
	if err != nil {
		if configRequired {
			log.Fatalf("No config file found. Aborting application...")
		}

		log.Println("No config file found. Using default configuration...")
		return DefaultConfig()
	}

	dec := json.NewDecoder(f)

	var cfg Config
	err = dec.Decode(&cfg)
	if err != nil {
		log.Fatalf("failed to decode config file: %s", err)
	}

	log.Printf("Successfully loaded config file: %s", configFile)
	return cfg
}
