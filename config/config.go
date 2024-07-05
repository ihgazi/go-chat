package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/pelletier/go-toml/v2"
	"log"
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
    SecretKey string    
    PostgresUrl string
}

// Using a singleton pattern to load the config only once and reduce read calls
var config *Config

func LoadConfig() Config {
	// returning config if already loaded
	if config != nil {
		return *config
	}

	// loading config if not already loaded
	config = &Config{}

	tomlFile, err := os.ReadFile("config.toml")
	if err != nil {
		panic(err)
	}

	if err = toml.Unmarshal(tomlFile, config); err != nil {
		panic(err)
	}

    err = godotenv.Load(".env")
    if err != nil {
        log.Fatalf("Error loading .env file: %s", err)
    }

    config.SecretKey = os.Getenv("SECRET_KEY")
    config.PostgresUrl = os.Getenv("POSTGRES_URL")	

    return *config
}

// Load .env configs
func LoadEnv() string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	secret_key := os.Getenv("SECRET_KEY")
	return secret_key
}
