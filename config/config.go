package config

import (
	"os"

	"github.com/pelletier/go-toml/v2"
)

type Config struct {
	GatewayConfig struct {
		Host string `toml:"HOST"`
		Port int    `toml:"PORT"`
	} `toml:"gateway_dev"`
	ServerConfig struct {
		Host string `toml:"HOST"`
		Port int    `toml:"PORT"`
	} `toml:"server_dev"`
}

// Using a singleton pattern to load the config only once and reduce read calls
var config *Config

func LoadConfig(path string) Config {
	// returning config if already loaded
	if config != nil {
		return *config
	}

	// loading config if not already loaded
	config = &Config{}

	tomlFile, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	if err = toml.Unmarshal(tomlFile, config); err != nil {
		panic(err)
	}

	return *config
}
