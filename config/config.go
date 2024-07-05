package config

import (
	"os"

	"github.com/joho/godotenv"
	"log"
)

type Config struct {
	ServerHost  string
	ServerPort  string
	SecretKey   string
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

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	config.ServerHost = os.Getenv("SERVER_HOST")
	config.ServerPort = os.Getenv("SERVER_PORT")
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
